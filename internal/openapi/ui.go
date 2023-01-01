package openapi

var swaggerUiDefaultParameters = map[string]string{
	"dom_id":               `"#swagger-ui"`,
	"deepLinking":          "true",
	"showExtensions":       "true",
	"showCommonExtensions": "true",
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

func makeRedocUiHtml(title, openapiUrl, jsUrl, faviconUrl string) string {
	indexPage := `
	<!DOCTYPE html>
	<html>
	<head>= 
		<title>` + title + ` </title>
		<!-- needed for adaptive design -->
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
		<link rel="shortcut icon" href="` + faviconUrl + `">
		<!--
			ReDoc doesn't change outer page styles
		-->
		<style>
			body {{
				margin: 0;
				padding: 0;
			}}
		</style>
	</head>`

	indexPage += `
	<body>
	<noscript>
		ReDoc requires Javascript to function. Please enable it to browse the documentation.
	</noscript>
	<redoc spec-url="` + openapiUrl + `"></redoc>
	<script src="` + jsUrl + `"> </script>
	</body>
	</html>`

	return indexPage
}
