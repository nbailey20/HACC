import os
import subprocess
import winreg

from console.hacc_console import console


## Function that builds Windows HACC executable using pyinstaller
## Returns absolute path of executable if successful build, None otherwise
def build_executable(source_dir):
    cwd = os.getcwd()
    os.chdir(source_dir)
    build_res = subprocess.run(['pyinstaller', 'hacc'],
                                stdout=subprocess.DEVNULL).returncode
    os.chdir(cwd)
    if build_res != 0:
        return None
    return source_dir+'\dist\hacc'


## Function to persistently set user PATH variable via Windows registry
## Returns True on success and False otherwise
def set_user_path(value):
    console.print(f'Updating path to be: {value}')
    try:
        path_update_res = subprocess.run(
            ['setx', 'PATH', value],
            stdout=subprocess.DEVNULL
        ).returncode

        if path_update_res == 0:
            return True
        return False
    except:
        return False


## Function to get user PATH env variable from Windows registry
## Returns value of PATH on success and None otherwise
def get_user_path():
    REG_PATH = 'Environment'
    try:
        with winreg.OpenKey(winreg.HKEY_CURRENT_USER, REG_PATH,
                            0, winreg.KEY_READ) as registry_key:
            value, _ = winreg.QueryValueEx(registry_key, 'Path')
        return value
    except WindowsError:
        return None


## Function to update user path with newly built executable location by manipulating Windows registry
## Returns True if successful, False otherwise
def update_user_path(executable_path):
    path_str = get_user_path()
    console.print(f'Current user path: {path_str}')

    ## Make sure not to overwrite whole PATH string, only HACC
    path_list = path_str.split(';')
    hacc_path_found = False
    for idx, path in enumerate(path_list):
        if '\\dist\\hacc' in path:
            path_list[idx] = executable_path
            hacc_path_found = True

    if not hacc_path_found:
        console.print('[red]Error finding existing Hacc location in environment PATH')
        return False

    updated_path_str = ';'.join(path_list)
    if set_user_path(updated_path_str):
        return True
    return False
