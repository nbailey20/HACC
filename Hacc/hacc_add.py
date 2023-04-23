from classes.vault import Vault
from classes.hacc_service import HaccService


def add_credential_for_service(display, vault_obj, service_name, user, passwd):
    service_obj = HaccService(service_name, vault=vault_obj)

    ## empty dicts evaluate to False
    if not bool(service_obj.credentials):
        display.update(
            display_type='text_append', 
            data={'text': f'Creating new service {service_name} in Vault'}
        )

    if not service_obj.add_credential(user, passwd):
        display.update(
            display_type='text_append', 
            data={'text': f'Username {user} already exists for service {service_name}'}
        )
        return

    service_obj.push_to_vault()
    display.update(
        display_type='text_append', 
        data={'text': f'Successfully added new username {user} for service {service_name}.'}
    )
    return



def add(display, args, config):
    vault = Vault(config)

    ## Import credentials from provided filename
    if args.file:
        display.update(
            display_type='text_new',
            data={'text': 'Importing credentials from backup file...'}
        )
        creds_list = vault.parse_import_file(args.file)

        if not creds_list:
            display.update(
                display_type='text_append',
                data={'text': 'Could not parse import file, provide valid file name created by Vault backup'}
            )
            return

        for cred in creds_list:
            add_credential_for_service(vault, cred['service'], cred['username'], cred['password'])
        display.update(
            display_type='text_append', 
            data={'text': 'Credential import complete.'}
        )


    ## Add single credential via CLI arguments
    else:
        display.update(
            display_type='text_new', 
            data={'text': 'Adding new credential...'}
        )

        service_name = args.service
        user = args.username
        password = args.password

        add_credential_for_service(display, vault, service_name, user, password)
    return