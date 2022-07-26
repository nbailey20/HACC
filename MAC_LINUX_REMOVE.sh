#!/bin/sh

echo "Removing all signs of HACC installation..."

sudo rm -rf /usr/local/bin/hacc
sudo rm -rf /usr/local/Hacc/
rm -rf ~/.hacc/

echo "HACC removal complete. Bye!"