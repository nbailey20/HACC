from hacc_input import parse_args, eval_args, validate_args_for_action
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




def main():
    args = parse_args()

    if not eval_args(args):
        return 

    if args.debug:
        logger.setLevel(logging.DEBUG)
    else:
        logger.setLevel(logging.INFO)

    logger.debug(f'Initial args provided: {args}')

    valid_args = validate_args_for_action(args)
    if not valid_args:
        return False

    logger.debug(f'Validated input args: {valid_args}')

    ## Call appropriate function for action
    globals()[valid_args.action](valid_args)


if __name__ == '__main__':
    main()