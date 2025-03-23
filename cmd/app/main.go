package main

import "github.com/coffee/support/internal/app"

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	a := app.New()

	a.Run()
}
