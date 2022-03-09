import hacc_vars
from classes.vault_installation import VaultInstallation
from hacc_add import add

    
## Setup IAM user, KMS CMK for Vault
## Optionally setup SCP for organizational account to lock down Vault access
def install(args):
    print('Installing Vault...')
    total_resources_to_create = 3 if hacc_vars.create_scp else 2

    install = VaultInstallation()

    install.create_cmk_with_alias()

    if not install.cmk:
        print('Continuing installation user even though encryption key creation failed')

    install.create_user_with_policy()


    if hacc_vars.create_scp:

        ## Make sure all Vault created by applying SCP
        if not install.cmk or not install.user:
            print('Vault resources must be fully setup before SCP can be created.')
            print('Enter command again to attempt to resume installation')
            return

        install.create_scp()

    else:
        print('SCP disabled for this installation, skipping')


    ## Determine how many resources were created during the install
    num_resources_created = len([x for x in [install.cmk, install.user, install.scp] if x != None])
    print(f'Successfully created {num_resources_created}/{total_resources_to_create} resources for Vault')

    if num_resources_created != total_resources_to_create:
        print()
        print('Vault installation completed but not all resources successfully created')
        print('Enter command again to attempt to finish installation')
        return


    print()
    print('Successfully installed Vault.')

    ## If import file provided, add credentials to new Vault
    if args.file:
        add(args)
    return