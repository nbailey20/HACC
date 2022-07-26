import subprocess, os


## Returns True if existing AWS credentials/config file contains a profile for Vault
## Return False otherwise
def hacc_profile_exists(profile_name):
    vault_user_region = subprocess.run(['aws', 'configure', 'get', 
        'region', '--profile', profile_name], 
        stdout=subprocess.DEVNULL
    )

    ## Return code 255 if profile not found
    if vault_user_region.returncode != 0:
        return False
    return True



## Returns AWS access key ID for Vault user from credentials file
## Returns False if no access key found
def get_hacc_access_key(profile_name):
    aws_access_key_obj = subprocess.run(['aws', 'configure', 'get', 
                                    'aws_access_key_id', '--profile', 
                                    profile_name],
                                    capture_output=True, text=True
                                )
    aws_access_key = aws_access_key_obj.stdout.strip()

    ## Return code 255 if credential not found
    if aws_access_key_obj.returncode != 0:
        return False
    return aws_access_key


## Returns AWS secret key ID for Vault user from credentials file
## Returns False if no access key found
def get_hacc_secret_key(profile_name):
    aws_secret_key_obj = subprocess.run(['aws', 'configure', 'get', 
                                    'aws_secret_access_key', '--profile', 
                                    profile_name],
                                    capture_output=True, text=True
                                )
    aws_secret_key = aws_secret_key_obj.stdout.strip()

    ## Return code 255 if credential not found
    if aws_secret_key_obj.returncode != 0:
        return False
    return aws_secret_key


## Sets AWS access key ID for Vault user in credentials file
## Returns True on success and False if no access key found
def set_hacc_access_key(access_key_id, profile_name):
    aws_access_key_obj = subprocess.run(['aws', 'configure', 'set', 
                                    'aws_access_key_id', access_key_id,
                                    '--profile', profile_name
                                ])
    ## Return code 255 if credential not found
    if aws_access_key_obj.returncode != 0:
        return False
    return True


## Sets AWS secret key ID for Vault user in credentials file
## Returns True on success and False if no secret key found
def set_hacc_secret_key(secret_access_key, profile_name):
    aws_secret_key_obj = subprocess.run(['aws', 'configure', 'set', 
                                    'aws_secret_access_key', secret_access_key,
                                    '--profile', profile_name
                                ])
    ## Return code 255 if credential not found
    if aws_secret_key_obj.returncode != 0:
        return False
    return True


## Adds a new profile to AWS credentials/config file
## Returns True if profile successfully saved, False otherwise
def create_hacc_profile(access_key_id, secret_access_key):
    access_key_set = set_hacc_access_key(access_key_id)
    secret_key_set = set_hacc_secret_key(secret_access_key)
    if access_key_set and secret_key_set:
        return True
    return False



## Removes an existing profile from AWS credentials/config file
## Returns True if profile successfully removed, False otherwise
## No way to completely remove AWS profile using aws configure, clean up manually
def delete_hacc_profile(profile_name):
    creds_file_location = ''
    config_file_location = ''

    ## check for windows system
    if 'USERPROFILE' in os.environ:
        creds_file_location = os.environ['USERPROFILE'] + '\.aws\credentials'
        config_file_location = os.environ['USERPROFILE'] + '\.aws\config'
    ## check for linux/mac
    elif 'HOME' in os.environ:
        creds_file_location = os.environ['HOME'] + '/.aws/credentials'
        config_file_location = os.environ['HOME'] + '/.aws/config'
    else:
        return False


    ## Check if credentials or config files located in non-default location
    if 'AWS_SHARED_CREDENTIALS_FILE' in os.environ:
        creds_file_location = os.environ['AWS_SHARED_CREDENTIALS_FILE']

    if 'AWS_CONFIG_FILE' in os.environ:
        config_file_location = os.environ['AWS_CONFIG_FILE']

    try:
        f = open(creds_file_location, 'r')
        creds_file_contents = f.readlines()
        f.close()

        f = open(config_file_location, 'r')
        config_file_contents = f.readlines()
        f.close()
    except:
        return False


    ## Find hacc profile lines in credentials file and remove
    start_profile = f'[{profile_name}]'
    deleting = False
    idx = 0
    while idx < len(creds_file_contents):
        if creds_file_contents[idx].rstrip() == start_profile:
            deleting = True
            creds_file_contents.pop(idx)

        ## reached start of next profile, done deleting
        if creds_file_contents[idx].startswith('[') and deleting:
            break

        if deleting:
            creds_file_contents.pop(idx)
        else:
            idx += 1

    
    ## Find hacc profile lines in config file and remove
    start_profile = f'[profile {profile_name}]'
    deleting = False
    idx = 0
    while idx < len(config_file_contents):
        if config_file_contents[idx].rstrip() == start_profile:
            deleting = True
            config_file_contents.pop(idx)

        ## reached start of next profile, done deleting
        if config_file_contents[idx].startswith('[') and deleting:
            break

        if deleting:
            config_file_contents.pop(idx)
        else:
            idx += 1


    ## Write updated content back to files
    try:
        f = open(creds_file_location, 'w')
        f.write(''.join(creds_file_contents))
        f.close()

        f = open(config_file_location, 'w')
        f.write(''.join(config_file_contents))
        f.close()
    except:
        return False

    return True
