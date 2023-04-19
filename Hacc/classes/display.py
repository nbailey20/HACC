import logging
import sys

try:
    from rich.panel import Panel
    from rich.padding import Padding
    from rich.layout import Layout
    from rich.text import Text
except:
    print('Python module "rich" required for HACC client. Install (pip install rich) and try again.')
    sys.exit()


NUM_CHOICES_PER_PAGE = 9
STARTUP_PANEL_SIZE = 5

logger=logging.getLogger()

class Display:
    ## Builds main panel to display on client startup
    ## Returns tuple of (rich.padding, int) where int is size (height) of panel
    def __build_startup_panel(self, client_version):
        panel = Panel(f'Starting up...',
            title=f'[steel_blue3]HACC {client_version}',
            expand=False)
        return (Padding(panel, 1), 5)


    ## Builds exit panel to display on ctrl-c before exiting
    ## Returns tuple of (rich.padding, int) where int is size (height) of panel
    def __build_exit_panel(self, client_version):
        panel = Panel(f'Ctrl-c received, exiting.',
            title=f'[steel_blue3]HACC {client_version}',
            expand=False)
        return (Padding(panel, 1), 5)


    ## Builds credential panel for display
    ## Returns tuple of (rich.padding, int) where int is size (height) of panel
    def __build_credential_panel(self, service, user, passwd):
        panel = Panel(f'[yellow]{user} [green]: [purple]{passwd}',
            title=f'[steel_blue3]{service}',
            expand=False)
        return (Padding(panel, 1), 5)

    ## Updates self.layout based on display_type and display_data
    ## Allowed display_types: startup, credential, service, vault
    ## display_data fields =
    ##  client_version, service, user, passwd
    ## Returns None
    def update(self, display_type, display_data):
        ## Default display configs
        panel_size = 5
        self.layout['footer'].size = 2
        footer = Text('Built in Python3 by Nick Bailey')
        footer.stylize('green')

        ## Handle different display_type options
        if display_type == 'startup':
            main_panel, panel_size = self.__build_startup_panel(display_data['client_version'])

        elif display_type == 'exit':
            main_panel, panel_size = self.__build_exit_panel(display_data['client_version'])

        elif display_type == 'credential':
            main_panel, panel_size = self.__build_credential_panel(display_data['service'], display_data['user'], display_data['passwd'])

        else:
            return

        self.layout['main'].size = panel_size
        self.layout['main'].update(main_panel)
        self.layout['footer'].update(footer)



    def __init__(self):
        self.MAX_NUM = NUM_CHOICES_PER_PAGE
        self.layout = Layout()

        self.layout.split_column(
            Layout(name='main'),
            Layout(name='footer')
        )