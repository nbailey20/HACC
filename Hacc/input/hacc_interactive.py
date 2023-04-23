import sys
from math import ceil

try:
    import keyboard
except:
    print('Python module "keyboard" required for HACC client interactive mode. Install (pip install keyboard) and try again or use flags to narrow your search results.')
    sys.exit()
    
from hacc_generate import generate_password



## Gets user input from paginated numbered list of acceptable choices
## User can provide choice number or prefix string to match against
def get_input_with_choices(display, choices, choice_type, service=None):
    num_choices_per_page = display.NUM_CHOICES_PER_PAGE
    total_pages = ceil((len(choices)*1.0) / num_choices_per_page)
    curr_page = 0
    curr_idx = 0
    display_data = {
        'choices': choices,
        'choice_type': choice_type,
        'curr_page': curr_page,
        'total_pages': total_pages
    }
    if service:
        display_data['service'] = service
    
    ## Begin interactive loop
    display.update(display_type='interactive', data=display_data)
    while True:
        event = keyboard.read_event()
        if event.event_type != keyboard.KEY_DOWN:
            continue

        if event.name == 'down':
            display_data['curr_idx'] += 1
        elif event.name == 'up':
            display_data['curr_idx'] -= 1
        elif event.name == 'right':
            display_data['curr_page'] += 1
            display_data['curr_idx'] = display_data['curr_page']*num_choices_per_page
        elif event.name == 'left':
            display_data['curr_page'] -= 1
            display_data['curr_idx'] = display_data['curr_page']*num_choices_per_page
        elif event.name in [str(x) for x in range(1,10)]:
            display_data['selection'] = int(event.name) + display_data['curr_page']*num_choices_per_page - 1
            display.update(display_type='interactive', data=display_data)
            return display_data['selection'] + 1

        display.update(display_type='interactive', data=display_data)



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


## Wait until user indicates they are done with the running client
## Returns when the user presses ending sequence q
def get_input_for_end():
    ## Begin interactive loop
    while True:
        event = keyboard.read_event()
        if event.event_type != keyboard.KEY_DOWN:
            continue

        if event.name == 'q':
            return
