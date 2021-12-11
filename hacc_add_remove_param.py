import hacc_vars
from hacc_install import aws_call
import boto3

def get_kms_id(debug, kms_client):
    # hacc_kms = kms_client.describe_key(
    #     KeyId='alias/{key}'.format(key=hacc_vars.aws_hacc_kms_alias)
    # )
    # if debug: print('Retrieving AWS KMS CMK details: ', hacc_kms['ResponseMetadata']['HTTPStatusCode'])
    hacc_kms = aws_call(
        kms_client, 'describe_key', debug, 
        KeyId='alias/{key}'.format(key=hacc_vars.aws_hacc_kms_alias)
    )

    return hacc_kms['KeyMetadata']['Arn']


def add(args):
    print('Adding new credential...')
    kms = boto3.client('kms', region_name=hacc_vars.aws_hacc_region)
    kms_id = get_kms_id(args.debug, kms)

    hacc_session = boto3.session.Session(profile_name=hacc_vars.aws_hacc_uname)
    ssm = hacc_session.client('ssm', region_name=hacc_vars.aws_hacc_region)

    if not args.username:
        args.username = input('Enter new username: ')

    if not args.password:
        args.password = input('Enter new password: ')

    if not args.service:
        args.service = input('Enter service that uses the credential: ')

    # creds_add = ssm.put_parameter(
    #     Name='/{path}/{service}'.format(
    #             path=hacc_vars.aws_hacc_param_path, 
    #             service=args.service
    #         ),
    #     Value='{0}:{1}'.format(args.username, args.password),
    #     Type='SecureString',
    #     KeyId=kms_id,
    #     Overwrite=True,
    #     Tier='Standard',
    #     DataType='text'
    # )
    # if args.debug: print('Encrypted parameter saved to AWS: ', creds_add['ResponseMetadata']['HTTPStatusCode'])
    aws_call(
        ssm, 'put_parameter', args.debug, 
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


    print('Credential added to vault.')
    return


def delete(args):
    print('Deleting credential...')
    hacc_session = boto3.session.Session(profile_name=hacc_vars.aws_hacc_uname)
    ssm = hacc_session.client('ssm', region_name=hacc_vars.aws_hacc_region)

    if not args.service:
        args.service = input('Enter service that should have credential removed: ')
    # cred_delete = ssm.delete_parameter(
    #     Name='/{path}/{service}'.format(
    #             path=hacc_vars.aws_hacc_param_path, 
    #             service=args.service
    #         )
    # )
    # if args.debug: print('Encrypted parameter deleted AWS: ', cred_delete['ResponseMetadata']['HTTPStatusCode'])
    aws_call(
        ssm, 'delete_parameter', args.debug, 
        Name='/{path}/{service}'.format(
                path=hacc_vars.aws_hacc_param_path, 
                service=args.service
            )
    )
    
    print('Credential deleted from vault.')
    return