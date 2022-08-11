package main

import (
	"os"

	"github.com/koki-develop/gogogo/backend/pkg/infrastructures/router"
)

// ローカル実行用
func main() {
	os.Setenv("IS_LOCAL", "true")

	r := router.New()
	if err := r.Run(); err != nil {
		panic(err)
	}
}
