package helper

import (
	base "dend-qrcode/internal/base"
	"flag"
	"fmt"
	"image"
	"os"
)

var (
	input  = flag.String("i", "logo.png", "Logo to be placed over QR code")
	output = flag.String("o", "qr.png", "Output filename")
	size   = flag.Int("size", 512, "Image size in pixels")
)

func Command() {
	flag.Usage = Usage
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	text := flag.Arg(0)
	file, err := os.Open(*input)
	errcheck(err, "Failed to open logo file")
	defer file.Close()

	logo, _, err := image.Decode(file)
	errcheck(err, "Failed to decode PNG logo file")

	qrcod, err := base.Encode(text, logo, *size)
	errcheck(err, "Failed to encode QR code")

	out, err := os.Create(*output)
	errcheck(err, "Failed to create output file")

	out.Write(qrcod.Bytes())
	out.Close()

	fmt.Println("QR code with logo created in", *output)
}

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage: go run . [options] text")
	flag.PrintDefaults()
}

func errcheck(err error, str string) {
	if err != nil {
		fmt.Println(str, err)
		os.Exit(1)
	}
}
