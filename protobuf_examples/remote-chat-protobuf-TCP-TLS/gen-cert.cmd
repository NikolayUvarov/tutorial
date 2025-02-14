# Generate key and certificate
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -config openssl.cnf
openssl x509 -req -in server.csr -signkey server.key -out server.crt -days 365 -extfile openssl.cnf -extensions req_ext


openssl x509 -in server.crt -text -noout | grep -A 1 "Subject Alternative Name"

