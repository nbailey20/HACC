import json

from classes.hacc_service import HaccService
from classes.vault import Vault


## Creates new backup file with entire Vault contents
def backup(args, config):
    print('Backing up Vault data...')
    vault = Vault(config)

    all_svcs = vault.get_all_services()
    # If no services in vault, can't look anything up
    if not all_svcs:
        print('No credentials in vault, nothing to backup.')
        return

    num_svcs = len(all_svcs)
    print(f'Found {num_svcs} services with credentials to back up')

    creds_list = []
    for svc in all_svcs:
        print(f'  Backing up service {svc}')
        svc_obj = HaccService(svc, vault=vault)

        for cred in svc_obj.credentials:
            creds_list.append(
                {
                    'service': svc, 
                    'username': cred, 
                    'password': svc_obj.credentials[cred]
                }
            )

    backup_content = {'creds_list': creds_list}

    try:
        f = open(args.file, 'w')
        f.write(json.dumps(backup_content))
        f.close()
    except:
        print('Failed to write Vault data to backup file, exiting')
        return
    
    print()
    print(f'Successfully created Vault backup file: {args.file}')
    return