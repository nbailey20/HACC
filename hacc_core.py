class HaccService:
    def __init__(self, service=None, user=None, password=None):
        self.servicename = service
        self.username = user
        self.password = password

    def print_service(self, user=False, password=False):
        print()
        print('#############################')
        print(f'Service {self.servicename}')
        if user:
            print(f'Username: {self.username}')
        if password:
            print(f'Password: {self.password}')
        print('#############################')
        print()