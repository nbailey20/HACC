#!/bin/sh

VERSION="v0.9"

echo "Starting setup for HACC ${VERSION} - Homemade Authentication Credential Client"

if [ ! -f "Hacc/hacc.conf" ]; then
    if [ ! -f "Hacc/hacc.conf.template" ]; then
            echo "Neither HACC configuration file or template file found, exiting." && exit 1
    fi

    read -p "Are you configuring the client with an existing vault? (y/n): " existing && [[ $existing == [yY] || $existing == [nN] ]] || exit 2
    if [[ $existing == [nN] ]]; then
        echo "To configure new Vault, please ensure all required values are present in Hacc/hacc.conf.template, then rename template to hacc.conf and rerun this script."
        exit 3
    else
        echo "Thanks for confirming :)"
        mkdir ~/.hacc
        cp Hacc/hacc.conf.template ~/.hacc/hacc.conf
        echo "Created empty configuration file ~/.hacc/hacc.conf that can be updated with hacc --configure"
    fi

else
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
