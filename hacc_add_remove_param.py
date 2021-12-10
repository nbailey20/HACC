import hacc_vars
import boto3

def get_kms_id(debug):
    kms = boto3.client('kms', region_name=hacc_vars.aws_hacc_region)
    hacc_kms = kms.describe_key(
        KeyId='alias/{key}'.format(key=hacc_vars.aws_hacc_kms_alias)
    )
    if debug: print('Retrieving AWS KMS CMK details: ', hacc_kms['ResponseMetadata']['HTTPStatusCode'])

    return hacc_kms['KeyMetadata']['Arn']


def add(args):
    hacc_session = boto3.session.Session(profile_name=hacc_vars.aws_hacc_uname)
    kms_id = get_kms_id(args.debug)
    ssm = hacc_session.client('ssm', region_name=hacc_vars.aws_hacc_region)

    creds_add = ssm.put_parameter(
        Name='/{path}/{service}'.format(
                path=hacc_vars.aws_hacc_param_path, 
                service=args.service
            ),
        Value='{0}:{1}'.format(args.username, args.password),
        Type='SecureString',
        KeyId=kms_id,
        Overwrite=True,
        Tier='Standard',
        DataType='text'
    )
    if args.debug: print('Encrypted parameter saved to AWS: ', creds_add['ResponseMetadata']['HTTPStatusCode'])

    print('Credential added.')
    return

def delete(args):
    hacc_session = boto3.session.Session(profile_name=hacc_vars.aws_hacc_uname)
    ssm = hacc_session.client('ssm', region_name=hacc_vars.aws_hacc_region)
    cred_delete = ssm.delete_parameter(
        Name='/{path}/{service}'.format(
                path=hacc_vars.aws_hacc_param_path, 
                service=args.service
            )
    )
    if args.debug: print('Encrypted parameter deleted AWS: ', cred_delete['ResponseMetadata']['HTTPStatusCode'])
    
    print('Credential deleted.')
    return