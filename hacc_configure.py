from install.hacc_credentials import get_hacc_access_key
from configure.hacc_encryption import encrypt_config_data, decrypt_config_data
from hacc_generate import generate_password
#import hacc_vars
import re, json

def get_all_configuration(config):
    config['aws_access_key_id'] = get_hacc_access_key()
    ## TODO, get aws_secret_access_key
    return json.dumps(config, indent=4, sort_keys=True)


def get_configuration(val, config):
    if val == 'all':
        print(get_all_configuration(config))

    elif val == 'aws_access_key_id':
        print(get_hacc_access_key())

    elif val in config and config[val] != '':
        print(config[val])

    elif val in config:
        print(f'No value set for {val}')

    else:
        print(f'{val} is not a known HACC configuration variable. See hacc_vars_template for expected variables.')
    return


def set_config_from_cli(inpt):
    ## Parse user input
    param_name, param_val = inpt.split('=')

    ## TODO: Update credential variables
    if param_name == 'aws_access_key_id':
        return
    elif param_name == 'aws_secret_access_key':
        return

    ## If not credential variable, must be in hacc_vars
    f = open('hacc_vars.py', 'r')
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

    ## Write updated data to hacc_vars
    f = open('hacc_vars.py', 'w')
    f.write(''.join(config))
    f.close()
    print(f'Updated {param_name}.')
    return


def set_config_from_file(param, config):
    if param == 'all':
        for p in config:
            set_config_from_file(p, config)
        return

    if not param in config:
        print(f'Variable {param} not found in imported configuration.')
    else:
        cli_expr = f'{param}={config[param]}'
        set_config_from_cli(cli_expr)
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
        print(f'Created encrypted file {f_name}.')
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
        set_type = get_set_type(args.set, args.file)

        if set_type == 'conflict':
            print('Usage: "--set" accepts either command line assignment with "param=value" or an import file with -f, not both!')
        elif set_type == 'cli':
            set_config_from_cli(args.set)
        elif set_type == 'file':
            config = import_configuration_file(args.file)
            if config:
                set_config_from_file(args.set, config)

    else:
        print('Usage: configure action requires "show", "set", or "export" argument')
    return