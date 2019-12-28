package parser

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

type Parser struct {
	Links []Link
}

func New() *Parser {
	return &Parser{
		Links: []Link{},
	}
}

func (p *Parser) GetResponseFor(url string) (string, error) {
	resp, err := http.Get("http://example.com")
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get virtual host: %s", err)
		fmt.Println(errMsg)
		return "", fmt.Errorf(errMsg)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to parse response body: %s", err)
		fmt.Println(errMsg)
		return "", fmt.Errorf(errMsg)
	}

	fmt.Println("Result ", string(body))
	return string(body), nil
}

func (p *Parser) ExtractAnchorTagsFrom(responseBody string) error {
	if responseBody == "" {
		return fmt.Errorf("The response body was empty")
	}
	r := strings.NewReader(responseBody)
	z := html.NewTokenizer(r)

	var temp *Link

	depth := 0
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			fmt.Println("error", z.Err().Error() == "EOF")
			err := z.Err()
			if err.Error() == "EOF" {
				fmt.Println("End of file")
				fmt.Println("Links", p.Links)
				return nil
			}

			return z.Err()
		case html.TextToken:
			if depth > 0 {
				trimmedText := strings.TrimSpace(string(z.Text()))
				if trimmedText != "" {
					fmt.Println("text and depth", trimmedText, depth)
					if temp == nil {
						fmt.Println("Link should not be nil when adding text")
						continue
					}

					temp.Text += fmt.Sprintf("%s ", trimmedText)
				}
			}
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()
			if len(tn) == 1 && tn[0] == 'a' {
				if tt == html.StartTagToken {
					_, val, extra := z.TagAttr()
					temp = &Link{
						Href: string(val),
					}
					fmt.Println("Starting anchor tag - creating link", string(val), extra)
					depth++
				} else {
					depth--
					if depth == 0 {
						fmt.Println("Closing anchor tag - adding to list")

						p.Links = append(p.Links, Link{
							Text: strings.TrimSpace(temp.Text),
							Href: temp.Href,
						})

						temp = nil
					}
				}
			}
		}
	}

	return nil
}
