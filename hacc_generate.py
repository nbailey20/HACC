import secrets ## random library unsuitable to cryptographic purposes

MAX_CHAR_SWAPS = 5 ## Max number of character substitions to make in password
NUM_WORDS_IN_PASS = 4 ## XKCD-style passwords
WORDLIST_FILE = 'wordlist.txt'

DIGIT_CHAR_MAP = {
    'b': '6',
    'e': '3',
    'i': '1',
    'o': '0',
    't': '7'
}

SPECIAL_CHAR_MAP = {
    'a': '@',
    'i': '!',
    's': '$'
}

## Returns number of lines in wordlist file
def get_num_words(filename):
    with open(filename) as f:
        for i, _ in enumerate(f):
            pass
    return i+1


## Perform random character substitions in given password
def sub_chars(password, char_map, num_subs):
    pass_chars = list(password)
    num_pass_chars = len(pass_chars)

    subs_made = 0
    max_attempts = 10000
    num_attempts = 0
    ## if we can't find a char to sub after 10000 attempts, give up
    while subs_made < num_subs and num_attempts < max_attempts:
        
        ## choose random char in password
        random_int = secrets.randbelow(num_pass_chars)
        random_char = pass_chars[random_int]
        if random_char in char_map:
            pass_chars[random_int] = char_map[random_char]
            subs_made += 1

        num_attempts += 1

    new_pass = ''.join(pass_chars)
    return {
        'password': new_pass,
        'subs_made': subs_made
    }
   


## Quick function to generate XKCD-style password,
##  4 random words joined with some char subs
def generate_password():
    pass_words = []
    wordlist_len = get_num_words(WORDLIST_FILE)

    for _ in range(NUM_WORDS_IN_PASS):
        f = open(WORDLIST_FILE)
        ## get random line number for wordlist
        line_num = secrets.randbelow(wordlist_len)

        ## add word at line_num to password
        for i, line in enumerate(f):
            if i == line_num:
                pass_words.append(line.strip())
            elif i > line_num:
                break
        f.close()

    ## CamelCase words
    pass_words = [x[0].upper()+x[1:] for x in pass_words]
    
    password = ''.join(pass_words)

    ## Sub random number of letters with special chars / digits
    num_char_swaps = secrets.randbelow(MAX_CHAR_SWAPS+1)
    num_digit_swaps = secrets.randbelow(num_char_swaps+1)
    num_special_swaps = num_char_swaps - num_digit_swaps

    sub_obj = sub_chars(password, DIGIT_CHAR_MAP, num_digit_swaps)
    ## If not enough digit subs available, try to make extra special subs
    if sub_obj['subs_made'] < num_digit_swaps:
        num_special_swaps += num_digit_swaps - sub_obj['subs_made']

    ## If not enough special subs available, oh well
    password = sub_chars(sub_obj['password'], SPECIAL_CHAR_MAP, num_special_swaps)['password']

    return password