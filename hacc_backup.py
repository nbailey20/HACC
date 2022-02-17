import hacc_vars
from hacc_core import get_all_services, HaccService
import boto3
import json
import logging

logger=logging.getLogger(__name__)


## Creates new backup file with entire Vault contents
def backup(args):
    hacc_session = boto3.session.Session(profile_name=hacc_vars.aws_hacc_uname)
    ssm_client = hacc_session.client('ssm', region_name=hacc_vars.aws_hacc_region)


    all_svcs = get_all_services(ssm_client)
    # If no services in vault, can't look anything up
    if not all_svcs:
        print('No credentials in vault, nothing to backup.')
        return

    backup_content = []
    for svc in all_svcs:
        logger.debug(f'Backing up service {svc}')
        svc_obj = HaccService(svc)
        backup_content.append({'service': svc, 'creds': svc_obj.credentials})

    f = open(args.outfile, 'w')
    f.write(json.dumps(backup_content))
    f.close()
    
    print(f'Successfully created Vault backup file: {args.outfile}')
    return