import argparse
import re
import time

from logger.hacc_logger import logger
from console.hacc_console import console
from classes.vault import Vault
from classes.hacc_service import HaccService

from input.hacc_interactive import get_input_with_choices, get_password_for_credential, get_input_string_for_subarg
from hacc_generate import generate_password


ALLOWED_FLAGS = {
    'install': {
        's_flag':  '-i',
        'action':  'install',
        'help':    'Create new authentication credential Vault',
    },
    'eradicate': {
        's_flag':  '-e',
        'action':  'eradicate',
        'help':    'Delete entire Vault - cannot be undone'
    },
    'add': {
        's_flag':  '-a',
        'action':  'add',
        'help':    'Add new set of credentials to Vault',
    },
    'delete': {
        's_flag':  '-d',
        'action':  'delete',
        'help':    'Delete credentials in Vault',
    },
    'rotate': {
        's_flag':  '-r',
        'action':  'rotate',
        'help':    'Rotate existing password in Vault',
    },
    'backup': {
        's_flag':  '-b',
        'action':  'backup',
        'help':    'Backup entire Vault'
    },
    'configure': {
        's_flag':   '-c',
        'action':   'configure',
        'help':     'View or modify client configuration'
    },
    'upgrade': {
        's_flag':   None,
        'action':   'upgrade',
        'help':     'Upgrade client software'
    },
    'subargs': [
        {
            's_flag': '-u',
            'name':   'username',
            'help':   'Username to perform action on',
            'type':   'string'
        },
        {
            's_flag': '-p',
            'name':   'password',
            'help':   'Password for new credentials, used with add action',
            'type':   'string'
        },
        {
            's_flag': '-g',
            'name':   'generate',
            'help':   'Generate random password for add operation',
            'type':   'bool'
        },
        {
            's_flag': '-f',
            'name':   'file',
            'help':   'File name for import / export operations',
            'type':   'string'
        },
        {
            's_flag': '-w',
            'name':   'wipe',
            'help':   'Wipe all existing credentials during Vault eradication',
            'type':   'bool'
        },
        {
            's_flag': None,
            'name':   'export',
            'help':   'Export existing client configuration as encrypted file',
            'type':   'bool'
        },
        {
            's_flag': None,
            'name':   'set',
            'help':   'Set client configuration parameter',
            'type':   'string'
        },
        {
            's_flag': None,
            'name':   'show',
            'help':   'Show client configuration parameter',
            'type':   'string'
        }
    ]
}

ACTION_ALLOWED_SUBARGS = {
    'install':   ['debug', 'file'],
    'eradicate': ['debug', 'wipe'],
    'add':       ['debug', 'service', 'username', 'password', 'generate', 'file'],
    'delete':    ['debug', 'service', 'username'],
    'rotate':    ['debug', 'service', 'username', 'password', 'generate'],
    'search':    ['debug', 'service', 'username'],
    'backup':    ['debug', 'file'],
    'configure': ['debug', 'export', 'set', 'show', 'file', 'password'],
    'upgrade':   ['debug']
}

ACTION_INCOMPATABLE_SUBARGS = {
    'add': [
        ['password', 'generate'], 
        ['username', 'file'],
        ['password', 'file'],
        ['service', 'file']
    ],
    'rotate': [
        ['password', 'generate']
    ],
    'configure': [
        ['export', 'set'],
        ['export', 'show'],
        ['set', 'show'],
        ['show', 'file']
    ]
}

ACTION_REQUIRED_SUBARGS = {
    'install':   [],
    'eradicate': [],
    'add':       ['service', 'username', 'password'],
    'delete':    ['service', 'username'],
    'rotate':    ['service', 'username', 'password'],
    'search':    ['service', 'username'],
    'backup':    ['file'],
    'configure': [],
    'upgrade':   []
}

## Arguments that are required together for action TODO use this
# ARGUMENT_REQUIRED_SUBARGS = {
#     'export':    ['file']
# }

## -v/--debug arg is always allowed, redundant to list here
## same for -y/--yes, ignored if not applicable
ALLOWED_ARG_INPUT_COMBOS = [
    ['install'],
    ['install', 'file'],

    ['eradicate'],
    ['eradicate', 'wipe'],

    ['upgrade'],

    #['add'],
    #['add', 'service'],
    #['add', 'username'],
    #['add', 'password'],
    #['add', 'generate'],
    ['add', 'file'],
    #['add', 'service', 'username'],
    #['add', 'service', 'password'],
    #['add', 'service', 'generate'],
    ['add', 'username', 'password'],
    ['add', 'username', 'generate'],
    ['add', 'service', 'username', 'password'],
    ['add', 'service', 'username', 'generate'],

    ['delete'],
    ['delete', 'service'],
    ['delete', 'username'],
    ['delete', 'service', 'username'],

    #['rotate'],
    #['rotate', 'service'],
    #['rotate', 'username'],
    ['rotate', 'password'],
    ['rotate', 'generate'],
    #['rotate', 'service', 'username'],
    ['rotate', 'service', 'password'],
    ['rotate', 'service', 'generate'],
    ['rotate', 'username', 'password'],
    ['rotate', 'username', 'generate'],
    ['rotate', 'service', 'username', 'password'],
    ['rotate', 'service', 'username', 'generate'],

    ['search'],
    ['search', 'service'],
    ['search', 'username'],
    ['search', 'service', 'username'],

    #['backup'],
    ['backup', 'file'],

    ['configure', 'show'],
    ['configure', 'set'],
    #['configure', 'export'],
    ['configure', 'set', 'file', 'password'],
    ['configure', 'export', 'file']
]

HACC_SAMPLE_USAGE = [
    'Sample Usage:',
    '  hacc --configure --set aws_hacc_region=us-east-2',
    '  hacc -c --show all',
    '  hacc --install -v',
    '  hacc -a --username example@gmail.com -p 1234 testService',
    '  hacc --add testService -u test@yahoo.com --generate',
    '  hacc testService',
    '  hacc --rotate test -u test@yahoo.com -g',
    '  hacc --backup -f test_backup.txt',
    '  hacc -d testService -u example',
    '  hacc --eradicate --wipe'
]


def parse_args():
    # Define all arguments for client
    parser = argparse.ArgumentParser(
        prog='hacc',
        description='Homemade Authentication Credential Client - HACC',
        epilog='\n'.join(HACC_SAMPLE_USAGE),
        formatter_class=argparse.RawTextHelpFormatter
    )
    # Set all actions to be mutually exclusive args
    action = parser.add_argument_group('Actions')
    action_types = action.add_mutually_exclusive_group()
    
    # Add main actions, set search as default action
    for a in [ALLOWED_FLAGS[x] for x in ALLOWED_FLAGS if x != 'subargs']:
        if a['s_flag']:
            action_types.add_argument(
                a['s_flag'], 
                '--' + a['action'], 
                dest='action', 
                default='search',
                action='store_const', 
                const= a['action'], 
                help= a['help']
            )
        else:
            action_types.add_argument(
                '--' + a['action'], 
                dest='action', 
                default='search',
                action='store_const', 
                const= a['action'], 
                help= a['help']
            )

    # Add default search action
    parser.add_argument(
        'service', nargs='?', default=None, 
        help='Service name, a folder that can hold multiple credentials'
    )

    # Add optional subargs, some without single-letter abbreviations
    for a in ALLOWED_FLAGS['subargs']:
        if a['type'] == 'string':
            if a['s_flag']:
                parser.add_argument(
                    a['s_flag'], 
                    '--' + a['name'], 
                    dest=a['name'], 
                    help=a['help']
                )
            else:
                parser.add_argument(
                    '--' + a['name'], 
                    dest=a['name'], 
                    help=a['help']
                )
        elif a['type'] == 'bool':
            if a['s_flag']:
                parser.add_argument(
                    a['s_flag'], 
                    '--' + a['name'], 
                    dest=a['name'],
                    action='store_true',
                    help=a['help']
                )
            else:
                parser.add_argument(
                    '--' + a['name'], 
                    dest=a['name'],
                    action='store_true',
                    help=a['help']
                )
    
    # Add verbosity flag
    parser.add_argument(
        '-v', 
        '--verbose', 
        dest='debug', 
        action='store_true', 
        help='Display verbose output'
    )

    args = parser.parse_args()
    return args


def allowed_subargs(args):
    # Example args:
    # Namespace(action='delete', service=None, username='asdf', password='qwer', debug=True)

    for subarg in [sa for sa in vars(args) if sa != 'action' and sa != 'debug']:
        arg_val = getattr(args, subarg)
        if not arg_val:
            continue

        valid_arg = True if subarg in ACTION_ALLOWED_SUBARGS[args.action] else False
        if not valid_arg:
            console.print(f'hacc --{args.action}: unknown option --{subarg}')
            return False

    return True



def compatable_subargs(args):
    ## If action has no incompatable subargs, we're done
    if args.action not in ACTION_INCOMPATABLE_SUBARGS:
        return True

    ## Make sure args are compatable with each other
    for inc_rule in ACTION_INCOMPATABLE_SUBARGS[args.action]:
        inc_args = []
        for subarg in inc_rule:
            if getattr(args, subarg):
                inc_args.append(subarg)
        
        if len(inc_args) > 1:
            console.print(f'Incompatable arguments: {"/".join(inc_args)}, aborting.')
            return False
    return True



## Evaluate all arguments are to ensure they are allowed and compatable for given action
def eval_args(progress, args):
    allowed_task = progress.add_task("[steel_blue3]Confirming flags are allowed for action...", total=1)
    if not allowed_subargs(args):
        return False
    progress.update(allowed_task, advance=1)
    time.sleep(1)

    compatible_task = progress.add_task("[steel_blue3]Confirming flags are compatible for action...", total=1)
    if not compatable_subargs(args):
        return False
    progress.update(compatible_task, advance=1)
    time.sleep(1)

    return True



## Compares input value against acceptable values
## Returns acceptable values with same prefix as input
## Returns empty list if no possible input matches
def get_possible_input_matches(input, choices):
    if not input:
        return choices

    r = re.compile(f'{input}.*')
    matches = list(filter(r.match, choices))
    return matches



## Takes provided input service name
## Input can either match service name exactly, or be a prefix
## Returns complete service name, False if no matches in Vault
def get_service_name_from_input(svc_input, vault):

    ## Check if exact service name provided
    if vault.service_exists(svc_input):
        return svc_input

    ## Check if service name prefix provided
    possible_svcs = get_possible_input_matches(svc_input, vault.get_all_services())

    ## No possible matches
    if len(possible_svcs) == 0:
        return False

    ## Found unique match for prefix, we're done!
    elif len(possible_svcs) == 1:
        return possible_svcs[0]

    
    else:
        ## Multiple possible input matches, interactively ask for clarification
        new_input = get_input_with_choices(possible_svcs, 'Service')

        ## Check if user entered choice line number instead of value
        try:
            choice_num = int(new_input)
            new_input = possible_svcs[choice_num-1]
        except:
            choice_num = None

        ## recurse on new input
        return get_service_name_from_input(new_input, vault)



## Takes provided input username for service
## Input can either match username exactly, or be a prefix
## Returns complete username, False if no matches in Vault
def get_service_user_from_input(service, user_input, vault):
    svc_obj = HaccService(service, vault=vault)

    ## Check if exact username provided
    if svc_obj.user_exists(user_input):
        return user_input

    ## Check if username prefix provided
    possible_users = get_possible_input_matches(user_input, svc_obj.get_users())

    ## No possible matches
    if len(possible_users) == 0:
        return False

    ## Found unique match for prefix, make sure user has some confirmation before performing action!
    elif len(possible_users) == 1 and user_input:
        return possible_users[0]

    
    else:
        ## Multiple possible input matches, interactively ask for clarification
        new_input = get_input_with_choices(possible_users, 'Username', service=service)

        ## Check if user entered choice line number instead of value
        try:
            choice_num = int(new_input)
            new_input = possible_users[choice_num-1]
        except ValueError:
            choice_num = None

        ## recurse on new input
        return get_service_user_from_input(service, new_input, vault)



## Returns a list of arguments needed for action that have not been provided
## If no args missing, empty list returned
def missing_args_for_action(args):
    action = args.action
    missing_args = []

    for req_subarg in ACTION_REQUIRED_SUBARGS[action]:
        if not getattr(args, req_subarg):
            missing_args.append(req_subarg)
    return missing_args



## Ensure all required args for action are provided and valid
## Expand any subarg prefix shortcuts
## Returns a set of updated args that are ready for the given action
## Returns False if args cannot be gathered or input not valid
def validate_args_for_action(args, config):
    action = args.action
    

    # ## Handle empty Vault scenario
    # if len(vault.get_all_services()) == 0:
    #     if action == 'search' or action == 'delete' or action == 'backup':
    #         print('Vault is curently empty, add a service credential with `hacc add`')
    #         return False


    ## Validate service and username input where possible
    if action == 'search' or action == 'delete':
        vault = Vault(config)
        if len(vault.get_all_services()) == 0:
            console.print('Vault is curently empty, add a service credential with `hacc add`')
            return False

        args.service = get_service_name_from_input(args.service, vault)
        if not args.service:
            console.print(f'Could not validate [steel_blue3]service, [white]check spelling or browse all services')
            return False

        args.username = get_service_user_from_input(args.service, args.username, vault)
        if not args.username:
            console.print(f'Could not validate username for service [steel_blue3]{args.service}, [white]check spelling or browse existing usernames')
            return False

        logger.debug('Successfully validated service and username input')

    ## If importing credentials from a file, provide temp args
    ##      to meet arg requirements for add action
    if action == 'add' and args.file:
        args.username = 'file'
        args.password = 'file'
        args.service = 'file'


    ## Gather any other missing args that can't be validated, if needed
    subargs_to_get = missing_args_for_action(args)
    for subarg in subargs_to_get:

        ## Check if we should generate password
        if subarg == 'password':
            args.password = get_password_for_credential(args.generate)
            continue
        else:
            subarg_val = get_input_string_for_subarg(subarg, action)
            if not subarg_val:
                logger.debug(f'Required subarg {subarg} not provided')
                return False

            setattr(args, subarg, subarg_val)
            logger.debug(f'Interactively gathered required argument {subarg}')

    logger.debug('All required arguments provided')
    return args
