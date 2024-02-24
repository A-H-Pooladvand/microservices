FROM hashicorp/vault

RUN apk add --no-cache bash

COPY app-policy.hcl                   /vault/config/dev-policy.hcl
COPY entrypoint.sh /vault/entrypoint.sh

RUN chmod +x /vault/entrypoint.sh

ENTRYPOINT ["/vault/entrypoint.sh"]