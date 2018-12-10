
# Copyright 2018. bolaxy.org authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/sh

# create self-signed server certificate:
# Generate RPC Server Key
echo "Create RPC Server Key..."
cd ../

PWD=$(cd "$(dirname "$0")";pwd)
CERTS_PATH=${PWD}/certs
CONFIGPATH=${PWD}/config
CER_PATH=${CONFIGPATH}/cer

if [[ ! -d "$CERTS_PATH" ]]; then
	mkdir -p $CERTS_PATH
else
	rm -rf $CERTS_PATH"/*"
fi

SUBJECT="/C=CN/ST=Shanghai/L=Earth/O=BOX/OU=DC/CN=box.la/emailAddress"

EMAIL=${2:-develop@2se.com}
DAYS=${3:-3650}

openssl req -new -nodes -x509 -out ${CERTS_PATH}/server.pem -keyout ${CERTS_PATH}/server.key -days ${DAYS} -subj "${SUBJECT}=${EMAIL}"

# Generate RPC Client Key
echo "Generate RPC Client Key..."
openssl req -new -nodes -x509 -out ${CERTS_PATH}/client.pem -keyout ${CERTS_PATH}/client.key -days ${DAYS} -subj "${SUBJECT}=${EMAIL}"

# Generate HTTPS Key
cd $CONFIGPATH
echo "Generate HTTPS Key..."

if [[ ! -d "$CER_PATH" ]]; then
	mkdir -p $CER_PATH
else
	rm -rf $CER_PATH"/*"
fi

cd $CER_PATH

openssl genrsa -des3 -out server.key 1024

echo "Create server certificate signing request..."

SUBJECT="/C=US/ST=Mars/L=iTranswarp/O=iTranswarp/OU=iTranswarp/CN=server"

openssl req -new -subj $SUBJECT -key server.key -out server.csr

echo "Remove password..."

mv server.key server.origin.key
openssl rsa -in server.origin.key -out server.key

echo "Sign SSL certificate..."

openssl x509 -req -days 3650 -in server.csr -signkey server.key -out server.crt

openssl x509 -in server.crt -out server.cer -outform der
