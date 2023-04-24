import logging

from classes.hacc_service import HaccService


logger=logging.getLogger()


## Print specific credential for service
def search(display, args, config):

    service_name = args.service
    user = args.username

    local_service = HaccService(display, service_name, config=config)
    if not bool(local_service.credentials):
        display.update(
            display_type='text_append', 
            data={'text': f'Service {service_name} does not exist in Vault, exiting.'}
        )
        return

    display_data = {
        'user': user,
        'passwd': local_service.get_credential(user),
        'service': service_name
    }
    display.update(display_type='credential', data=display_data)
    return
