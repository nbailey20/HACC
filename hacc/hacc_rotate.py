from hacc_add import add
from hacc_delete import delete
import time

def rotate(args, config):
    delete(args, config)

    print()
    print('Sleeping for 3 seconds to fully delete before updating...')
    time.sleep(3)
    print()

    add(args, config)
    return
