import hacc_vars
import boto3, subprocess, json

def install(args):
    account = boto3.client('sts').get_caller_identity().get('Account')
    kms = boto3.client('kms', region_name=hacc_vars.aws_hacc_region)
    iam = boto3.client('iam')

    hacc_kms = kms.create_key(
        KeyUsage='ENCRYPT_DECRYPT',
        KeySpec='SYMMETRIC_DEFAULT'
    )
    if args.debug: print('Created new KMS CMK: ', hacc_kms['ResponseMetadata']['HTTPStatusCode'])

    hacc_alias = kms.create_alias(
        AliasName='alias/{key}'.format(key=hacc_vars.aws_hacc_kms_alias),
        TargetKeyId=hacc_kms['KeyMetadata']['Arn']
    )
    if args.debug: print('Created new KMS alias: ', hacc_alias['ResponseMetadata']['HTTPStatusCode'])

    hacc_user = iam.create_user(
        UserName=hacc_vars.aws_hacc_uname
    )
    if args.debug: print('Created new AWS user: ', hacc_user['ResponseMetadata']['HTTPStatusCode'])

    hacc_creds = iam.create_access_key(
        UserName=hacc_vars.aws_hacc_uname
    )
    if args.debug: print('Created new AWS creds for user: ', hacc_creds['ResponseMetadata']['HTTPStatusCode'])

    aws_config_region = subprocess.run(['aws', 'configure', 'set', 
        'region', hacc_vars.aws_hacc_region, 
        '--profile', hacc_vars.aws_hacc_uname]
    )
    if args.debug: print('Set vault AWS region: ', aws_config_region.returncode)

    aws_config_access = subprocess.run(['aws', 'configure', 'set', 
        'aws_access_key_id', hacc_creds['AccessKey']['AccessKeyId'], 
        '--profile', hacc_vars.aws_hacc_uname]
    )
    if args.debug: print('Set AWS user access key: ', aws_config_access.returncode)

    aws_config_secret = subprocess.run(['aws', 'configure', 'set', 
        'aws_secret_access_key', hacc_creds['AccessKey']['SecretAccessKey'], 
        '--profile', hacc_vars.aws_hacc_uname]
    )
    if args.debug: print('Set AWS user secret key: ', aws_config_secret.returncode)

    vault_perms = {
        'Version': '2012-10-17',
        'Statement': [
            {
                'Effect': 'Allow',
                'Action': [
                    'ssm:DescribeParameters',
                    'ssm:GetParameter*',
                    'ssm:DeleteParameter*',
                    'ssm:PutParameter'
                ],
                'Resource': 'arn:aws:ssm:{region}:{account}:parameter/{path}/*'.format(
                                region=hacc_vars.aws_hacc_region, 
                                account=account,
                                path=hacc_vars.aws_hacc_param_path
                            )
            },
            {
                'Effect': 'Allow',
                'Action': [
                    'kms:Encrypt',
                    'kms:Decrypt'
                ],
                'Resource': hacc_kms['KeyMetadata']['Arn']
            }
        ]
    }
    hacc_policy = iam.put_user_policy(
        UserName=hacc_vars.aws_hacc_uname,
        PolicyName=hacc_vars.aws_hacc_policy,
        PolicyDocument=json.dumps(vault_perms)
    )
    if args.debug: print('Created IAM policy for user: ', hacc_policy['ResponseMetadata']['HTTPStatusCode'])

    print('Vault setup complete.')
    return