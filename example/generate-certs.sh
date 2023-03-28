
#!/usr/bin/env bash

CERTS_DIR=certs

if [ -d $CERTS_DIR ]
then
    echo "Delete existing $CERTS_DIR directory before running this script"
    exit 1
fi

mkdir -p $CERTS_DIR

# Create OpenSSL CA configuration file
cat > ca.cnf <<EOF
# OpenSSL CA configuration file
[ ca ]
default_ca         = CA_default
[ CA_default ]
default_days       = 365
database           = index.txt
serial             = serial.txt
default_md         = sha256
copy_extensions    = copy 

# Used to create the CA certificate
[ req ]
prompt             =no
distinguished_name = distinguished_name
x509_extensions    = extensions
[ distinguished_name ]
organizationName   = Redpanda
commonName         = Redpanda CA
[ extensions ]
keyUsage           = critical,digitalSignature,nonRepudiation,keyEncipherment,keyCertSign
basicConstraints   = critical,CA:true,pathlen:1

# Common policy for nodes and users.
[ signing_policy ]
organizationName   = supplied
commonName         = optional