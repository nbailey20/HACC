from classes.aws_client import AwsClient
import hacc_vars


## VaultComponent Object:
##      checks existence of Vault AWS components
##
## Attributes:
##      aws_client, AwsClient object
##      aws_account_id, string
##      cmk, ARN of HACC CMK/False
##      user, ARN of HACC user/False
##      scp, ARN of HACC SCP/False
##
## Methods:
##      __cmk_exists
##      __user_exists
##      __scp_exists
##      required
##      active
##
class VaultComponents:

    def __init__(self):
        # region = hacc_vars.aws_hacc_region
        # path = hacc_vars.aws_hacc_param_path
        # self.aws_ssm_path_arn = f'arn:aws:ssm:{region}:{self.aws_account_id}:parameter/{path}'

        ## Configure AWS clients based on whether SCP (multi-account) enabled
        self.aws_client = AwsClient(client_type='mgmt', create_scp=hacc_vars.create_scp)

        identity_info = self.aws_client.call('sts', 'get_caller_identity')
        self.aws_account_id = identity_info['Account']

        scp = self.__scp_exists() if hacc_vars.create_scp else False
        self.scp = scp if scp else None

        cmk = self.__cmk_exists()
        self.cmk = cmk if cmk else None

        user = self.__user_exists()
        self.user = user if user else None



    ## Checks if KMS key exists with expected alias from config file
    ## Returns KMS ARN if key exists, otherwise False
    def __cmk_exists(self):
        kms_alias = f'alias/{hacc_vars.aws_hacc_kms_alias}'

        hacc_kms_arn = self.aws_client.call(
                            'kms', 'describe_key',
                            KeyId = kms_alias
                        )
        if hacc_kms_arn:
            return hacc_kms_arn['KeyMetadata']['Arn']
        return False



    ## Checks if expected IAM user exists for Vault
    ## Returns user ARN if exists, False otherwise
    def __user_exists(self):
        username = hacc_vars.aws_hacc_uname

        hacc_user = self.aws_client.call(
                        'iam', 'get_user',
                        UserName = username
                    )
        if hacc_user:
            return hacc_user['User']['Arn']
        return False



    ## Returns SCP ID if Vault SCP exists for account, False otherwise
    def __scp_exists(self):
        hacc_scp_name = hacc_vars.aws_hacc_scp

        vault_account_scp_obj = self.aws_client.call(
                                    'org', 'list_policies_for_target',
                                    TargetId = self.aws_account_id,
                                    Filter = 'SERVICE_CONTROL_POLICY'
                                )

        if not vault_account_scp_obj:
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
                return False

            for scp in vault_account_scp_obj['Policies']:
                name = scp['Name']
                if name == hacc_scp_name:
                    return scp['Id']
        return False


    ## Function that returns list of required components for Vault based on config params
    def required(self):
        components = ['cmk', 'user']
        if hacc_vars.create_scp:
            components.append('scp')
        return components


    ## Function that returns list of active (fully setup) Vault components
    def active(self):
        components = []
        if self.cmk:
            components.append('cmk')
        if self.user:
            components.append('user')
        if hacc_vars.create_scp and self.scp:
            components.append('scp')
        return components