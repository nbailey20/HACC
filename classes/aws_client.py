import hacc_vars
import boto3
import logging

logger=logging.getLogger()



## AwsClient Object:
##      Creates boto3 clients for interaction with Vault services
## Attributes:
##      ssm, boto3 client obj
##      kms, boto3 client obj
##      iam, boto3 client obj
##      sts, boto3 client obj
##      org, boto3 client obj
##
class AwsClient:

    def __init__(self, client_type='data', create_scp=True):

        ## If API client is for data operations, only need SSM and KMS clients
        if client_type == 'data':
            hacc_session = boto3.session.Session(profile_name=hacc_vars.aws_hacc_uname)

            self.ssm = hacc_session.client('ssm', region_name=hacc_vars.aws_hacc_region)
            self.kms = hacc_session.client('kms', region_name=hacc_vars.aws_hacc_region)

        ## Install/eradicate action clients
        if client_type == 'mgmt':
            self.sts = boto3.client('sts')

            ## Assume cross account role for Vault resources, create SCP in root account
            if create_scp:
                assumed_role_object = self.call('sts', 'assume_role',
                                                RoleArn=hacc_vars.aws_member_role,
                                                RoleSessionName="HaccInstallSession"
                                            )
                role_creds = assumed_role_object['Credentials']
                aws_access_key_id = role_creds['AccessKeyId']
                aws_secret_access_key = role_creds['SecretAccessKey']
                aws_session_token = role_creds['SessionToken']

                self.iam = boto3.client('iam',
                                        aws_access_key_id = aws_access_key_id,
                                        aws_secret_access_key = aws_secret_access_key,
                                        aws_session_token = aws_session_token
                                    )
                self.kms = boto3.client('kms', region_name = hacc_vars.aws_hacc_region,
                                            aws_access_key_id = aws_access_key_id,
                                            aws_secret_access_key = aws_secret_access_key,
                                            aws_session_token = aws_session_token
                                        )
                self.sts = boto3.client('sts', region_name = hacc_vars.aws_hacc_region,
                                            aws_access_key_id = aws_access_key_id,
                                            aws_secret_access_key = aws_secret_access_key,
                                            aws_session_token = aws_session_token
                                        )
                self.org = boto3.client('organizations')


            ## Only one account, use current AWS creds
            else:
                self.iam = boto3.client('iam')
                self.kms = boto3.client('kms', region_name=hacc_vars.aws_hacc_region)



    ## Execute generic AWS API call with basic error handling
    ## Returns result of api call, or False if call failed
    def call(self, client_type, api_name, **kwargs):
        client = None
        if client_type == 'ssm':
            client = self.ssm
        if client_type == 'kms':
            client = self.kms
        if client_type == 'sts':
            client = self.sts
        if client_type == 'iam':
            client = self.iam
        if client_type == 'org':
            client = self.org
    
        try:
            method_to_call = getattr(client, api_name)
        except:
            print(f'API {api_name} not known for {client_type} client, exiting')
            return False

        try:
            result = method_to_call(**kwargs)
            ## add leading spaces for cleaner debug output
            logger.debug(f'    {api_name} API execution successful') 
            return result

        except Exception as e:
            print(f'API call failed, exiting: {e}')
            return False
