import hacc_vars
import boto3
import logging, sys

logger=logging.getLogger()


## Execute generic AWS API call with basic error handling
## Returns result of api call, or False if call failed
def aws_call(client, api_name, **kwargs):
    try:
        method_to_call = getattr(client, api_name)
    except:
        print(f'ERROR: API {api_name} not known, exiting')
        return False

    try:
        result = method_to_call(**kwargs)
        ## add leading spaces for cleaner debug output
        logger.debug(f'    {api_name} API execution successful') 
        return result

    except Exception as e:
        print(f'ERROR: API call failed, exiting: {e}')
        return False



## Return all service names stored in Vault
def get_all_services(ssm_client):
    all_svc_list = []

    logger.debug('Retrieving service names from Vault')

    curr_params = aws_call(
        ssm_client, 'get_parameters_by_path', 
        Path = '/'+hacc_vars.aws_hacc_param_path,
        WithDecryption = False
    )
    all_svc_list += curr_params['Parameters']
    
    # check for NextToken='string' in response, only returns <= 10 parameters per call
    while 'NextToken' in curr_params:

        more_params = aws_call(
            ssm_client, 'get_parameters_by_path', 
            Path = '/'+hacc_vars.aws_hacc_param_path,
            WithDecryption = False,
            NextToken = curr_params['NextToken']
        )
        all_svc_list += more_params['Parameters']
        curr_params = more_params

    logger.debug('Gathered all service names from Vault')
    
    ## Example service Name: /hacc-vault/test
    formatted_svcs_list = list(set([x['Name'].split('/')[-1] for x in all_svc_list]))
    return formatted_svcs_list



## Return True if service exists in Vault, False otherwise
def service_exists(service, ssm_client):
    all_svcs = get_all_services(ssm_client)
    if service in all_svcs:
        return True
    return False



## Return True if user exists for service, False otherwise
def user_exists_for_service(user, service, ssm_client):
    service_obj = HaccService(service, ssm_client=ssm_client)

    if user in service_obj.get_users():
        return True
    return False



## Returns KMS ARN used to encrypt Vault credentials
def get_kms_arn(kms_client):
    kms_alias = hacc_vars.aws_hacc_kms_alias

    logger.debug(f'Retrieving Vault encryption key')

    hacc_kms_arn = aws_call(
        kms_client, 'describe_key', 
        KeyId=f'alias/{kms_alias}'
    )['KeyMetadata']['Arn']

    return hacc_kms_arn




## ApiClient Object:
##      Creates boto3 clients for interaction with Vault
## Attributes:
##      ssm, boto3 client obj
##      kms, boto3 client obj
##      iam, boto3 client obj
class ApiClient:

    def __init__(self, ssm=False, kms=False, iam=False):
        hacc_session = boto3.session.Session(profile_name=hacc_vars.aws_hacc_uname)

        if ssm:
            self.ssm = hacc_session.client('ssm', region_name=hacc_vars.aws_hacc_region)

        if kms:
            self.kms = hacc_session.client('kms', region_name=hacc_vars.aws_hacc_region)




## HaccService Object:
##      Used for interaction with service credentials in Vault
##
## Attributes:
##      service_name, string
##      ssm_client, boto3 obj
##      kms_id, string
##      credentials, dict of user:passwd
## 
## Methods:
##      pull_from_vault
##      push_to_vault
##      get_users
##      get_credential
##      add_credential
##      remove_credential
##      print_credential
class HaccService:

    ## Pulls creds from Vault and updates self.credentials
    def pull_from_vault(self):
        ssm_client = self.ssm_client
        vault_path = hacc_vars.aws_hacc_param_path

        try:
            creds_string = aws_call(
                ssm_client, 'get_parameter', 
                Name = f'/{vault_path}/{self.service_name}',
                WithDecryption = True
            )['Parameter']['Value']

            logger.debug('Successfully pulled credential data from Vault')

            # Parse credential string of form 'user1:pass1,user2:pass2,etc'
            creds_dict = {}
            creds_list = creds_string.split(',')
            for cred in creds_list:
                user, passwd = cred.split(':')
                creds_dict[user] = passwd

            self.credentials = creds_dict

        except Exception as e:
            print(f'Unexpected error retrieving credential, exiting: {e}')
            sys.exit(99)


    ## Pushes self.credentials to Vault if at least one credential present
    ## Otherwise delete service from Vault if all creds removed
    def push_to_vault(self):
        ssm_client = self.ssm_client
        svc_creds = self.credentials
        vault_path = hacc_vars.aws_hacc_param_path
        
        ## Create parameter string from dict of form {user1:pass1,user2:pass2,etc}
        param_string = ''
        for user in svc_creds:
            param_string += f'{user}:{svc_creds[user]},'

        if not param_string:
            aws_call(
                ssm_client, 'delete_parameter', 
                Name=f'/{vault_path}/{self.service_name}'
            )
            logger.debug('Successfully removed service with no more credentials from Vault')

        else:
            aws_call(
                ssm_client, 'put_parameter',
                Name=f'/{vault_path}/{self.service_name}',
                Value=param_string[:-1], ## remove final comma
                Type='SecureString',
                KeyId=self.kms_id,
                Overwrite=True,
                Tier='Standard',
                DataType='text'
            )
            logger.debug('Successfully pushed updated credential data to Vault')


    ## Returns list of all usernames associated with service
    def get_users(self):
        return list(self.credentials)

    
    ## Returns password for given username, False if it doesn't existS
    def get_credential(self, user):
        if user in self.credentials:
            logger.debug(f'Found credential for user "{user}"')
            return self.credentials[user]
        logger.debug(f'No credential for user {user}')
        return False


    ## Adds new user:passwd to self.credentials
    ## Returns False if user already exists, True if success
    def add_credential(self, user, passwd):
        if not user in self.credentials:
            self.credentials[user] = passwd
            return True
        else:
            return False


    ## Removes user from self.credentials
    ## Returns False if user doesn't exist, True if success
    def remove_credential(self, user):
        if user in self.credentials:
            self.credentials.pop(user)
            return True
        else:
            return False


    ## Upon object init, pull existing service data if it exists
    def __init__(self, service_name, ssm_client=None, kms_id=None):

        ## Minimize the number of boto3 sessions/clients created
        if not ssm_client and not kms_id:
            api_obj = ApiClient(ssm=True, kms=True)
            self.ssm_client = api_obj.ssm
            self.kms_id = get_kms_arn(api_obj.kms)
        elif not ssm_client:
            api_obj = ApiClient(ssm=True)
            self.ssm_client = api_obj.ssm
            self.kms_id = kms_id
        elif not kms_id:
            api_obj = ApiClient(kms=True)
            self.ssm_client = ssm_client
            self.kms_id = get_kms_arn(api_obj.kms)
        else:
            self.ssm_client = ssm_client
            self.kms_id = kms_id


        self.service_name = service_name

        if service_exists(service_name, self.ssm_client):
            logger.debug(f'Found existing service "{service_name}" in Vault, pulling data')
            self.pull_from_vault()
            
        else:
            self.credentials = {}
            logger.debug('Did not find existing service in Vault, creating new')

    
    ## Prints username, and optionally password for service
    def print_credential(self, user, print_password=False):
        if print_password:
            password = self.get_credential(user)
        print()
        print('#############################')
        print(f'Service {self.service_name}')
        print(f'Username: {user}')
        if print_password:
            print(f'Password: {password}')
        print('#############################')
        print()






## VaultInstallation Object:
##      idempotently install or eradicate Vault AWS components
##
## Attributes:
##      mgmt_org, boto3 org object in mgmt account (if SCP)
##      kms, boto3 kms object
##      iam, boto3 iam object
##      
## 
## Methods:
##      create_cmk_with_alias
##      delete_cmk
##      
##      
##      
##      
##      
class VaultInstallation:

    ## Accepts alias or key id to delete
    ## Returns True if delete successful or key doesn't exist
    ## Returns False if failed to delete key/alias
    def delete_cmk(self, cmk):
        hacc_key_id = cmk

        if cmk.startswith('alias/'):
            ## Get key ID to delete after alias
            hacc_key_id = get_kms_arn(self.kms)

            alias_delete_res = aws_call(
                                self.kms, 'delete_alias', 
                                AliasName = f'alias/{cmk}'
                            )
            if not alias_delete_res:
                logger.debug('Failed to delete alias for Vault KMS key')
                return False
            logger.debug('Successfully deleted alias for Vault KMS key')


        key_delete_res = aws_call(
                            self.kms, 'schedule_key_deletion', 
                            KeyId = hacc_key_id,
                            PendingWindowInDays = 7
                        )
        if not key_delete_res:
            logger.debug('Failed to delete KMS key')
            return False

        logger.debug('Successfully deleted Vault KMS key')
        return True


    ## Creates or confirms KMS exists for Vault 
    ## Returns True if key is present
    ## Returns False if key/alias failed to create
    def create_cmk_with_alias(self):

        ## Check if key already exists
        arn = get_kms_arn(self.kms)
        if arn:
            logger.debug('Found existing KMS key and alias for Vault')
            return True

        ## Otherwise create new key and alias
        hacc_kms = aws_call(
                    self.kms, 'create_key', 
                    KeyUsage = 'ENCRYPT_DECRYPT', 
                    KeySpec = 'SYMMETRIC_DEFAULT'
                )
        if not hacc_kms:
            logger.debug('Failed to create new symmetric KMS key')
            return False
        logger.debug('Successfully created new symmetric KMS key for Vault')

        key_name = hacc_vars.aws_hacc_kms_alias
        key_id = hacc_kms['KeyMetadata']['Arn']

        hacc_alias = aws_call(
            self.kms, 'create_alias', 
            AliasName = f'alias/{key_name}',
            TargetKeyId = key_id
        )

        ## If alias creates properly we're done
        if hacc_alias:
            logger.debug('Successfully created alias for KMS key')
            return True

        ## If alias fails but key exists, clean up first
        logger.debug('Failed to create alias for KMS key')
        logger.debug('Cleaning up KMS key with no alias')
        self.delete_cmk(key_id)
        return False
            

    def __init__(self):

        ## SCP is setup in mgmt account, don't use member role
        self.mgmt_org = boto3.client('organizations')

        ## Assume member role for all other Vault resources
        sts = boto3.client('sts')
        assumed_role_object=sts.assume_role(
            RoleArn=hacc_vars.aws_member_role,
            RoleSessionName="HaccInstallSession"
        )
        role_creds = assumed_role_object['Credentials']

        ## arn:aws:iam::account:role/name
        account = hacc_vars.aws_member_role.split(':')[4]

        self.kms = boto3.client('kms', region_name=hacc_vars.aws_hacc_region,
                            aws_access_key_id=role_creds['AccessKeyId'],
                            aws_secret_access_key=role_creds['SecretAccessKey'],
                            aws_session_token=role_creds['SessionToken']
                            )
        self.iam = boto3.client('iam',
                            aws_access_key_id=role_creds['AccessKeyId'],
                            aws_secret_access_key=role_creds['SecretAccessKey'],
                            aws_session_token=role_creds['SessionToken']
                            )

    

        