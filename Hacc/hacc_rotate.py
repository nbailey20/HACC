import time

from hacc_add import add
from hacc_delete import delete


def rotate(display, args, config):
    delete(display, args, config)
    display.update(
        display_type='text_append',
        data={'text': 'Sleeping for 3 seconds to fully delete before updating...'}
    )
    time.sleep(3)
    add(display, args, config)
    return
