
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

# Used to sign node certificates.
[ signing_node_req ]
keyUsage           = critical,digitalSignature,keyEncipherment
extendedKeyUsage   = serverAuth,clientAuth

# Used to sign client certificates.
[ signing_client_req ]
keyUsage           = critical,digitalSignature,keyEncipherment
extendedKeyUsage   = clientAuth
EOF

# Generate RSA private key for the CA
openssl genrsa -out $CERTS_DIR/ca.key 2048
chmod 400 $CERTS_DIR/ca.key

# Generate self-signed CA certificate
openssl req -new -x509 -config ca.cnf -key $CERTS_DIR/ca.key -out $CERTS_DIR/ca.crt -days 365 -batch

rm -f index.txt serial.txt
touch index.txt
echo '01' > serial.txt

# Create OpenSSL node and client configuration file
cat > node.cnf <<EOF
# OpenSSL node configuration file
[ req ]
prompt             = no
distinguished_name = distinguished_name
req_extensions     = extensions
[ distinguished_name ]
organizationName   = Redpanda
[ extensions ]
subjectAltName     = critical,@alt_names
[ alt_names ]
DNS  = 127.0.0.1
IP.1 = 127.0.0.1
IP.2 = 172.24.1.10
IP.3 = 172.24.1.20
EOF

# Generate RSA private key for the nodes