package service

import "fmt"

func NewService(port int) string {
	return fmt.Sprintf("Port %d\n", port)
}
