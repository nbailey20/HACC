from console.hacc_console import console
from classes.vault import Vault
from classes.hacc_service import HaccService


def add_credential_for_service(vault_obj, service_name, user, passwd):
    service_obj = HaccService(service_name, vault=vault_obj)

    ## empty dicts evaluate to False
    if not bool(service_obj.credentials):
        console.print(f'Creating new service [steel_blue3]{service_name} [white]in Vault...')

    if not service_obj.add_credential(user, passwd):
        console.print(f'Username [yellow]{user} [white]already exists for service [steel_blue3]{service_name}.')
        return

    service_obj.push_to_vault()
    console.print(f'Successfully added new username [yellow]{user} [white]for service [steel_blue3]{service_name}.')
    return



def add(args, config):
    vault = Vault(config)

    ## Import credentials from provided filename
    if args.file:
        console.print(f'Importing credentials from file [green]{args.file}...')
        creds_list = vault.parse_import_file(args.file)

        if not creds_list:
            console.print(f'Could not parse import file [green]{args.file}, [white]provide valid file name created by Vault backup')
            return

        for cred in creds_list:
            add_credential_for_service(vault, cred['service'], cred['username'], cred['password'])
        console.print('Credential import complete.')


    ## Add single credential via CLI arguments
    else:
        console.print('Adding new credential...')

        service_name = args.service
        user = args.username
        password = args.password

        add_credential_for_service(vault, service_name, user, password)
    return