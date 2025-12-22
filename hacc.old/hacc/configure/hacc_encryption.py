## Code courtesy of David Sorkin, shared 2021

import io

## Cleverly map any binary data to set of 64 ASCII chars in a reversible manner
def bin_to_ascii(data_bytes):
    chars = "abcdefghijklmnopqrstuvwxyz"
    chars += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    chars += "0123456789"
    chars += "/:"
    char_list = [c for c in chars]

    ## grabbing 6 bits at time, total len including padding at end should divide
    padding = ((6 - len(data_bytes)) % 6) * b'X'

    bc = 0
    acc = 0
    f = io.StringIO()

    ## line up bytes of binary data
    for b in padding + data_bytes:
        acc |= b * 2 ** bc
        bc += 8
        ## grab 6 bits at a time and map to one of 64 chars
        while bc >= 6:
            f.write(char_list[acc & 63])
            acc >>= 6
            bc -= 6

    return f.getvalue()


## Reverse mapping from ASCII chars to binary data
def ascii_to_binary(ascii_in):
    chars = "abcdefghijklmnopqrstuvwxyz"
    chars += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    chars += "0123456789"
    chars += "/:"

    char_dict = dict()

    for i in range(len(chars)):
        char_dict[chars[i]] = i

    acc = 0
    bc = 0
    b_out = io.BytesIO()

    ## 64 possible chars takes up 6 bits, line the bits out
    for c in ascii_in:
        ## bitwise OR assignment
        acc |= char_dict[c] * 2 ** bc
        bc += 6
        ## grab 8 bits from the line and form a byte of binary
        while bc >= 8:
            b_out.write(bytes([acc & 255]))
            acc >>= 8
            bc -= 8

    ret = b_out.getvalue()

    ## Strip off leading X's
    p = 0
    while ret[p] == ord(b"X"):
        p += 1

    return ret[p:]


## XOR binary text with password
def encrypt(clear_bytes, c_str):
    cipher_bytes = bytearray(clear_bytes, 'utf-8')
    for i in range(len(clear_bytes)):
        cipher_bytes[i] ^= ord(c_str[i % len(c_str)])

    return cipher_bytes


def decrypt(cipher_bytes, c_str):
    clear_bytes = bytearray(cipher_bytes)
    for i in range(len(clear_bytes)):
        clear_bytes[i] ^= ord(c_str[i % len(c_str)])

    return clear_bytes



def decrypt_config_data(data, passwd):
    data = ascii_to_binary(data)
    data = decrypt(data, passwd)
    ## convert result to string
    return data.decode()


def encrypt_config_data(data, passwd):
    data = encrypt(data, passwd)
    data = bin_to_ascii(data)
    return data