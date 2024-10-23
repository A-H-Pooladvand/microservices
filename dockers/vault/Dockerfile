FROM hashicorp/vault

RUN apk add --no-cache bash

COPY app-policy.hcl                   /vault/config/dev-policy.hcl
COPY entrypoint.sh /vault/entrypoint.sh

RUN chmod +x /vault/entrypoint.sh

RUN dos2unix /vault/entrypoint.sh

ENTRYPOINT ["/bin/bash", "/vault/entrypoint.sh"]