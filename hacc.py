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
        },
        {
            's_flag': '-p',
            'name':   'password',
            'help':   'Password for new credentials, used with add action'
        }
    ]
}

ACTION_VALID_SUBARGS = {
    'install': ['debug'],
    'eradicate': ['debug'],
    'add': ['debug', 'service', 'username', 'password'],
    'delete': ['debug', 'service', 'username'],
    'search': ['debug', 'service', 'username']
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
        parser.add_argument(
            a['s_flag'], 
            '--' + a['name'], 
            dest=a['name'], 
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


def eval_args(args):
    # Example args:
    # Namespace(action='delete', service=None, username='asdf', password='qwer', debug=True)

    # Determine if provided subarguments are valid for action
    action = args.action
    for subarg in vars(args):
        arg_val = getattr(args, subarg)
        # Looping over all args, so skip arguments that weren't provided or aren't subargs
        if subarg == 'action' or subarg == 'debug' or not arg_val:
            continue
        valid_arg = True if subarg in ACTION_VALID_SUBARGS[action] else False
        if not valid_arg:
            print('hacc --{0}: unknown option --{1}'.format(action, subarg))
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