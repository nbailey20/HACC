import json

from logger.hacc_logger import logger
from classes.aws_client import AwsClient


## Vault Object:
##
## Attributes:
##      aws_client, AwsClient obj
##      kms_arn, string
##      param_path, string
## 
## Methods:
##      get_all_services
##      service_exists
##      get_kms_arn
##      parse_import_file
##          
class Vault:
    def __init__(self, config, aws_client=None):
        if aws_client:
            self.aws_client = aws_client
        else:
            self.aws_client = AwsClient(config, client_type='data')
        self.kms_arn = self.get_kms_arn(config['aws_hacc_kms_alias'])
        self.param_path = config['aws_hacc_param_path']


    ## Return all service names stored in Vault
    ## Returns False if failure
    def get_all_services(self):
        logger.debug('Retrieving service names from Vault')

        curr_params = self.aws_client.call(
                        'ssm', 'get_parameters_by_path', 
                        Path = '/'+self.param_path,
                        WithDecryption = False
                    )
        
        if not curr_params:
            logger.debug('Failed to pull Vault services')
            return False

        all_svc_list = curr_params['Parameters']
        
        # check for NextToken='string' in response, only returns <= 10 parameters per call
        while 'NextToken' in curr_params:

            more_params = self.aws_client.call(
                            'ssm', 'get_parameters_by_path', 
                            Path = '/'+self.param_path,
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
    def service_exists(self, service):
        all_svcs = self.get_all_services()
        if service in all_svcs:
            return True
        return False



    ## Returns KMS ARN used to encrypt Vault credentials
    ## Returns False if failure to get key
    def get_kms_arn(self, kms_alias):
        logger.debug(f'Retrieving Vault encryption key')
        try:
            hacc_kms_arn = self.aws_client.call(
                                'kms', 'describe_key', 
                                KeyId = f'alias/{kms_alias}'
                            )['KeyMetadata']['Arn']
            return hacc_kms_arn
        except:
            return False


    ## Returns a list of credentials read from a Vault backup output file
    ## If unable to parse, returns False
    def parse_import_file(self, filename):
        try:
            f = open(filename, 'r')
            creds_list = json.loads(f.read())['creds_list']
            f.close()
            return creds_list
        except:
            return False
