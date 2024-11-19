package main

import (
	"fmt"

	"github.com/dlbarduzzi/guestbook/internal/service"
)

func main() {
	port := 8080

	svc := service.NewService(port)
	fmt.Println(svc)
}
