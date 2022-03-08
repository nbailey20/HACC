import subprocess
import hacc_vars


## Returns True if existing AWS credentials/config file contains a profile for Vault
## Return False otherwise
def hacc_profile_exists():
    vault_user_region = subprocess.run(['aws', 'configure', 'get', 
        'region', '--profile', hacc_vars.aws_hacc_uname], 
        stdout=subprocess.DEVNULL
    )

    ## Return code 255 if profile not found
    if vault_user_region.returncode != 0:
        return False
    return True


## Adds a new profile to AWS credentials/config file
## Returns True if profile successfully saved, False otherwise
def create_hacc_profile(access_key_id, secret_access_key):

    try:
        aws_config_region = subprocess.run(['aws', 'configure', 'set', 
            'region', hacc_vars.aws_hacc_region, 
            '--profile', hacc_vars.aws_hacc_uname]
        )
        
        aws_config_access = subprocess.run(['aws', 'configure', 'set', 
            'aws_access_key_id', access_key_id, 
            '--profile', hacc_vars.aws_hacc_uname]
        )

        aws_config_secret = subprocess.run(['aws', 'configure', 'set', 
            'aws_secret_access_key', secret_access_key, 
            '--profile', hacc_vars.aws_hacc_uname]
        )

        return True
    
    except:
        return False



## Removes an existing profile from AWS credentials/config file
## Returns True if profile successfully removed, False otherwise
def delete_hacc_profile():
    return True
