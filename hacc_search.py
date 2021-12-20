import hacc_vars
import boto3
import re


def get_all_services(ssm_client, debug):
    vault_params = ssm_client.get_parameters_by_path(
        Path='/'+hacc_vars.aws_hacc_param_path,
        WithDecryption=False
    )
    if debug: print('INFO: retrieved known services')
    # TODO: check for NextToken='string' in response, only returns <= 10 parameters per call
    while 'NextToken' in vault_params:
        print('more params to get!')
    
    ## Example service Name: /hacc-vault/test
    svcs_list = list(set([x['Name'].split('/')[-1] for x in vault_params['Parameters']]))
    if debug: print('Services found: ', svcs_list)
    return svcs_list


def get_service_name(svc, ssm_client, debug):
    svc_list = get_all_services(ssm_client, debug)
    existing_service = False
    for s in svc_list:
        if svc.lower() == s.lower():
            existing_service = True
            if svc != s:
                print('Found similar existing service: {}'.format(s))
                if input('Use this service (y/n)? ') == 'y':
                    return s
            else:
                return s
    return False

def get_service_creds(name, ssm_client, debug):
    try:
        param = ssm_client.get_parameter(
            Name='/{path}/{service}'.format(
                    path=hacc_vars.aws_hacc_param_path, 
                    service=name
                ),
            WithDecryption=True
        )['Parameter']['Value']
        if debug: print('INFO: successfully retrieved and decrypted parameter')
        return param

    except ssm_client.exceptions.ParameterNotFound:
        print('No credential exists for service {}, exiting'.format(name))
        exit(2)
    except Exception as e:
        print('Unexpected error retrieving credential, exiting: ', e)
        exit(99)


def user_exists_for_service(svc_creds, user, debug):
    creds_list = svc_creds.split(',')
    for cred in creds_list:
        name,_ = cred.split(':')
        if name == user:
            if debug: print('INFO: found existing username {}'.format(name))
            return True
    return False
    

def search(args):
    hacc_session = boto3.session.Session(profile_name=hacc_vars.aws_hacc_uname)
    ssm = hacc_session.client('ssm', region_name=hacc_vars.aws_hacc_region)

    if not args.service:
        all_svcs = get_all_services(ssm, args.debug)
        # If no services in vault, can't look anything up
        if not all_svcs:
            print('No credentials in vault, add credential before searching.')
            return
        # Print available services if at least one exists
        print()
        print('Available services:')
        for svc in all_svcs:
            print(svc)
        print()
        args.service = input('Enter service name from above list to retrieve credential: ')

    name = get_service_name(args.service, ssm, args.debug)
    if not name:
        print('Service {} not known, exiting.'.format(args.service))
        exit(2)
    creds_list = get_service_creds(name, ssm, args.debug).split(',')
    print()
    print('Credential(s) found for service {}'.format(name))
    for cred in creds_list:
        user,passwd = cred.split(':')
        if args.username:
            if args.username == user:
                print('Username: {}'.format(user))
                print('Password: {}'.format(passwd))
        else:
            print('Username: {}'.format(user))
            print('Password: {}'.format(passwd))
            print()
    
    return
