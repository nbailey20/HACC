from classes.vault_components import VaultComponents
import hacc_vars

MGMT_ACTIONS = ['install', 'eradicate']
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


## Function to check all required config vars are set when client is invoked
def required_config_set_for_action(args):
    all_vars_exist = True
    missing_vars = []

    if args.action in DATA_ACTIONS:
        for v in REQUIRED_DATA_VARS:
            if not hasattr(hacc_vars, v) or getattr(hacc_vars, v) == '':
                missing_vars.append(v)
                all_vars_exist = False

    elif args.action in MGMT_ACTIONS:
        for v in REQUIRED_MGMT_VARS:
            if not hasattr(hacc_vars, v) or getattr(hacc_vars, v) == '':
                missing_vars.append(v)
                all_vars_exist = False

        ## conditionally check for SCP vars if create_scp == True
        if not 'create_scp' in missing_vars and hacc_vars.create_scp == True:
            for v in REQUIRED_SCP_VARS:
                if not hasattr(hacc_vars, v) or getattr(hacc_vars, v) == '':
                    missing_vars.append(v)
                    all_vars_exist = False

    for m in missing_vars:
        print(f'Required configuration parameter {m} not set in hacc_vars.py, please set value and try again.')

    return all_vars_exist


## Function to confirm all required Vault components for action are setup
def vault_components_exist_for_action(args):
    if args.action in DATA_ACTIONS:
        components = VaultComponents()
        active = components.active()
        required = components.required()

        if len(active) == 0:
            print('No Vault detected, install before attempting this command.')
            return False

        elif len(active) != len(required):
            print('Vault is not fully setup, complete installation before attempting this command.')
            return False

    ## If wipe flag provided, make sure enough components exist to do so
    elif args.action == 'eradicate' and args.wipe:
        components = VaultComponents()
        if not components.user or not components.cmk:
            print('Missing Vault components needed for wipe (-w) flag')
            print('  If you are resuming a previous Vault eradication, try again without wipe flag.')
            return False
    return True


def startup(args):
    return required_config_set_for_action(args) and vault_components_exist_for_action(args)