import os
import time

from console.hacc_console import console
from classes.vault_components import VaultComponents

from versions.hacc_versions import check_for_upgrades, check_for_old_versions, cleanup_old_versions




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

REQUIRED_COMMON_VARS = [
    'check_for_upgrades',
    'cleanup_old_versions'
]


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
        console.print(f'Unable to read required configuration file {hacc_config_location}, aborting.')
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


## Function to compare list of required vars against loaded vars
## Returns list containing any missing required vars, empty if none
def get_missing_vars(required, loaded):
    missing_vars = []
    for v in required:
        if not v in loaded or loaded[v] == '':
            missing_vars.append(v)
    return missing_vars


## Function to check all required config vars are set when client is invoked
def required_config_set_for_action(args, config):
    missing_vars = []

    if args.action in DATA_ACTIONS:
        missing_vars = get_missing_vars(REQUIRED_DATA_VARS, config)

    elif args.action in MGMT_ACTIONS:
        missing_vars = get_missing_vars(REQUIRED_MGMT_VARS, config)

        ## conditionally check for SCP vars if create_scp == True
        if not 'create_scp' in missing_vars and config['create_scp'] == True:
            missing_vars.append(get_missing_vars(REQUIRED_SCP_VARS, config))

    missing_vars.append(get_missing_vars(REQUIRED_COMMON_VARS, config))
    missing_vars = [item for sublist in missing_vars for item in sublist]

    for m in missing_vars:
        console.print(f'Required configuration parameter {m} not set in hacc_vars.py, please set value with hacc configure --set {m}=... (or edit manually) and try again.')

    if len(missing_vars):
        return False
    return True


## Function to confirm all required Vault components for action are setup
def vault_components_exist_for_action(args, config):
    console.print('hi')
    if args.action in DATA_ACTIONS:
        components = VaultComponents(config)
        active = components.active()
        required = components.required()

        if len(active) == 0:
            console.print('No Vault detected, execute [salmon1]hacc --install [white]or [salmon1]hacc --configure [white]before attempting this command.')
            return False

        elif len(active) != len(required):
            missing = [x for x in required if x not in active]

            console.print('Vault is not fully setup, complete installation with [salmon1]hacc --install [white]before attempting this command.')
            console.print(f'Active components: {active}')
            console.print(f'Missing components: {missing}')
            return False

    ## If wipe flag provided, make sure enough components exist to do so
    elif args.action == 'eradicate' and args.wipe:
        components = VaultComponents(config)
        if not components.user or not components.cmk:
            console.print('Missing Vault components needed for wipe (-w) flag')
            console.print('If you are resuming a previous Vault eradication, try again without wipe flag.')
            return False
    return True



## Initialization function for HACC client
## 1. Loads config variables from configuration file
## 2. Confirms all required variables for action exist
## 3. Checks for software upgrades (if check_for_upgrades == True)
## 4. Cleans up old software versions (if cleanup_old_versions == True)
## 5. Confirms all required Vault components properly setup for action
##
## Returns configuration variable dict required by client, None if error
def startup(progress, args, current_version):
    ## Get all saved config vars from config file
    get_config_task = progress.add_task("[steel_blue3]Retrieving client configuration...", total=1)
    config = get_config_params()
    progress.update(get_config_task, advance=1)
    time.sleep(0.1)

    ## Check for upgrades automatically if indicated in config file
    if config['check_for_upgrades']:
        upgrade_task = progress.add_task("[steel_blue3]Checking for available HACC upgrades...", total=1)
        new_version = check_for_upgrades(current_version)
        if new_version:
            console.print(f'New HACC version {new_version} available for installation. Upgrade with hacc --upgrade')
        progress.update(upgrade_task, advance=1)
        time.sleep(0.1)

    ## Cleanup previous HACC installations automatically if indicated in config file
    if config['cleanup_old_versions']:
        cleanup_task = progress.add_task("[steel_blue3]Cleaning up all previous HACC installations...", total=1)
        old_versions = check_for_old_versions(current_version)
        if old_versions:
            console.print(f'Found previous versions: {old_versions}')
            cleanup_old_versions(old_versions)
        progress.update(cleanup_task, advance=1)
        time.sleep(0.1)

    ## Ensure all required config vars for action exist
    required_config_task = progress.add_task("[steel_blue3]Confirming client config is valid for action...", total=1)
    if not required_config_set_for_action(args, config):
        progress.stop()
        return None
    progress.update(required_config_task, advance=1)
    time.sleep(0.1)

    ## Ensure Vault is setup correctly for action
    vault_task = progress.add_task("[steel_blue3]Confirming vault components are setup for action...", total=1)
    if not vault_components_exist_for_action(args, config):
        console.print('Vault components not setup, try [salmon1]hacc --install [white]first')
        return None
    progress.update(vault_task, advance=1)
    time.sleep(0.1)

    config['version'] = current_version
    return config