import hacc_vars
from classes.vault_installation import VaultInstallation

    

# Setup IAM user KMS CMK for Vault in member account
# Setup SCP for member account in mgmt account to lock down
def install(args):
    print('Installing new vault...')

    vault = VaultInstallation()

    vault.create_cmk_with_alias()

    # # Assume role in member account with mgmt account creds
    # sts = boto3.client('sts')
    # assumed_role_object=sts.assume_role(
    #     RoleArn=hacc_vars.aws_member_role,
    #     RoleSessionName="HaccInstallSession"
    # )
    # role_creds = assumed_role_object['Credentials']

    # # arn:aws:iam::account:role/name
    # account = hacc_vars.aws_member_role.split(':')[4]

    # kms = boto3.client('kms', region_name=hacc_vars.aws_hacc_region,
    #                     aws_access_key_id=role_creds['AccessKeyId'],
    #                     aws_secret_access_key=role_creds['SecretAccessKey'],
    #                     aws_session_token=role_creds['SessionToken']
    #                     )
    # iam = boto3.client('iam',
    #                     aws_access_key_id=role_creds['AccessKeyId'],
    #                     aws_secret_access_key=role_creds['SecretAccessKey'],
    #                     aws_session_token=role_creds['SessionToken']
    #                     )
    # # SCP is setup in mgmt account, don't use member role
    # orgs = boto3.client('organizations')
    # debug = args.debug

    # hacc_kms = aws_call(
    #             kms, 'create_key', debug, 
    #             KeyUsage='ENCRYPT_DECRYPT', 
    #             KeySpec='SYMMETRIC_DEFAULT'
    #             )

    # aws_call(
    #     kms, 'create_alias', debug, 
    #     AliasName='alias/{key}'.format(key=hacc_vars.aws_hacc_kms_alias), 
    #     TargetKeyId=hacc_kms['KeyMetadata']['Arn']
    # )

    vault.create_user_with_policy()

    # aws_call(
    #     iam, 'create_user', debug, 
    #     UserName=hacc_vars.aws_hacc_uname
    # )

    # hacc_creds = aws_call(
    #                 iam, 'create_access_key', debug, 
    #                 UserName=hacc_vars.aws_hacc_uname
    #             )

    # create_hacc_profile(
    #     hacc_creds['AccessKey']['AccessKeyId'], 
    #     hacc_creds['AccessKey']['SecretAccessKey'], 
    #     debug
    # )

    # vault_path_arn = 'arn:aws:ssm:{region}:{account}:parameter/{path}'.format(
    #                     region=hacc_vars.aws_hacc_region, 
    #                     account=account,
    #                     path=hacc_vars.aws_hacc_param_path
    #                 )
    # vault_key_arn = hacc_kms['KeyMetadata']['Arn']

    # user_perms = json.loads(VAULT_IAM_PERMS % (vault_path_arn, vault_path_arn+'/*', vault_key_arn))
    # aws_call(
    #     iam, 'put_user_policy', debug, 
    #     UserName=hacc_vars.aws_hacc_uname,
    #     PolicyName=hacc_vars.aws_hacc_iam_policy,
    #     PolicyDocument=json.dumps(user_perms)
    # )

    if not hacc_vars.create_scp:
        print('Single-account Vault setup complete.')
        return

    ## TODO: create SCP 
    # iam_user_arn = 'arn:aws:iam::*:user/{}'.format(hacc_vars.aws_hacc_uname)
    # scp = json.loads(VAULT_SCP % (vault_path_arn, 
    #                             vault_path_arn+'/*',
    #                             iam_user_arn,
    #                             vault_key_arn,
    #                             iam_user_arn)
    #                 )
    # hacc_scp = aws_call(
    #             orgs, 'create_policy', debug,
    #             Content=json.dumps(scp),
    #             Name=hacc_vars.aws_hacc_scp,
    #             Description='SCP for HACC',
    #             Type='SERVICE_CONTROL_POLICY'
    #         )

    # aws_call(
    #     orgs, 'attach_policy', debug,
    #     PolicyId=hacc_scp['Policy']['PolicySummary']['Id'],
    #     TargetId=account
    # )
    

    print('Organizational Vault setup complete with SCP.')
    return