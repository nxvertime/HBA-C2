#!/bin/bash


CERTS_DIR="./certs"


mkdir -p "$CERTS_DIR"


openssl req -x509 -newkey rsa:4096 -keyout "$CERTS_DIR/server.key" -out "$CERTS_DIR/server.crt" -days 365 -nodes


if [ $? -eq 0 ]; then
    echo "[+] Certificates succesfully generated in $CERTS_DIR/"
else
    echo "[-] Certificates generation failed, please install openssl"
    exit 1
fi

