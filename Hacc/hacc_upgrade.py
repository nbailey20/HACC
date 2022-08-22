import os
import urllib.request

from versions.hacc_versions import check_for_upgrades

SOURCE_DOWNLOAD = 'https://api.github.com/repos/nbailey20/HACC/zipball/refs/tags/'


def get_download_location(version):
    ## check for windows system
    if 'USERPROFILE' in os.environ:
        return f'%USERPROFILE%\Downloads\Hacc{version}.gz'
    ## check for linux/mac
    if 'HOME' in os.environ:
        return f'/tmp/Hacc{version}.gz'


def upgrade(_, config):
    new_version = check_for_upgrades(config['version'])
    if not new_version:
        print(f'Current HACC version {config["version"]} is up to date.')
        return

    url = SOURCE_DOWNLOAD + new_version
    try:
        urllib.request.urlretrieve(url, get_download_location())
    except Exception as e:
        print(f'Failed to download new HACC version, aborting: {e}')
    return