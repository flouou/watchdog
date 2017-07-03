package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var root = flag.String("root", "", "Directory containing the original images")
var width = flag.Int64("width", 0, "Desired width of the derivate")
var height = flag.Int64("height", 0, "Height of the desired derivate")

func main() {
	flag.Parse()

	err := filepath.Walk(*root, findFiles)
	if err != nil {
		fmt.Printf("Error while walking: %v\n", err)
	}
}

func findFiles(path string, fi os.FileInfo, err error) error {
	//fmt.Printf("Converted: %s\n", path)
	switch mode := fi.Mode(); {
	case mode.IsDir():
	case mode.IsRegular():
		resize(path)
	}

	return nil
}

//converts the image
func resize(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	//fmt.Printf("Found a file: %s\n", file.Name())
	//fmt.Println(canConvert(*file))
	convertible, extension := canConvert(*file)
	if convertible {
		fmt.Printf("Converting %s\n", extension)
	} else {
		fmt.Printf("Conversion not doable: %s\n", extension)
	}
	return nil
}

//checks, if the file is even convertible by filetype
func canConvert(file os.File) (bool, string) {
	extension := filepath.Ext(file.Name())

	if extension == ".jpg" || extension == ".gif" || extension == ".png" {
		//fmt.Printf("%s is a %s", file.Name(), extension)
		return true, extension
	}
	return false, extension
}
