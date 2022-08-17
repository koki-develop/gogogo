package main

import (
	"github.com/koki-develop/gogogo/backend/pkg/infrastructures/router"
)

// ローカル実行用
func main() {
	r := router.New()
	if err := r.Run(); err != nil {
		panic(err)
	}
}
