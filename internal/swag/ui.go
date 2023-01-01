package swag

var swaggerUiDefaultParameters = map[string]string{
	"dom_id":               `"#swagger-ui"`,
	"deepLinking":          "true",
	"showExtensions":       "true",
	"showCommonExtensions": "true",
	//"layout":               "BaseLayout",
}

func makeSwaggerUiHtml(title, openapiUrl, jsUrl, cssUrl, faviconUrl string) string {
	indexPage := `
	<!DOCTYPE html>
	<html>
	<head>
		<link type="text/css" rel="stylesheet" href="` + cssUrl + `">
		<link rel="shortcut icon" href="` + faviconUrl + `">
		<title>` + title + `</title>
	</head>
	<body>
		<div id="swagger-ui">
		</div>
		<script src="` + jsUrl + `"></script>
		<script>
		const ui = SwaggerUIBundle({
		url: './` + openapiUrl + `',
	`
	for k, v := range swaggerUiDefaultParameters {
		indexPage = indexPage + `"` + k + `": ` + v + ",\n"
	}

	indexPage = indexPage + "oauth2RedirectUrl: window.location.origin + '/docs/oauth2-redirect',\n"

	indexPage += ` presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIBundle.SwaggerUIStandalonePreset
        ],
    })
	</script>
    </body>
    </html>
	`
	return indexPage
}
