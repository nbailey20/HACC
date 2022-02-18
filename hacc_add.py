from hacc_core import HaccService
import sys


def add(args):
    print('Adding new credential...')

    service_name = args.service
    user = args.username
    passwd = args.password

    local_service = HaccService(service_name)

    ## empty dicts evaluate to False
    if not bool(local_service.credentials):
        print(f'Creating new service {service_name} in Vault')

    if not local_service.add_credential(user, passwd):
        print(f'Username {user} already exists for service {service_name}, exiting.')
        sys.exit(1)

    local_service.push_to_vault()
    print('Successfully added credential to Vault.')
    return