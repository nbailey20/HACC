from classes.hacc_service import HaccService
from classes.vault import Vault


def delete_all_creds():
    vault = Vault()

    for svc_name in vault.get_all_services():
        svc_obj = HaccService(svc_name, vault=vault)

        for user in svc_obj.get_users():
            svc_obj.remove_credential(user)

        svc_obj.push_to_vault()
        print(f'  All credentials for service {svc_name} have been deleted')

    return



def delete(args):

    ## Wipe all Vault creds before eradication
    if args.wipe:
        print('Wiping all credentials from Vault...')
        delete_all_creds()
        return


    ## Delete single credential via input arguments
    print('Deleting credential...')
   
    service_name = args.service
    user = args.username

    svc_obj = HaccService(service_name)
    if not svc_obj.remove_credential(user):
        print(f'  Username {user} does not exist for service {service_name}, exiting.')
        return

    svc_obj.push_to_vault()
    print(f'Successfully deleted username {user} from service {service_name}.')

    if not bool(svc_obj.credentials):
        print(f'  No more credentials for service {service_name}.')
    return