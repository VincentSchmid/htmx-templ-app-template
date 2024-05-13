//go:build dev

package main

import (
	"fmt"
	"net/http"
)

func publicHandler() http.Handler {
	fmt.Println("dev mode")
	return http.StripPrefix("/public/", http.FileServer(http.Dir("public")))
}
