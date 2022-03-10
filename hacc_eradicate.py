import hacc_vars
from classes.vault_installation import VaultInstallation
from classes.vault import Vault
from hacc_delete import delete
import time



def eradicate(args):

    vault = Vault()
    if len(vault.get_all_services()) != 0:

        ## Print scary message to prevent accidental deletion
        if args.wipe:
            print('This operation will delete ALL Vault credentials.')
            print('If you continue the credentials will be gone forever!')
        else:
            print('This operation schedules master key deletion for all credentials.')
            print('If you continue, you have up to 7 days to manually disable deletion or credentials can never be decrypted!')
        
        proceed = True if input('Are you sure you want to proceed (y/n)? ') == 'y' else False
        if not proceed:
            print('Aborting.')
            return

        ## Wipe all credentials before Vault deletion
        if args.wipe:
            delete(args)
            time.sleep(5)
            if len(vault.get_all_services()) != 0:
                print('Failed to delete all credentials from Vault, aborting eradication due to wipe argument')
                return


    ## Vault already empty, less scary of a warning :)
    else:
        proceed = True if input('Are you sure you want to remove empty Vault (y/n)? ') == 'y' else False
        if not proceed:
            print('Aborting.')
            return


    ## Delete Vault
    print()
    print('Eradicating Vault...')
    total_resources_to_destroy = 3 if hacc_vars.create_scp else 2

    eradicate = VaultInstallation()

    if hacc_vars.create_scp:
        eradicate.delete_scp()

    if eradicate.scp:
        print('Cannot continue eradicating vault until SCP is removed')
        print('Retry to attempt to resume eradication')
        return

    eradicate.delete_cmk()
    eradicate.delete_user()


    ## Determine how many resources were destroyed during the eradication
    num_resources_destroyed = len([x for x in [eradicate.cmk, eradicate.user, eradicate.scp] if x == None])
    if not hacc_vars.create_scp:
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