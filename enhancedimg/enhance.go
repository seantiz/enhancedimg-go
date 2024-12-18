package enhancedimg

import (
	"bytes"
	"fmt"
	"image"
	"math"
	"os"
	"path/filepath"
	"strings"

	"image/gif"
	"image/jpeg"
	"image/png"
)

var fallbacks map[string]selectedEncoder

type enhancedImg struct {
	sourceImagePath string
	sources         map[string][]string
	img             struct {
		src    string
		width  int
		height int
	}
}

type selectedEncoder struct {
	name    string
	encode  func(img image.Image, quality int) ([]byte, error)
	quality int
}

type sizeVariant struct {
	width int
	label string
}

func (ei enhancedImg) aspectRatio() float64 {
	return float64(ei.img.height) / float64(ei.img.width)
}

func init() {
	fallbacks = map[string]selectedEncoder{
		".avif": {"png", encodePNG, 0},
		".gif":  {"gif", encodeGIF, 0},
		".heif": {"jpg", encodeJPEG, 85},
		".jpeg": {"jpg", encodeJPEG, 85},
		".jpg":  {"jpg", encodeJPEG, 85},
		".png":  {"png", encodePNG, 0},
		".tiff": {"jpg", encodeJPEG, 85},
		".webp": {"png", encodePNG, 0},
	}
}

func encodeJPEG(img image.Image, quality int) ([]byte, error) {
	var buf bytes.Buffer
	options := &jpeg.Options{
		Quality: quality,
	}
	if err := jpeg.Encode(&buf, img, options); err != nil {
		return nil, fmt.Errorf("unexpected error encoding JPEG: %w", err)
	}
	return buf.Bytes(), nil
}

func encodePNG(img image.Image, _ int) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encodeGIF(img image.Image, _ int) ([]byte, error) {
	var buf bytes.Buffer
	if err := gif.Encode(&buf, img, nil); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func selectEncoder(inputFormat string) selectedEncoder {
	if encoder, ok := fallbacks[strings.ToLower(inputFormat)]; ok {
		return encoder
	}
	return selectedEncoder{"jpg", encodeJPEG, 85}
}

func enhanceImage(src string) (enhancedImg, error) {
	src = strings.TrimPrefix(src, "/")

	f, err := os.Open(src)
	if err != nil {
		return enhancedImg{}, fmt.Errorf("couldn't open file %s: %w", src, err)
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		return enhancedImg{}, fmt.Errorf("couldn't decode image %s (format: %s): %w", src, format, err)
	}

	inputExt := filepath.Ext(src)
	encoder := selectEncoder(inputExt)

	formats := []selectedEncoder{
		encoder,
	}

	bounds := img.Bounds()
	enhancedImage := enhancedImg{
		sourceImagePath: src,
		sources:         make(map[string][]string),
		img: struct {
			src    string
			width  int
			height int
		}{
			src:    src,
			width:  bounds.Dx(),
			height: bounds.Dy(),
		},
	}

	// This helper ensures we don't return larger sizes than source image
	sizes := calculateSizeVariants(bounds.Dx())

	baseFileName := filepath.Base(src)
	baseFileName = strings.TrimSuffix(baseFileName, filepath.Ext(baseFileName))

	for _, format := range formats {
		if len(formats) == 0 {
			return enhancedImg{}, fmt.Errorf("couldn't find any valid images")
		}
		srcsets := []string{}

		for _, size := range sizes {
			height := int(math.Round(float64(size.width) * enhancedImage.aspectRatio()))
			resized := resizeImage(img, size.width, height)
			processed, err := format.encode(resized, format.quality)
			if err != nil {
				continue
			}

			processedFileName := fmt.Sprintf("%s-%s-%d.%s",
				baseFileName,
				size.label,
				size.width,
				format.name,
			)

			outPath := filepath.Join("static", "processed", processedFileName)
			if err := os.WriteFile(outPath, processed, 0644); err != nil {
				continue
			}

			srcsets = append(srcsets, fmt.Sprintf("/static/processed/%s %dw",
				processedFileName,
				size.width,
			))
		}

		if len(srcsets) > 0 {
			enhancedImage.sources[format.name] = srcsets
		}
	}

	return enhancedImage, nil
}

func calculateSizeVariants(originalWidth int) []sizeVariant {
	if originalWidth <= 0 {
		return []sizeVariant{}
	}
	// Sizes beyond 2xl follow SvelteJS's reasoning and https://screensiz.es/ common device sizes
	sizes := []sizeVariant{
		{width: 640, label: "sm"},
		{width: 768, label: "md"},
		{width: 1024, label: "lg"},
		{width: 1280, label: "xl"},
		{width: 1536, label: "2xl"},
		{width: 1920, label: "3xl"}, // 1080p
		{width: 2560, label: "4xl"}, // 2K
		{width: 3000, label: "5xl"},
		{width: 4096, label: "6xl"}, // 4K
		{width: 5120, label: "7xl"},
	}

	var result []sizeVariant
	for _, size := range sizes {
		if size.width > originalWidth {
			break
		}
		result = append(result, size)
	}

	if len(result) == 0 || result[len(result)-1].width != originalWidth {
		result = append(result, sizeVariant{
			width: originalWidth,
			label: "original",
		})
	}

	return result
}