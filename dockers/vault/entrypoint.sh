#!/bin/bash

set -e

export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_FORMAT='json'

vault server -dev -dev-listen-address="0.0.0.0:8200" &
sleep 2s

vault login -no-print "${VAULT_DEV_ROOT_TOKEN_ID}"

vault policy write "${APPROLE_ROLE_NAME}"-policy /vault/config/dev-policy.hcl

vault kv put -mount=secret database host="127.0.0.1" port=5432 username=root password=1 db=app
vault kv put -mount=secret app port=8000 debug=true
vault kv put -mount=secret logstash address="http://127.0.0.1:50000"

vault auth enable approle

vault write auth/approle/role/"${APPROLE_ROLE_NAME}" \
    token_policies="${APPROLE_ROLE_NAME}"-policy \
    secret_id_ttl=0 \
    token_num_uses=0 \
    token_ttl=0 \
    token_max_ttl=0 \
    secret_id_num_uses=0

# Overwrite our role id with a known value to simplify our demo
vault write auth/approle/role/"${APPROLE_ROLE_NAME}"/role-id role_id="${APPROLE_ROLE_ID}"

vault write -f auth/approle/role/"${APPROLE_ROLE_NAME}"/secret-id

tail -f /dev/null & trap 'kill %1' TERM ; wait