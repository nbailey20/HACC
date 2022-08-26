import time

from classes.vault_eradicator import VaultEradicator
from classes.vault import Vault
from hacc_delete import delete


def eradicate(args, config):
    ## Print scary message to prevent accidental deletion
    if args.wipe:
        print('This operation will delete ALL Vault credentials.')
        print('If you continue the credentials will be gone forever!')
    else:
        print('This operation schedules master key deletion for all credentials.')
        print('If you continue, you have up to 7 days to manually disable KMS key deletion or credentials can never be decrypted!')
    
    proceed = True if input('Are you sure you want to proceed (y/n)? ') == 'y' else False
    if not proceed:
        print('Aborting, close one ;)')
        return

    ## Wipe all credentials before Vault deletion
    if args.wipe:
        vault = Vault(config)
        delete(args)
        time.sleep(5)
        if len(vault.get_all_services()) != 0:
            print('Failed to delete all credentials from Vault, aborting eradication due to user providing wipe argument')
            return

    ## Delete Vault
    print()
    print('Eradicating Vault...')
    total_resources_to_destroy = 3 if config['create_scp'] else 2

    eradicator = VaultEradicator(config)

    if config['create_scp']:
        eradicator.delete_scp()

    eradicator.delete_iam_user_with_policy(config['aws_hacc_uname'], config['aws_hacc_iam_policy'])

    if eradicator.scp:
        print('Cannot delete Vault master key until SCP is removed')
    else:
        eradicator.delete_cmk_with_alias(config['aws_hacc_kms_alias'])


    ## Determine how many resources were destroyed during the eradication
    num_resources_destroyed = len([x for x in [eradicator.cmk, eradicator.user, eradicator.scp] if x == None])
    if not config['create_scp']:
        num_resources_destroyed -= 1
        
    print()
    print(f'{num_resources_destroyed}/{total_resources_to_destroy} Vault components destroyed')

    if num_resources_destroyed != total_resources_to_destroy:
        print('Vault eradication finished but not all resources successfully destroyed')
        print('  Retry to attempt to complete eradication')
        return

    print()
    print('Successfully completed Vault eradication.')
    return