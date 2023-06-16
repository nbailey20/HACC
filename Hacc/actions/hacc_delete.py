import sys

try:
    from rich.panel import Panel
    from rich.padding import Padding
except:
    print('Python module "rich" required for HACC. Install (pip install rich) and try again')
    sys.exit()

from helpers.console.hacc_console import console

from classes.hacc_service import HaccService
from classes.vault import Vault


def delete_all_creds(config):
    vault = Vault(config)

    for svc_name in vault.get_all_services():
        svc_obj = HaccService(svc_name, vault=vault)

        for user in svc_obj.get_users():
            svc_obj.remove_credential(user)

        svc_obj.push_to_vault()
        console.print(f'All credentials for service [steel_blue3]{svc_name} [white]have been deleted')
    return



def delete(args, config):
    ## Wipe all Vault creds before eradication
    if args.wipe:
        console.print('Wiping all credentials from Vault...')
        delete_all_creds(config)
        return

    ## Delete single credential via input arguments
    console.print('Deleting credential...')
    service_name = args.service
    user = args.username

    svc_obj = HaccService(service_name, config=config)
    if not svc_obj.remove_credential(user):
        console.print(f'Username [yellow]{user} [white]does not exist for service [steel_blue3]{service_name}, [white]exiting.')
        return

    svc_obj.push_to_vault()
    panel = Panel(
        f'Deleted [yellow]{user} [white]from [steel_blue3]{service_name}.',
        title='[green]Success',
        expand=False,
    )
    console.print(Padding(panel, (1,0,0,0)))

    if not bool(svc_obj.credentials):
        console.print(f'No more credentials for service [steel_blue3]{service_name}.')
    return