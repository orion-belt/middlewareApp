#!/usr/bin/env bash

# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# bootstrap.bash simulates a registered dev gateway going through the bootstrap
# process, outputting generated session secrets as /tmp/magma_protos/gateway.{key,crt}.
#
# Defaults to targeting dev environment bootstrapper, but can also target
# production deployments.
#
# NOTE: target gateway must
#   - already be provisioned
#   - have been provisioned with the (debug) device key type lte_gateway.device.key.key_type=ECHO

set -e
hwid=${1}
# hwid=$(uuidgen)
# hwid="c29f9ded-0d34-4e64-9ee5-c66d202081d6"
echo "Hardware ID :- "$hwid
CERT_DIR=$PWD/magmanbi/.certs
echo "CERT_DIR :- "$CERT_DIR

[[ ${hwid} == "" ]] && usage
bootstrapper_location=${2:-localhost:7444}

# ${MAGMA_ROOT:-~/magma}/orc8r/tools/scripts/consolidate_protos.bash

# protoc \
#   -I /tmp/magma_protos \
#   -I /tmp/magma_protos/orc8r/protos/prometheus/ \
#   --descriptor_set_out /tmp/magma_protos/bootstrapper.protoset \
#   --include_imports \
#   /tmp/magma_protos/orc8r/protos/bootstrapper.proto

challenge=$(grpcurl \
  -insecure \
  -authority bootstrapper-controller.magma.test \
  -protoset $PWD/magmanbi/scripts/bootstrapper.protoset \
  -d "{\"id\": \"${hwid}\"}" \
  ${bootstrapper_location} \
  magma.orc8r.Bootstrapper/GetChallenge \
  | jq -r .challenge
)

echo "<-- challenge -->"
echo " $challenge"
echo "<-- challenge -->"

openssl genrsa -out $CERT_DIR/gateway.key 2048
openssl req -new -key $CERT_DIR/gateway.key -out $CERT_DIR/gateway.csr.der -outform DER -subj "/C=US/CN=${hwid}"
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
  csr_bytes=$(base64 -i -w0 $CERT_DIR/gateway.csr.der)
else
  csr_bytes=$(base64 -i $CERT_DIR/gateway.csr.der)
fi

req="
{
  'echoResponse': {
    'response': '${challenge}'
  },
  'csr': {
    'validTime': '345600s',
    'csrDer': '${csr_bytes}',
    'id': {
      'gateway': {
        'hardwareId': '${hwid}'
      }
    }
  },
  'hwId': {
    'id': '${hwid}'
  },
  'challenge': '${challenge}'
}"
req=$(echo ${req} | tr "'" '"')

cert_der=$(grpcurl \
  -insecure \
  -authority bootstrapper-controller.magma.test \
  -protoset $PWD/magmanbi/scripts/bootstrapper.protoset \
 -d "${req}" \
  ${bootstrapper_location} \
  magma.orc8r.Bootstrapper/RequestSign \
  | jq -r .certDer
)
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
  echo -n ${cert_der} | base64 -d | openssl x509 -inform der -out $CERT_DIR/gateway.crt
else
  echo -n ${cert_der} | base64 -D | openssl x509 -inform der -out $CERT_DIR/gateway.crt
fi

echo ''
echo 'Success'
echo 'Session secrets output to' $CERT_DIR 'gateway.{key,crt}'

