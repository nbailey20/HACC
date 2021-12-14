import hacc_vars
from hacc_install import aws_call
from hacc_search import get_service_name, get_service_creds
import boto3

def get_kms_id(debug, kms_client):
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

    svc = get_service_name(args.service, ssm, args.debug)
    param = ''
    if svc:
        if args.debug: print('INFO: Existing service found, adding new credential')
        param = get_service_creds(svc, ssm, args.debug)
        ## TODO: create multiple params for same service if more than 4KB creds
        if len(param)> 4000:
            print('Too many creds for this service!')
            return

        param += ',{0}:{1}'.format(args.username, args.password)
    else:
        param = '{0}:{1}'.format(args.username, args.password)
        if args.debug: print('INFO No existing service found, creating new')

    aws_call(
        ssm, 'put_parameter', args.debug, 
        Name='/{path}/{service}'.format(
                path=hacc_vars.aws_hacc_param_path, 
                service=args.service
            ),
        Value=param,
        Type='SecureString',
        KeyId=kms_id,
        Overwrite=True,
        Tier='Standard',
        DataType='text'
    )
    print('Credential added to {} vault service.'.format('existing' if svc else 'new'))
    return


def delete(args):
    print('Deleting credential...')
    hacc_session = boto3.session.Session(profile_name=hacc_vars.aws_hacc_uname)
    ssm = hacc_session.client('ssm', region_name=hacc_vars.aws_hacc_region)

    if not args.service:
        args.service = input('Enter service that should have credential removed: ')
   
    aws_call(
        ssm, 'delete_parameter', args.debug, 
        Name='/{path}/{service}'.format(
                path=hacc_vars.aws_hacc_param_path, 
                service=args.service
            )
    )
    
    print('Credential deleted from vault.')
    return