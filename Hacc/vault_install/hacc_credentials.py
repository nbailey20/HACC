import subprocess
import os

from logger.hacc_logger import logger


## Returns True if existing AWS credentials/config file contains a profile for Vault
## Return False otherwise
def hacc_profile_exists(profile_name):
    vault_user_region = subprocess.run(['aws', 'configure', 'get', 
        'region', '--profile', profile_name], 
        stdout=subprocess.DEVNULL
    )

    ## Return code 255 if profile not found
    if vault_user_region.returncode != 0:
        logger.debug(f'Subprocess returned error code {vault_user_region.returncode}, 255 means profile not found')
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
        logger.debug(f'Subprocess returned error code {aws_access_key_obj.returncode}, 255 means credential not found')
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
        logger.debug(f'Subprocess returned error code {aws_secret_key_obj.returncode}, 255 means credential not found')
        return False
    return aws_secret_key


## Sets AWS access key ID for Vault user in credentials file
## Returns True on success and False if no access key found
def set_hacc_access_key(access_key_id, profile_name):
    aws_access_key_obj = subprocess.run(['aws', 'configure', 'set', 
                                    'aws_access_key_id', access_key_id,
                                    '--profile', profile_name
                                ])
    ## Return code 255 if cr not found
    if aws_access_key_obj.returncode != 0:
        logger.debug(f'Subprocess returned error code {aws_access_key_obj.returncode}, 255 means profile not found')
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
        logger.debug(f'Subprocess returned error code {aws_secret_key_obj.returncode}, 255 means profile not found')
        return False
    return True


## Adds a new profile to AWS credentials/config file
## Returns True if profile successfully saved, False otherwise
def create_hacc_profile(access_key_id, secret_access_key, profile_name):
    access_key_set = set_hacc_access_key(access_key_id, profile_name)
    secret_key_set = set_hacc_secret_key(secret_access_key, profile_name)
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
        logger.debug('Detected Windows environment')
        creds_file_location = os.environ['USERPROFILE'] + '\.aws\credentials'
        config_file_location = os.environ['USERPROFILE'] + '\.aws\config'
    ## check for linux/mac
    elif 'HOME' in os.environ:
        logger.debug('Detected Linux environment')
        creds_file_location = os.environ['HOME'] + '/.aws/credentials'
        config_file_location = os.environ['HOME'] + '/.aws/config'
    else:
        logger.debug(f'Could not identify OS type from: {os.environ}')
        return False

    logger.debug(f'Default AWS client credential file location: {creds_file_location}')
    logger.debug(f'Default AWS client config file location: {config_file_location}')

    ## Check if credentials or config files located in non-default location
    if 'AWS_SHARED_CREDENTIALS_FILE' in os.environ:
        creds_file_location = os.environ['AWS_SHARED_CREDENTIALS_FILE']
        logger.debug(f'Updated credential location to {creds_file_location} from env var AWS_SHARED_CREDENTIALS_FILE')

    if 'AWS_CONFIG_FILE' in os.environ:
        config_file_location = os.environ['AWS_CONFIG_FILE']
        logger.debug(f'Updated config location to {config_file_location} from env var AWS_CONFIG_FILE')

    try:
        with open(creds_file_location, 'r') as f:
            creds_file_contents = f.readlines()
            logger.debug('Read credential file contents')
        with open(config_file_location, 'r') as f:
            config_file_contents = f.readlines()
            logger.debug('Read config file contents')
    except:
        logger.debug('Failed to read AWS client files')
        return False


    ## Find hacc profile lines in credentials file and remove
    start_profile = f'[{profile_name}]'
    deleting = False
    idx = 0
    while idx < len(creds_file_contents):
        if creds_file_contents[idx].rstrip() == start_profile:
            logger.debug(f'Found start of profile {profile_name} in AWS client credentials file')
            deleting = True
            creds_file_contents.pop(idx)

        ## reached start of next profile, done deleting
        logger.debug(f'Finished deleting profile {profile_name} from AWS client credentials file')
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
            logger.debug(f'Found start of profile {profile_name} in AWS client config file')
            deleting = True
            config_file_contents.pop(idx)

        ## reached start of next profile, done deleting
        logger.debug(f'Finished deleting profile {profile_name} from AWS client config file')
        if config_file_contents[idx].startswith('[') and deleting:
            break

        if deleting:
            config_file_contents.pop(idx)
        else:
            idx += 1


    ## Write updated content back to files
    try:
        with open(creds_file_location, 'w') as f:
            f.write(''.join(creds_file_contents))
            logger.debug(f'Successfully wrote updated AWS client credentials file to {creds_file_location}')
        with open(config_file_location, 'w') as f:
            f.write(''.join(config_file_contents))
            logger.debug(f'Successfully wrote updated AWS client config file to {config_file_location}')

    except:
        logger.debug('Failed to write updated AWS client files')
        return False

    return True
