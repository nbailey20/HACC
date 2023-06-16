import time
import sys

try:
    from rich.prompt import Prompt
    from rich.panel import Panel
    from rich.padding import Padding
except:
    print('Python module "rich" required for HACC. Install (pip install rich) and try again')
    sys.exit()

from helpers.console.hacc_console import console

from classes.vault_installer import VaultInstaller

from actions.hacc_add import add


## Setup IAM user, KMS CMK for Vault
## Optionally setup SCP for organizational account to lock down Vault access
def install(args, config):
    console.print('Installing new Vault...')

    total_resources_to_create = 3 if config['create_scp'] else 2
    installer = VaultInstaller(config)

    if installer.cmk or installer.user or installer.scp:
        console.print('Previous installation detected, resuming')

    if not installer.cmk:
        installer.create_cmk_with_alias(config['aws_hacc_kms_alias'])

    if not installer.user:
        installer.create_user_with_policy(config['aws_hacc_uname'], config['aws_hacc_iam_policy'])

    if config['create_scp']:
        if not installer.scp:
            installer.create_scp(config['aws_hacc_scp'])

    else:
        console.print('SCP disabled for Vault installation, skipping setup')
    
    ## Determine how many resources were created during the install
    num_resources_created = len([x for x in [installer.cmk, installer.user, installer.scp] if x])
    console.rule(style='salmon1')
    console.print(f'{num_resources_created}/{total_resources_to_create} Vault components ready')

    if num_resources_created != total_resources_to_create:
        console.print('Vault installation finished but not all resources successfully created')
        console.print('Retry to attempt to resume installation')
        return

    panel = Panel(
        '[green]Vault installation completed successfully.',
        expand=False,
    )
    console.print(Padding(panel, (1,0,0,0)))

    ## If import file provided, add credentials to new Vault
    if args.file:
        console.print('Pausing for 15 seconds for Vault components to become active before importing credentials...')
        time.sleep(15)

        add(args, config)
        console.print('Credentials successfully imported into new Vault.')
    return