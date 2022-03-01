from math import ceil
from hacc_generate import generate_password

NUM_CHOICES_PER_PAGE = 10

## Prints single page of numbered choices 
def print_input_choices(choices, page_num, total_pages):
    start_idx = page_num * NUM_CHOICES_PER_PAGE

    print()
    print(f'__________Page {page_num+1} / {total_pages}__________')
    print('##############################')

    for idx in range(start_idx, min(len(choices), start_idx+NUM_CHOICES_PER_PAGE)):
        print(f'[{idx+1}] {choices[idx]}')
        
    print('##############################')
    print()


## Gets user input from paginated numbered list of acceptable choices
## User can provide choice number or prefix string to match against
def get_input_with_choices(choices, choice_type):
    num_choice_pages = ceil((len(choices)*1.0) / NUM_CHOICES_PER_PAGE)
    input_val = None

    for page_num in range(num_choice_pages-1):
        print_input_choices(choices, page_num, num_choice_pages)
        input_val = input(f'Select {choice_type} name/number from above options, or hit enter to display more: ')
        if input_val:
            break
    
    ## display final page if no choice yet
    if not input_val:
        print_input_choices(choices, num_choice_pages-1, num_choice_pages)
        input_val = input(f'Select {choice_type} name/number from above options: ')

    return input_val


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
                break
            print()

    else:
        input_password = input(f'Enter password: ')
        if not input_password:
            print(f'Value for password not provided, exiting.')
            return False

        credential_password = input_password
        
    return credential_password
