package parser

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseHtml(t *testing.T) {
	html := `<!doctype html>
		<html>
			<head>
				<link rel="stylesheet" href="/styles/main.css">
			</head>
			<body>
				<img src="/img/logo.png">
				<a href="/about">About</a>
			</body>
		</html>`
	base, _ := url.Parse("https://example.com")

	links, newBody := ParseHtml([]byte(html), base)

	require.Len(t, links, 3)

	bodyStr := string(newBody)
	require.Contains(t, bodyStr, "/styles/main.css")
	require.Contains(t, bodyStr, "/img/logo.png")
	require.Contains(t, bodyStr, "/about")
}
