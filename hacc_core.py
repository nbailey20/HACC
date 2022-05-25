import hacc_vars

MGMT_ACTIONS = ['install', 'eradicate']
REQUIRED_MGMT_VARS = [
    'aws_hacc_region',
    'aws_hacc_uname',
    'aws_hacc_iam_policy',
    'aws_hacc_kms_alias',
    'aws_hacc_param_path',
    'create_scp'
]

DATA_ACTIONS = ['add', 'delete', 'search', 'rotate', 'backup']
REQUIRED_DATA_VARS = [
    'aws_hacc_kms_alias',
    'aws_hacc_uname',
    'aws_hacc_region',
    'aws_hacc_param_path'
]

REQUIRED_SCP_VARS = ['aws_hacc_scp', 'aws_member_role']


## Function to check all required config vars are set
##   when client is invoked
def startup(action):
    all_vars_exist = True
    missing_vars = []

    ## TODO: validate variables to ensure sensible values, not just checking if empty

    if action in DATA_ACTIONS:
        for v in REQUIRED_DATA_VARS:
            if not hasattr(hacc_vars, v) or getattr(hacc_vars, v) == '':
                missing_vars.append(v)
                all_vars_exist = False

    elif action in MGMT_ACTIONS:
        for v in REQUIRED_MGMT_VARS:
            if not hasattr(hacc_vars, v) or getattr(hacc_vars, v) == '':
                missing_vars.append(v)
                all_vars_exist = False

        ## conditionally check for SCP vars if create_scp == True
        if not 'create_scp' in missing_vars and hacc_vars.create_scp == True:
            for v in REQUIRED_SCP_VARS:
                if not hasattr(hacc_vars, v) or getattr(hacc_vars, v) == '':
                    missing_vars.append(v)
                    all_vars_exist = False

    for m in missing_vars:
        print(f'Required configuration parameter {m} not set in hacc_vars.py, please set value and try again.')

    return all_vars_exist