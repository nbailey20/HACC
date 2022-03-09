from classes.aws_client import AwsClient
from install.hacc_credentials import create_hacc_profile, delete_hacc_profile
import install.hacc_policies
import hacc_vars
import json, time

import logging

logger=logging.getLogger()



## VaultInstallation Object:
##      idempotently install or eradicate Vault AWS components
##
## Attributes:
##      aws_client, AwsClient object
##      cmk, string
##      user, string
##      scp, string
##      aws_account_id, string
##      aws_user_arn, string (for SCP only)
##      aws_ssm_path_arn, string (for SCP only)
## 
## Methods:
##      cmk_exists, private
##      create_cmk_with_alias
##      delete_cmk
##      user_exists, private
##      create_user_with_policy
##      delete_user
##      scp_exists, private
##      create_vault_scp
##      delete_vault_scp
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
        hacc_kms_res = self.aws_client.call('kms', 'create_key', 
                                        KeyUsage = 'ENCRYPT_DECRYPT', 
                                        KeySpec = 'SYMMETRIC_DEFAULT'
                                    )
        if not hacc_kms_res:
            logger.info('Failed to create new symmetric KMS key')
            return False
        logger.info('Successfully created new symmetric KMS key for Vault')


        key_name = hacc_vars.aws_hacc_kms_alias
        key_id = hacc_kms_res['KeyMetadata']['Arn']

        hacc_alias_res = self.aws_client.call('kms', 'create_alias', 
                                        AliasName = f'alias/{key_name}',
                                        TargetKeyId = key_id
                                    )

        ## If alias creates properly we're done
        if hacc_alias_res:
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
                            )
            return hacc_user['User']['Arn']
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

        logger.info('Successfully deleted Vault IAM user')
        self.user = None

        ## Delete saved user credentials locally
        deleted_creds = delete_hacc_profile()

        if not deleted_creds:
            logger.info('User deleted but failed to clean up deprecated Hacc profile in AWS credentials/config files')
            return False
        
        logger.info('Successfully deleted saved Vault user credentials')
        self.user = None
        return True



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


        ## Idempotently put user policy for new/existing user
        user_perms = json.loads(install.hacc_policies.VAULT_IAM_PERMS % (self.aws_ssm_path_arn, self.aws_ssm_path_arn+'/*', self.cmk))
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


    
    ## Returns SCP ID if Vault SCP exists for account, False otherwise
    def __scp_exists(self):
        hacc_scp_name = hacc_vars.aws_hacc_scp

        vault_account_scp_obj = self.aws_client.call(
                                        'org', 'list_policies_for_target',
                                        TargetId = self.aws_account_id,
                                        Filter = 'SERVICE_CONTROL_POLICY'
                                    )

        if not vault_account_scp_obj:
            logger.info('Failed to retrieve existing SCPs for Vault account')
            return False

        ## Check if policy with expected name exists
        for scp in vault_account_scp_obj['Policies']:
            name = scp['Name']
            if name == hacc_scp_name:
                return scp['Id']

        ## If more than 10 SCPs, get more to check
        while 'NextToken' in vault_account_scp_obj:

            vault_account_scp_obj = self.aws_client.call(
                                            'org', 'list_policies_for_target',
                                            TargetId = self.aws_account_id,
                                            Filter = 'SERVICE_CONTROL_POLICY',
                                            NextToken = vault_account_scp_obj['NextToken']
                                        )
            if not vault_account_scp_obj:
                logger.info('Failed to retrieve all SCPs for Vault account')
                return False

            for scp in vault_account_scp_obj['Policies']:
                name = scp['Name']
                if name == hacc_scp_name:
                    return scp['Id']

        return False



    ## Returns True if delete successful or SCP doesn't exist and unsets scp attribute
    ## Returns False if failed to delete SCP
    def delete_scp(self):
        hacc_scp_id = self.__scp_exists()

        if not hacc_scp_id:
            logger.info('No existing SCP found')
            return True

        detach_scp_res = self.aws_client.call(
                                    'org', 'detach_policy',
                                    PolicyId = hacc_scp_id,
                                    TargetId = self.aws_account_id
                                )

        if not detach_scp_res:
            logger.info('Failed to detach SCP from Vault account')
            return False

        delete_scp_res = self.aws_client.call(
                                    'org', 'delete_policy',
                                    PolicyId = hacc_scp_id
                                )

        if not delete_scp_res:
            logger.info('Failed to delete SCP in Vault account')
            return False

        # Give SCP deletion time to take effect ('immediate' per AWS docs is not good enough :)
        logger.info('Waiting 10 seconds for SCP to fully delete')
        time.sleep(10)

        logger.info('Successfully deleted SCP from Vault account')
        self.scp = None
        return True
        
        


    ## Creates or confirms SCP exists for Vault 
    ## Returns True if SCP is present
    ## Sets scp attribute with SCP ID
    ## Returns False if SCP failed to create/attach to Vault account
    def create_scp(self):
        hacc_scp_id = self.__scp_exists()

        if hacc_scp_id:
            logger.info('Existing SCP found for Vault account')
            self.scp = hacc_scp_id
            return True

        scp_policy = json.loads(install.hacc_policies.VAULT_SCP % 
                                                (
                                                    self.aws_ssm_path_arn, 
                                                    self.aws_ssm_path_arn+'/*',
                                                    self.aws_user_arn,
                                                    self.cmk,
                                                    self.aws_user_arn)
                                                )
        create_scp_res = self.aws_client.call(
                                    'org', 'create_policy',
                                    Content = json.dumps(scp_policy),
                                    Name = hacc_vars.aws_hacc_scp,
                                    Description = 'SCP for HACC',
                                    Type = 'SERVICE_CONTROL_POLICY'
                                )
        if not create_scp_res:
            logger.debug('Failed to create SCP in AWS organization')
            return False


        attach_scp_res = self.aws_client.call(
                                    'org', 'attach_policy',
                                    PolicyId = hacc_scp_id,
                                    TargetId = self.aws_account_id
                                )
        if not attach_scp_res:
            logger.debug('Failed to attach SCP to Vault account')
            logger.debug('Cleaning up unattached SCP in AWS organization')
            self.delete_scp()
            return False

        logger.debug('Successfully attached SCP to Vault account')
        self.scp = hacc_scp_id
        return True

               

    def __init__(self):
        self.cmk = self.__cmk_exists()
        self.user = self.__user_exists()

        ## If no SCP (single-account setup), use current AWS creds to install Vault
        if not hacc_vars.create_scp:
            self.aws_client = AwsClient(client_type='mgmt', create_scp=False)
            self.aws_account_id = self.aws_client.call('sts', 'get_caller_identity')['Account']

        else:
            self.scp = self.__scp_exists()
            self.aws_client = AwsClient(client_type='mgmt', create_scp=True)

            curr_account_obj = self.aws_client.call('sts', 'get_caller_identity')
            self.aws_account_id = curr_account_obj['Account']
            self.aws_user_arn = curr_account_obj['Arn']

            region = hacc_vars.aws_hacc_region
            path = hacc_vars.aws_hacc_param_path
            self.aws_ssm_path_arn = f'arn:aws:ssm:{region}:{self.aws_account_id}:parameter/{path}'
                    