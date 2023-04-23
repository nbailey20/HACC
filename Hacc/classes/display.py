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
BASE_PANEL_HEIGHT = 5
PANEL_WIDTH_PADDING = 5
FOOTER_HEIGHT = 2
DEFAULT_FOOTER_TEXT = "Just HACC it."

logger=logging.getLogger()

class Display:
    ## Builds main panel to display on client startup
    ## Returns tuple of (rich.padding, int, int, rich.Text)
    ## Ints are height and width of panel
    ## Text is footer of layout
    def __build_text_layout(self, text, append=False):
        if append:
            text = self.panel_text + '\n' + text

        panel = Panel(
            text,
            title=f'[steel_blue3]HACC {self.client_version}',
            expand=False
        )

        max_line_width = 0
        for line in text.split('\n'):
            line_len = len(line)
            if line_len > max_line_width:
                max_line_width = line_len

        self.panel_text = text
        self.panel_width = max_line_width + PANEL_WIDTH_PADDING
        panel_height = text.count('\n') + BASE_PANEL_HEIGHT

        footer = Text(DEFAULT_FOOTER_TEXT)
        footer.stylize('green')
        footer.align('center', self.panel_width)

        return (Padding(panel, 1), panel_height, self.panel_width, footer)


    ## Builds credential panel for display
    ## Returns tuple of (rich.padding, int, int) where ints are the height and width of panel
    def __build_credential_layout(self, service, user, passwd):
        panel_text = f'[yellow]{user} [green]: [purple]{passwd}'
        panel = Panel(
            panel_text,
            title=f'[steel_blue3]{service}',
            expand=False
        )

        self.panel_text = panel_text
        self.panel_width = len(user+' : '+passwd)+PANEL_WIDTH_PADDING

        footer = Text(DEFAULT_FOOTER_TEXT)
        footer.stylize('red')
        footer.align('center', self.panel_width)

        return (Padding(panel, 1), BASE_PANEL_HEIGHT, self.panel_width, footer)


    ## Creates single page of numbered choices with current idx highlighted
    def __build_interactive_layout(self, choices, choice_type, page_num, total_pages, service=None, selection=None):
        start_idx = page_num * NUM_CHOICES_PER_PAGE
        panel_text = ""
        max_line_len = 0

        for idx in range(start_idx, min(len(choices), start_idx+NUM_CHOICES_PER_PAGE)):
            display_idx = idx % NUM_CHOICES_PER_PAGE + 1

            line_len = len(f'[{display_idx}] {choices[idx]}')
            if line_len > max_line_len:
                max_line_len = line_len

            if selection and selection == idx:
                panel_text += f'[gold1][{display_idx}] [gold1]{choices[idx]}\n'
            else:
                panel_text += f'[purple3][{display_idx}] [green]{choices[idx]}\n'

        self.panel_text = panel_text
        self.panel_width = max_line_len+PANEL_WIDTH_PADDING

        if service:
            panel = Panel(
                panel_text,
                title=f'[steel_blue3]{service}',
                subtitle=f'[steel_blue3]{page_num+1} / {total_pages}',
                expand=False
            )
        else:
            panel = Panel(
                panel_text,
                title=f'[steel_blue3]HACC {self.client_version}',
                subtitle=f'[steel_blue3]{page_num+1} / {total_pages}',
                expand=False
            )

        if total_pages > 1:
            footer_text = f'Use left/right arrows to browse {choice_type}s and # keys to make a selection.'
        else:
            footer_text = f'Use # keys to select a {choice_type}.'
        if len(footer_text) > max_line_len:
            self.panel_width = len(footer_text)

        footer = Text(footer_text)
        footer.stylize('gold1')
        footer.align('center', self.panel_width)

        return (Padding(panel, 1), NUM_CHOICES_PER_PAGE+BASE_PANEL_HEIGHT, self.panel_width, footer)


    ## Updates self.layout based on display_type and display_data
    ## Allowed display_types: text_new, text_append, credential, interactive, end
    ## display_data fields =
    ##  text, service, user, passwd
    ## Returns None
    def update(self, display_type, data):
        ## Default display configs
        self.layout['footer'].size = FOOTER_HEIGHT
        footer = Text('Just HACC it.')

        ## Handle different display_type options
        if display_type == 'text_new':
            main_panel, panel_height, self.panel_width, footer = self.__build_text_layout(data['text'])
            self.layout['main'].size = panel_height
            self.layout['main'].update(main_panel)

        elif display_type == 'text_append':
            main_panel, panel_height, self.panel_width, footer = self.__build_text_layout(data['text'], append=True)
            self.layout['main'].size = panel_height
            self.layout['main'].update(main_panel)

        elif display_type == 'credential':
            main_panel, panel_height, self.panel_width, footer = self.__build_credential_layout(
                data['service'],
                data['user'],
                data['passwd']
            )
            self.layout['main'].size = panel_height
            self.layout['main'].update(main_panel)

        elif display_type == 'interactive':
            main_panel, panel_height, self.panel_width, footer = self.__build_interactive_layout(
                data['choices'],
                data['choice_type'],
                data['curr_page'],
                data['total_pages'],
                service = None if not 'service' in data else data['service'],
                selection = None if not 'selection' in data else data['selection']
            )
            self.layout['main'].size = panel_height
            self.layout['main'].update(main_panel)

        elif display_type == 'end':
            footer = Text('Press q to exit.')
            footer.stylize('green')
            self.layout['footer'].update(footer)
            footer.align('center', self.panel_width)

        self.layout['footer'].update(footer)



    def __init__(self, client_version):
        self.layout = Layout()
        self.NUM_CHOICES_PER_PAGE = NUM_CHOICES_PER_PAGE
        self.client_version = client_version
        self.panel_text = None
        self.panel_width = None

        self.layout.split_column(
            Layout(name='main'),
            Layout(name='footer')
        )