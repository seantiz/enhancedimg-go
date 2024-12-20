# Enhanced Images (Go)

Use this for serving optimised images in your Go + HTML/X web apps. Please remember to run this as a pre-processor before any `go build` command.

### What It Actually Does

__Enhanced Images Go__ parses your app's source HTML templates and replaces all optimisable `<img>` elements with a `<picture>` element containing the responsive-friendly images inside.

It's still in alpha because I've yet to implement ability to apply transformation effects to target images.

## Quick Start

Installation options:

1. Install globally:
   ```bash
   go install github.com/seantiz/enhancedimg-go@latest
   ```

   Then run with `enhancedimg-go` command

2. Via Makefile: Create a Makefile to run `enhanceddimg-go` programatically before your 'go build' command.

## Features

- Responsive web design supported but no data transforms handled yet.
- All common device sizes supported - ranging from Tailwind's small ("sm") device breakpoint and up to 4K (and slightly beyond).
- AVIF source images not fully supported and will be converted to PNG.
- HEIF and TIFF source images will be converted to JPEG.

## Known Limitations

Right now this library (as of v0.6.0) deals purely with injecting picture elements and responsive image variants only. Actual image enhancement (like transforms) are yet to be implemented.

## Long-Term Performance?

Pre-processing time increases with content size - especially image content, less so with HTML structure holding the images. Here's a performance test with just 6 input images (and this is without transforms as part of the job!):

```bash
HTML parsing: 132.917µs (0.132917 milliseconds)
Image processing: 15.704458ms (15.704458 milliseconds)
```
Under bulk testing, processing 512 images in one run took just over 40 seconds.

Image processing takes around 118 times longer than parsing the HTML.

One option is to implement a version of enhancedimg-go that uses `bimg` for performance. But I decided against going down that route, for now, because it means:

1. Any developer would have to install `libvips` as a system-level dependency
2. I don't know any medium or entreprise-scale web apps published with a Go + HTML/X stack as of yet, so content size (and processing time) seems like it has a modest upper bounds in most cases right now.

Feel free to get in touch if you disagree or you're running a Go + HTML/X codebase where the content size demands greater performance.

## Credit Where It's Due

The logic behind this package was heavily lifted from the SvelteJS enhanced-img preprocessor (being a big SvelteJS user myself) and Tailwind.