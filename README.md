# Homemade Authentication Credential Client - HACC

HACC is an open-source credential manager command-line tool that uses your personal AWS account to securely store your secrets, so you don't have to wonder who's secretly selling your data behind the scenes!

* Built with Golang
* Credentials stored in AWS Parameter Store
* Encrypted with AWS KMS (AWS-managed or CMEK) Key
* Can store up to 40MB of encrypted data

Current Version: v2.0.0

## Current Features

* Add, search, rotate, delete, backup, and generate credentials with your own Hacc Vault
* Organize credentials by services e.g. gmail, spectrum, etc
* Interactively view and retrieving existing credentials in the Vault
    * Credential is auto-copied to clipboard with visibility toggle for security
* CLI input is autocompleted where possible to minimize keystrokes required for usage
* Generate hard-to-guess password for new credentials offline via built-in wordlist
    * Includes min/max password length and special / numeric character requirement
* Backup a credential, all credentials in a service, or the entire Vault to a file
* Add credentials to the Vault in bulk by supplying "--file" subarg to the add action
    * Example: importing a backup file to a different Vault

## TODO
* Optionally check for client version updates in github, --upgrade keyword if not on latest version
* Better sample usage / -h output
* Backup option --no-creds for usernames/services only
* Support for different password types, e.g. base64, random, xkcd
* Configuration bool to indicate whether previous HACC client versions should be cleaned up from system
* Add MIT license


## hacc -h
```
hacc is a credential manager backed by AWS SSM

Usage:
  hacc [flags]
  hacc [command]

Available Commands:
  add         Add a credential for a service
  backup      backup credential for all or one service
  delete      Delete a credential for a service
  hacc        View credentials for a service
  help        Help about any command
  rotate      rotate a credential for a service

Flags:
  -h, --help   help for hacc

Use "hacc [command] --help" for more information about a command.

Sample Usage:
  hacc a --username example@gmail.com -p 1234 testService
  hacc add testService -u test@yahoo.com --generate
  hacc testService
  hacc rotate test -u test@yahoo.com -g
  hacc backup -f test_backup.txt
  hacc d testService -u example
```
