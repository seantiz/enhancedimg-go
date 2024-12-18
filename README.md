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
- WEBP, AVIF not yet supported but are handled by falling back to PNG on conversion.
- HEIF and TIFF source images are converted to JPEG.

## Dependencies and Limitations (PRs welcome)

1. I leaned on the standard Go HTML parser in this package. I'll look at any way to shave down this dependency, if possible.

2. I couldn't find a clean way to encode WebP images in Go just yet, despite the fact `optimisable()` utility tries to match the `.webp` format during parsing. The solutions I find meant pulling in more dependencies.

3. Right now this library (as of v0.3.0) deals purely with resizing images for responsive design but actual image enhancement still needs to be implement by reading all metadata for transformation attributes.

## Long-Term Performance?

Pre-processing time increases with content size, so I initially built a version of enhancedimg-go that uses `bimg` for a more performant image-processing step.

I decided against going down that route because it means:

1. Any developer would have to install `libvips` as a system-level dependency
2. I don't know any medium or entreprise-scale web apps published with a Go + HTML/X stack as of yet, so content size (and processing time) seems like it has a modest upper bounds in most cases right now.

Feel free to get in touch if you disagree or you're running a Go + HTML/X codebase where the content size demands greater performance.

## Credit Where It's Due

The logic behind this package was heavily lifted from the SvelteJS enhanced-img preprocessor (being a big Svelte user myself) and Tailwind.