import hacc_vars
import boto3, subprocess, json

VAULT_IAM_PERMS = """
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ssm:DescribeParameters",
                "ssm:GetParameter*",
                "ssm:DeleteParameter*",
                "ssm:PutParameter"
            ],
            "Resource": "%s"
        },
        {
            "Effect": "Allow",
            "Action": [
                "kms:Encrypt",
                "kms:Decrypt"
            ],
            "Resource": "%s"
        }
    ]
}
"""

def create_hacc_profile(access_key_id, secret_access_key, debug):
    # Adds a new profile to AWS credentials/config file
    aws_config_region = subprocess.run(['aws', 'configure', 'set', 
        'region', hacc_vars.aws_hacc_region, 
        '--profile', hacc_vars.aws_hacc_uname]
    )
    if debug: print('INFO: Successfully wrote vault AWS region to profile')

    aws_config_access = subprocess.run(['aws', 'configure', 'set', 
        'aws_access_key_id', access_key_id, 
        '--profile', hacc_vars.aws_hacc_uname]
    )
    if debug: print('INFO: Successfully wrote vault AWS access key to profile')

    aws_config_secret = subprocess.run(['aws', 'configure', 'set', 
        'aws_secret_access_key', secret_access_key, 
        '--profile', hacc_vars.aws_hacc_uname]
    )
    if debug: print('INFO: Successfully wrote vault AWS secret key to profile')
    return


def aws_call(client, api, debug, **kwargs):
    # Execute generic AWS API call with error handling
    try:
        method_to_call = getattr(client, api)
        if debug: print('INFO: About to call API {}'.format(api))
    except:
        if debug: print('ERROR: API {} not known'.format(api))

    try:
        result = method_to_call(**kwargs)
        if debug: print('INFO: API execution successful')
        return result
    except Exception as e:
        if debug: print('ERROR: API call failed, exiting.: {}'.format(e))
        exit(1)
    

def install(args):
    # Setup IAM user and KMS CMK for Vault
    account = boto3.client('sts').get_caller_identity().get('Account')
    kms = boto3.client('kms', region_name=hacc_vars.aws_hacc_region)
    iam = boto3.client('iam')
    debug = args.debug

    hacc_kms = aws_call(
                kms, 'create_key', debug, 
                KeyUsage='ENCRYPT_DECRYPT', 
                KeySpec='SYMMETRIC_DEFAULT'
                )

    aws_call(
        kms, 'create_alias', debug, 
        AliasName='alias/{key}'.format(key=hacc_vars.aws_hacc_kms_alias), 
        TargetKeyId=hacc_kms['KeyMetadata']['Arn']
    )

    aws_call(
        iam, 'create_user', debug, 
        UserName=hacc_vars.aws_hacc_uname
    )

    hacc_creds = aws_call(
                    iam, 'create_access_key', debug, 
                    UserName=hacc_vars.aws_hacc_uname
                )

    create_hacc_profile(
        hacc_creds['AccessKey']['AccessKeyId'], 
        hacc_creds['AccessKey']['SecretAccessKey'], 
        debug
    )

    vault_path_arn = 'arn:aws:ssm:{region}:{account}:parameter/{path}/*'.format(
                        region=hacc_vars.aws_hacc_region, 
                        account=account,
                        path=hacc_vars.aws_hacc_param_path
                    )
    vault_key_arn = hacc_kms['KeyMetadata']['Arn']
    perms = json.loads(VAULT_IAM_PERMS % (vault_path_arn, vault_key_arn))
    aws_call(
        iam, 'put_user_policy', debug, 
        UserName=hacc_vars.aws_hacc_uname,
        PolicyName=hacc_vars.aws_hacc_policy,
        PolicyDocument=json.dumps(perms)
    )

    print('Vault setup complete.')
    return