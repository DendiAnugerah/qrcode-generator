package handler

import (
	base "dend-qrcode"
	"flag"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
)

var (
	outputDirectory = flag.String("outputDirectory", "output", "Directory to save the generated QR codes")
)

func HandlerQRCode(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to read image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	decoder := base.DefaultEncoder

	text := r.FormValue("text")
	logo, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Failed to decode image file", http.StatusInternalServerError)
		return
	}

	qrCode, err := decoder.Encode(text, logo, 512)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	if _, err := os.Stat(*outputDirectory); os.IsNotExist(err) {
		err = os.Mkdir(*outputDirectory, 0755)
		if err != nil {
			http.Error(w, "Failed to create output directory", http.StatusInternalServerError)
			return
		}
	}

	outputPath := fmt.Sprintf("%s/qr.png", *outputDirectory)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		http.Error(w, "Failed to create output file", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, qrCode)
	if err != nil {
		http.Error(w, "Failed to save QR code", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "QR code with logo created and saved at %s\n", outputPath)
}
