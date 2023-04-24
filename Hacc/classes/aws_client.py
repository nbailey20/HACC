import sys
import logging

try:
    import boto3
    from botocore.exceptions import ClientError
except:
    print('Python module boto3 required for HACC client. Install (pip install boto3) and try again.')
    sys.exit()


logger=logging.getLogger()



## AwsClient Object:
##      Creates boto3 clients for interaction with Vault services
## Attributes:
##      ssm, boto3 client obj
##      kms, boto3 client obj
##      iam, boto3 client obj
##      sts, boto3 client obj
##      org, boto3 client obj
##      display, classes.Display object
##
class AwsClient:

    def __init__(self, display, config, client_type='data'):
        self.display = display
        
        region = config['aws_hacc_region']
        create_scp = config['create_scp']

        ## If API client is for data operations, only need SSM and KMS clients
        if client_type == 'data':
            hacc_session = boto3.session.Session(profile_name=config['aws_hacc_uname'])

            self.ssm = hacc_session.client('ssm', region_name=region)
            self.kms = hacc_session.client('kms', region_name=region)

        ## Install/eradicate action clients
        if client_type == 'mgmt':
            self.sts = boto3.client('sts')

            ## Assume cross account role for Vault resources, create SCP in root account
            if create_scp:
                assumed_role_object = self.call('sts', 'assume_role',
                                                RoleArn=config['aws_member_role'],
                                                RoleSessionName="HaccSession"
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
                self.kms = boto3.client('kms', region_name = region,
                                            aws_access_key_id = aws_access_key_id,
                                            aws_secret_access_key = aws_secret_access_key,
                                            aws_session_token = aws_session_token
                                        )
                self.sts = boto3.client('sts', region_name = region,
                                            aws_access_key_id = aws_access_key_id,
                                            aws_secret_access_key = aws_secret_access_key,
                                            aws_session_token = aws_session_token
                                        )
                self.org = boto3.client('organizations')


            ## Only one account, use current AWS creds
            else:
                self.iam = boto3.client('iam')
                self.kms = boto3.client('kms', region_name=region)



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
            self.display.update(
                display_type='text_append', 
                data={'text': f'API {api_name} not known for {client_type} client, exiting'}
            )
            return False

        try:
            result = method_to_call(**kwargs)
            ## add leading spaces for cleaner debug output
            logger.debug(f'    {api_name} API execution successful') 
            return result

        except ClientError as e:
            if e.response['ResponseMetadata']['HTTPStatusCode'] == 403:
                self.display.update(
                    display_type='text_append', 
                    data={'text': f'HACC is not authorized to perform {client_type} {api_name}, exiting'}
                )
            else:
                logger.debug(f'AWS client error: {e}')
            return False

        except Exception as e:
            self.display.update(
                display_type='text_append', 
                data={'text': f'Unknown API error, exiting: {e}'}
            )
            return False
