from hacc_classes.aws_client import AwsClient
import hacc_vars
import logging

logger=logging.getLogger()



## VaultInstallation Object:
##      idempotently install or eradicate Vault AWS components
##
## Attributes:
##      AwsClient object
##      
## 
## Methods:
##      cmk_exists
##      create_cmk_with_alias
##      delete_cmk
##      
##      
##      
##      
##      
class VaultInstallation:

    ## Checks if KMS key exists with expected alias from config file
    ## Returns KMS ARN if key exists, otherwise False
    def cmk_exists(self):
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
            return False



    ## Accepts alias or key id to delete
    ## Returns True if delete successful or key doesn't exist
    ## Returns False if failed to delete key/alias
    def delete_cmk(self, key_id=None):
        kms = self.aws_client.kms

        ## Look for existing key using expected alias
        if not key_id:
            key_id = self.cmk_exists()
            key_alias = hacc_vars.aws_hacc_kms_alias

            if key_id:
                alias_delete_res = self.aws_client.call(kms, 'delete_alias',
                                                    AliasName = f'alias/{key_alias}'
                                                )
            if not alias_delete_res:
                logger.debug('Failed to delete alias for Vault KMS key')
                return False
            logger.debug('Successfully deleted alias for Vault KMS key')


        key_delete_res = self.aws_client.call(kms, 'schedule_key_deletion', 
                                            KeyId = key_id,
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
        kms = self.aws_client.kms

        ## Check if key already exists
        arn = self.cmk_exists()
        if arn:
            logger.debug('Found existing KMS key and alias for Vault')
            return True

        ## Otherwise create new key and alias
        hacc_kms = self.aws_client.call(kms, 'create_key', 
                                        KeyUsage = 'ENCRYPT_DECRYPT', 
                                        KeySpec = 'SYMMETRIC_DEFAULT'
                                    )
        if not hacc_kms:
            logger.debug('Failed to create new symmetric KMS key')
            return False
        logger.debug('Successfully created new symmetric KMS key for Vault')


        key_name = hacc_vars.aws_hacc_kms_alias
        key_id = hacc_kms['KeyMetadata']['Arn']

        hacc_alias = self.aws_client.call(kms, 'create_alias', 
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
        self.delete_cmk(key_id=key_id)
        return False
            

    def __init__(self):

        ## Assume member role to access account for first time (if SCP)
        aws_client = AwsClient(sts=True)
        assumed_role_object = aws_client.call('sts', 'assume_role',
                                                RoleArn=hacc_vars.aws_member_role,
                                                RoleSessionName="HaccInstallSession"
                                            )
        role_creds = assumed_role_object['Credentials']

        ## Create API client for installation
        self.aws_client = AwsClient(kms=True, iam=True, org=True,
                                aws_access_key_id = role_creds['AccessKeyId'],
                                aws_secret_access_key = role_creds['SecretAccessKey'],
                                aws_session_token = role_creds['SessionToken']
                            )


    