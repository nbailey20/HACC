import hacc_vars
from classes.vault_installation import VaultInstallation
from hacc_add import add
import time

    
## Setup IAM user, KMS CMK for Vault
## Optionally setup SCP for organizational account to lock down Vault access
def install(args):
    print('Installing Vault...')
    total_resources_to_create = 3 if hacc_vars.create_scp else 2

    install = VaultInstallation()

    if install.cmk or install.user:
        print('  Previous installation detected, resuming')

    if not install.cmk:
        install.create_cmk_with_alias()

    if not install.user:
        install.create_user_with_policy()

    if hacc_vars.create_scp:
        # ## Make sure all Vault resources created before applying SCP
        # if not install.cmk or not install.user:
        #     print('Vault resources must be fully setup before SCP can be created.')
        #     print('Retry to attempt to resume installation')
        #     return
        
        if not install.scp:
            install.create_scp()

    else:
        print('SCP disabled for Vault installation, skipping')


    ## Determine how many resources were created during the install
    num_resources_created = len([x for x in [install.cmk, install.user, install.scp] if x != None])
    print()
    print(f'{num_resources_created}/{total_resources_to_create} Vault components ready')

    if num_resources_created != total_resources_to_create:
        print('Vault installation finished but not all resources successfully created')
        print('  Retry to attempt to resume installation')
        return

    print('Vault installation completed successfully.')
    print()

    ## If import file provided, add credentials to new Vault
    if args.file:
        print('Pausing for 15 seconds for Vault components to become active before importing credentials...')
        time.sleep(15)
        add(args)
        print()
        print('Vault install completed and credentials imported.')

    return