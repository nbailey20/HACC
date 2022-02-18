from hacc_core import get_all_services, get_kms_arn, HaccService, ApiClient
import json
import logging

logger=logging.getLogger()


## Creates new backup file with entire Vault contents
def backup(args):
    api_obj = ApiClient(ssm=True, kms=True)
    ssm_client = api_obj.ssm
    kms_id = get_kms_arn(api_obj.kms)

    all_svcs = get_all_services(ssm_client)
    # If no services in vault, can't look anything up
    if not all_svcs:
        print('No credentials in vault, nothing to backup.')
        return

    backup_content = []
    for svc in all_svcs:
        logger.debug(f'Backing up service {svc}')
        svc_obj = HaccService(svc, ssm_client=ssm_client, kms_id=kms_id)
        backup_content.append({'service': svc, 'creds': svc_obj.credentials})

    f = open(args.outfile, 'w')
    f.write(json.dumps(backup_content))
    f.close()
    
    print(f'Successfully created Vault backup file: {args.outfile}')
    return