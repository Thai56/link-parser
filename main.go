package main

import (
	"fmt"
	"log"

	"github.com/html-link-parser/parser"
)

func main() {
	p := parser.New()
	s, err := p.GetResponseFor("http://example.com")
	if err != nil {
		log.Fatal(err)
	}
	// s := `<p>Links:</p><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul>`

	err = p.ExtractAnchorTagsFrom(s)

	fmt.Println("error", err)
}
