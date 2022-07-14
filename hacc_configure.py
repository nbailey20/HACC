from install.hacc_credentials import get_hacc_access_key
from configure.hacc_encryption import encrypt_config_data, decrypt_config_data
from hacc_generate import generate_password
#import hacc_vars
import re, json

def get_all_configuration(config):
    config['aws_access_key_id'] = get_hacc_access_key()
    ## TODO, get aws_secret_access_key
    return json.dumps(config)


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


def set_configuration_variable(val):
    ## Parse user input
    try:
        param_name, param_val = val.split('=')
        #var_name = val_list[0].strip()
        #var_val = val_list[1]
    except:
        print('Malformed input, expecting string of form variable=value')
        return

    ## TODO: Update credential variables
    if param_name == 'aws_access_key_id':
        return
    elif param_name == 'aws_secret_access_key':
        return

    ## If not credential variable, must be in hacc_vars
    f = open('hacc_vars.py', 'r')
    config = f.readlines()
    f.close()

    ## Update data locally
    value_updated = False
    for i in range(len(config)):
        line_pattern = param_name + r'\s+=\s+\S+'
        line_repl = param_name + ' = ' + '\''+param_val+'\''
        updated_line, num_subs = re.subn(line_pattern, line_repl, config[i])

        if num_subs > 0:
            config[i] = updated_line
            value_updated = True
            break
    if not value_updated:
        print(f'Variable {param_name} not known, see hacc_vars_template for acceptable configuration variables.')
        return

    ## Write updated data to hacc_vars
    f = open('hacc_vars.py', 'w')
    f.write(''.join(config))
    f.close()
    print(f'Updated {param_name}.')
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


def import_configuration_file(f_name):
    ## Read encrypted file
    try:
        f = open(f_name, 'r')
        enc_config = f.readlines()
        f.close()
    except:
        print('Could not open file {f_name}.')

    ## Decrypt file
    passwd = input('Enter password used to encrypt configuration file: ')
    config = decrypt_config_data(enc_config, passwd)
    if not config:
        print(f'Failed to decrypt configuration file {f_name}.')

    ## Set hacc_vars variables
    f = open('hacc_vars.py', 'w')
    for a in config:
        if a != 'aws_access_key_id' and a != 'aws_secret_access_key':
            f.write(f'{a} = {a.value()}')
    f.close()
    print('Updated variables in hacc_vars file.')

    ## TODO, set AWS credential vars
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
        if args.file:
            import_configuration_file(args.file)
        else:
            set_configuration_variable(args.set)
    else:
        print('Usage: configure action requires "show", "set", or "export" argument')
    return