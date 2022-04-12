package main

import (
	"errors"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/wii-tools/GoNUSD"
	"github.com/wii-tools/lz11"
	. "github.com/wii-tools/powerpc"
	"github.com/wii-tools/wadlib"
	"io/fs"
	"log"
	"os"
)

// mainDol holds the main DOL - our content at index 1.
var mainDol []byte

// filePresent returns whether the specified path is present on disk.
func filePresent(path string) bool {
	_, err := os.Stat(path)
	return errors.Is(err, fs.ErrNotExist) == false
}

// createDir creates a directory at the given path if it is not already present.
func createDir(path string) {
	if !filePresent(path) {
		os.Mkdir(path, 0755)
	}
}

func main() {
	fmt.Println("=================================")
	fmt.Println("=       Wii-no-Ma Patcher       =")
	fmt.Println("=================================")

	// Create directories we may need later.
	createDir("./output")
	createDir("./cache")

	var originalWad *wadlib.WAD
	var err error

	// Determine whether Wii no Ma is present on disk.
	if !filePresent("./cache/original.wad") {
		log.Println("Downloading a copy of the original Wii no Ma, please wait...")
		originalWad, err = GoNUSD.Download(0x00010001_4843494a, 1025, true)
		check(err)

		// Cache this downloaded WAD to disk.
		contents, err := originalWad.GetWAD(wadlib.WADTypeCommon)
		check(err)

		os.WriteFile("./cache/original.wad", contents, 0755)
	} else {
		originalWad, err = wadlib.LoadWADFromFile("./cache/original.wad")
		check(err)
	}

	// Load main DOL
	mainDol, err = originalWad.GetContent(1)
	check(err)

	// The DOL is lz11 compressed. Decompress it.
	mainDol, err = lz11.Decompress(mainDol)
	check(err)

	// Apply all DOL patches
	fmt.Println(aurora.Green("Applying DOL patches..."))
	mainDol, err = ApplyPatchSet(DetermineLanguageCodePatch, mainDol)
	check(err)

	// Save main DOL
	err = originalWad.UpdateContent(1, mainDol)
	check(err)

	// Generate a patched WAD with our changes
	output, err := originalWad.GetWAD(wadlib.WADTypeCommon)
	check(err)

	fmt.Println(aurora.Green("Done! Install ./output/patched.wad, sit back, and enjoy."))
	writeOut("patched.wad", output)
}

// check well checks if there is an error
func check(err error) {
	if err != nil {
		panic(any(err))
	}
}

// writeOut writes a file with the given name and contents to the output folder.
func writeOut(filename string, contents []byte) {
	os.WriteFile("./output/"+filename, contents, 0755)
}
