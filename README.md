# Enhanced Image Go

Use this for serving optimised images in your HTML/X web apps. Please remember to run this as a pre-processor before any build command.

### What It Actually Does

__Enhanced Image Go__ parses your app's source HTML templates and replaces all optimisable `<img>` elements with a `<picture>` element containing the responsive-friendly images inside.

It's still in alpha because I've yet to implement ability to apply transformation effects to target images.

## Quick Start

Installation options:

1. Install globally:
   ```bash
   go install github.com/seantiz/enhancedimg-go@latest
   ```

   Then run with `enhancedimg-go` command

2. Via Makefile: Create a Makefile to run `enhanceddimg-go` programatically before your build command.

## Features

- Responsive web design supported but no data transforms handled yet.
- All common device sizes supported - potentially up to ten breakpoints handled (but for efficiency it will never generate srcset images that were bigger than your original image).
- AVIF source images not fully supported and will be converted to PNG.
- HEIF and TIFF source images will be converted to JPEG.

## IMPORTANT NOTE: Do You Really Need Responsive Image Sets?

I discovered a practical kink in the road that I probably should've known before I built this package: Modern browsers use device pixel ratio (DPR) to select from responsive image sets.

This can lead to unexpected behavior: High-DPR devices (most modern smartphones and laptops) apply a 2x or 3x multiplier when selecting the right image from your responsive image srcset.

This means I could be browsing your optimised website on an iPhone or Android with a very modest screen size, and my browser will still download your largest, ultimate image variant from within your sets.

It may be a better strategy to aggressively convert all your images to WEBP with data cache attributes and, in the words of @CaptainCodeman, "cache the beeejezus out of them."

### When To Use An Aggressive WEBP Strategy

If you know the likelihood is high that your userbase will be browsing your site from high-DPR devices most of the time, then arguably your server should only store the largest image available in WEBP format (with the appropriate data cache attributes) and let the browser + browser cache handle the rest.

Save yourself pre-processing time and storage space going down this road, rather than storing smaller image variants in a set that will never be matched by the browser.

### When to Use Full Responsive HTML Strategy (Enhanced Image Go)

The full responsive approach remains valuable when:

- You have diverse users or you don't know your userbase
- You need to support legacy browsers or devices
- Bandwidth optimization is critical for your use case

## Known Limitations

Right now this library (as of v0.6.0) deals purely with injecting responsive picture elements only; actual image-processing enhancements (like transforms) are yet to be implemented.

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
2. I don't know any medium or entreprise-scale web apps published through an "HTML by wire" model that this tool is built for (e.g. optimising images in a Go + HTML/X web stack), so content size (and processing time) seems like it has a modest upper bounds in most cases right now.

Feel free to get in touch if you disagree or you're running a Go + HTML/X codebase where the content size demands greater performance.

## Credit Where It's Due

The logic behind this package was heavily lifted from the SvelteJS enhanced-img preprocessor (being a big SvelteJS user myself) and Tailwind.