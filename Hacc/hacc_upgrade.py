import tempfile
import os
import shutil
import zipfile
import urllib.request

from versions.hacc_versions import check_for_upgrades

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
        from upgrade.windows_upgrade import build_executable, update_user_path

        client_entrypoint = build_executable(install_dir)
        if not client_entrypoint:
            print('ERROR building Windows executable from Python source, aborting.')
            return
        path_updated = update_user_path(client_entrypoint)
        if not path_updated:
            print(f'ERROR updating user PATH with value {client_entrypoint}, manually set value to complete upgrade.')

    if os_type == 'linux':
        from upgrade.mac_linux_upgrade import update_bin_symlink

        symlink_created = update_bin_symlink(install_dir)
        if not symlink_created:
            print(f'ERROR updating symlink /usr/local/bin/hacc to point to {install_dir}, manual assistance required to complete upgrade.')
    return


def upgrade(_, config):
    current_version = config['version']
    new_version = check_for_upgrades(current_version)
    if not new_version:
        print(f'Current HACC version {current_version} is up to date.')
        return

    ## Get OS type
    os_type = get_os_type()
    if not os_type:
        print('ERROR determining OS type, aborting.')
        return

    ## Download new version of client
    url = SOURCE_DOWNLOAD + new_version
    tempdir = get_download_dir(os_type)
    download_dest = tempdir + 'source.gz'
    try:
        urllib.request.urlretrieve(url, download_dest)
    except Exception as e:
        print(f'Failed to download new HACC version, aborting: {e}')
        return

    ## Unzip download file, clean up temporary download dir
    install_dir = get_install_dir(os_type, new_version)
    try:
        with zipfile.ZipFile(download_dest, 'r') as zf:
            zf.extractall(install_dir)
        shutil.rmtree(tempdir)

        ## zipfile contains Hacc folder, only want its contents in install_dir
        hacc_files = os.listdir(os.path.join(install_dir, 'Hacc'))
        for file_name in hacc_files:
            shutil.move(os.path.join(install_dir, 'Hacc', file_name), install_dir)
        
    except Exception as e:
        print(f'ERROR extracting new version sourcecode to installation directory {install_dir}, aborting: {e}')
        return

    ## Complete upgrade with OS-dependent commands
    upgrade_success = complete_upgrade(os_type, install_dir)
    if not upgrade_success:
        print('Failed to install upgraded version of client, aborting.')
    else:
        print(f'Successfully upgraded HACC software to version {new_version}')
    return