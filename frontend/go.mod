module github.com/koki-develop/gogogo/frontend

go 1.19

require (
	github.com/hexops/vecty v0.6.0
	github.com/koki-develop/gogogo/backend v0.0.0-20230615232836-b787e366d61d
	github.com/tdewolff/minify/v2 v2.12.4
)

require github.com/tdewolff/parse/v2 v2.6.4 // indirect

replace github.com/koki-develop/gogogo/backend => ../backend
