package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// u can use path/filepath, io/ioutil, strings, io, os
const (
	line       = "│"
	lineMiddle = "├───"
	lineEnd    = "└───"
)

func main() {
	// out := os.Stdout
	// if !(len(os.Args) == 2 || len(os.Args) == 3) {
	// 	panic("usage go run main.go . [-f]")
	// }
	// path := os.Args[1]
	// printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	// err := dirTree(out, path, printFiles)
	// if err != nil {
	// 	panic(err.Error())
	// }
	var doPrintFiles bool
	flag.BoolVar(&doPrintFiles, "f", false, "whether to print file size info")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("usage go run main.go . [-f]")
	}

	if err := dirTree(os.Stdout, args[0], doPrintFiles); err != nil {
		log.Fatal(err)
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var counter int
	var prefix string
	dirTreeRec(out, path, counter, prefix, printFiles)
	return nil
}

func dirTreeRec(out io.Writer, path string, counter int, prefix string, printFiles bool) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	counter++

	onlyDir := []fs.DirEntry{}
	if !printFiles {
		for _, file := range files {
			if file.IsDir() {
				onlyDir = append(onlyDir, file)
			}
		}
		files = onlyDir
	}

	for i, file := range files {
		prefixNew := prefix
		if i != len(files)-1 {
			if file.IsDir() {
				prefixNew = lineGenerate(counter, false, prefix)
				fmt.Fprintln(out, prefixNew+lineMiddle+file.Name())
				prefixNew = lineGenerate(counter, true, prefix)
				dirTreeRec(out, filepath.Join(path, file.Name()), counter, prefixNew, printFiles)
			} else {
				prefixNew = lineGenerate(counter, false, prefixNew)
				fmt.Fprint(out, prefixNew+lineMiddle+file.Name())
				countBytes(filepath.Join(path, file.Name()), out)
				fmt.Fprintln(out)
			}
		} else {
			if file.IsDir() {
				prefixNew = lineGenerate(counter, false, prefix)
				fmt.Fprintln(out, prefixNew+lineEnd+file.Name())
				dirTreeRec(out, filepath.Join(path, file.Name()), counter, prefixNew, printFiles)
			} else {
				prefixNew = lineGenerate(counter, false, prefix)
				fmt.Fprint(out, prefixNew+lineEnd+file.Name())
				countBytes(filepath.Join(path, file.Name()), out)
				fmt.Fprintln(out)
			}
		}
	}
}

func countBytes(path string, out io.Writer) {
	fs, _ := os.Stat(path)
	bytes := fs.Size()
	if bytes == 0 {
		fmt.Fprint(out, " (empty)")
	} else {
		fmt.Fprintf(out, " (%vb)", bytes)
	}
}

func lineGenerate(counter int, chekcLast bool, prefix string) string {
	if counter == 1 {
		if chekcLast {
			prefix += line
		}
		return prefix
	}
	prefix += "\t"
	if chekcLast {
		prefix += line
	}
	return prefix
}
