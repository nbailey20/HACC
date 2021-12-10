import hacc_vars
import boto3

def search(args):
    hacc_session = boto3.session.Session(profile_name=hacc_vars.aws_hacc_uname)
    ssm = hacc_session.client('ssm', region_name=hacc_vars.aws_hacc_region)
    
    param = ssm.get_parameter(
        Name='/{path}/{service}'.format(
                path=hacc_vars.aws_hacc_param_path, 
                service=args.service
            ),
        WithDecryption=True
    )['Parameter']['Value']

    user,passwd = param.split(':')

    print('Credential found for service {}'.format(args.service))
    print('Username: {}'.format(user))
    print('Password: {}'.format(passwd))
    return