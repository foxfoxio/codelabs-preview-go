package main

import (
	"bytes"
	"fmt"
	"github.com/graemephi/goldmark-qjs-katex"
	"github.com/yuin/goldmark"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	md := goldmark.New(
		goldmark.WithExtensions(&qjskatex.Extension{}),
		goldmark.WithRendererOptions(
			gmhtml.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	in := []byte(`$$a_{n}=\dfrac{1}{\pi }\int_{-\pi }^{\pi }f(x)\cos nx\;dx\qquad n=0,1,2,3,...,$$`)
	if err := md.Convert(in, &buf); err != nil {
		fmt.Printf("Failed to convert %s: %s", in, err)
	} else {
		html := strings.Replace(htmlTemplate, "{{content}}", buf.String(), 1)
		ioutil.WriteFile("katex_out.html", []byte(html), os.ModePerm)
	}
}

const htmlTemplate = `
<html>
<head>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.12.0/katex.min.css"
          integrity="sha512-h7nl+xz8wgDlNM4NqKEM4F1NkIRS17M9+uJwIGwuo8vGqIl4BhuCKdxjWEINm+xyrUjNCnK5dCrhM0sj+wTIXw=="
          crossorigin="anonymous"/>
</head>
<body>
{{content}}
</body>
</html>
`
