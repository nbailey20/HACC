import sys
import os
import requests
import shutil

try:
    from packaging import version
except:
    print('Python module "packaging" required for HACC client. Install (pip install packaging) and try again.')
    sys.exit()

from helpers.console.hacc_console import console



## Function that takes in a list of Hacc versions + current version and returns either:
## 1. newest version newer than current version (can be empty) OR
## 2. list of versions older than the current version (can be empty)
def compare_hacc_versions(version_list, current_version, operation='newest'):
    current_version_obj = version.parse(current_version)

    if operation == 'newest':
        newest = {'obj': current_version_obj, 'str': current_version}
        for v in version_list:
            try:
                v_obj = version.parse(v)
            except:
                console.print('[red]Error parsing version from remote list.')
                continue

            if v_obj > newest['obj']:
                newest['obj'] = v_obj
                newest['str'] = v
        ## If we're running the current version, there's nothing newer
        if newest['str'] == current_version:
            return ""
        return newest['str']

    elif operation == 'older':
        older_versions = []
        for v in version_list:
            try:
                v_obj = version.parse(v)
            except:
                console.print(f'[red]Error parsing local source folder version {v} while checking for older versions.')
                continue
            if v_obj < current_version_obj:
                older_versions.append(v)
        return older_versions


## Returns appropriate directory where HACC versions are located
## Based on OS type
def get_hacc_source_dir():
    ## check for windows system
    if 'USERPROFILE' in os.environ:
        return '%USERPROFILE%\AppData\Local\Hacc\\'
    ## check for linux/mac
    if 'HOME' in os.environ:
        return '/usr/local/Hacc/'


## Function to remove previous installations of Hacc
def cleanup_old_versions(old_versions):
    source_dir = get_hacc_source_dir()
    for v in old_versions:
        shutil.rmtree(source_dir+v)
    return


## Function to look for previously installed versions of HACC
## Returns list of previous versions found, empty if none
def check_for_old_versions(current_version):
    hacc_dir = get_hacc_source_dir()
    try:
        versions = os.listdir(hacc_dir)
    except:
        console.print('[red]Could not list contents in HACC source directory, unable to check for old versions.')
        return []

    old_versions = compare_hacc_versions(versions, current_version, operation='older')
    return old_versions


## Function to check github for newer versions of HACC
## Returns newer version if it exists, None otherwise
def check_for_upgrades(current_version):
    try:
        all_version_data = requests.get('https://api.github.com/repos/nbailey20/HACC/tags').json()
        all_versions = [tag['name'] for tag in all_version_data]
    except Exception as e:
        console.print(f'[red]Error checking for potential upgrades: {e}')
        return None

    newer_version = compare_hacc_versions(all_versions, current_version)
    if newer_version:
        return newer_version
    return None
