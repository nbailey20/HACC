from classes.hacc_service import HaccService
import logging

logger=logging.getLogger()


## Print specific credential for service
def search(args):

    service_name = args.service
    user = args.username

    local_service = HaccService(service_name)
    if not bool(local_service.credentials):
        print(f'Service {service_name} does not exist in Vault, exiting.')
        return

    local_service.print_credential(user, print_password=True)
    return
