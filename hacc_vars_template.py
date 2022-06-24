## Rename this file to hacc_vars.py once all values provided

# Region where Vault should live
aws_hacc_region = 'us-east-1'

# IAM user name in member account for all Vault operations
aws_hacc_uname = 'hacc-user'

# IAM policy name to apply to above user
aws_hacc_iam_policy = 'hacc-policy'

# KMS key alias name for Vault
aws_hacc_kms_alias = 'hacc-key'

# SSM Parameter Store path where credentials will be kept (the vault itself)
aws_hacc_param_path = 'hacc-vault'



# Boolean to indicate whether optional SCP should be created to further lock down Vault account
create_scp = False

# Organizational SCP name that locks down Vault so only client can access
## Mandatory if create_scp == True, UNCOMMENT LINE BELOW
# aws_hacc_scp = 'hacc-scp'

# Cross account role for member (vault) account accessible via mgmt account creds
## Mandatory if create_scp == True, UNCOMMENT LINE BELOW
# aws_member_role = ''