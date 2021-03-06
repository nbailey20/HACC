from classes.vault_eradicator import VaultEradicator
from install.hacc_credentials import create_hacc_profile
import install.hacc_policies
import json


## VaultInstaller Object:
##      idempotently install Vault AWS components
##
## Methods:
##      __create_cmk
##      __create_alias
##      create_cmk_with_alias
##      __create_iam_user
##      __put_iam_policy
##      __create_iam_access_key
##      create_iam_user_with_policy
##      __create_org_policy
##      __attach_org_policy
##      create_scp
##
class VaultInstaller(VaultEradicator):
    def __init__(self, config):
        super().__init__()

        if config['create_scp']:
            cross_account_user = config['aws_hacc_uname']
            self.aws_user_arn = f'arn:aws:iam::{self.aws_account_id}:user/{cross_account_user}'

        region = config['aws_hacc_region']
        path = config['aws_hacc_param_path']
        self.aws_ssm_path_arn = f'arn:aws:ssm:{region}:{self.aws_account_id}:parameter/{path}' 



    ## Creates new CMK and returns ID, False on fail
    def __create_cmk(self):
        hacc_kms_res = self.aws_client.call(
                            'kms', 'create_key',
                            KeyUsage = 'ENCRYPT_DECRYPT',
                            KeySpec = 'SYMMETRIC_DEFAULT'
                        )
        if not hacc_kms_res:
            return False
        return hacc_kms_res['KeyMetadata']['Arn']


    ## Creates new key alias with provided value, returns False on fail
    def __create_alias(self, name, id):
        hacc_alias_res = self.aws_client.call(
                                'kms', 'create_alias',
                                AliasName = f'alias/{name}',
                                TargetKeyId = id
                            )
        return hacc_alias_res


    ## Creates or confirms KMS exists for Vault 
    ## Returns True if key is present and sets CMK component parent attribute with ARN
    ## Returns False if key/alias failed to create
    def create_cmk_with_alias(self, kms_alias):
        print('Checking for existing Vault key')

        if self.cmk:
            print('Found existing Vault key with alias')
            return True
        print('No existing Vault key found, creating...')

        ## Create new key
        key_id = self.__create_cmk()
        if not key_id:
            print('Failed to create new symmetric KMS key')
            return False

        ## Create new alias
        hacc_alias_res = self.__create_alias(kms_alias, key_id)
        if hacc_alias_res:
            self.cmk = key_id
            return True

        ## If alias fails but key exists, clean up for future
        print('Failed to create alias for Vault key')
        print('  Cleaning up Vault key component for future install attempt')
        self.delete_cmk_with_alias()
        return False



    ## Create new IAM user with provided name, returns user ARN or False on fail
    def __create_iam_user(self, name):
        hacc_user = self.aws_client.call(
                        'iam', 'create_user',
                        UserName = name
                    )
        if not hacc_user:
            return False
        return hacc_user['User']['Arn']


    ## Idempotently put IAM policy with provided name, returns False on fail
    def __put_iam_policy(self, user, policy_name, policy_doc):
        hacc_policy = self.aws_client.call(
                            'iam', 'put_user_policy',
                            UserName = user,
                            PolicyName = policy_name,
                            PolicyDocument = policy_doc
                        )
        return hacc_policy


    ## Create new IAM access key for provided user, returns tuple key data or False on fail
    def __create_iam_access_key(self, user):
        hacc_creds = self.aws_client.call(
                        'iam', 'create_access_key',
                        UserName = user
                    )
        if not hacc_creds:
            return (False, False)
        return ( hacc_creds['AccessKey']['AccessKeyId'],
                 hacc_creds['AccessKey']['SecretAccessKey'] )


    ## Creates or confirms IAM user exists for Vault 
    ## Returns True if user is present
    ## Sets user component parent attribute with IAM ARN
    ## Returns False if user/policy failed to create
    def create_user_with_policy(self, username, policy_name):
        print('Checking for existing IAM user')

        if self.user:
            print('Existing IAM user found')
            return True
        print('No existing user found, creating...')

        ## Create IAM user
        user_arn = self.__create_iam_user(username)
        if not user_arn:
            print('Failed to create IAM user for Vault')
            return False

        ## Put IAM policy for user
        user_perms = json.loads(install.hacc_policies.VAULT_IAM_PERMS % 
                                    (self.aws_ssm_path_arn,
                                    self.aws_ssm_path_arn+'/*',
                                    self.cmk)
                        )
        hacc_policy = self.__put_iam_policy(
                            username,
                            policy_name,
                            json.dumps(user_perms)
                        )
        if not hacc_policy:
            print('Failed to create IAM user policy') 
            print('  Cleaning up Vault IAM component for future install attempt')
            self.delete_iam_user_with_policy()
            return False

        ## Create IAM access key for user
        access_key, secret_key = self.__create_iam_access_key(username)
        if not access_key:
            print('Failed to create credentials for user')
            print('  Cleaning up Vault user with no credentials')
            self.delete_iam_user_with_policy()
            return False

        ## Create new HACC profile in local AWS credentials files
        created_profile = create_hacc_profile(access_key, secret_key)
        if not created_profile:
            print('Failed to save Vault user credentials')
            print('  Cleaning up user without saved credentials')
            self.delete_iam_user_with_policy()
            return False

        print('Successfully created Vault IAM user and saved credentials')
        self.user = user_arn
        return True



    ## Create new organizational security policy, returns policy ID or False on fail
    def __create_org_policy(self, name, content):
        create_scp_res = self.aws_client.call(
                                'org', 'create_policy',
                                Name = name,
                                Content = content,
                                Description = 'SCP for HACC',
                                Type = 'SERVICE_CONTROL_POLICY'
                            )
        if not create_scp_res:
            return False
        return create_scp_res['Policy']['PolicySummary']['Id']


    ## Attach provided organizational policy to target account, returns False on fail
    def __attach_org_policy(self, policy, account):
        attach_scp_res = self.aws_client.call(
                                'org', 'attach_policy',
                                PolicyId = policy,
                                TargetId = account
                            )
        return attach_scp_res


    ## Creates or confirms SCP exists for Vault 
    ## Returns True if SCP is present
    ## Sets scp component parent attribute with SCP ID
    ## Returns False if SCP failed to create/attach to Vault account
    def create_scp(self, scp_name):
        print('Checking for existing Vault SCP')

        if self.scp:
            print('Existing SCP found for Vault account')
            return True
        print('No existing SCP found, creating Vault SCP...')

        ## Create policy in org
        scp_policy = json.loads(install.hacc_policies.VAULT_SCP % 
                                    (self.aws_ssm_path_arn,
                                     self.aws_ssm_path_arn+'/*',
                                     self.aws_user_arn,
                                     self.cmk,
                                     self.aws_user_arn)
                                )
        policy_id = self.__create_org_policy(scp_name, json.dumps(scp_policy))
        if not policy_id:
            print('Failed to create SCP in AWS organization')
            return False

        attach_scp_res = self.__attach_org_policy(policy_id, self.aws_account_id)
        if not attach_scp_res:
            print('Failed to attach SCP to Vault account')
            print('  Cleaning up unattached SCP in AWS organization')
            self.delete_scp()
            return False

        print('Successfully attached SCP to Vault account')
        self.scp = policy_id
        return True