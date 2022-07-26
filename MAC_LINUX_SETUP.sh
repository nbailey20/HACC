#!/bin/sh

echo "Starting setup for HACC - Homemade Authentication Credential Client"

if [ ! -f "hacc.conf" ]; then
    if [ ! -f "hacc.conf.template" ]; then
            echo "Neither HACC configuration file or template file found, exiting." && exit 1
    fi

    read -p "Are you configuring the client with an existing vault? (y/n): " existing && [[ $existing == [yY] || $existing == [nN] ]] || echo "Unknown input" && exit 1
    if [[ $existing == [nN] ]]; then
        echo "Please ensure all required values are present in hacc.conf.template and then rename to hacc.conf and rerun this script."
        exit 2
    else
        echo "Thanks for confirming :)"
        mkdir ~/hacc
        cp hacc.conf.template ~/hacc/hacc.conf
        echo "Created empty configuration file ~/hacc/hacc.conf that can be updated with hacc --configure"

else
    mv hacc.conf ~/hacc/hacc.conf
    echo "Moved configuration file to ~/hacc/hacc.conf"
fi

sudo cp -r Hacc/ /usr/local/ && sudo ln -s /usr/local/Hacc/hacc /usr/local/bin/hacc
if [[ $? == 0 ]]; then
    echo "Installed client at /usr/local/bin/hacc, please start new terminal and test with `hacc`"
else
    echo "Failed to install client in /usr/local/bin/hacc, aborting."
    exit 3
fi