#!/usr/bin/env bash

echo "Remove existing *.pem files."

rm *.pem
rm client/*.pem
rm server/*.pem

echo "----- Generate Certificate Authority's private key and self-signed certificate."
openssl req \
    -days 365 \
    -keyout certificate-authority/private_key.pem \
    -newkey rsa:4096 \
    -noenc \
    -out certificate-authority/certificate.pem \
    -subj "/C=US/ST=NV/L=Las Vegas/O=Senzing/OU=Test CA/CN=senzing.com" \
    -x509

openssl x509 \
    -in certificate-authority/certificate.pem \
    -noout \
    -text

# Generate web server's private key and certificate signing request (CSR)

echo "----- Generate server certificate."

openssl req \
    -keyout server/private_key.pem \
    -newkey rsa:4096 \
    -noenc \
    -out server/certificate_request.pem \
    -subj "/C=US/ST=NV/L=Las Vegas/O=Senzing/OU=Test Server/CN=senzing.com"

# Use CA's private key to sign web server's CSR and get back the signed certificate.

openssl x509 \
    -CA certificate-authority/certificate.pem \
    -CAcreateserial \
    -CAkey certificate-authority/private_key.pem \
    -days 360 \
    -extfile server/ext.cnf \
    -in server/certificate_request.pem \
    -out server/certificate.pem \
    -req

openssl x509 \
    -in server/certificate.pem \
    -noout \
    -text

# Generate client's private key and certificate signing request (CSR)

echo "----- Generate client certificate."

openssl req \
    -keyout client/private_key.pem \
    -newkey rsa:4096 \
    -noenc \
    -out client/certificate_request.pem \
    -subj "/C=US/ST=NV/L=Las Vegas/O=Senzing/OU=Test Client/CN=senzing.com"

# Use CA's private key to sign client's CSR and get back the signed certificate

openssl x509 \
    -CA certificate-authority/certificate.pem \
    -CAcreateserial \
    -CAkey certificate-authority/private_key.pem \
    -days 360 \
    -extfile client/ext.cnf \
    -in client/certificate_request.pem \
    -out client/certificate.pem \
    -req

openssl x509 \
    -in client/certificate.pem \
    -noout \
    -text