import argparse
import hacc_interactive
from hacc_search import search
from hacc_add import add
from hacc_delete import delete
from hacc_install import install
from hacc_uninstall import eradicate
from hacc_backup import backup
import logging


 ## Configure logger
logging.basicConfig()

## suppress lots of noisy logs
logging.getLogger('boto3').setLevel(logging.CRITICAL) 
logging.getLogger('botocore').setLevel(logging.CRITICAL)
logging.getLogger('nose').setLevel(logging.CRITICAL)
logging.getLogger('urllib3').setLevel(logging.CRITICAL)

logger=logging.getLogger()
logger.handlers[0].setFormatter(logging.Formatter(
    'DEBUG: %(message)s'
))


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
    'backup': {
        's_flag':  '-b',
        'action':  'backup',
        'help':    'Backup entire Vault'
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
            's_flag': '-o',
            'name':   'outfile',
            'help':   'File name to use for backup operation',
            'type':   'string'
        }
    ]
}

ACTION_ALLOWED_SUBARGS = {
    'install': ['debug'],
    'eradicate': ['debug'],
    'add': ['debug', 'service', 'username', 'password', 'generate'],
    'delete': ['debug', 'service', 'username'],
    'search': ['debug', 'service', 'username'],
    'backup': ['debug', 'outfile']
}

ACTION_INCOMPATABLE_SUBARGS = {
    'add': [['password', 'generate']]
}

ACTION_REQUIRED_SUBARGS = {
    'install': [],
    'eradicate': [],
    'add': ['service', 'username', 'password'],
    'delete': ['service', 'username'],
    'search': ['service', 'username'],
    'backup': ['outfile']
}


def parse_args():
    # Define all arguments for client
    parser = argparse.ArgumentParser(
        prog='hacc',
        description='Homemade Authentication Credential Client - HACC',
        epilog='Sample Usage:\n  hacc -iv\n  hacc -a -u example@gmail.com -p 1234 testService\n  hacc testService\n  hacc -d testService\n  hacc -e -v',
        formatter_class=argparse.RawTextHelpFormatter
    )
    # Set all actions to be mutually exclusive args
    action = parser.add_argument_group('Actions')
    action_types = action.add_mutually_exclusive_group()
    
    # Add main actions, set search as default action
    for a in [ALLOWED_FLAGS[x] for x in ALLOWED_FLAGS if x != 'subargs']:
        action_types.add_argument(
            a['s_flag'], 
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

    # Add optional subargs
    for a in ALLOWED_FLAGS['subargs']:
        if a['type'] == 'string':
            parser.add_argument(
                a['s_flag'], 
                '--' + a['name'], 
                dest=a['name'], 
                help=a['help']
            )
        elif a['type'] == 'bool':
            parser.add_argument(
                a['s_flag'], 
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
            print('hacc --{0}: unknown option --{1}'.format(args.action, subarg))
            return False

    return True



def compatable_subargs(args):
    ## check if action has incompatable flags
    if args.action not in ACTION_INCOMPATABLE_SUBARGS:
        return True

    ## Check each incompatable rule for given action
    for inc_rule in ACTION_INCOMPATABLE_SUBARGS[args.action]:
        inc_args = []
        for subarg in inc_rule:
            if getattr(args, subarg):
                inc_args.append(subarg)
        
        if len(inc_args) > 1:
            print(f'Incompatable arguments: {"/".join(inc_args)}, aborting.')
            return False
    return True



## Evaluate all arguments are to ensure they allowed and compatable for given action
def eval_args(args):
    if not allowed_subargs(args):
        return False

    if not compatable_subargs(args):
        return False
    
    return True


## If not all required args provided, interactively ask user
## Return validated arguments object for action logic
## Return False if unable to gather required input
def collect_missing_args(args):
    action = args.action

    for subarg in ACTION_REQUIRED_SUBARGS[action]:
        if not getattr(args, subarg):
            
            if subarg == 'service':
                subarg_val = hacc_interactive.get_service()
            ## if subarg == user, display possible users with numbering
            ## if subarg == password, ask to generate one

            else:
                subarg_val = input(f'Enter {subarg} to {action}: ')
                if not subarg_val:
                    print(f'Value for {subarg} not provided, exiting.')
                    return False
            
            logger.debug(f'Interactively gathered required argument {subarg}')
            setattr(args, subarg, subarg_val)

    logger.debug('All required arguments provided')
    return args



def main():
    args = parse_args()

    if args.debug:
        logger.setLevel(logging.DEBUG)
    else:
        logger.setLevel(logging.INFO)

    logger.debug(f'Args provided: {args}')


    if not eval_args(args):
        return 

    required_args = collect_missing_args(args)
    
    globals()[required_args.action](required_args)


if __name__ == '__main__':
    main()