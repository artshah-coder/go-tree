package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

var indent = ""

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

// Check is directory empty or not
func isEmptyDir(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	files, _ := f.ReadDir(0)
	count := len(files)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()

	})
	switch printFiles {
	case false:
		dirs := []os.DirEntry{}
		for _, file := range files {
			if file.IsDir() {
				dirs = append(dirs, file)
			}
		}
		count = len(dirs)
		for i, dir := range dirs {
			output := ""
			if i != count-1 {
				if len(indent) > 0 {
					output += indent + "├───%s\n"
				} else {
					output += "├───%s\n"
				}
				fmt.Fprintf(out, output, dir.Name())

				if isEmpty, _ := isEmptyDir(path + "/" + dir.Name()); !isEmpty {
					indent += "│\t"
					dirTree(out, path+"/"+dir.Name(), printFiles)
					indent, _ = strings.CutSuffix(indent, "│\t")
				}
			} else {
				if len(indent) > 0 {
					output = indent + "└───%s\n"
				} else {
					output = "└───%s\n"
				}
				fmt.Fprintf(out, output, dir.Name())

				if isEmpty, _ := isEmptyDir(path + "/" + dir.Name()); !isEmpty {
					indent += "\t"
					dirTree(out, path+"/"+dir.Name(), printFiles)
					indent, _ = strings.CutSuffix(indent, "\t")
				}
			}
		}
	case true:
		for i, file := range files {
			output := ""
			if i != count-1 {
				if len(indent) > 0 {
					output += indent + "├───%s"
				} else {
					output += "├───%s"
				}
				if file.Type().IsRegular() {
					fileInfo, err := file.Info()
					if err != nil {
						return err
					}
					fileSize := fileInfo.Size()
					if fileSize == 0 {
						output += " (empty)\n"
						fmt.Fprintf(out, output, file.Name())
					} else {
						output += " (%db)\n"
						fmt.Fprintf(out, output, file.Name(), fileSize)
					}
				} else {
					output += "\n"
					fmt.Fprintf(out, output, file.Name())
				}

				if isEmpty, _ := isEmptyDir(path + "/" + file.Name()); !isEmpty {
					indent += "│\t"
					dirTree(out, path+"/"+file.Name(), printFiles)
					indent, _ = strings.CutSuffix(indent, "│\t")
				}
			} else {
				if len(indent) > 0 {
					output = indent + "└───%s"
				} else {
					output = "└───%s"
				}
				if file.Type().IsRegular() {
					fileInfo, err := file.Info()
					if err != nil {
						return err
					}
					fileSize := fileInfo.Size()
					if fileSize == 0 {
						output += " (empty)\n"
						fmt.Fprintf(out, output, file.Name())
					} else {
						output += " (%db)\n"
						fmt.Fprintf(out, output, file.Name(), fileSize)
					}
				} else {
					output += "\n"
					fmt.Fprintf(out, output, file.Name())
				}

				if isEmpty, _ := isEmptyDir(path + "/" + file.Name()); !isEmpty {
					indent += "\t"
					dirTree(out, path+"/"+file.Name(), printFiles)
					indent, _ = strings.CutSuffix(indent, "\t")
				}
			}
		}
	}
	return nil
}
