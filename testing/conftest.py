import pytest

## Contains shared fixtures usable by any tests, no import required :)

@pytest.fixture(params=["scp", "no-scp"])
def config(request):
    c = {
        "aws_hacc_region":     "us-east-1",
        "aws_hacc_uname":      "test_user",
        "aws_hacc_kms_alias":  "testalias",
        "aws_hacc_param_path": "testpath",
        "create_scp":          False,
    }
    if request.param == "scp":
        c["create_scp"]      = True
        c["aws_member_role"] = "arn:aws:iam::00000:role/test_role"
    return c