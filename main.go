package main

import (
	"os"
	"po/cmd"
)

// Todo:: Move apm integration to dedicated place
// Todo:: Try to implement vault
// Todo:: Add dockers such as vault and postgres and...
// Todo::
func main() {
	_ = os.Setenv("ELASTIC_APM_SERVER_URL", "http://127.0.0.1:8200")
	_ = os.Setenv("ELASTIC_APM_SERVICE_NAME", "app")
	_ = os.Setenv("ELASTIC_APM_SECRET_TOKEN", "")
	//_ = os.Setenv("ELASTIC_APM_SERVER_CA_CERT_FILE", "docker/tls/certs/ca/ca.crt")
	_ = os.Setenv("ELASTIC_APM_LOG_LEVEL", "debug")
	_ = os.Setenv("ELASTIC_APM_LOG_FILE", "stderr")
	_ = os.Setenv("ELASTIC_APM_ENVIRONMENT", "staging")
	_ = os.Setenv("ELASTIC_APM_VERIFY_SERVER_CERT", "false")

	cmd.Execute()
}
