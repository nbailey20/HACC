import sys
from math import ceil

try:
    from rich.table import Table
    from rich.padding import Padding
    from rich.prompt import Prompt
    from rich.panel import Panel
except:
    print('Python module "rich" required for HACC. Install (pip install rich) and try again')
    sys.exit()

from console.hacc_console import console, NUM_CHOICES_PER_TABLE
from hacc_generate import generate_password


## Gets user input from paginated numbered list of acceptable choices
## User can provide choice number or prefix string to match against
def get_input_with_choices(choices, choice_type, service=None):
    total_pages = ceil((len(choices)*1.0) / NUM_CHOICES_PER_TABLE)
    total_choices = len(choices)

    default_input = 'more...'
    if total_pages == 1:
        default_input = None

    ## Begin interactive loop
    user_input = default_input
    while user_input == default_input:
        choice_idx = 0
        for _ in range(total_pages):
            table = Table()
            table.add_column('#', justify='right', style='cyan', no_wrap=True)
            table.add_column(choice_type, style='magenta')

            for _ in range(min(NUM_CHOICES_PER_TABLE, total_choices-choice_idx)):
                table.add_row(str(choice_idx+1), choices[choice_idx])
                choice_idx += 1
            console.print(Padding(table, (1, 0, 0, 0)))

            default_input = 'more...'
            if choice_idx == total_choices:
                default_input = None

            if service:
                user_input = Prompt.ask(f'Enter a {choice_type}/# for service {service}', default=default_input)
            else:
                user_input = Prompt.ask(f'Enter a {choice_type}/#', default=default_input)
            
            if user_input != default_input:
                break
    return user_input



## If user_requested_generate, generate password and return it
## If user did not specify, ask and then either return generated password or user input
## Returns False if cannot gather password input
def get_password_for_credential(user_requested_generate):
    gen_password = True
    credential_password = None

    if not user_requested_generate:
        gen_password = False if Prompt.ask('Would you like to generate a [green]password?', default='y').lower() != 'y' else True

    if gen_password:
        need_passwd = True
        while need_passwd:
            proposed_password = generate_password()
            console.print(Panel(f'[steel_blue3]Generated password: [purple]{proposed_password}', expand=False))
            if Prompt.ask('Use this [green]password?', default='y').lower() == 'y':
                credential_password = proposed_password
                break

    else:
        input_password = Prompt.ask(f'Enter [green]password')
        if not input_password:
            console.print(f'[red]Value for [green]password [red]not provided, exiting.')
            return False

        credential_password = input_password
    return credential_password



## Gets free-form user input for subarg
def get_input_string_for_subarg(subarg, action):
    subarg_val = Prompt.ask(f'Enter [green]{subarg} [white]for {action}', console=console)

    if not subarg_val:
        console.print(f'[red]Value for [green]{subarg} [red]not provided, exiting.')
        return False
    return subarg_val