import sys

try:
    from rich import print
    from rich.live import Live
    from rich.panel import Panel
    from rich.padding import Padding
    from rich.prompt import Prompt
    from rich.layout import Layout
    from rich.text import Text
except:
    print('Python module "rich" required for HACC client. Install (pip install rich) and try again.')
    sys.exit()

from math import ceil
from hacc_generate import generate_password

NUM_CHOICES_PER_PAGE = 9


## Creates single page of numbered choices with current idx highlighted
def build_input_choices_panel(choices, choice_type, curr_idx, page_num, total_pages):
    start_idx = page_num * NUM_CHOICES_PER_PAGE
    panel_text = ""

    for idx in range(start_idx, min(len(choices), start_idx+NUM_CHOICES_PER_PAGE)):
        display_idx = idx % NUM_CHOICES_PER_PAGE + 1
        if idx == curr_idx:
            panel_text += f'[gold1][{display_idx}] [gold1]{choices[idx]}\n'
        else:
            panel_text += f'[purple3][{display_idx}] [green]{choices[idx]}\n'
   # panel_text += Prompt.ask(f'Select {choice_type} name/number from above options, or hit enter to display more:')

    panel_type = choice_type[0].upper()+choice_type[1:]
    if len(choices) > 1:
        panel_type += 's'
    panel = Panel(panel_text,
                  title=f'[steel_blue3]{panel_type} {page_num+1} / {total_pages}',
                  expand=False)
    return Padding(panel, 1)


## Gets user input from paginated numbered list of acceptable choices
## User can provide choice number or prefix string to match against
def get_input_with_choices(choices, choice_type):
    try:
        import keyboard
    except:
        print('Python module "keyboard" required for HACC client interactive mode. Install (pip install keyboard) and try again or use flags to narrow your search results.')
        sys.exit()

    num_choice_pages = ceil((len(choices)*1.0) / NUM_CHOICES_PER_PAGE)
    #input_val = ""

    ## add 4 to account for panel space
    service_size = NUM_CHOICES_PER_PAGE+4
    if len(choices) < NUM_CHOICES_PER_PAGE:
        service_size = len(choices)+4

    layout = Layout()
    layout.split(
        
        Layout(name='services', size=service_size),
        Layout(name='instruction', size=2)
    )

    instruction_text = Text(f'Select {choice_type} with arrow/# keys and press enter.')
    instruction_text.stylize('steel_blue3')
    layout['instruction'].update(instruction_text)

    with Live(layout):
        curr_page = 0
        curr_idx = 0
        padded_panel = build_input_choices_panel(choices, choice_type, curr_idx, curr_page, num_choice_pages)
        layout['services'].update(padded_panel)
        layout['instruction'].update(instruction_text)

        ## Begin interactive loop
        while True:
            event = keyboard.read_event()
            if event.event_type != keyboard.KEY_DOWN:
                continue

            if event.name == 'down':
                curr_idx += 1
            elif event.name == 'up':
                curr_idx -= 1
            elif event.name == 'right':
                curr_page += 1
                curr_idx = curr_page*NUM_CHOICES_PER_PAGE
            elif event.name == 'left':
                curr_page -= 1
                curr_idx = curr_page*NUM_CHOICES_PER_PAGE
            elif event.name == 'enter':
                return curr_idx+1

            padded_panel = build_input_choices_panel(choices, choice_type, curr_idx, curr_page, num_choice_pages)
            layout["services"].update(padded_panel)



## Gets free-form user input for subarg
def get_input_string_for_subarg(subarg, action):
    subarg_val = input(f'Enter {subarg} for {action}: ')

    if not subarg_val:
        print(f'Value for {subarg} not provided, exiting.')
        return False
    return subarg_val


## If user_requested_generate, generate password and return it
## If user did not specify, ask and then either return generated password or user input
## Returns False if cannot gather password input
def get_password_for_credential(user_requested_generate):
    gen_password = True
    credential_password = None

    if not user_requested_generate:
        gen_password = False if input('Would you like to generate a password (y/n)? ').lower() != 'y' else True

    if gen_password:
        need_passwd = True
        while need_passwd:
            proposed_password = generate_password()
            print(f'Generated password: {proposed_password}')
            if input('Use this password (y/n)? ').lower() == 'y':
                credential_password = proposed_password
                print()
                break
            print()

    else:
        input_password = input(f'Enter password: ')
        if not input_password:
            print(f'Value for password not provided, exiting.')
            return False

        credential_password = input_password
        
    return credential_password
