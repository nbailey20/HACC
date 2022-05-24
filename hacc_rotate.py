from hacc_add import add
from hacc_delete import delete
import time

def rotate(args):
    delete(args)

    print()
    print('Sleeping for 3 seconds before updating password...')
    time.sleep(3)
    print()

    add(args)
    return
