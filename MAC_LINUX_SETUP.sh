#!/bin/sh

VERSION="v0.9"

echo "Starting setup for HACC ${VERSION} - Homemade Authentication Credential Client"

if [ ! -f "Hacc/hacc.conf" ]; then
    if [ ! -f "Hacc/hacc.conf.template" ]; then
            echo "ERROR: Could not find hacc.conf or hacc.conf.template, exiting." && exit 1
    fi

    read -p "Are you configuring the client with a new or existing Vault? (new/existing): " res && [[ $res == "new" || $res == "existing" ]] || echo "Unknown response" && exit 2
    if [[ $res == "new" ]]; then
        echo "To configure client for a new Vault, please ensure all required values are present in Hacc/hacc.conf.template, then rename template to hacc.conf and rerun this script."
        exit 3
    else
        echo "Welcome back!"
        mkdir ~/.hacc
        cp Hacc/hacc.conf.template ~/.hacc/hacc.conf
        echo "Created empty configuration, connect to your existing Vault with:"
        echo "  hacc --configure --set all --file <exported_config_file> --password <config_encryption_password>"
    fi

else
    echo "Found hacc.conf file..."
    mkdir ~/.hacc
    mv Hacc/hacc.conf ~/.hacc/hacc.conf
    if [[ $? == 0 ]]; then
        echo "Configuration file moved to ~/.hacc/hacc.conf, values can up be updated with hacc --configure"
    else
        echo "Could not move hacc.conf to ~/.hacc/, please ensure you have permission to install on this system."
        exit 4
    fi
fi

sudo mkdir -p /usr/local/Hacc/${VERSION}
sudo chown -R $USER /usr/local/Hacc
sudo cp -r Hacc/* /usr/local/Hacc/${VERSION} && sudo ln -s /usr/local/Hacc/${VERSION}/hacc /usr/local/bin/hacc
sudo chmod +x /usr/local/Hacc/${VERSION}/hacc
if [[ $? == 0 ]]; then
    echo "Installed client at /usr/local/bin/hacc, test with command 'hacc' (terminal restart may be needed)"
else
    echo "Failed to install client in /usr/local/bin/hacc, aborting."
    exit 5
fi
