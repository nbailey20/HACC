# Homemade Authentication Credential Client - HACC

HACC is an open-source credential manager command-line tool that uses your personal AWS account to securely store your secrets, so you don't have to wonder what's happening behind the scenes.

* Built with Python3
* Credentials stored in AWS Parameter Store
* Encrypted with AWS KMS Customer-Managed Key so only you can decrypt
* Up to 40MB of encrypted data costs $1/month (only cost is the key)

Current Version: v0.4

## Current Features

* Add new credential to vault with username/password for a service
* Store multiple credentials per service
* Provided arguments with flags or interactively
* Generate hard-to-guess password for new credentials
* Delete service and associated credentials from vault
* Return specific or all credentials for service by 'haccing' it
* Use existing AWS CLI credentials to create new least-privilege vault user and KMS CMK, saved credentials as new HACC profile
* Eradicate action to remove IAM, KMS resources - does not yet remove any parameters but schedules master key deletion - credentials manually recoverable for up to 7 days
* Checks for existing services and usernames before adding/deleting credentials
* Completely locks down operations on credential parameters to hacc-user so nobody else can read/delete by unintentionally - via SCP applied to member account
* Backup entire vault to file 


## hacc -h
```
HACC v0.4
usage: hacc [-h] [-i | -e | -a | -d] [-u USERNAME] [-p PASSWORD] [-g] [-b BACKUP] [-v] [service]

Homemade Authentication Credential Client - HACC

positional arguments:
  service               Service name, a folder that can hold multiple credentials

optional arguments:
  -h, --help            show this help message and exit
  -u USERNAME, --username USERNAME
                        Username to perform action on
  -p PASSWORD, --password PASSWORD
                        Password for new credentials, used with add action
  -g, --generate        Generate random password for operation
  -b BACKUP, --backup BACKUP
                        Backup entire Vault and write to file name: hacc -b out_file
  -v, --verbose         Display verbose output

Actions:
  -i, --install         Create new authentication credential vault
  -e, --eradicate       Delete entire vault - cannot be undone
  -a, --add             Add new set of credentials to vault
  -d, --delete          Delete credentials in vault

Sample Usage:
  hacc -iv
  hacc -a -u example@gmail.com -p 1234 testService
  hacc testService
  hacc -d testService
  hacc -e -v
```

## Future Needs
* Clean up code with classes, logging instead of prints
* Add 'list' keyword with pagination, print default HACC ascii if nothing provided - TBD
* Ability to rotate credential passwords
* Fully wipe any sign of hacc AWS profile from credentials/config file once vault eradicated
* Support for services with more than 4KB of credentials via multiple parameters per service
* Error handling for install/eradicate and idempotency
* 'configure' (in addition to install) keyword to grant additional devices access to Vault
* Make SCP optional for single-account deployments not part of an org
* Ability to provide backup file for new vault install or add to existing vault, "--file" subarg