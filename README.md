# Homemade Authentication Credential Client - HACC

HACC is an open-source credential manager command-line tool that uses your personal AWS account to securely store your secrets, so you don't have to wonder what's happening behind the scenes.

* Credentials stored in AWS Parameter Store
* Encrypted with AWS KMS Customer-Managed Key so only you can decrypt
* Up to 40MB of encrypted data costs $1/month (only cost is the key)

Current Version: v0.1

## Current Features

* Ability to add new credential to vault with username/password via -u, -p subargs to specified service
* Ability to delete service and associated credential from vault
* Ability to return credential for service by 'haccing' it
* Uses current AWS CLI credentials to create new least-privilege vault user and KMS CMK, saved credentials as new HACC profile
* Eradicate action to remove IAM, KMS resources - does not yet remove any parameters


## hacc -h
```
HACC v0.1
usage: hacc [-h] [-i | -e | -a | -d] [-u USERNAME] [-p PASSWORD] [-v] [service]

Homemade Authentication Credential Client - HACC

positional arguments:
  service               Service name, a folder that can hold multiple credentials

optional arguments:
  -h, --help            show this help message and exit
  -u USERNAME, --username USERNAME
                        Username to perform action on
  -p PASSWORD, --password PASSWORD
                        Password for new credentials, used with add action
  -v, --verbose         Display verbose output

Actions:
  -i, --install         Create new authentication credential vault
  -e, --eradicate       Delete entire vault - cannot be undone
  -a, --add             Add new set of credentials to vault
  -d, --delete          Delete credentials in vault

Sample Usage:
  hacc -i
  hacc -a -u example@gmail.com -p 1234 testService
  hacc testService
  hacc -d testService
  hacc -e -v
```

## Future Needs

* Error handling
* Interactive prompts if subargs not provided, list existing services
* Store multiple credentials per service
* Ability to generate hard-to-guess password for new credentials
* Ability to update/rotate credential passwords
* Fully wipe any sign of hacc AWS profile from credentials/config file once vault eradicated
* Lock-down Vault with SCP so only HACC AWS user can read/write