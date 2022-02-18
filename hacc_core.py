import hacc_vars
import boto3
import logging, sys

logger=logging.getLogger()


## Execute generic AWS API call with basic error handling
def aws_call(client, api_name, **kwargs):
    try:
        method_to_call = getattr(client, api_name)
    except:
        print(f'ERROR: API {api_name} not known, exiting')

    try:
        result = method_to_call(**kwargs)
        ## add leading spaces for cleaner debug output
        logger.debug(f'    {api_name} API execution successful') 
        return result
    except Exception as e:
        print(f'ERROR: API call failed, exiting: {e}')
        exit(1)



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

        if iam:
            self.iam = hacc_session.client('iam', region_name=hacc_vars.aws_hacc_region)





## HaccService Object:
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
        return self.credentials.keys()

    
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
