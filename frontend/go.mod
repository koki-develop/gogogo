module github.com/koki-develop/gogogo/frontend

go 1.19

require (
	github.com/hexops/vecty v0.6.0
	github.com/koki-develop/gogogo/backend v0.0.0-00010101000000-000000000000
	github.com/tdewolff/minify/v2 v2.20.15
)

require github.com/tdewolff/parse/v2 v2.7.10 // indirect

replace github.com/koki-develop/gogogo/backend => ../backend
