from hacc_core import get_all_services, HaccService, ApiClient


def get_service():
    api_obj = ApiClient(ssm=True)
    svc_list = get_all_services(api_obj.ssm)



def get_user(svc_name):
    return


def get_password():
    return

