import hacc_vars
from hacc_add_remove_param import get_kms_id
import boto3, subprocess

def eradicate(args):
    iam = boto3.client('iam')
    kms = boto3.client('kms', region_name=hacc_vars.aws_hacc_region)

    hacc_key_id = get_kms_id(args.debug, kms)
    kms_delete = kms.schedule_key_deletion(
        KeyId=hacc_key_id,
        PendingWindowInDays=7
    )
    if args.debug: print('Deleting AWS KMS CMK: ', kms_delete['ResponseMetadata']['HTTPStatusCode'])

    alias_delete = kms.delete_alias(
        AliasName='alias/{key}'.format(key=hacc_vars.aws_hacc_kms_alias)
    )
    if args.debug: print('Deleting AWS KMS alias: ', alias_delete['ResponseMetadata']['HTTPStatusCode'])

    policy_delete = iam.delete_user_policy(
        UserName=hacc_vars.aws_hacc_uname,
        PolicyName=hacc_vars.aws_hacc_policy
    )
    if args.debug: print('Deleting AWS IAM policy: ', policy_delete['ResponseMetadata']['HTTPStatusCode'])

    aws_access_key = subprocess.run(['aws', 'configure', 'get', 
        'aws_access_key_id', '--profile', hacc_vars.aws_hacc_uname],
        capture_output=True, text=True
    )
    aws_access_key = aws_access_key.stdout.strip()
    if args.debug: print('Retrieved saved AWS user access key: ', aws_access_key)

    creds_delete = iam.delete_access_key(
        UserName=hacc_vars.aws_hacc_uname,
        AccessKeyId=aws_access_key
    )
    if args.debug: print('Removing access key from AWS user: ', creds_delete['ResponseMetadata']['HTTPStatusCode'])

    user_delete = iam.delete_user(
        UserName=hacc_vars.aws_hacc_uname
    )
    if args.debug: print('Deleting AWS user: ', user_delete['ResponseMetadata']['HTTPStatusCode'])

    ## Remove profile from ~/.aws/config and credentials! read/regex/write Python
    print('Vault has been eradicated.')
    return