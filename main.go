package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	root := flag.Arg(0)
	err := filepath.Walk(root, findFiles)
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

func convert(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	fmt.Printf("Found a file: %s\n", file.Name())
	fmt.Println(canConvert(*file))
	convertible, convertibleFile := canConvert(*file)
	if convertible {
		fmt.Printf("Converting %s\n", convertibleFile)
	} else {
		fmt.Printf("Conversion not doable: %s\n", convertibleFile)
	}
	return err
}

//checks, if the file is even convertible by filetype
func canConvert(file os.File) (bool, string) {
	switch extension := filepath.Ext(file.Name()); {
	case extension == "jpg":
		fmt.Printf("%s is a %s", file.Name(), extension)
		return true, extension
	case extension == "gif":
		fmt.Printf("%s is a %s", file.Name(), extension)
		return true, extension
	}
	return false, file.Name()
}
