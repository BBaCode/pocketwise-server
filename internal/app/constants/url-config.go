package constants

var LocalHost = "http://localhost:4200"
var TestURL = "https://deploy-preview-13--pocketwise.netlify.app"
var ProdURL = "https://pocketwise.netlify.app"

var AllowedOrigins = map[string]bool{
	LocalHost: true,
	TestURL:   true,
	ProdURL:   true,
}
