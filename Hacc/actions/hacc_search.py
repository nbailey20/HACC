import sys

try:
    from rich.padding import Padding
    from rich.panel import Panel
except:
    print('Python module "rich" required for HACC. Install (pip install rich) and try again')
    sys.exit()

from helpers.console.hacc_console import console

from classes.hacc_service import HaccService

BASE_PANEL_HEIGHT = 5
PANEL_WIDTH_PADDING = 5


## Print specific credential for service
def search(args, config):
    console.print("Retrieving credential...")
    service_name = args.service
    user = args.username

    local_service = HaccService(service_name, config=config)
    if not bool(local_service.credentials):
        console.print(f'Service [steel_blue3]{service_name} [white]does not exist in Vault, exiting.')
        return

    passwd = local_service.get_credential(user)

    panel_text = f'[yellow]{user} [green]: [purple]{passwd}'
    panel_width = len(user+' : '+passwd)+PANEL_WIDTH_PADDING
    panel = Panel(
        panel_text,
        title=f'[steel_blue3]{service_name}',
        expand=False,
        width=panel_width
    )

    console.print(Padding(panel, (1,0,0,0)))
    return
