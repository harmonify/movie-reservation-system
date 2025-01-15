#!/bin/bash

root="$(realpath "$(dirname "$(dirname "$0")")")"
echo $root
cd $root/cli
make build
chmod +x $root/cli/.dist/mrs-cli
mv $root/cli/.dist/mrs-cli $root/bin/mrs-cli
