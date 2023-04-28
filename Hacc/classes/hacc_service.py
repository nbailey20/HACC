import sys

from logger.hacc_logger import logger
from classes.vault import Vault


## HaccService Object:
##      Used for interaction with service credentials in Vault
##
## Attributes:
##      service_name, string
##      vault, Vault obj
##      credentials, dict of user:passwd
## 
## Methods:
##      pull_from_vault
##      push_to_vault
##      get_users
##      get_credential
##      user_exists
##      add_credential
##      remove_credential
class HaccService:

    ## Pulls creds from Vault and updates self.credentials
    def pull_from_vault(self):
        vault_path = self.vault.param_path

        try:
            creds_string = self.vault.aws_client.call(
                'ssm', 'get_parameter', 
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
        svc_creds = self.credentials
        vault_path = self.vault.param_path
        kms_id = self.vault.kms_arn
        
        ## Create parameter string from dict of form {user1:pass1,user2:pass2,etc}
        param_string = ''
        for user in svc_creds:
            param_string += f'{user}:{svc_creds[user]},'

        if not param_string:
            self.vault.aws_client.call(
                'ssm', 'delete_parameter', 
                Name=f'/{vault_path}/{self.service_name}'
            )
            logger.debug('Successfully removed service with no more credentials from Vault')

        else:
            self.vault.aws_client.call(
                'ssm', 'put_parameter',
                Name = f'/{vault_path}/{self.service_name}',
                Value = param_string[:-1], ## remove final comma
                Type = 'SecureString',
                KeyId = kms_id,
                Overwrite = True,
                Tier = 'Standard',
                DataType = 'text'
            )
            logger.debug('Successfully pushed updated credential data to Vault')


    ## Returns list of all usernames associated with service
    def get_users(self):
        return list(self.credentials)

    
    ## Returns password for given username, False if it doesn't exist
    def get_credential(self, user):
        if user in self.credentials:
            logger.debug(f'Found credential for user "{user}"')
            return self.credentials[user]
        logger.debug(f'No credential for user {user}')
        return False


    ## Returns True if username exists in service, False otherwise
    def user_exists(self, user):
        return True if user in self.get_users() else False


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
    ## Either provide config to initialize new vault, or existing vault object
    def __init__(self, display, service_name, config=None, vault=None):
        self.vault = Vault(display, config) if vault == None else vault
        self.service_name = service_name

        if self.vault.service_exists(service_name):
            logger.debug(f'Found existing service "{service_name}" in Vault, pulling data')
            self.pull_from_vault()
            
        else:
            self.credentials = {}
            logger.debug(f'Did not find existing service {service_name} in Vault, creating new')