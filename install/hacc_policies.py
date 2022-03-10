VAULT_IAM_PERMS = """
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ssm:DescribeParameters",
                "ssm:GetParameter",
                "ssm:GetParametersByPath",
                "ssm:DeleteParameter*",
                "ssm:PutParameter"
            ],
            "Resource": [
                "%s",
                "%s"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "kms:Encrypt",
                "kms:Decrypt",
                "kms:DescribeKey"
            ],
            "Resource": "%s"
        }
    ]
}
"""

VAULT_SCP = """
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Deny",
            "Action": [
                "ssm:DescribeParameters",
                "ssm:GetParameter*",
                "ssm:GetParametersByPath",
                "ssm:DeleteParameter*",
                "ssm:PutParameter"
            ],
            "Resource": [
                "%s",
                "%s"
            ],
            "Condition": {
                "StringNotLike": {
                    "aws:PrincipalARN": "%s"
                }
            }
        },
        {
            "Effect": "Deny",
            "Action": [
                "kms:Encrypt",
                "kms:Decrypt",
                "kms:ScheduleKeyDeletion",
                "kms:DisableKey",
                "kms:DeleteAlias",
                "kms:GenerateDataKey*",
                "kms:CreateGrant"
            ],
            "Resource": "%s",
            "Condition": {
                "StringNotLike": {
                    "aws:PrincipalARN": "%s"
                }
            }
        }
    ]
}
"""