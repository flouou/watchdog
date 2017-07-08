package main

import (
	"flag"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"image/gif"

	"github.com/flouou/watchdog/config"
	"github.com/nfnt/resize"
)

//Configuration resembles every configuration made in config.json

var configuration = config.LoadConfig("config.json")
var source = flag.String("source", "", "Directory containing the original images")
var width = flag.Uint("width", 150, "Desired width of the derivate")
var height = flag.Uint("height", 150, "Height of the desired derivate")
var destination = flag.String("destination", "", "Destination directory for the created derivates (default {source}_derivates)")
var logDir = configuration.LogDir
var logFile = filepath.Join(logDir, configuration.LogFile)

func main() {
	flag.Parse()

	if *source == "" {
		log.Fatalf("You have to define a source directory\n")
	}
	if *destination == "" {

		fi, err := os.Stat(*source)
		if err != nil {
			log.Fatalf("Could not scan source directory\n")
		}
		switch mode := fi.Mode(); {
		case mode.IsDir():
			*destination = *source + "_derivates"
		case mode.IsRegular():
			*destination = filepath.Dir(*source) + "_derivates"
		}
		log.Printf("No destination directory given. Using default destination: %s\n", *destination)
	}

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Could not create log directory %s\n", logDir)
		}
	}
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Could not start logger: %s\n", err)
	}
	defer f.Close()

	log.SetOutput(f)

	walkErr := filepath.Walk(*source, findFiles)
	if walkErr != nil {
		log.Fatalf("Error while walking: %sv\n", err)
	}
}

func findFiles(path string, fi os.FileInfo, err error) error {
	switch mode := fi.Mode(); {
	case mode.IsDir():
	case mode.IsRegular():
		err := convert(path)
		if err != nil {
			log.Panicf("Error while converting the file: %s -> %s\n", err, path)
			return err
		}
	}

	return nil
}

//converts the image
func convert(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Panicf("Error while opening file: %s -> %s\n", err, path)
	}

	convertible, extension := canConvert(*file)
	if convertible {
		log.Printf("Converting %s -> %s\n", extension, file.Name())
		image, _, err := image.Decode(file)
		if err != nil {
			log.Printf("Error while decoding file: %s -> Did not decode %s\n", err, file.Name())
			return err
		}
		file.Close()

		thumbnailImage := resize.Thumbnail(*width, *height, image, resize.Lanczos3)

		saveErr := saveImage(thumbnailImage, file.Name(), extension)
		if saveErr != nil {
			log.Panicf("Error while saving file: %s -> Did not save %s\n", saveErr, file.Name())
			return err
		}
	} else {
		file.Close()
		log.Printf("Conversion not doable for type %s -> %s\n", extension, file.Name())
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

func saveImage(img image.Image, sourcePath, extension string) error {
	baseName := filepath.Base(sourcePath)
	destination := *destination
	createPath(destination)
	destinationFile := filepath.Join(destination, baseName)
	out, err := os.Create(destinationFile)
	if err != nil {
		log.Panicf("Error while creating destination file: %s", destinationFile)
		return err
	}
	defer out.Close()

	switch extension {
	case ".png":
		encErr := png.Encode(out, img)
		if encErr != nil {
			out.Close()
			log.Panicf("Error while encoding file: %s -> %s", encErr, out.Name())
			return encErr
		}
	case ".jpg":
		encErr := jpeg.Encode(out, img, nil)
		if encErr != nil {
			out.Close()
			log.Panicf("Error while encoding file: %s -> %s", encErr, out.Name())
			return encErr
		}
	case ".gif":
		encErr := gif.Encode(out, img, nil)
		if encErr != nil {
			out.Close()
			log.Panicf("Error while encoding file: %s -> %s", encErr, out.Name())
			return encErr
		}
	}
	return nil
}

func createPath(path string) {
	_ = os.MkdirAll(path, os.ModePerm)
}
