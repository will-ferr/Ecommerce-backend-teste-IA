package docs

// SwaggerInfo holds exported Swagger Info so clients can import it
var SwaggerInfo = struct {
	Version        string
	Title          string
	Description    string
	TermsOfService string
	Contact        struct {
		Name  string
		URL   string
		Email string
	}
	License struct {
		Name string
		URL  string
	}
	Host     string
	BasePath string
}{
	Version:        "1.0.0",
	Title:          "Smart Choice E-commerce API",
	Description:    "A comprehensive e-commerce backend with authentication, products, orders, and payment processing",
	TermsOfService: "http://swagger.io/terms/",
	Contact: struct {
		Name  string
		URL   string
		Email string
	}{
		Name:  "Smart Choice Team",
		URL:   "https://github.com/will-ferr/Ecommerce-backend-teste-IA",
		Email: "support@smartchoice.com",
	},
	License: struct {
		Name string
		URL  string
	}{
		Name: "MIT",
		URL:  "https://opensource.org/licenses/MIT",
	},
	Host:     "localhost:8080",
	BasePath: "/",
}
