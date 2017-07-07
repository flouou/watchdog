package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"

	"path"

	"github.com/nfnt/resize"
)

var source = flag.String("source", "", "Directory containing the original images")
var width = flag.Uint("width", 1000, "Desired width of the derivate")
var height = flag.Uint("height", 1000, "Height of the desired derivate")
var destination = flag.String("destination", "", "Destination directory for the created derivates")

func main() {
	flag.Parse()

	err := filepath.Walk(*source, findFiles)
	if err != nil {
		fmt.Printf("Error while walking: %v\n", err)
	}
}

func findFiles(path string, fi os.FileInfo, err error) error {
	//fmt.Printf("Converted: %s\n", path)
	switch mode := fi.Mode(); {
	case mode.IsDir():
	case mode.IsRegular():
		convert(path)
	}

	return nil
}

//converts the image
func convert(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	//fmt.Printf("Found a file: %s\n", file.Name())
	//fmt.Println(canConvert(*file))
	convertible, extension := canConvert(*file)
	if convertible {
		fmt.Printf("Converting %s\n", extension)
		fmt.Print(file)
		image, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err)
			return err
		}
		file.Close()

		thumbnailImage := resize.Thumbnail(*width, *height, image, resize.Lanczos3)

		saveErr := saveImage(thumbnailImage, file.Name())
		if saveErr != nil {
			return saveErr
		}
	} else {
		fmt.Printf("Conversion not doable for type %s\n", extension)
	}
	return nil
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
	baseName := path.Base(sourcePath)
	destination := *destination
	createPath(destination)
	destinationFile := destination + baseName
	fmt.Printf("something even stranger")
	out, err := os.Create(destinationFile)
	if err != nil {
		return err
	}
	defer out.Close()
	//png.Encode(out, img, nil)
	return nil
}

func createPath(path string) {
	fmt.Printf("something strange")
	_ = os.MkdirAll(path, os.ModePerm)
}
