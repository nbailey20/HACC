import hacc_vars
import boto3

def search(args):
    hacc_session = boto3.session.Session(profile_name=hacc_vars.aws_hacc_uname)
    ssm = hacc_session.client('ssm', region_name=hacc_vars.aws_hacc_region)

    if not args.service:
        args.service = input('Enter service to retrieve credential: ')
    
    try:
        param = ssm.get_parameter(
            Name='/{path}/{service}'.format(
                    path=hacc_vars.aws_hacc_param_path, 
                    service=args.service
                ),
            WithDecryption=True
        )['Parameter']['Value']
        if args.debug: print('INFO: successfully retrieved and decrypted parameter')

        user,passwd = param.split(':')

        print('Credential found for service {}'.format(args.service))
        print('Username: {}'.format(user))
        print('Password: {}'.format(passwd))
        return
    except ssm.exceptions.ParameterNotFound:
        print('No credential exists for service {}, exiting'.format(args.service))
        exit(2)
    except Exception as e:
        print('Unexpected error retrieving credential, exiting: ', e)
        exit(99)
