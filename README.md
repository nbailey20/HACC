# Homemade Authentication Credential Client - HACC

HACC is an open-source credential manager command-line tool that uses your personal AWS account to securely store your secrets, so you don't have to wonder who's secretly selling your data behind the scenes!

* Built with Python3
* Credentials stored in AWS Parameter Store
* Encrypted with AWS KMS Customer-Managed Key so only you can decrypt
* Can store up to 40MB of encrypted data

Current Version: v0.7

## Current Features

* Add new credential to vault with username/password for a service
* Store multiple credentials per service
* Provided arguments with flags or interactively
* Pagination for interactive input, search services/usernames by name prefix or line number
* Generate hard-to-guess password for new credentials offline via built-in wordlist
* Delete service and associated credentials from vault
* Return specific or all credentials for service by 'haccing' it
* Use existing AWS CLI credentials to create new least-privilege vault user and KMS CMK, saved credentials as new HACC profile
* Eradicate action to remove IAM, KMS resources and optionally all Vault credentials
* Fully wipe any sign of hacc AWS profile from local credentials/config files once Vault eradicated
* Error handling for install/eradicate and idempotency, other actions if Vault not setup
* Checks for existing services and usernames before adding/deleting credentials
* AWS Organization mode locks down operations on credential parameters to hacc-user so nobody else can read/delete by unintentionally - via optional SCP applied to member account containing Vault (requires adequate role in org to apply SCPs)
* Backup entire vault to file
* Ability to provide backup file for new vault install or add to existing vault, "--file" subarg
* Ability to rotate existing credentials
* Cleanly exit at any point with ctrl-c
* 'configure' keyword to set/show/export client configuration parameters
    * export config as encrypted file to grant multiple users access to single Vault


## hacc -h
```
HACC v0.7

usage: hacc [-h] [-i | -e | -a | -d | -r | -b | -c] [-u USERNAME] [-p PASSWORD] [-g] [-f FILE] [-w] [--export] [--set SET] [--show SHOW] [-v] [service]

Homemade Authentication Credential Client - HACC

positional arguments:
  service               Service name, a folder that can hold multiple credentials

optional arguments:
  -h, --help            show this help message and exit
  -u USERNAME, --username USERNAME
                        Username to perform action on
  -p PASSWORD, --password PASSWORD
                        Password for new credentials, used with add action
  -g, --generate        Generate random password for add operation
  -f FILE, --file FILE  File name for import / export operations
  -w, --wipe            Wipe all existing credentials during Vault eradication
  --export              Export existing client configuration as encrypted file
  --set SET             Set client configuration parameter
  --show SHOW           Show client configuration parameter
  -v, --verbose         Display verbose output

Actions:
  -i, --install         Create new authentication credential Vault
  -e, --eradicate       Delete entire Vault - cannot be undone
  -a, --add             Add new set of credentials to Vault
  -d, --delete          Delete credentials in Vault
  -r, --rotate          Rotate existing password in Vault
  -b, --backup          Backup entire Vault
  -c, --configure       View or modify client configuration

Sample Usage:
  hacc --configure --set aws_hacc_region=us-east-2
  hacc -c --show all
  hacc --install -v
  hacc -a --username example@gmail.com -p 1234 testService
  hacc --add testService -u test@yahoo.com --generate
  hacc testService
  hacc --rotate test -u test@yahoo.com -g
  hacc --backup -f test_backup.txt
  hacc -d testService -u example
  hacc --eradicate --wipe
```

## Creating executable file from Python source
```pyinstaller hacc```
* Will create dist\hacc folder, add this directory to PATH env variable
* Note: if rebuilding, first delete build and dist folders or errors may occur, and restart terminal after running command

## Future Needs
* Support for services with more than 4KB of credentials via multiple parameters per service
* Confirm user doesn't already exist for add action before asking for password if not provided
* Test that all required config variables for action are valid
* Set max password generation length for dumb sites --max-len
* --no-specials and --no-numbers for password generation
* Backup option --no-creds for usernames/services only
* -y option to automate interactive questions
* Support service/name prefixes properly for rotate action
* Support for different password types, e.g. base64, random, xkcd
* Add confirmation before installing vault, but show how many components are setup first
* Optional (via hacc var) check for latest client version in github, upgrade if not - use brew/windows thing instead?
* Reduce args passed to funcs to only necessary values
* Option for AWS-managed SSM KMS key for completely free Vault
* Add MIT license