import argparse
from argparse import RawTextHelpFormatter

from hacc_search import search
from hacc_add_remove_param import add, delete
from hacc_install import install
from hacc_uninstall import eradicate


ALLOWED_FLAGS = {
    'install': {
        's_flag':  '-i',
        'action':  'install',
        'help':    'Create new authentication credential vault',
    },
    'eradicate': {
        's_flag':  '-e',
        'action':  'eradicate',
        'help':    'Delete entire vault - cannot be undone'
    },
    'add': {
        's_flag':  '-a',
        'action':  'add',
        'help':    'Add new set of credentials to vault',
    },
    'delete': {
        's_flag':  '-d',
        'action':  'delete',
        'help':    'Delete credentials in vault',
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
            'help':   'Generate random password for operation',
            'type':   'bool'
        },
        {
            's_flag': '-b',
            'name':   'backup',
            'help':   'Backup entire Vault and write to file name: hacc -b out_file',
            'type':   'string'
        }
    ]
}

ACTION_VALID_SUBARGS = {
    'install': ['debug'],
    'eradicate': ['debug'],
    'add': ['debug', 'service', 'username', 'password', 'generate'],
    'delete': ['debug', 'service', 'username'],
    'search': ['debug', 'service', 'username', 'backup'],
}

INCOMPATABLE_SUBARGS = {
    'add': [['password', 'generate']]
}

def parse_args():
    # Define all arguments for client
    parser = argparse.ArgumentParser(
        prog='hacc',
        description='Homemade Authentication Credential Client - HACC',
        epilog='Sample Usage:\n  hacc -iv\n  hacc -a -u example@gmail.com -p 1234 testService\n  hacc testService\n  hacc -d testService\n  hacc -e -v',
        formatter_class=RawTextHelpFormatter
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


def valid_subargs(args):
    # Example args:
    # Namespace(action='delete', service=None, username='asdf', password='qwer', debug=True)

    for subarg in [sa for sa in vars(args) if sa != 'action' and sa != 'debug']:
        arg_val = getattr(args, subarg)
        if not arg_val:
            continue

        valid_arg = True if subarg in ACTION_VALID_SUBARGS[args.action] else False
        if not valid_arg:
            print('hacc --{0}: unknown option --{1}'.format(args.action, subarg))
            return False

    return True



def compatable_subargs(args):
    ## check if action has incompatable flags
    if args.action not in INCOMPATABLE_SUBARGS:
        return True

    ## Check each incompatable rule for given action
    for inc_rule in INCOMPATABLE_SUBARGS[args.action]:
        inc_args = []
        for subarg in inc_rule:
            if getattr(args, subarg):
                inc_args.append(subarg)
        
        if len(inc_args) > 1:
            print(f'Incompatable arguments: {"/".join(inc_args)}, aborting.')
            return False
    return True


## Evaluate if parsed arguments are valid and compatable
def eval_args(args):
    if not valid_subargs(args):
        return False

    if not compatable_subargs(args):
        return False
    return True


def main():
    args = parse_args()
    if args.debug: print('Args provided: ', args)

    if not eval_args(args):
        return 
    
    globals()[args.action](args)


if __name__ == '__main__':
    main()