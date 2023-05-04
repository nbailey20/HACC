import time
import sys

try:
    from rich.prompt import Prompt
except:
    print('Python module "rich" required for HACC. Install (pip install rich) and try again')
    sys.exit()

from console.hacc_console import console
from classes.vault import Vault
from classes.vault_eradicator import VaultEradicator

from hacc_delete import delete


def eradicate(args, config):
    ## Print scary message to prevent accidental deletion
    if args.wipe:
        console.print('This operation will delete ALL Vault credentials.')
        console.print('If you continue the credentials will be gone forever!')
    else:
        console.print('This operation schedules master key deletion for all credentials.')
        console.print('If you continue, you have up to 7 days to manually disable KMS key deletion or credentials can never be decrypted!')

    proceed = Prompt.ask('Are you sure you want to proceed?', default='n')
    if proceed.lower() != 'y' and proceed.lower() != 'yes':
        console.print('Aborting, close one ;)')
        return

    ## Wipe all credentials before Vault deletion
    if args.wipe:
        vault = Vault(config)
        delete(args, config)
        time.sleep(5)

        if len(vault.get_all_services()) != 0:
            console.print('[red]Failed to delete all credentials from Vault, aborting eradication')
            return

    ## Delete Vault
    console.print('Eradicating Vault...')

    total_resources_to_destroy = 3 if config['create_scp'] else 2
    eradicator = VaultEradicator(config)

    if config['create_scp']:
        eradicator.delete_scp()

    eradicator.delete_iam_user_with_policy(config['aws_hacc_uname'], config['aws_hacc_iam_policy'])

    if eradicator.scp:
        console.print('Cannot delete Vault master key until SCP is removed, skipping')
    else:
        eradicator.delete_cmk_with_alias(config['aws_hacc_kms_alias'])


    ## Determine how many resources were destroyed during the eradication
    num_resources_destroyed = len([x for x in [eradicator.cmk, eradicator.user, eradicator.scp] if x == None])
    if not config['create_scp']:
        num_resources_destroyed -= 1

    console.print(f'{num_resources_destroyed}/{total_resources_to_destroy} Vault components destroyed')

    if num_resources_destroyed != total_resources_to_destroy:
        console.print('Vault eradication finished but not all resources successfully destroyed')
        console.print('Retry to attempt to complete eradication')
        return

    console.print('Successfully completed Vault eradication.')
    return