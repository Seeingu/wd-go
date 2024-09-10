package main

import (
	"golang.org/x/net/html"
	"os"
	"strings"
)

type HtmlRef struct {
	scripts []string
	styles  []string
}

func htmlBody(doc *html.Node) (HtmlRef, error) {
	htmlRef := HtmlRef{
		scripts: []string{},
		styles:  []string{},
	}
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "script" {
			for _, attr := range node.Attr {
				if attr.Key == "src" && !strings.HasPrefix(attr.Val, "http") {
					var fileUrl = attr.Val
					if strings.HasPrefix(fileUrl, "/") {
						fileUrl = fileUrl[1:]
					}
					htmlRef.scripts = append(htmlRef.scripts, fileUrl)
				}
			}
			return
		}
		if node.Type == html.ElementNode && node.Data == "link" {
			for _, attr := range node.Attr {
				// TODO: should check rel
				if attr.Key == "href" && !strings.HasPrefix(attr.Val, "http") {
					var fileUrl = attr.Val
					if strings.HasPrefix(fileUrl, "/") {
						fileUrl = fileUrl[1:]
					}
					htmlRef.styles = append(htmlRef.styles, fileUrl)
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	return htmlRef, nil
}

func HtmlCrossReference(filePath string) (HtmlRef, error) {
	r, err := os.Open(filePath)
	if err != nil {
		return HtmlRef{}, err
	}
	doc, err := html.Parse(r)
	if err != nil {
		return HtmlRef{}, err
	}
	ref, err := htmlBody(doc)
	return ref, err
}
