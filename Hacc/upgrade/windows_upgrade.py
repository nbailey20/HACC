import os
import subprocess
import winreg


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
    return source_dir+'dist\hacc'


## Function to persistently set user PATH variable via Windows registry
## Returns True on success and False otherwise
def update_user_path(value):
    REG_PATH = 'Environment'
    try:
        with winreg.OpenKey(winreg.HKEY_CURRENT_USER, REG_PATH,
                            0, winreg.KEY_WRITE) as registry_key:
            winreg.SetValueEx(registry_key, 'Path', 0, winreg.REG_SZ, value)
        return True
    except WindowsError:
        return False


## Function to get user PATH env variable from Windows registry
## Returns value of PATH on success, None otherwise
def get_user_path():
    REG_PATH = 'HKEY_CURRENT_USER\Environment'
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
    current_path = get_user_path()
    print(f'Current path: {current_path}')
    print(f'Path to add: {executable_path}')
    return True