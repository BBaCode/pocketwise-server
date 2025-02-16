package constants

var TestURL = "https://deploy-preview-13--pocketwise.netlify.app"
var ProdURL = "https://pocketwise.netlify.app"

var AllowedOrigins = map[string]bool{
	TestURL: true,
	ProdURL: true,
}
