import os

## Function to update symlink in /usr/bin to point to new HACC version in executable_dir
## Returns True on sucess and False otherwise
def update_bin_symlink(executable_dir):
    dest = '/usr/local/bin/hacc'
    try:
        os.symlink(executable_dir, dest)
    except:
        return False
    return True