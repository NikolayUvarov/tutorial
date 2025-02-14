# Создаем закрытый ключ
#openssl genrsa -out server.key 2048

# Создаем сертификат
#openssl req -new -x509 -key server.key -out server.crt -days 365  -subj "/C=US/ST=State/L=City/O=Example/OU=IT Department/CN=localhost"



# Generate a new private key
openssl genrsa -out server.key 2048

# Generate a Certificate Signing Request (CSR)
openssl req -new -key server.key -out server.csr -config openssl_silent.cnf

# Generate the final certificate with SANs
openssl x509 -req -in server.csr -signkey server.key -out server.crt -days 365 -extfile openssl_silent.cnf -extensions req_ext


:: check
openssl x509 -in server.crt -text -noout | grep -A 1 "Subject Alternative Name"
