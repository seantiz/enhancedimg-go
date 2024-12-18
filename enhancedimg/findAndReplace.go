// Package enhancedimg provides image optimisation for HTML template content.
// The package supports JPEG, PNG and GIF formats. WEBP, AVIF, HEIF and TIFF are not fully supported
// but will fallback to supported formats upon conversion to target image files.
// The package can generate size variants for responsive web design and common-device breakpoints
// while maintaining aspect ratios rounded to the nearest pixel.
package enhancedimg

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// FindAllImageElements tracks difference between source and target HTML before saving all changes to target.
// The map avoids any problems with hot reloading in a dev environment.
func FindAllImageElements(templatesDir string) error {
	processedFiles := make(map[string]bool)

	return filepath.Walk(templatesDir, func(path string, info fs.FileInfo, walkErr error) error {
		if !strings.HasSuffix(path, "html") {
			return nil
		}

		absPath, _ := filepath.Abs(path)
		if processedFiles[absPath] {
			return nil
		}
		processedFiles[absPath] = true

		if walkErr != nil {
			return walkErr
		}

		content, err := parseHTMLFile(path)
		if err != nil {
			return err
		}

		originalContent := content
		doc, err := html.Parse(strings.NewReader(content))
		if err != nil {
			return err
		}

		var imageTags []*html.Node
		findImageTags(doc, &imageTags)

		modified := content
		changesMade := false

		for _, img := range imageTags {
			src := getAttr(img, "src")
			fmt.Printf("Found image: %s in file: %s\n", src, path)

			_ = strings.ReplaceAll(getOriginalTag(img), " ", "")
			picture, err := replaceImgWithPicture(img)
			if err != nil {
				fmt.Printf("Error processing image %s: %v\n", src, err)
				continue
			}

			if picture != nil {
				var buf strings.Builder
				html.Render(&buf, picture)
				newTag := buf.String()

				regexPattern := fmt.Sprintf(`<img[^>]*?src="%s"[^>]*?>`, regexp.QuoteMeta(src))
				re := regexp.MustCompile(regexPattern)
				modified = re.ReplaceAllString(modified, newTag)

				if modified != content {
					changesMade = true
					fmt.Printf("Successfully replaced %s with picture element\n", src)
				} else {
					fmt.Printf("Failed to replace %s - no changes made\n", src)
				}
			}
		}

		if changesMade && modified != originalContent {
			return saveHTMLFile(modified, path)
		}

		return nil
	})
}

func replaceImgWithPicture(img *html.Node) (*html.Node, error) {
	src := getAttr(img, "src")
	if !optimisable(src) {
		return nil, nil
	}

	enhanced, err := enhanceImage(src)
	if err != nil {
		return nil, err
	}

	picture := createPictureElement(enhanced, img)
	return picture, nil
}

func createPictureElement(ei enhancedImg, originalImg *html.Node) *html.Node {
	picture := &html.Node{
		Type: html.ElementNode,
		Data: "picture",
	}

	for format, srcsets := range ei.sources {
		source := &html.Node{
			Type: html.ElementNode,
			Data: "source",
		}
		setAttr(source, "srcset", strings.Join(srcsets, ", "))
		setAttr(source, "type", "image/"+format)
		setAttr(source, "sizes", "(max-width: 640px) 640px, (max-width: 1024px) 1024px, 1920px")
		picture.AppendChild(source)
	}

	img := &html.Node{
		Type: html.ElementNode,
		Data: "img",
	}
	setAttr(img, "src", "/"+strings.TrimPrefix(ei.img.src, "/"))
	setAttr(img, "width", fmt.Sprintf("%d", ei.img.width))
	setAttr(img, "height", fmt.Sprintf("%d", ei.img.height))
	setAttr(img, "loading", "lazy")
	setAttr(img, "class", getAttr(originalImg, "class"))
	picture.AppendChild(img)

	return picture
}
