from classes.vault_components import VaultComponents
import os, subprocess

MGMT_ACTIONS = ['install', 'eradicate'] ## configure is mgmt action but doesn't have any required vars
REQUIRED_MGMT_VARS = [
    'aws_hacc_region',
    'aws_hacc_uname',
    'aws_hacc_iam_policy',
    'aws_hacc_kms_alias',
    'aws_hacc_param_path',
    'create_scp'
]

DATA_ACTIONS = ['add', 'delete', 'search', 'rotate', 'backup']
REQUIRED_DATA_VARS = [
    'aws_hacc_kms_alias',
    'aws_hacc_uname',
    'aws_hacc_region',
    'aws_hacc_param_path'
]

REQUIRED_SCP_VARS = ['aws_hacc_scp', 'aws_member_role']


## Function to read config parameters from hacc_vars file into object
def get_config_params():
    ## check for windows system
    if 'USERPROFILE' in os.environ:
        hacc_config_location = os.environ['USERPROFILE'] + '\.hacc\hacc.conf'
    ## check for linux/mac
    elif 'HOME' in os.environ:
        hacc_config_location = os.environ['HOME'] + '/.hacc/hacc.conf'

    try:
        f = open(hacc_config_location, 'r')
        hacc_vars = f.readlines()
        f.close()
    except:
        print(f'Unable to read required configuration file {hacc_config_location}, aborting.')
        return False

    config = {}
    for line in hacc_vars:
        line = line.strip()
        try:
            ## ignore comments in config file, strip quotes
            if line[0] != '#':
                param, value = [x.strip() for x in line.split('=')]
                value = value.replace('"', '').replace("'", "")
                if value == 'True':
                    config[param] = True
                elif value == 'False':
                    config[param] = False
                else:
                    config[param] = value
        except:
            continue

    return config



## Function to check all required config vars are set when client is invoked
def required_config_set_for_action(args, config):
    all_vars_exist = True
    missing_vars = []

    if args.action in DATA_ACTIONS:
        for v in REQUIRED_DATA_VARS:
            if not v in config or config[v] == '':
                missing_vars.append(v)
                all_vars_exist = False

    elif args.action in MGMT_ACTIONS:
        for v in REQUIRED_MGMT_VARS:
            if not v in config or config[v] == '':
                missing_vars.append(v)
                all_vars_exist = False

        ## conditionally check for SCP vars if create_scp == True
        if not 'create_scp' in missing_vars and config['create_scp'] == True:
            for v in REQUIRED_SCP_VARS:
                if not v in config  or config[v] == '':
                    missing_vars.append(v)
                    all_vars_exist = False

    for m in missing_vars:
        print(f'Required configuration parameter {m} not set in hacc_vars.py, please set value with hacc configure --set {m}=... (or edit manually) and try again.')

    return all_vars_exist


## Function to confirm all required Vault components for action are setup
def vault_components_exist_for_action(args, config):
    if args.action in DATA_ACTIONS:
        components = VaultComponents(config)
        active = components.active()
        required = components.required()

        if len(active) == 0:
            print('No Vault detected, install before attempting this command.')
            return False

        elif len(active) != len(required):
            print('Vault is not fully setup, complete installation before attempting this command.')
            print(f'Active components: {active}')
            missing = [x for x in required if x not in active]
            print(f'Missing components: {missing}')
            return False

    ## If wipe flag provided, make sure enough components exist to do so
    elif args.action == 'eradicate' and args.wipe:
        components = VaultComponents(config)
        if not components.user or not components.cmk:
            print('Missing Vault components needed for wipe (-w) flag')
            print('  If you are resuming a previous Vault eradication, try again without wipe flag.')
            return False
    return True


def startup(args, config):
    return required_config_set_for_action(args, config) and vault_components_exist_for_action(args, config)