from install.hacc_credentials import get_hacc_access_key
from configure.hacc_encryption import encrypt_config_data, decrypt_config_data
from hacc_generate import generate_password
import hacc_vars
import re

def get_all_configuration():
    setattr(hacc_vars, 'aws_access_key_id', get_hacc_access_key())
    ## TODO, get aws_secret_access_key
    return hacc_vars



def get_configuration(val):
    if val == '*':
        print(get_all_configuration())

    elif val == 'aws_access_key_id':
        print(get_hacc_access_key())

    elif hasattr(hacc_vars, val) and getattr(hacc_vars, val) != '':
        print(getattr(hacc_vars, val))

    elif hasattr(hacc_vars, val):
        print(f'No value set for {val}')

    else:
        print(f'{val} is not a known HACC configuration variable. See hacc_vars_template for expected variables.')
    return


def set_configuration_variable(val):
    ## Parse user input
    try:
        val_list = val.split('=')
        var_name = val_list[0]
        var_val = val_list[1]
    except:
        print('Malformed input, expecting string of form variable=value')
        return

    ## Update credential variables
    if var_name == 'aws_access_key_id':
        return
    elif var_name == 'aws_secret_access_key':
        return

    ## If not credential variable, must be in hacc_vars
    f = open('hacc_vars.py', 'r')
    config = f.readlines()
    f.close()

    ## Update data locally
    value_updated = False
    for line in config:
        line_pattern = var_name + r'\s+=\s+(\S+)'
        line_repl = r'\0' + var_val
        line, num_subs = re.subn(line_pattern, line_repl, line)
        if num_subs > 0:
            value_updated = True
            break
    if not value_updated:
        print('Variable {var_name} not known, see hacc_vars_template for acceptable configuration variables.')
        return

    ## Write updated data to hacc_vars
    f = open('hacc_vars.py', 'w')
    f.write(config)
    f.close()
    print(f'Updated {var_name}.')
    return


def export_configuration(f_name, generate_passwd=True):
    print('Exporting HACC configuration as file...')
    config = get_all_configuration()
    if generate_passwd:
        passwd = generate_password()
    else:
        passwd = input('Enter password to encrypt configuration data file: ')

    enc_config = encrypt_config_data(config, passwd)
    if not enc_config:
        print('Could not encrypt configuration data.')
        return

    try:
        f = open(f_name, 'w')
        f.write(enc_config)
        f.close()
        print(f'Created encrypted file {f_name}.')
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


def configure(args):
    if args.export:
        export_configuration(args.file)
    elif args.show:
        get_configuration(args.show)
    elif args.set:
        if args.file:
            import_configuration_file(args.file)
        else:
            set_configuration_variable(args.set)
    return