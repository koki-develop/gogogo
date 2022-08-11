package util

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

func NewLoading() *vecty.HTML {
	root := WithClasses(
		elem.Div(
			newLoadingDot(),
			WithClasses(newLoadingDot(), "mx-4"),
			newLoadingDot(),
		),
		"flex", "justify-center", "py-12",
	)

	return root
}

func newLoadingDot() *vecty.HTML {
	return WithClasses(elem.Div(), "animate-ping", "h-2", "w-2", "bg-blue-600", "rounded-full")
}
