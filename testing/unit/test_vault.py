import pytest

from classes.vault import Vault


class MockAwsClient:
    def call(self, client_type, api_name, **kwargs):
        if client_type == "kms" and api_name == "describe_key":
            if kwargs["KeyId"] == "alias/testalias":
                return {
                    "KeyMetadata": {
                        "Arn": "arn:aws:kms::00000:key/12345"
                    }
                }

        if client_type == "ssm" and api_name == "get_parameters_by_path":
            if kwargs["Path"] == "/testpath":
                if "NextToken" in kwargs:
                    return {
                        "Parameters": [
                            {
                                "Name": "testpath/p5"
                            },
                            {
                                "Name": "testpath/p6"
                            }
                        ]
                    }
                return {
                    "Parameters": [
                        {
                            "Name": "testpath/p1"
                        },
                        {
                            "Name": "testpath/p2"
                        },
                        {
                            "Name": "testpath/p3"
                        },
                        {
                            "Name": "testpath/p4"
                        }
                    ],
                    "NextToken": "1"
                }

class MockFile:
    def read(self):
        return '{"creds_list": ["a", "b", "c"]}'
    def close(self):
        pass

def open(fname, access_type):
    return MockFile()


## BEGIN TESTS
## init calls get_kms_arn method, test here
def test_vault_init(config):
    v = Vault(config, aws_client=MockAwsClient())
    assert v.aws_client is not None
    assert v.kms_arn == "arn:aws:kms::00000:key/12345"
    assert v.param_path == "testpath"

def test_get_all_services(config):
    v = Vault(config, aws_client=MockAwsClient())
    all_svcs = v.get_all_services()
    all_svcs.sort()
    assert all_svcs == [
        "p1",
        "p2",
        "p3",
        "p4",
        "p5",
        "p6"
    ]

def test_service_exists(config):
    v = Vault(config, aws_client=MockAwsClient())
    assert v.service_exists("p3") is True
    assert v.service_exists("p7") is False

def test_parse_import_file(config, mocker):
    mocker.patch("builtins.open", new=open)
    v = Vault(config, aws_client=MockAwsClient())
    assert v.parse_import_file("test") == ["a", "b", "c"]