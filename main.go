package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

var source = flag.String("source", "", "Directory containing the original images")
var width = flag.Uint("width", 1000, "Desired width of the derivate")
var height = flag.Uint("height", 1000, "Height of the desired derivate")
var destination = flag.String("destination", "", "Destination directory for the created derivates")
var pathSeparator = fmt.Sprintf("%c", os.PathSeparator)
var logDir = "logs"
var logFile = logDir + "/converter.log"

func main() {
	flag.Parse()

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Could not create log directory\n")
		}
	}
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Could not start logger\n")
	}
	defer f.Close()

	log.SetOutput(f)

	walkErr := filepath.Walk(*source, findFiles)
	if walkErr != nil {
		log.Fatalf("Error while walking: %v\n", err)
	}
}

func findFiles(path string, fi os.FileInfo, err error) error {
	switch mode := fi.Mode(); {
	case mode.IsDir():
	case mode.IsRegular():
		convert(path)
	}

	return nil
}

//converts the image
func convert(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Panicln(err)
	}

	convertible, extension := canConvert(*file)
	if convertible {
		log.Printf("Converting %s -> %s\n", extension, file.Name())
		image, _, err := image.Decode(file)
		if err != nil {
			log.Panicln(err)
		}
		file.Close()

		thumbnailImage := resize.Thumbnail(*width, *height, image, resize.Lanczos3)

		saveErr := saveImage(thumbnailImage, file.Name())
		if saveErr != nil {
			log.Panic(err)
		}
	} else {
		log.Printf("Conversion not doable for type %s -> %s\n", extension, file.Name())
	}
}

//checks, if the file is even convertible by filetype
func canConvert(file os.File) (bool, string) {
	extension := filepath.Ext(file.Name())

	if extension == ".jpg" || extension == ".gif" || extension == ".png" {
		return true, extension
	}
	return false, extension
}

func saveImage(img image.Image, sourcePath string) error {
	baseName := filepath.Base(sourcePath)
	destination := *destination
	createPath(destination)
	destinationFile := destination + pathSeparator + baseName
	out, err := os.Create(destinationFile)
	if err != nil {
		log.Panicln(err)
		return err
	}
	defer out.Close()
	encErr := png.Encode(out, img)
	if encErr != nil {
		out.Close()
		log.Panicln(encErr)
	}
	return nil
}

func createPath(path string) {
	_ = os.MkdirAll(path, os.ModePerm)
}
