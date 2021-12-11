import hacc_vars
from hacc_add_remove_param import get_kms_id
from hacc_install import aws_call
import boto3, subprocess

def eradicate(args):
    print('Eradicating Vault...')
    iam = boto3.client('iam')
    kms = boto3.client('kms', region_name=hacc_vars.aws_hacc_region)
    debug = args.debug

    hacc_key_id = get_kms_id(args.debug, kms)
    aws_call(
        kms, 'schedule_key_deletion', debug, 
        KeyId=hacc_key_id,
        PendingWindowInDays=7
    )

    aws_call(
        kms, 'delete_alias', debug, 
        AliasName='alias/{key}'.format(key=hacc_vars.aws_hacc_kms_alias)
    )

    aws_call(
        iam, 'delete_user_policy', debug, 
        UserName=hacc_vars.aws_hacc_uname,
        PolicyName=hacc_vars.aws_hacc_policy
    )

    aws_access_key = subprocess.run(['aws', 'configure', 'get', 
        'aws_access_key_id', '--profile', hacc_vars.aws_hacc_uname],
        capture_output=True, text=True
    )
    aws_access_key = aws_access_key.stdout.strip()
    if args.debug: print('INFO: Retrieved saved AWS user access key')
    aws_call(
        iam, 'delete_access_key', debug, 
        UserName=hacc_vars.aws_hacc_uname,
        AccessKeyId=aws_access_key
    )

    aws_call(
        iam, 'delete_user', debug, 
        UserName=hacc_vars.aws_hacc_uname,
    )

    ## Remove profile from ~/.aws/config and credentials! read/regex/write Python
    print('Vault has been eradicated.')
    return