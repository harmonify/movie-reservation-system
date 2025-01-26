#!/bin/bash

# This script is a simple wrapper around the CLI class. It is used to run the CLI
# from the command line. It is not intended to be used in a production environment.

cd ""$(dirname "$(dirname "$0")")"/cli"

go run main.go "$@"
