package main

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"regexp"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
)

const tmpl = `
<html>
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>GoGoGo</title>
    <script src="https://cdn.tailwindcss.com"></script>
  </head>

  <body>
    <div class="flex justify-center py-12">
      <div class="animate-ping h-2 w-2 bg-blue-600 rounded-full"></div>
      <div class="animate-ping h-2 w-2 bg-blue-600 rounded-full mx-4"></div>
      <div class="animate-ping h-2 w-2 bg-blue-600 rounded-full"></div>
    </div>

    <script src="wasm_exec.js"></script>
    <script>
      if (!WebAssembly.instantiateStreaming) {
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
          const source = await (await resp).arrayBuffer();
          return await WebAssembly.instantiate(source, importObject);
        };
      }

      const go = new Go();

      (async () => {
        const result = await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
        await go.run(result.instance);
        await WebAssembly.instantiate(result.module, go.importObject);
      })()
    </script>
  </body>
</html>
`

func main() {
	tpl, err := template.New("index.html").Parse(tmpl)
	if err != nil {
		log.Fatalln(err)
	}

	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, nil); err != nil {
		log.Fatalln(err)
	}

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

	if err := os.MkdirAll("dist", os.ModePerm); err != nil {
		log.Fatalln(err)
	}
	f, err := os.Create("dist/index.html")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	if err := m.Minify("text/html", f, buf); err != nil {
		log.Fatalln(err)
	}
}
