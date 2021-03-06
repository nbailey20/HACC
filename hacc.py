from hacc_core import get_config_params, startup
from input.hacc_input import parse_args, eval_args, validate_args_for_action

from hacc_search import search
from hacc_add import add
from hacc_delete import delete
from hacc_rotate import rotate
from hacc_install import install
from hacc_eradicate import eradicate
from hacc_backup import backup
from hacc_configure import configure

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


HACC_VERSION = 'v0.7'


def main():
    print(f'HACC {HACC_VERSION}')
    print()

    try:
        ## Parse args, ensure required config vars for action are set
        ##  and Vault setup properly for given action
        args = parse_args()

        ## Read config params from hacc_vars file into object
        config = get_config_params()
        if not config:
            return

        if not startup(args, config):
            return

        ## Ensure args are valid for action
        if not eval_args(args):
            return 

        ## Setting logging level
        if args.debug:
            logger.setLevel(logging.DEBUG)
        else:
            logger.setLevel(logging.INFO)

        logger.debug(f'Initial args provided: {args}')

        ## Validate input/gather any remaining args before passing to action
        valid_args = validate_args_for_action(args, config)
        if not valid_args:
            return False
        logger.debug(f'Validated input args: {valid_args}')

        ## Call appropriate function for action
        globals()[valid_args.action](valid_args, config)
    
    ## cleanly exit without errors
    except KeyboardInterrupt:
        print('\n\nReceived ctrl-c, exiting.')


if __name__ == '__main__':
    main()