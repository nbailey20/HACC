import pytest

from classes.aws_client import AwsClient


class Session:
    def client(self, *args, **kwargs):
        return 1

def call_res(_, client_type, api_name, **kwargs):
    if client_type == "sts" and api_name == "assume_role":
        return {
            "Credentials": {"AccessKeyId": 1, "SecretAccessKey": 2, "SessionToken": 3}
        }
    elif client_type == "kms" and api_name == "create_key":
        if kwargs["KeyUsage"] == 'ENCRYPT_DECRYPT' and kwargs['KeySpec'] == 'SYMMETRIC_DEFAULT':
            return "success"
        return "fail"
    return "fail"

## Mocking boto3 and AwsClient.call functions
@pytest.fixture(params=["data", "mgmt"])
def client(request, config, mocker):
    mocker.patch("classes.aws_client.AwsClient.call", new=call_res)
    mocker.patch("boto3.client", return_value=1)
    mocker.patch("boto3.session.Session", return_value=Session)
    return AwsClient(config, client_type=request.param)


## BEGIN TESTS
def test_client_empty_init():
    with pytest.raises(TypeError):
        AwsClient()

def test_client_init(client):
    assert client is not None and client.kms == 1

def test_client_call_args(client):
    kms_res = client.call(
        'kms', 'create_key',
        KeyUsage = 'ENCRYPT_DECRYPT',
        KeySpec = 'SYMMETRIC_DEFAULT'
    )
    assert kms_res == "success"