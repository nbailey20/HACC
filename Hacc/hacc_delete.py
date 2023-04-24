from classes.hacc_service import HaccService
from classes.vault import Vault


def delete_all_creds(display, config):
    vault = Vault(display, config)

    for svc_name in vault.get_all_services():
        svc_obj = HaccService(display, svc_name, vault=vault)

        for user in svc_obj.get_users():
            svc_obj.remove_credential(user)

        svc_obj.push_to_vault()
        display.update(
            display_type='text_append', 
            data={'text': f'All credentials for service {svc_name} have been deleted'}
        )
    return



def delete(display, args, config):

    ## Wipe all Vault creds before eradication
    if args.wipe:
        display.update(
            display_type='text_new', 
            data={'text': 'Wiping all credentials from Vault...'}
        )
        delete_all_creds(display, config)
        return


    ## Delete single credential via input arguments
    display.update(
        display_type='text_new', 
        data={'text': 'Deleting credential...'}
    )

    service_name = args.service
    user = args.username

    svc_obj = HaccService(display, service_name, config=config)
    if not svc_obj.remove_credential(user):
        display.update(
            display_type='text_append', 
            data={'text': f'Username {user} does not exist for service {service_name}, exiting.'}
        )
        return

    svc_obj.push_to_vault()
    display.update(
        display_type='text_append', 
        data={'text': f'Successfully deleted username {user} from service {service_name}.'}
    )

    if not bool(svc_obj.credentials):
        display.update(
            display_type='text_append', 
            data={'text': f'No more credentials for service {service_name}.'}
        )
    return