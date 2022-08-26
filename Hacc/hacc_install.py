import time

from classes.vault_installer import VaultInstaller
from hacc_add import add


## Setup IAM user, KMS CMK for Vault
## Optionally setup SCP for organizational account to lock down Vault access
def install(args, config):
    print('Installing Vault...')
    total_resources_to_create = 3 if config['create_scp'] else 2

    installer = VaultInstaller(config)

    if installer.cmk or installer.user or installer.scp:
        print('  Previous installation detected, resuming')

    if not installer.cmk:
        installer.create_cmk_with_alias(config['aws_hacc_kms_alias'])

    if not installer.user:
        installer.create_user_with_policy(config['aws_hacc_uname'])

    if config['create_scp']:
        if not installer.scp:
            installer.create_scp(config['aws_hacc_scp'])

    else:
        print('SCP disabled for Vault installation, skipping')


    ## Determine how many resources were created during the install
    num_resources_created = len([x for x in [installer.cmk, installer.user, installer.scp] if x != None])
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