package main

import (
	"po/cmd"
	_ "po/docs" // Swag CLI generates docs, you have to import it.
)

//	@title			Application API
//	@version		1.0
//	@description	Application description.
//	@termsOfService	https://example.com/terms

//	@contact.name	API Support
//	@contact.url	https://example.com/support
//	@contact.email	a.h.pooladvand@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8000
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	//docs.SwaggerInfo.Title = "Application API"
	//docs.SwaggerInfo.Description = "Application description."
	//docs.SwaggerInfo.Version = "1.0"

	cmd.Execute()
}
