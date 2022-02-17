import hacc_vars
from hacc_install import aws_call
from hacc_search import get_service_name, get_service_creds, user_exists_for_service
import boto3
from hacc_generate import generate

def get_kms_id(debug, kms_client):
    hacc_kms = aws_call(
        kms_client, 'describe_key', debug, 
        KeyId='alias/{key}'.format(key=hacc_vars.aws_hacc_kms_alias)
    )
    if debug: print('INFO: retrieved kms arn {}'.format(hacc_kms['KeyMetadata']['Arn']))
    return hacc_kms['KeyMetadata']['Arn']


def remove_cred_from_credlist(svc_creds, user, debug):
    # Returns an updated credlist with credential for user removed
    creds_list = svc_creds.split(',')
    i = 0
    while i < len(creds_list): 
        name,_ = creds_list[i].split(':')
        if name == user:
            creds_list = creds_list[:i] + creds_list[i+1:]
            if debug: print('INFO: updated credentials for service')
            return ",".join(creds_list)
        i += 1


def add(args):
    print('Adding new credential...')
    hacc_session = boto3.session.Session(profile_name=hacc_vars.aws_hacc_uname)
    ssm = hacc_session.client('ssm', region_name=hacc_vars.aws_hacc_region)
    kms = hacc_session.client('kms', region_name=hacc_vars.aws_hacc_region)
    kms_id = get_kms_id(args.debug, kms)

    # if args.password and args.generate:
    #     print('Can either generate a new password or provide one, aborting.')
    #     return

    if not args.username:
        args.username = input('Enter new username: ')
        if not args.username:
            print('No username provided, aborting.')
            return

    if args.generate:
        good_pass = False
        while not good_pass:
            args.password = generate(args.debug)
            print(f'Generated password: {args.password}')
            gp = input('Is this acceptable (y/n)? ')
            good_pass = True if gp.lower() == 'y' else False

    if not args.password:
        args.password = input('Enter new password: ')
        if not args.password:
            print('No password provided, aborting.')
            return

    if not args.service:
        args.service = input('Enter service that uses the credential: ')
        if not args.service:
            print('No service provided, aborting.')
            return

    svc = get_service_name(args.service, ssm, args.debug)
    svc_creds = ''
    if svc:
        if args.debug: print('INFO: Existing service found')
        svc_creds = get_service_creds(svc, ssm, args.debug)

        # Check if username already exists for service
        if user_exists_for_service(svc_creds, args.username, args.debug):
            print('Username already exists for service {}, aborting.'.format(svc))
            return

        ## TODO: create multiple params for same service if more than 4KB creds
        if len(svc_creds)> 4000:
            print('Too many creds for this service!')
            return

        # Add new credential to service
        svc_creds += ',{0}:{1}'.format(args.username, args.password)
    else:
        svc_creds = '{0}:{1}'.format(args.username, args.password)
        if args.debug: print('INFO No existing service found')

    if args.debug: print('INFO: Creating new credential for service')
    aws_call(
        ssm, 'put_parameter', args.debug, 
        Name='/{path}/{service}'.format(
                path=hacc_vars.aws_hacc_param_path, 
                service=args.service
            ),
        Value=svc_creds,
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
        args.service = input('Enter service that should have credential deleted: ')
    
    # Make sure service exists in vault
    svc = get_service_name(args.service, ssm, args.debug)
    svc_creds = ''
    if not svc:
        print('No service {} found, aborting.'.format(args.service))
        return
    if args.debug: print('INFO: Existing service found')

    if not args.username:
        args.username = input('Enter username that should be deleted: ')

    # Check if username exists for service
    svc_creds = get_service_creds(svc, ssm, args.debug)
    if not user_exists_for_service(svc_creds, args.username, args.debug):
        print('Username does not exist for service {}, aborting.'.format(svc))
        return
   
    # Remove credential from service and update parameter if other creds remaining
    updated_creds = remove_cred_from_credlist(svc_creds, args.username, args.debug)
    if updated_creds:
        kms = hacc_session.client('kms', region_name=hacc_vars.aws_hacc_region)
        kms_id = get_kms_id(args.debug, kms)

        if args.debug: print('INFO: Deleting credential from vault')
        aws_call(
            ssm, 'put_parameter', args.debug, 
            Name='/{path}/{service}'.format(
                    path=hacc_vars.aws_hacc_param_path, 
                    service=args.service
                ),
            Value=updated_creds,
            Type='SecureString',
            KeyId=kms_id,
            Overwrite=True,
            Tier='Standard',
            DataType='text'
        )
        print('Credential deleted from vault.')

    else:
        # If no other creds for service, remove parameter
        aws_call(
            ssm, 'delete_parameter', args.debug, 
            Name='/{path}/{service}'.format(
                    path=hacc_vars.aws_hacc_param_path, 
                    service=args.service
                )
        )
        print('Credential deleted from vault, no more credentials for service {}.'.format(args.service))
    return