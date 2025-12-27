package collector

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/M-kos/wb_level2/task_16/internal/parser"
)

var mimeTypeToExt = map[string]string{
	"text/html":              ".html",
	"application/xhtml+xml":  ".xhtml",
	"image/jpeg":             ".jpg",
	"image/png":              ".png",
	"application/javascript": ".js",
	"text/css":               ".css",
	"application/json":       ".json",
	"application/xml":        ".xml",
	"application/pdf":        ".pdf",
	"image/gif":              ".gif",
	"image/svg+xml":          ".svg",
	"audio/mpeg":             ".mp3",
	"video/mp4":              ".mp4",
}

type Semaphore interface {
	Acquire()
	Release()
}

type Collector struct {
	urls      []*url.URL
	resultDir string
	visited   map[string]struct{}
	client    *http.Client
	mu        sync.RWMutex
	semaphore Semaphore
}

func NewCollector(urls []*url.URL, semaphore Semaphore, resultDir string, requestTimeout time.Duration) *Collector {
	return &Collector{
		urls:      urls,
		resultDir: resultDir,
		visited:   make(map[string]struct{}),
		client: &http.Client{
			Timeout: requestTimeout,
		},
		semaphore: semaphore,
	}
}

func (c *Collector) Start(ctx context.Context) error {
	var wg sync.WaitGroup

	for _, u := range c.urls {
		wg.Add(1)
		c.semaphore.Acquire()

		go func() {
			defer c.semaphore.Release()
			defer wg.Done()
			c.worker(ctx, u)
		}()
	}

	wg.Wait()

	return nil
}

func (c *Collector) worker(ctx context.Context, targetUrl *url.URL) {
	queue := []*url.URL{targetUrl}

	for len(queue) > 0 {
		select {
		case <-ctx.Done():
			return
		default:
			task := queue[0]
			queue = queue[1:]

			urlStr := task.String()

			c.mu.Lock()
			if _, ok := c.visited[urlStr]; ok {
				c.mu.Unlock()
				continue
			}

			c.visited[urlStr] = struct{}{}
			c.mu.Unlock()

			body, mimeType, err := c.download(ctx, task)
			if err != nil {
				fmt.Printf("download error: %v\n", err)
				continue
			}

			localPath := c.localPath(targetUrl, task, mimeType)

			if mimeType == "text/html" {
				links, htmlBody := parser.ParseHtml(body, task)
				body = htmlBody
				queue = append(queue, links...)
			}

			err = c.saveFile(localPath, body)
			if err != nil {
				fmt.Println("save error:", err)
				continue
			}

		}
	}
}

func (c *Collector) download(ctx context.Context, targetUrl *url.URL) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", targetUrl.String(), nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Println("error closing response body:", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download %s: status code %d", targetUrl.String(), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	mimeType := strings.Split(resp.Header.Get("Content-Type"), ";")[0]

	return body, mimeType, err
}

func (c *Collector) localPath(targetUrl *url.URL, currentUrl *url.URL, mimeType string) string {
	urlPath := currentUrl.Path
	if urlPath == "" || strings.HasSuffix(urlPath, "/") {
		urlPath += "index"
	}

	ext, ok := mimeTypeToExt[mimeType]
	if ok && ext != path.Ext(urlPath) {
		urlPath += ext
	}

	return filepath.Join(c.resultDir, targetUrl.Host, urlPath)
}

func (c *Collector) saveFile(path string, payload []byte) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println("error closing file:", err)
		}
	}()

	_, err = file.Write(payload)
	if err != nil {
		return err
	}

	return nil
}
