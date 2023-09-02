package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
    htmlFileName := flag.String("html", "index.html", "html file name")
	flag.Parse()

	file, err := os.ReadFile(*htmlFileName)
	if err != nil {
		fmt.Printf("Failed to open html file: %s\n ", *htmlFileName)
	    os.Exit(1)
	}

	r := strings.NewReader(string(file))
	links, err := parse(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", links)
}

type link struct {
	Href string
	Text string
}

func parse(r io.Reader) ([]link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linkNodes(doc)
	var links []link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}

func buildLink(n *html.Node) link {
	var ret link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = text(n)
	return ret
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += text(c) + ""
	}
	return strings.Join(strings.Fields(ret), "")
}

func linkNodes(n *html.Node) []*html.Node  {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}

