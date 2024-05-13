package main

import (
	"github.com/VincentSchmid/htmx-templ-app-template/cmd"
)

func main() {
	cmd.PublicHandler = publicHandler()
	cmd.Execute()
}
