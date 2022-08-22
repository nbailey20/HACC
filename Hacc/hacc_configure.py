from vault_install.hacc_credentials import get_hacc_access_key, get_hacc_secret_key, set_hacc_access_key, set_hacc_secret_key
from configure.hacc_encryption import encrypt_config_data, decrypt_config_data
from hacc_generate import generate_password
import re, json, os

def get_all_configuration(config):
    profile_name = config['aws_hacc_uname']
    config['aws_access_key_id'] = get_hacc_access_key(profile_name)
    config['aws_secret_access_key'] = get_hacc_secret_key(profile_name)
    return json.dumps(config, indent=4, sort_keys=True)


def get_configuration(val, config):
    profile_name = config['aws_hacc_uname']

    if val == 'all':
        print(get_all_configuration(config))

    elif val == 'aws_access_key_id':
        print(get_hacc_access_key(profile_name))

    elif val == 'aws_secret_access_key':
        print(get_hacc_secret_key(profile_name))

    elif val in config and config[val] != '':
        print(config[val])

    elif val in config:
        print(f'No value set for {val}')

    else:
        print(f'{val} is not a known HACC configuration variable. See hacc_vars_template for expected variables.')
    return


def set_config_from_cli(inpt, config, config_location):
    ## Parse user input
    param_name, param_val = inpt.split('=')

    ## Handle credential variable updates
    if param_name == 'aws_access_key_id' or param_name == 'aws_secret_access_key':
        profile_name = config['aws_hacc_uname']

        if param_name == 'aws_access_key_id':
            credential_set = set_hacc_access_key(param_val, profile_name)
        elif param_name == 'aws_secret_access_key':
            credential_set = set_hacc_secret_key(param_val, profile_name)

        if not credential_set:
            print(f'Error setting {param_name}, manual configuration of Vault AWS credentials required.')
        else:
            print(f'Successfully set parameter {param_name} in AWS credentials file for profile {profile_name}.')
        return

    ## If not credential variable, must be in config file
    f = open(config_location, 'r')
    config = f.readlines()
    f.close()

    ## Update data locally, also check for commented config
    value_updated = False
    for i in range(len(config)):
        line_pattern = r'\s*#*\s*' + param_name + r'\s*=\s*\S+'
        line_repl = param_name + ' = ' + '\''+param_val+'\''
        updated_line, num_subs = re.subn(line_pattern, line_repl, config[i])

        if num_subs > 0:
            config[i] = updated_line
            value_updated = True
            break
    if not value_updated:
        print(f'Variable {param_name} not known, cannot set value.')
        return

    ## Write updated data to config file
    f = open(config_location, 'w')
    f.write(''.join(config))
    f.close()
    print(f'Updated {param_name}.')
    return


## Set param in file config_location to value in imported config object
def set_config_from_file(param, config, config_location):
    if param == 'all':
        for p in config:
            set_config_from_file(p, config, config_location)
        return

    if not param in config:
        print(f'Variable {param} not found in imported configuration.')
    else:
        cli_expr = f'{param}={config[param]}'
        set_config_from_cli(cli_expr, config, config_location)
    return


def get_set_type(inpt, f_name):
    if f_name and '=' in inpt:
        return 'conflict'
    elif f_name:
        return 'file'
    elif '=' in inpt:
        return 'cli'



def import_configuration_file(f_name):
    ## Read encrypted file
    try:
        f = open(f_name, 'r')
        enc_config = f.read()
        f.close()
    except:
        print('Could not open file {f_name}.')

    ## Decrypt file
    passwd = input('Enter password used to encrypt configuration file: ')
    config = decrypt_config_data(enc_config, passwd)
    try:
        config = json.loads(config)
        return config
    except:
        print(f'Failed to decrypt configuration file {f_name}.')
        return



def export_configuration(f_name, config, generate_passwd=True):
    print('Exporting HACC configuration as file...')
    config_str = get_all_configuration(config)
    if generate_passwd:
        passwd = generate_password()
    else:
        passwd = input('Enter password to encrypt configuration data file: ')

    enc_config = encrypt_config_data(config_str, passwd)
    if not enc_config:
        print('Could not encrypt configuration data.')
        return

    try:
        f = open(f_name, 'w')
        f.write(enc_config)
        f.close()
        print(f'Created encrypted file {f_name} in {os.getcwd()}.')
        if generate_passwd:
            print()
            print('Password for future decryption:')
            print(f'{passwd}')
    except:
        print(f'Could not create configuration file {f_name}.')
    return



def configure(args, config):
    if args.export:
        if args.password:
            export_configuration(args.file, config, generate_passwd=False)
        else:
            export_configuration(args.file, config)

    elif args.show:
        get_configuration(args.show, config)

    elif args.set:
        ## Get configuration file location
        if 'USERPROFILE' in os.environ:
            hacc_config_location = os.environ['USERPROFILE'] + '\.hacc\hacc.conf'
        ## check for linux/mac
        elif 'HOME' in os.environ:
            hacc_config_location = os.environ['HOME'] + '/.hacc/hacc.conf'

        set_type = get_set_type(args.set, args.file)

        if set_type == 'conflict':
            print('Usage: "--set" accepts either command line assignment with "param=value" or an import file with -f, not both!')
        elif set_type == 'cli':
            set_config_from_cli(args.set, config, hacc_config_location)
        elif set_type == 'file':
            config = import_configuration_file(args.file)
            if config:
                set_config_from_file(args.set, config, hacc_config_location)

    else:
        print('Usage: configure action requires "show", "set", or "export" argument')
    return