#!/usr/bin/env bash

echo "Remove existing *.pem files."

rm certificate-authority/*.pem
rm client/*.pem
rm server/*.pem

echo "----- Generate Certificate Authority's self-signed certificate and private key."

openssl req \
    -days 365 \
    -keyout certificate-authority/private_key.pem \
    -newkey rsa:4096 \
    -noenc \
    -out certificate-authority/certificate.pem \
    -subj "/C=US/ST=NV/L=Las Vegas/O=Senzing/OU=Test CA/CN=senzing.com" \
    -x509

# Create encrypted private key.

openssl rsa \
    -aes256 \
    -in certificate-authority/private_key.pem \
    -out certificate-authority/private_key_encrypted.pem \
    -passout pass:Passw0rd

# View certificate.

openssl x509 \
    -in certificate-authority/certificate.pem \
    -noout \
    -text

# Generate server's private key and certificate signing request (CSR).

echo "----- Generate server certificate."

openssl req \
    -keyout server/private_key.pem \
    -newkey rsa:4096 \
    -noenc \
    -out server/certificate_request.pem \
    -subj "/C=US/ST=NV/L=Las Vegas/O=Senzing/OU=Test Server/CN=senzing.com"

# Use CA's private key to sign server's CSR and get back the signed certificate.

openssl x509 \
    -CA certificate-authority/certificate.pem \
    -CAcreateserial \
    -CAkey certificate-authority/private_key.pem \
    -days 360 \
    -extfile server/ext.cnf \
    -in server/certificate_request.pem \
    -out server/certificate.pem \
    -req

# Create encrypted private key.

openssl rsa \
    -aes256 \
    -in server/private_key.pem \
    -out server/private_key_encrypted.pem \
    -passout pass:Passw0rd

# View certificate.

openssl x509 \
    -in server/certificate.pem \
    -noout \
    -text

# Generate client's private key and certificate signing request (CSR).

echo "----- Generate client certificate."

openssl req \
    -keyout client/private_key.pem \
    -newkey rsa:4096 \
    -noenc \
    -out client/certificate_request.pem \
    -subj "/C=US/ST=NV/L=Las Vegas/O=Senzing/OU=Test Client/CN=senzing.com"

# Use CA's private key to sign client's CSR and get back the signed certificate.

openssl x509 \
    -CA certificate-authority/certificate.pem \
    -CAcreateserial \
    -CAkey certificate-authority/private_key.pem \
    -days 360 \
    -extfile client/ext.cnf \
    -in client/certificate_request.pem \
    -out client/certificate.pem \
    -req

# Create encrypted private key.

openssl rsa \
    -aes256 \
    -in client/private_key.pem \
    -out client/private_key_encrypted.pem \
    -passout pass:Passw0rd

# View certificate.

openssl x509 \
    -in client/certificate.pem \
    -noout \
    -text