import time

from console.hacc_console import console

from hacc_add import add
from hacc_delete import delete


def rotate(args, config):
    delete(args, config)
    console.print('Sleeping for 3 seconds to fully delete before updating...')
    time.sleep(3)
    add(args, config)
    return
