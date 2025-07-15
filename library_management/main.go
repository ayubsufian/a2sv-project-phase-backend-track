package main

import (
	"fmt"
	"library_management/controllers"
	"library_management/services"
)

func main() {
	libService := services.NewLibrary()
	controllers.StartCLI(libService)
	fmt.Println("Thank you for using the Library Management System!")
}
