package enhancedimg

import (
	"image"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"

	"golang.org/x/net/html"
)

func parseHTMLFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func saveHTMLFile(content string, path string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func optimisable(src string) bool {
	src = strings.TrimPrefix(src, "/")
	ext := strings.ToLower(filepath.Ext(src))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp"
}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func setAttr(n *html.Node, key, value string) {
	for i, attr := range n.Attr {
		if attr.Key == key {
			n.Attr[i].Val = value
			return
		}
	}
	n.Attr = append(n.Attr, html.Attribute{Key: key, Val: value})
}

func findImageTags(n *html.Node, tags *[]*html.Node) {
	if n.Type == html.ElementNode && n.Data == "img" {
		*tags = append(*tags, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findImageTags(c, tags)
	}
}

func resizeImage(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
	return dst
}

func getOriginalTag(n *html.Node) string {
	var buf strings.Builder
	html.Render(&buf, n)
	return buf.String()
}
