from hacc_core import HaccService, service_exists
from sys import exit



def add(args):
    print('Adding new credential...')

    service_name = args.service
    user = args.username
    passwd = args.password

    local_service = HaccService(service_name)

    ## empty dicts evaluate to False
    if not bool(local_service.credentials):
        print(f'Creating new service {service_name} in Vault')

    if not local_service.add_credential(user, passwd):
        print(f'Username {user} already exists for service {service_name}, exiting.')
        exit(1)

    local_service.push_to_vault()
    print('Successfully added credential to Vault.')

    # hacc_session = boto3.session.Session(profile_name=hacc_vars.aws_hacc_uname)
    # ssm = hacc_session.client('ssm', region_name=hacc_vars.aws_hacc_region)
    # kms = hacc_session.client('kms', region_name=hacc_vars.aws_hacc_region)
    # kms_id = get_kms_id(args.debug, kms)

    # if args.generate:
    #     good_pass = False
    #     while not good_pass:
    #         args.password = generate(args.debug)
    #         print(f'Generated password: {args.password}')
    #         gp = input('Is this acceptable (y/n)? ')
    #         good_pass = True if gp.lower() == 'y' else False


    # svc = get_service_name(args.service, ssm, args.debug)
    # svc_creds = ''
    # if svc:
    #     if args.debug: print('INFO: Existing service found')
    #     svc_creds = get_service_creds(svc, ssm, args.debug)

    #     # Check if username already exists for service
    #     if user_exists_for_service(svc_creds, args.username, args.debug):
    #         print('Username already exists for service {}, aborting.'.format(svc))
    #         return

    #     ## TODO: create multiple params for same service if more than 4KB creds
    #     if len(svc_creds)> 4000:
    #         print('Too many creds for this service!')
    #         return

    #     # Add new credential to service
    #     svc_creds += ',{0}:{1}'.format(args.username, args.password)
    # else:
    #     svc_creds = '{0}:{1}'.format(args.username, args.password)
    #     if args.debug: print('INFO No existing service found')

    # if args.debug: print('INFO: Creating new credential for service')
    # aws_call(
    #     ssm, 'put_parameter', args.debug, 
    #     Name='/{path}/{service}'.format(
    #             path=hacc_vars.aws_hacc_param_path, 
    #             service=args.service
    #         ),
    #     Value=svc_creds,
    #     Type='SecureString',
    #     KeyId=kms_id,
    #     Overwrite=True,
    #     Tier='Standard',
    #     DataType='text'
    # )
    # print('Credential added to {} vault service.'.format('existing' if svc else 'new'))
    return


def delete(args):
    print('Deleting credential...')
    # hacc_session = boto3.session.Session(profile_name=hacc_vars.aws_hacc_uname)
    # ssm = hacc_session.client('ssm', region_name=hacc_vars.aws_hacc_region)

    service_name = args.service
    user = args.username

    local_service = HaccService(service_name)
    if not local_service.remove_credential(user):
        print(f'Username {user} does not exist for service {service_name}, exiting.')
        exit(1)
    local_service.push_to_vault()

    print('Successfully deleted credential from Vault.')
    if not service_exists(service_name, local_service.ssm_client):
        print(f'No more credentials for service {service_name}.')
    
    # # Make sure service exists in vault
    # svc = get_service_name(args.service, ssm, args.debug)
    # svc_creds = ''
    # if not svc:
    #     print('No service {} found, aborting.'.format(args.service))
    #     return
    # if args.debug: print('INFO: Existing service found')

    # # Check if username exists for service
    # svc_creds = get_service_creds(svc, ssm, args.debug)
    # if not user_exists_for_service(svc_creds, args.username, args.debug):
    #     print('Username does not exist for service {}, aborting.'.format(svc))
    #     return
   
    # # Remove credential from service and update parameter if other creds remaining
    # updated_creds = remove_cred_from_credlist(svc_creds, args.username, args.debug)
    # if updated_creds:
    #     kms = hacc_session.client('kms', region_name=hacc_vars.aws_hacc_region)
    #     kms_id = get_kms_id(args.debug, kms)

    #     if args.debug: print('INFO: Deleting credential from vault')
    #     aws_call(
    #         ssm, 'put_parameter', args.debug, 
    #         Name='/{path}/{service}'.format(
    #                 path=hacc_vars.aws_hacc_param_path, 
    #                 service=args.service
    #             ),
    #         Value=updated_creds,
    #         Type='SecureString',
    #         KeyId=kms_id,
    #         Overwrite=True,
    #         Tier='Standard',
    #         DataType='text'
    #     )
    #     print('Credential deleted from vault.')

    # else:
    #     # If no other creds for service, remove parameter
    #     aws_call(
    #         ssm, 'delete_parameter', args.debug, 
    #         Name='/{path}/{service}'.format(
    #                 path=hacc_vars.aws_hacc_param_path, 
    #                 service=args.service
    #             )
    #     )
    #     print('Credential deleted from vault, no more credentials for service {}.'.format(args.service))
    return