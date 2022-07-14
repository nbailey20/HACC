from classes.vault import Vault
from classes.hacc_service import HaccService


def add_credential_for_service(vault_obj, service_name, user, passwd):
    service_obj = HaccService(service_name, vault=vault_obj)

    ## empty dicts evaluate to False
    if not bool(service_obj.credentials):
        print(f'Creating new service {service_name} in Vault')

    if not service_obj.add_credential(user, passwd):
        print(f'  Username {user} already exists for service {service_name}')
        return

    service_obj.push_to_vault()
    print(f'Successfully added new username {user} for service {service_name}.')
    return



def add(args, config):
    vault = Vault()

    ## Import credentials from provided filename
    if args.file:
        print('Importing credentials from backup file...')
        creds_list = vault.parse_import_file(args.file)

        if not creds_list:
            print('Could not parse import file, provide valid file name created by Vault backup')
            return

        for cred in creds_list:
            add_credential_for_service(vault, cred['service'], cred['username'], cred['password'])
        print('\nCredential import complete.')


    ## Add single credential via CLI arguments
    else:
        print('Adding new credential...')

        service_name = args.service
        user = args.username
        password = args.password

        add_credential_for_service(vault, service_name, user, password)
    return