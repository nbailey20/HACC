from classes.vault_components import VaultComponents
from install.hacc_credentials import delete_hacc_profile, get_hacc_access_key
import hacc_vars
import time


## VaultEradicator Object:
##      idempotently eradicate Vault AWS components
##
## Methods:
##      __delete_alias
##      __delete_cmk
##      delete_cmk_with_alias
##      __delete_iam_policy
##      __delete_iam_access_key
##      __delete_iam_user
##      delete_iam_user_with_policy
##      __detach_org_policy
##      __delete_org_policy
##      delete_scp
##
class VaultEradicator(VaultComponents):
    def __init__(self):
        super().__init__()



    ## Deletes expected CMK alias, returns False on fail
    def __delete_alias(self):
        key_alias = hacc_vars.aws_hacc_kms_alias
        alias_delete_res = self.aws_client.call(
                                'kms', 'delete_alias',
                                AliasName = f'alias/{key_alias}'
                            )
        return alias_delete_res


    ## Schedules earliest deletion for CMK, returns False on fail
    def __delete_cmk(self):
        key_delete_res = self.aws_client.call(
                                'kms', 'schedule_key_deletion',
                                KeyId = self.cmk,
                                PendingWindowInDays = 7
                            )
        return key_delete_res


    ## Returns True if delete successful or key doesn't exist
    ## Unsets cmk attribute of parent class
    ## Returns False if failed to delete key/alias
    def delete_cmk_with_alias(self):
        print('Checking for existing Vault key')

        if not self.cmk:
            print('No existing Vault key found')
            return True

        print('Existing Vault key found, deleting...')

        alias_delete_res = self.__delete_alias()
        ## If we can't delete alias, leave CMK component for next time
        if not alias_delete_res:
            print('Failed to delete alias for Vault key')
            return False

        key_delete_res = self.__delete_cmk()

        if not key_delete_res:
            print('Failed to schedule Vault key deletion')
            print(f'  Vault key {self.cmk} will need to be manually deleted')
            self.cmk = None
            return False

        print('Successfully deleted Vault key')
        self.cmk = None
        return True



    ## Deletes expected HACC IAM user policy, returns False on fail
    def __delete_iam_policy(self, username, policy):
        delete_policy_res = self.aws_client.call(
                                'iam', 'delete_user_policy',
                                UserName = username,
                                PolicyName = policy
                            )
        return delete_policy_res


    ## Deletes provided HACC IAM access key ID, returns False on fail
    def __delete_iam_access_key(self, username, access_key):
        delete_aws_access_res = self.aws_client.call(
                                    'iam', 'delete_access_key',
                                    UserName = username,
                                    AccessKeyId = access_key
                                )
        return delete_aws_access_res


    ## Deletes expected HACC IAM user, returns False on fail
    def __delete_iam_user(self, username):
        user_delete_res = self.aws_client.call(
                                'iam', 'delete_user',
                                UserName = username
                            )
        return user_delete_res


    ## Returns True if delete successful or user doesn't exist
    ## Unsets user attribute of parent class
    ## Returns False if failed to delete user/policy
    def delete_iam_user_with_policy(self):
        print('Checking for Vault IAM user')

        if not self.user:
            print('No existing IAM user found')
            return True
        print('Existing Vault user found, deleting...')

        username = hacc_vars.aws_hacc_uname

        delete_policy_res = self.__delete_iam_policy(username, hacc_vars.aws_hacc_iam_policy)
        ## If we can't delete policy, leave IAM component for next time
        if not delete_policy_res:
            print('Failed to delete user policy')
            return False

        ## Need local credential ID to delete AWS credential
        aws_access_key = get_hacc_access_key()
        if aws_access_key:
            delete_aws_access_res = self.__delete_iam_access_key(username, aws_access_key)
            if not delete_aws_access_res:
                print('Failed to delete Vault user access key')
                return False

        ## If AWS access key non-existent or deleted, remove local cred
        delete_local_access_res = delete_hacc_profile()
        if not delete_local_access_res:
            print('Failed to delete local Vault credentials')
            print('  AWS credentials/config files will need manual cleanup to remove deprecated profile')

        ## Delete user
        user_delete_res = self.__delete_iam_user(username)
        if not user_delete_res:
            print('Failed to delete Vault IAM user')
            return False

        print('Successfully deleted Vault IAM user')
        self.user = None
        return True



    ## Detaches provided SCP policy ID from org, returns False on fail
    def __detach_org_policy(self, policy_id):
        detach_scp_res = self.aws_client.call(
                                'org', 'detach_policy',
                                PolicyId = policy_id,
                                TargetId = self.aws_account_id
                            )
        return detach_scp_res


    ## Deletes provided SCP policy ID from org, returns False on fail
    def __delete_org_policy(self, policy_id):
        delete_scp_res = self.aws_client.call(
                                'org', 'delete_policy',
                                PolicyId = policy_id
                            )
        return delete_scp_res


    ## Returns True if delete successful or SCP doesn't exist and unsets scp attribute
    ## Returns False if failed to delete SCP
    def delete_scp(self):

        print('Checking for existing Vault SCP')
        if not self.scp:
            print('No existing SCP found')
            return True

        ## Detach SCP from Vault account
        print('Found Vault SCP, deleting...')
        detach_scp_res = self.__detach_org_policy(self.scp)
        if not detach_scp_res:
            print('Failed to detach SCP from Vault account')
            return False

        ## Delete SCP from org
        delete_scp_res = self.__delete_org_policy(self.scp)
        if not delete_scp_res:
            print('Failed to delete SCP from AWS organization')
            print(f'  Unattached SCP {self.scp} will need manual cleanup')
            self.scp = None
            return False

        # Give SCP deletion time to take effect ('immediate' per AWS docs is not good enough :)
        print('Successfully deleted SCP from Vault account, waiting 10 seconds to ensure completion')
        self.scp = None
        time.sleep(10)
        return True