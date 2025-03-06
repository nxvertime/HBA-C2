#!/bin/bash


echo "[$] Updating packets"
sudo apt update

echo "[$] Installing mysql-server"
sudo apt install -y mysql-server

echo "[$] Starting mysql"
sudo systemctl enable mysql

echo "[$] Installing golang"
sudo apt install -y golang-go

echo "[$] Installing OpenSSL"
sudo apt install -y openssl



MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASS="${MYSQL_PASS:-}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_HOST="${MYSQL_HOST:-localhost}"
MYSQL_FILE="${MYSQL_FILE:-hba_db.sql}"

echo "[$] Importing hba_db from $MYSQL_FILE"

sudo mysql -u "$MYSQL_USER" -p"$MYSQL_PASS" < "$MYSQL_FILE"



echo "[+] Done"

echo "[$] Installing project dependencies"
go mod tidy





echo "[$] Generating SSL Certificates"
chmod +x gen_certs.sh
./gen_certs.sh
echo "[+] Done"
