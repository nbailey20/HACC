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
class AwsClient:

    def __init__(self, 
                    ssm=False, 
                    kms=False, 
                    iam=False,
                    org=False,
                    sts=False,
                    aws_access_key_id=None, 
                    aws_secret_access_key=None,
                    aws_session_token=None
                ):

        hacc_session = None if aws_access_key_id else boto3.session.Session(profile_name=hacc_vars.aws_hacc_uname)

        if ssm:
            self.ssm = hacc_session.client('ssm', region_name = hacc_vars.aws_hacc_region)


        if kms:
            ## Check if we're given role creds for Vault install
            if not hacc_session:
                self.kms = boto3.client('kms', region_name=hacc_vars.aws_hacc_region,
                                        aws_access_key_id=aws_access_key_id,
                                        aws_secret_access_key=aws_secret_access_key,
                                        aws_session_token=aws_session_token
                                    )
            else:
                self.kms = hacc_session.client('kms', region_name=hacc_vars.aws_hacc_region)

        ## IAM client only used for Vault install with role creds
        if iam:
            self.iam = boto3.client('iam',
                                    aws_access_key_id=aws_access_key_id,
                                    aws_secret_access_key=aws_secret_access_key,
                                    aws_session_token=aws_session_token
                                )
        ## Used for install to setup SCP
        if org:
            self.org = boto3.client('organizations')

        ## Used for install to first access Vault account (if SCP)
        if sts:
            self.sts = boto3.client('sts')



    ## Execute generic AWS API call with basic error handling
    ## Returns result of api call, or False if call failed
    def call(self, client_type, api_name, **kwargs):
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
