package parser

import (
	"bytes"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

func ParseHtml(body []byte, baseURL *url.URL) ([]*url.URL, []byte) {
	htmlDocument, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, body
	}

	var links []*url.URL

	nodes := []*html.Node{htmlDocument.FirstChild}

	for len(nodes) != 0 {
		node := nodes[0]
		nodes = nodes[1:]

		if node.NextSibling != nil {
			nodes = append(nodes, node.NextSibling)
		}

		if node.FirstChild != nil {
			nodes = append(nodes, node.FirstChild)
		}

		if node.Type == html.ElementNode {
			for i, attr := range node.Attr {
				if isLinkAttr(node.Data, attr.Key) {
					link, err := baseURL.Parse(attr.Val)

					if err == nil {
						linkDomain, err := publicsuffix.EffectiveTLDPlusOne(link.Host)
						if err != nil {
							continue
						}
						baseDomain, err := publicsuffix.EffectiveTLDPlusOne(baseURL.Host)
						if err != nil || linkDomain != baseDomain {
							continue
						}

						links = append(links, link)
						node.Attr[i].Val = relativeLocalPath(link)
					}
				}
			}
		}
	}

	var buf bytes.Buffer
	err = html.Render(&buf, htmlDocument)
	if err != nil {
		return nil, body
	}

	return links, buf.Bytes()
}

func isLinkAttr(tag, attr string) bool {
	switch tag {
	case "a", "link":
		return attr == "href"
	case "img", "script":
		return attr == "src"
	}

	return false
}

func relativeLocalPath(u *url.URL) string {
	uPath := u.Path
	if uPath == "" || strings.HasSuffix(uPath, "/") {
		uPath += "index.html"
	}

	return uPath
}
