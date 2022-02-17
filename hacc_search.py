from hacc_core import HaccService, service_exists
import logging

logger=logging.getLogger(__name__)



# def get_service_name(svc, ssm_client, debug):
#     svc_list = get_all_services(ssm_client, debug)
#     existing_service = False
#     for s in svc_list:
#         if svc.lower() == s.lower():
#             existing_service = True
#             if svc != s:
#                 print('Found similar existing service: {}'.format(s))
#                 if input('Use this service (y/n)? ') == 'y':
#                     return s
#             else:
#                 return s
#     return False



## Print specific credential for service
def search(args):

    service_name = args.service
    user = args.username

    local_service = HaccService(service_name)
    if not service_exists(service_name, local_service.ssm_client):
        print(f'Service {service_name} does not exist in Vault, exiting.')
        return

    local_service.print_credential(user, password=True)
    

    # if not args.service:
    #     all_svcs = get_all_services(ssm, args.debug)
    #     # If no services in vault, can't look anything up
    #     if not all_svcs:
    #         print('No credentials in vault, add credential before searching.')
    #         return
    #     # Print available services if at least one exists
    #     print()
    #     print('Available services:')
    #     for svc in all_svcs:
    #         print(svc)
    #     print()
    #     args.service = input('Enter service name from above list to retrieve credential: ')
    
    return
