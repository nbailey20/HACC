import urllib.request
import json, random


MIN_WORD_SIZE = 7
MAX_WORD_SIZE = 10
NUM_CHAR_SWAPS = 3
RANDOM_WORDS_REQUESTED = 100

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

def sub_chars(password, char_map, num_subs):
    pass_chars = list(password)
    rand_seq = random.sample(range(len(pass_chars)), len(pass_chars))

    subs_made = 0
    idx = 0
    while subs_made < num_subs and idx < len(pass_chars):
        ## iterate over chars in random order to spread out subs
        curr_char = pass_chars[rand_seq[idx]]
        if curr_char in char_map:
            pass_chars[rand_seq[idx]] = char_map[curr_char]
            subs_made += 1
        idx += 1

    new_pass = ''.join(pass_chars)
    return {
        'password': new_pass,
        'subs_made': subs_made
    }
   


## Quick function to generate XKCD-style password,
##  4 random words joined with some char subs
def generate_password():
    random_words = json.loads(urllib.request.urlopen(f'https://random-word-api.herokuapp.com/word?number={RANDOM_WORDS_REQUESTED}').read())
    pass_words = []

    idx = 0
    while len(pass_words) < 4:
        word = random_words[idx]
        if len(word) > MIN_WORD_SIZE and len(word) < MAX_WORD_SIZE:
            pass_words.append(word)
        idx += 1

    ## CamelCase words
    pass_words = [x[0].upper()+x[1:] for x in pass_words]
    
    password = ''.join(pass_words)

    ## Sub some letters with special chars / digits
    num_digit_swaps = random.randint(0,NUM_CHAR_SWAPS)
    num_special_swaps = 3 - num_digit_swaps

    sub_obj = sub_chars(password, DIGIT_CHAR_MAP, num_digit_swaps)
    ## If not enough digit subs available, try to make extra special subs
    if sub_obj['subs_made'] < num_digit_swaps:
        num_special_swaps += num_digit_swaps - sub_obj['subs_made']

    ## If not enough special subs available, oh well
    password = sub_chars(sub_obj['password'], SPECIAL_CHAR_MAP, num_special_swaps)['password']

    return password