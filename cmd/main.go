package main

import (
	"fmt"
	"github.com/joelewaldo/go-micro-service/pkg/hello"
	"net/http"
)

func main() {
	fmt.Println("Running from cmd/go")
	hello.Test()
}
