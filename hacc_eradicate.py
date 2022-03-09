import hacc_vars
from classes.vault_installation import VaultInstallation
import boto3, subprocess, time

def eradicate(args):
    print('WARNING, this operation schedules master key deletion for all credentials.')
    print('If you continue, you have up to 7 days to manually disable deletion or credentials can never be decrypted')
    proceed = True if input('Are you sure you want to proceed (y/n)? ') == 'y' else False
    if not proceed:
        print('Aborting.')
        return

    ## TODO: check for args.wipe and delete all creds if so
        
    print('Eradicating Vault...')

    eradicate = VaultInstallation()

    if hacc_vars.create_scp:
        eradicate.delete_scp()

    eradicate.delete_cmk()

    eradicate.delete_user()
    # # arn:aws:iam::account:role/name
    # account = hacc_vars.aws_member_role.split(':')[4]                
    # orgs = boto3.client('organizations')
    # debug = args.debug

    # # Find and remove SCP before deleting vault resources
    # scp_list = aws_call(
    #     orgs, 'list_policies_for_target', debug,
    #     TargetId=account,
    #     Filter='SERVICE_CONTROL_POLICY'
    # )['Policies']

    # hacc_scp_id = [x['Id'] for x in scp_list if x['Name'] == hacc_vars.aws_hacc_scp][0]
    # if debug: print('INFO: found SCP for Vault with id {}'.format(hacc_scp_id))

    # aws_call(
    #     orgs, 'detach_policy', debug,
    #     PolicyId=hacc_scp_id,
    #     TargetId=account
    # )

    # aws_call(
    #     orgs, 'delete_policy', debug,
    #     PolicyId=hacc_scp_id
    # )

    # # Give SCP deletion time to take effect ('immediate' per AWS docs is not good enough :)
    # if debug: print('INFO: Waiting 10 seconds for SCP to fully delete')
    # time.sleep(10)

    # # Assume role in member account with mgmt account creds
    # sts = boto3.client('sts')
    # assumed_role_object=sts.assume_role(
    #     RoleArn=hacc_vars.aws_member_role,
    #     RoleSessionName="HaccEradicateSession"
    # )
    # role_creds = assumed_role_object['Credentials']

    # iam = boto3.client('iam',
    #                     aws_access_key_id=role_creds['AccessKeyId'],
    #                     aws_secret_access_key=role_creds['SecretAccessKey'],
    #                     aws_session_token=role_creds['SessionToken']
    #                     )
    # kms = boto3.client('kms', region_name=hacc_vars.aws_hacc_region,
    #                             aws_access_key_id=role_creds['AccessKeyId'],
    #                             aws_secret_access_key=role_creds['SecretAccessKey'],
    #                             aws_session_token=role_creds['SessionToken']
    #                     )

    # hacc_key_id = get_kms_arn(kms)
    # aws_call(
    #     kms, 'schedule_key_deletion', debug, 
    #     KeyId=hacc_key_id,
    #     PendingWindowInDays=7
    # )

    # aws_call(
    #     kms, 'delete_alias', debug, 
    #     AliasName='alias/{key}'.format(key=hacc_vars.aws_hacc_kms_alias)
    # )

    aws_call(
        iam, 'delete_user_policy', debug, 
        UserName=hacc_vars.aws_hacc_uname,
        PolicyName=hacc_vars.aws_hacc_iam_policy
    )

    aws_access_key = subprocess.run(['aws', 'configure', 'get', 
        'aws_access_key_id', '--profile', hacc_vars.aws_hacc_uname],
        capture_output=True, text=True
    )
    aws_access_key = aws_access_key.stdout.strip()
    if args.debug: print('INFO: Retrieved saved AWS user access key')
    aws_call(
        iam, 'delete_access_key', debug, 
        UserName=hacc_vars.aws_hacc_uname,
        AccessKeyId=aws_access_key
    )

    aws_call(
        iam, 'delete_user', debug, 
        UserName=hacc_vars.aws_hacc_uname,
    )

    ## Remove profile from ~/.aws/config and credentials! read/regex/write Python
    print('Vault has been eradicated.')
    return