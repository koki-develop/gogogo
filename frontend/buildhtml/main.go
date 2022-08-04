package main

import (
	"html/template"
	"log"
	"os"
)

const tmpl = `
<html>
  <head>
    <meta charset="utf-8" />
    <title>GoGoGo</title>
  </head>

  <body>
    <p>Loading...</p>

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

	f, err := os.Create("dist/index.html")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	if err := tpl.Execute(f, nil); err != nil {
		log.Fatalln(err)
	}
}
