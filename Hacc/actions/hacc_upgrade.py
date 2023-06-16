import tempfile
import os
import shutil
import zipfile
import urllib.request

from helpers.console.hacc_console import console
from helpers.versions.hacc_versions import check_for_upgrades

SOURCE_DOWNLOAD = 'https://api.github.com/repos/nbailey20/HACC/zipball/refs/tags/'


def get_os_type():
    if 'USERPROFILE' in os.environ:
        return 'windows'
    if 'HOME' in os.environ:
        return 'linux'
    return None


def get_download_dir(os_type):
    dirpath = tempfile.mkdtemp()
    if os_type == 'windows':
        dirpath += '\\'
    if os_type == 'linux':
        dirpath += '/'
    return dirpath


def get_install_dir(os_type, new_version):
    ## check for windows system
    if os_type == 'windows':
        return os.path.join(os.environ['userprofile'],
                        'AppData',
                        'Local',
                        'Hacc',
                        new_version)
    ## check for linux/mac
    if os_type == 'linux':
        return os.path.join('/usr', 'local', 'Hacc', new_version)


def complete_upgrade(os_type, install_dir):
    if os_type == 'windows':
        from helpers.upgrade.windows_upgrade import build_executable, update_user_path

        client_entrypoint = build_executable(install_dir)
        if not client_entrypoint:
            console.print('[red]Error building Windows executable from Python source, aborting.')
            return False

        path_updated = update_user_path(client_entrypoint)
        if not path_updated:
            console.print(f'[red]Error updating user PATH with value {client_entrypoint}, manually set value to complete upgrade.')

    if os_type == 'linux':
        from helpers.upgrade.mac_linux_upgrade import update_bin_symlink

        symlink_created = update_bin_symlink(install_dir)
        if not symlink_created:
            console.print(f'[red]Error updating symlink /usr/local/bin/hacc to point to {install_dir}, manual assistance required to complete upgrade.')
    return True


def upgrade(_, config):
    console.print('Upgrading HACC...')

    current_version = config['version']
    new_version = check_for_upgrades(current_version)
    if not new_version:
        console.print(f'Current HACC version {current_version} is up to date.')
        return

    ## Get OS type
    os_type = get_os_type()
    if not os_type:
        console.print('[red]Error determining OS type, aborting.')
        return

    ## Download new version of client
    url = SOURCE_DOWNLOAD + new_version
    tempdir = get_download_dir(os_type)
    download_dest = tempdir + 'source.gz'
    try:
        urllib.request.urlretrieve(url, download_dest)
    except Exception as e:
        console.print(f'[red]Error downloading new HACC version, aborting: {e}')
        return

    ## Unzip download file, clean up temporary download dir
    install_dir = get_install_dir(os_type, new_version)
    try:
        with zipfile.ZipFile(download_dest, 'r') as zf:
            zf.extractall(install_dir)
        shutil.rmtree(tempdir)
    except Exception as e:
        console.print(f'[red]Error extracting new version sourcecode to installation directory {install_dir}, aborting: {e}')
        return

    try:
        ## zipfile contains <hacc_setup_dir>/Hacc/, need to move files from both dirs to right place
        hacc_setup_dir = os.path.join(install_dir, os.listdir(install_dir)[0]) ## only thing in install directory
        console.print(f'Hacc setup dir {hacc_setup_dir}')

        for file_name in os.listdir(hacc_setup_dir):
            if file_name != 'Hacc':
                console.print(f'Moving setup file name {file_name}')
                shutil.move(os.path.join(hacc_setup_dir, file_name), install_dir)

        ## don't want additional Hacc directory in path, move its files up a level
        hacc_source_dir = os.path.join(hacc_setup_dir, 'Hacc')
        console.print(f'Moving setup file name {file_name}')

        for file_name in os.listdir(hacc_source_dir):
            console.print(f'Source file name {file_name}')
            shutil.move(os.path.join(hacc_source_dir, file_name), install_dir)
        shutil.rmtree(hacc_setup_dir)

    except Exception as e:
        console.print(f'[red]Error organizing files into expected path structure, aborting: {e}')
        return

    ## Complete upgrade with OS-dependent commands
    upgrade_success = complete_upgrade(os_type, install_dir)
    if not upgrade_success:
        console.print('[red]Failed to install upgraded version of client.')
    else:
        console.print(f'Successfully upgraded HACC software to version {new_version}')
        console.print('You may need to restart your terminal for the upgrade to take effect')
    return