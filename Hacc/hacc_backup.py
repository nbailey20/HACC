import json

from classes.hacc_service import HaccService
from classes.vault import Vault


## Creates new backup file with entire Vault contents
def backup(display, args, config):
    display.update(
        display_type='text_new',
        data={'text': 'Backing up Vault data...'}
    )
    vault = Vault(display, config)

    all_svcs = vault.get_all_services()
    # If no services in vault, can't look anything up
    if not all_svcs:
        display.update(
            display_type='text_append',
            data={'text': 'No credentials in vault, nothing to backup.'}
        )
        return

    num_svcs = len(all_svcs)
    display.update(
        display_type='text_append',
        data={'text': f'Found {num_svcs} services with credentials to back up'}
    )

    creds_list = []
    for svc in all_svcs:
        display.update(
            display_type='text_append',
            data={'text': f'Backing up service {svc}'}
        )
        svc_obj = HaccService(display, svc, vault=vault)

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
        with open(args.file, 'w') as f:
            f.write(json.dumps(backup_content))
    except:
        display.update(
            display_type='text_append',
            data={'text': 'Failed to write Vault data to backup file, exiting'}
        )
        return

    display.update(
        display_type='text_append',
        data={'text': f'Successfully created Vault backup file: {args.file}'}
    )
    return