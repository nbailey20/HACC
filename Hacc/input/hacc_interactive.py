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

        if event.name == 'right':
            display_data['curr_page'] = (display_data['curr_page']+1) % display_data['total_pages']

        elif event.name == 'left':
            display_data['curr_page'] = (display_data['curr_page']-1) % display_data['total_pages']

        elif event.name in [str(x) for x in range(1,10)]:
            ## highlight selection to let user know it was successfully chosen before moving on
            display_data['selection'] = int(event.name) + display_data['curr_page']*num_choices_per_page - 1
            display.update(display_type='interactive', data=display_data)
            return display_data['selection'] + 1

        display.update(display_type='interactive', data=display_data)


## Waits for user to confirm or deny the current confirmation panel
## Returns True or False based on user y/n input
def get_user_confirmation(display, prompt):
    display.update(
        display_type='confirm',
        data={'text': prompt, 'selection': None}
    )

    while True:
        event = keyboard.read_event()
        if event.event_type != keyboard.KEY_DOWN:
            continue

        if event.name == 'n' or event.name == 'y':
            keyboard.send('backspace')
            display.update(
                display_type='confirm',
                data={'text': prompt, 'selection': event.name}
            )
            if event.name == 'y':
                return True
            return False


## If user_requested_generate, generate password and return it
## If user did not specify, ask and then either return generated password or user input
## Returns False if cannot gather password input
def get_password_for_credential(display):
    need_passwd = True
    while need_passwd:
        proposed_password = generate_password()
        display.update(
            display_type='text_new',
            data={'text': f'Generated password: {proposed_password}'}
        )

        if get_user_confirmation(display, prompt='Use this password?'):
            return proposed_password


## Wait until user indicates they are done with the running client
## Returns when the user presses ending sequence q
def get_input_for_end():
    keyboard.wait('q')
    keyboard.send('backspace')
    return
