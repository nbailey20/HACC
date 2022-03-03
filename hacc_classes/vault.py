import hacc_vars
from aws_client import AwsClient
from hacc_service import HaccService
import logging

logger=logging.getLogger()


class Vault:
    ## Return all service names stored in Vault
    def get_all_services(self):
        ssm = self.aws_client.ssm
        all_svc_list = []

        logger.debug('Retrieving service names from Vault')

        curr_params = self.aws_client.call(
                        ssm, 'get_parameters_by_path', 
                        Path = '/'+hacc_vars.aws_hacc_param_path,
                        WithDecryption = False
                    )
        all_svc_list += curr_params['Parameters']
        
        # check for NextToken='string' in response, only returns <= 10 parameters per call
        while 'NextToken' in curr_params:

            more_params = self.aws_client.call(
                            ssm, 'get_parameters_by_path', 
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
    def service_exists(self, service):
        ssm = self.aws_client.ssm
        all_svcs = self.get_all_services(ssm)
        if service in all_svcs:
            return True
        return False



    ## Return True if user exists for service, False otherwise
    def user_exists_for_service(self, user, service):
        ssm = self.aws_client.ssm
        service_obj = HaccService(service, ssm_client=ssm)

        if user in service_obj.get_users():
            return True
        return False



    ## Returns KMS ARN used to encrypt Vault credentials
    ## Returns None if no KMS key exists
    def get_kms_arn(self):
        kms = self.aws_client.kms
        kms_alias = hacc_vars.aws_hacc_kms_alias

        logger.debug(f'Retrieving Vault encryption key')

        try:
            hacc_kms_arn = self.aws_client.call(
                                kms, 'describe_key', 
                                KeyId = f'alias/{kms_alias}'
                            )['KeyMetadata']['Arn']
            return hacc_kms_arn
        except:
            return None


    def __init__(self):
        self.aws_client = AwsClient(ssm=True, kms=True)
        self.kms_arn = self.get_kms_arn()
        self.param_path = hacc_vars.aws_hacc_param_path