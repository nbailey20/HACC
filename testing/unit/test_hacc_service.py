import pytest

from classes.vault import Vault
from classes.hacc_service import HaccService


class MockAwsClient:
    def call(self, client_type, api_name, **kwargs):
        if client_type == "ssm" and api_name == "get_parameters_by_path":
            return {
                "Parameters": [
                    {
                        "Name": "testpath/testservice"
                    }
                ]
            }

        if client_type == "ssm" and api_name == "get_parameter":
            if kwargs["Name"] == "/testpath/testservice":
                return {
                    "Parameter": {
                        "Value": "testuser1:testpass1,testuser2:testpass2,testuser3:testpass3"
                    }
                }

        if client_type == "ssm" and api_name == "put_parameter":
            if kwargs["Name"] == "/testpath/testservice":
                if kwargs["Value"] == "testuser1:testpass1,testuser2:testpass2,testuser3:testpass3":
                    return True
                return False
            if kwargs["Name"] == "/testpath/abc":
                if kwargs["Value"] == "user1:pass1":
                    return True
                return False

        if client_type == "ssm" and api_name == "delete_parameter":
            if kwargs["Name"] == "/testpath/abc":
                return True
            return False


## BEGIN TESTS
## init method calls pull_from_vault, test here
def test_hacc_service_init(config):
    v = Vault(config, aws_client=MockAwsClient())
    good_svc = HaccService("testservice", vault=v)
    assert good_svc.vault is not None
    assert good_svc.service_name == "testservice"
    assert good_svc.credentials["testuser1"] == "testpass1"
    assert good_svc.credentials["testuser2"] == "testpass2"
    assert good_svc.credentials["testuser3"] == "testpass3"

    bad_svc = HaccService("abc", vault=v)
    assert bad_svc.credentials == {}

def test_get_users(config):
    v = Vault(config, aws_client=MockAwsClient())
    good_svc = HaccService("testservice", vault=v)
    assert good_svc.get_users() == ["testuser1", "testuser2", "testuser3"]

    bad_svc = HaccService("abc", vault=v)
    assert bad_svc.get_users() == []

def test_get_credential(config):
    v = Vault(config, aws_client=MockAwsClient())
    s = HaccService("testservice", vault=v)
    assert s.get_credential("testuser2") == "testpass2"
    assert s.get_credential("baduser") is False

def test_user_exists(config):
    v = Vault(config, aws_client=MockAwsClient())
    s = HaccService("testservice", vault=v)
    assert s.user_exists("testuser3") is True
    assert s.user_exists("baduser") is False

def test_add_credential(config):
    v = Vault(config, aws_client=MockAwsClient())
    existing_svc = HaccService("testservice", vault=v)
    existing_svc.add_credential("testuser4", "testpass4")
    assert existing_svc.credentials == {
        "testuser1": "testpass1",
        "testuser2": "testpass2",
        "testuser3": "testpass3",
        "testuser4": "testpass4"
    }

    new_svc = HaccService("abc", vault=v)
    new_svc.add_credential("user1", "pass1")
    assert new_svc.credentials == {"user1": "pass1"}

def test_remove_credential(config):
    v = Vault(config, aws_client=MockAwsClient())
    s = HaccService("testservice", vault=v)
    s.remove_credential("testuser1")
    assert s.credentials == {
        "testuser2": "testpass2",
        "testuser3": "testpass3"
    }
    assert s.remove_credential("testuser4") is False
    assert s.credentials == {
        "testuser2": "testpass2",
        "testuser3": "testpass3"
    }

def test_push_to_vault(config):
    ## test no change to creds
    v = Vault(config, aws_client=MockAwsClient())
    existing_svc = HaccService("testservice", vault=v)
    assert existing_svc.push_to_vault() is True

    ## test new cred
    new_svc = HaccService("abc", vault=v)
    new_svc.add_credential("user1", "pass1")
    assert new_svc.push_to_vault() is True
    ## API call expects only user1 for testing
    new_svc.add_credential("user2", "pass2")
    assert new_svc.push_to_vault() is False

    ## test delete_parameter logic
    new_svc.remove_credential("user1")
    new_svc.remove_credential("user2")
    assert new_svc.push_to_vault() is True