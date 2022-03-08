from classes.aws_client import AwsClient
from install.hacc_credentials import create_hacc_profile, delete_hacc_profile
import install.hacc_policies
import hacc_vars
import json

import logging

logger=logging.getLogger()



## VaultInstallation Object:
##      idempotently install or eradicate Vault AWS components
##
## Attributes:
##      AwsClient object
##      
## 
## Methods:
##      cmk_exists, private
##      create_cmk_with_alias
##      delete_cmk
##      user_exists, private
##      create_user_with_policy
##      delete_user
##      
##      
##      
class VaultInstallation:

    ## Checks if KMS key exists with expected alias from config file
    ## Returns KMS ARN if key exists, otherwise False
    def __cmk_exists(self):
        kms_alias = hacc_vars.aws_hacc_kms_alias

        try:
            hacc_kms_arn = self.aws_client.call(
                                'kms', 'describe_key', 
                                KeyId = f'alias/{kms_alias}'
                            )['KeyMetadata']['Arn']
            return hacc_kms_arn
        except:
            return False



    ## Accepts key id to delete, or uses expected alias by default
    ## Returns True if delete successful or key doesn't exist
    ## Unsets cmk attribute
    ## Returns False if failed to delete key/alias
    def delete_cmk(self, key_id=None):

        ## Look for existing key using expected alias
        if not key_id:
            key_id = self.__cmk_exists()
            key_alias = hacc_vars.aws_hacc_kms_alias

            if key_id:
                alias_delete_res = self.aws_client.call('kms', 'delete_alias',
                                                    AliasName = f'alias/{key_alias}'
                                                )
            if not alias_delete_res:
                logger.info('Failed to delete alias for Vault KMS key')
                return False
            logger.info('Successfully deleted alias for Vault KMS key')


        key_delete_res = self.aws_client.call('kms', 'schedule_key_deletion', 
                                            KeyId = key_id,
                                            PendingWindowInDays = 7
                                        )
        if not key_delete_res:
            logger.info('Failed to delete KMS key')
            return False

        logger.info('Successfully deleted Vault KMS key')
        self.cmk = None
        return True


    ## Creates or confirms KMS exists for Vault 
    ## Returns True if key is present and sets cmk attribute with ARN
    ## Returns False if key/alias failed to create
    def create_cmk_with_alias(self):
        logger.info('Checking for existing KMS key')

        ## Check if key already exists
        arn = self.__cmk_exists()
        if arn:
            self.cmk = arn
            logger.info('Found existing KMS key and alias for Vault')
            return True

        ## Otherwise create new key and alias
        hacc_kms = self.aws_client.call('kms', 'create_key', 
                                        KeyUsage = 'ENCRYPT_DECRYPT', 
                                        KeySpec = 'SYMMETRIC_DEFAULT'
                                    )
        if not hacc_kms:
            logger.info('Failed to create new symmetric KMS key')
            return False
        logger.info('Successfully created new symmetric KMS key for Vault')


        key_name = hacc_vars.aws_hacc_kms_alias
        key_id = hacc_kms['KeyMetadata']['Arn']

        hacc_alias = self.aws_client.call('kms', 'create_alias', 
                                        AliasName = f'alias/{key_name}',
                                        TargetKeyId = key_id
                                    )

        ## If alias creates properly we're done
        if hacc_alias:
            logger.info('Successfully created alias for KMS key')
            self.cmk = key_id
            return True

        ## If alias fails but key exists, clean up first
        logger.info('Failed to create alias for KMS key')
        logger.info('Cleaning up KMS key with no alias')
        self.delete_cmk(key_id=key_id)
        return False


    ## Checks if expected IAM user name with appropriate policy exists for Vault
    ## Returns user ARN if exists, False otherwise
    def __user_exists(self):
        username = hacc_vars.aws_hacc_uname

        try:
            hacc_user = self.aws_client.call(
                                'iam', 'get_user', 
                                UserName = f'{username}'
                            )['User']['Arn']
            return hacc_user
        except:
            return False



    ## Returns True if delete successful or key doesn't exist and unsets user attribute
    ## Returns False if failed to delete user and remove creds from AWS credentials file
    def delete_user(self):
        if not self.__user_exists():
            logger.info('No existing IAM user found')
            return True

        ## Delete user from AWS
        username = hacc_vars.aws_hacc_uname
        user_delete_res = self.aws_client.call(
                                'iam', 'delete_user', 
                                UserName = f'{username}'
                            )
        
        if not user_delete_res:
            logger.info('Failed to delete Vault IAM user')
            return False

        logger.debug('Successfully deleted Vault IAM user')
        self.user = None

        ## Delete saved user credentials locally
        deleted_creds = delete_hacc_profile()

        if deleted_creds:
            logger.info('Successfully deleted saved Vault user credentials')
            return True
        
        logger.info('User deleted but failed to clean up old .aws/credentials file')
        return False


    ## Creates or confirms IAM user exists for Vault 
    ## Returns True if user is present
    ## Sets user attribute with IAM ARN
    ## Returns False if user/policy failed to create
    def create_user_with_policy(self):
        logger.info('Checking for existing IAM user')

        ## Check if user already exists
        user_arn = self.__user_exists()
        if user_arn:
            logger.info('Existing IAM user found')

        else:
            ## Create IAM user if not
            hacc_user = self.aws_client.call(
                                'iam', 'create_user',
                                UserName=hacc_vars.aws_hacc_uname
                            )
            if not hacc_user:
                logger.info('Failed to create IAM user for Vault')
                return False

            logger.info('Successfully created IAM user for Vault')
            user_arn = hacc_user['User']['Arn']


        region = hacc_vars.aws_hacc_region
        account = self.aws_account_id
        path = hacc_vars.aws_hacc_param_path
        vault_path_arn = f'arn:aws:ssm:{region}:{account}:parameter/{path}'
        vault_key_arn = self.cmk

        ## Idempotently put user policy for new/existing user
        user_perms = json.loads(install.hacc_policies.VAULT_IAM_PERMS % (vault_path_arn, vault_path_arn+'/*', vault_key_arn))
        hacc_policy = self.aws_client.call(
                                'iam', 'put_user_policy',
                                UserName=hacc_vars.aws_hacc_uname,
                                PolicyName=hacc_vars.aws_hacc_iam_policy,
                                PolicyDocument=json.dumps(user_perms)
                            )   

        ## If policy failed to create, clean up
        if not hacc_policy:
            logger.info('Failed to create IAM user policy') 
            logger.info('Cleaning up Vault user with no policy')
            self.delete_user()
            return False

        logger.info('Successfully updated IAM user policy for Vault')


        hacc_creds = self.aws_client.call(
                        'iam', 'create_access_key',
                        UserName=hacc_vars.aws_hacc_uname
                    )

        ## If credentials failed to create, clean up
        if not hacc_creds:
            logger.info('Failed to create credentials for user')
            logger.info('Cleaning up Vault user with no credentials')
            self.delete_user()
            return False

        created_profile = create_hacc_profile(
                            hacc_creds['AccessKey']['AccessKeyId'], 
                            hacc_creds['AccessKey']['SecretAccessKey']
                        )

        ## If credentials failed to save, clean up
        if not created_profile:
            logger.info('Failed to save Vault user credentials')
            logger.info('Cleaning up user without saved credentials')
            self.delete_user()
            return False

        logger.info('Successfully saved Vault IAM user credentials')
        self.user = user_arn
        return True


               

    def __init__(self):
        self.cmk = None
        self.user = None

        ## If no SCP (single-account setup), use current AWS creds to install Vault
        if not hacc_vars.create_scp:
            self.aws_client = AwsClient(client_type='mgmt', create_scp=False)
            self.aws_account_id = self.aws_client.call('sts', 'get_caller_identity')['Account']

        else:
            self.aws_client = AwsClient(client_type='mgmt', create_scp=True)
            self.aws_account_id = self.aws_client.call('sts', 'get_caller_identity')['Account']
                