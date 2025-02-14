# Create OpenSSL config file (automatic, no prompts)
cat > openssl.cnf <<EOF
[req]
default_bits       = 2048
default_md         = sha256
distinguished_name = req_distinguished_name
req_extensions     = req_ext
prompt             = no

[req_distinguished_name]
C  = US
ST = California
L  = San Francisco
O  = MyCompany
CN = localhost

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
EOF

# Generate key and certificate
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -config openssl.cnf
openssl x509 -req -in server.csr -signkey server.key -out server.crt -days 365 -extfile openssl.cnf -extensions req_ext

