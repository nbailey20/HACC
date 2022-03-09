#from hacc_core import HaccService
from classes.hacc_service import HaccService
import sys


def delete(args):
    print('Deleting credential...')
   
    service_name = args.service
    user = args.username

    local_service = HaccService(service_name)
    if not local_service.remove_credential(user):
        print(f'Username {user} does not exist for service {service_name}, exiting.')
        return

    local_service.push_to_vault()
    print(f'Successfully deleted username {user} from service {service_name}.')

    if not bool(local_service.credentials):
        print(f'No more credentials for service {service_name}.')
    return