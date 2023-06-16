import time

from helpers.console.hacc_console import console

from actions.hacc_add import add
from actions.hacc_delete import delete


def rotate(args, config):
    delete(args, config)
    console.print('Sleeping for 3 seconds to fully delete before updating...')
    time.sleep(3)
    add(args, config)
    return
