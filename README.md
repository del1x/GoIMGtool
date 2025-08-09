# GoIMGtool 🖼️

**GoIMGtool** is an image processing tool that adds watermarks to images.
- 🔍 **Size Optimization** - Resizes images to 1200×1200 pixels if larger.
- 💧 **Watermark Overlay** - Applies a semi-transparent watermark.
- 🚀 **Format Support** - Works with PNG, JPG, JPEG, and WebP.

Perfect for:
- Batch processing photos before uploading to a website.
- Adding logos to images.
- Adding routine image tasks.

## Requirements
- Go 1.24.6
- MSYS2 with `mingw-w64-x86_64-libwebp` for WebP support.

## Installation
1. Install Go 1.24.6[](https://golang.org/dl/).
2. Install MSYS2[](https://www.msys2.org/) and run:
   ```bash
   pacman -S mingw-w64-x86_64-libwebp
3. Set CGO environment variables (Windows cmd):
   ```bash
    set CGO_CFLAGS=-IC:/MSYS2/mingw64/include
    set CGO_LDFLAGS=-LC:/MSYS2/mingw64/lib -lwebp
    set CGO_ENABLED=1
    set PATH=С:\MSYS2\mingw64\bin;%PATH% 
4. Clone the repository and initialize:
   ```bash
   git clone https://github.com/del1x/GoIMGtool.git
    cd GoIMGtool
    go mod init
    go mod tidy

## Usage
1. Place images in the Images/ folder.
2. Add watermark.png to the same folder.
3. Run program
   ```bash
    go run -tags "desktop" main.go
4. Select the images folder, choose output format (webp/png), and process.
    - Processed files will appear in Images_watermarked/.

## Notes
 - Ensure watermark.png is a valid image.
 - For WebP support, CGO settings are required.