package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	flagLanguage = flag.String("lang", "go", "language you want to generate")
	rootFolders  = []string{"attack", "discovery", "regex", "wordlists-misc", "wordlists-user-passwd"}
)

func processDir(fpath []string) error {
	fmt.Println("-->", fpath)
	// list files/folders from directory
	dir, err := ioutil.ReadDir(path.Join(fpath...))
	if err != nil {
		return err
	}

	// check if file or folder
	for i := range dir {
		isDir := dir[i].IsDir()
		fname := dir[i].Name()

		if isDir {
			tmpPath := fpath
			tmpPath = append(tmpPath, fname)
			err := processDir(tmpPath)
			if err != nil {
				fmt.Println(err)
				return err
			}
		} else {
			if path.Ext(fname) == ".txt" {
				err = processFile(fpath, fname)
				if err != nil {
					fmt.Println(err)
					return err
				}
			}
		}
	}
	return nil
}

func processFile(fpath []string, fname string) error {
	fext := path.Ext(fname)

	source := path.Join(path.Join(fpath...), fname)
	d, err := ioutil.ReadFile(source)
	if err != nil {
		fmt.Println("Read source error:", err)
		return err
	}
	dLines := strings.Split(string(d), "\n")

	// cleaned up data
	data := make([]string, 0)
	for i := 0; i < len(dLines); i++ {
		// check if last character is a carrigan return
		tmp := dLines[i]
		tmpLength := len(tmp)
		if tmpLength > 2 {
			if string(tmp[tmpLength-1:]) == "\r" {
				dLines[i] = string(tmp[0 : tmpLength-1])
			}
			data = append(data, dLines[i])
		}
	}

	var generatedCode []byte
	fpathLastItem := fpath[len(fpath)-1]
	fnameWithoutExt := fname[0 : len(fname)-len(fext)]
	fnameWithoutExtToUpper := strings.ToUpper(fnameWithoutExt)
	fnameWithoutExtToUpper = strings.Replace(fnameWithoutExtToUpper, "-", "_", 100)

	fmt.Println("-->   generate", fnameWithoutExt)
	switch *flagLanguage {
	case "go":
		generatedCode = generateGo(data, source, fpathLastItem, fnameWithoutExtToUpper)
		break
	case "js":
		generatedCode = generateJs(data, source, fnameWithoutExtToUpper)
		break
	case "c":
	case "h":
	case "cpp":
	case "hpp":
		fmt.Println("c/c++ not supported")
		break
	case "py":
		fmt.Println("python not supported")
		break
	case "rb":
		fmt.Println("ruby not supported")
		break
	default:
		fmt.Printf("language %q not supported\n", fname)
	}

	perm := os.FileMode(0666)

	target := path.Join(path.Join(fpath...), fnameWithoutExt) + "." + *flagLanguage
	err = ioutil.WriteFile(target, generatedCode, perm)
	if err != nil {
		return err
	}
	return nil

	// targetTest := path.Join(path.Join(fpath...), fnameWithoutExt) + "_test." + *flagLanguage
	// generatedTestCode := generateGoTest(fpathLastItem, []string{"fofof", "dsd"})
	// return ioutil.WriteFile(targetTest, generatedTestCode, perm)
}

func main() {
	flag.Parse()
	fmt.Println("==> fuzzdb code generator")
	fmt.Println("==> language:", *flagLanguage)
	for i := range rootFolders {
		err := processDir([]string{rootFolders[i]})
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(127)
		}
	}
}

func header(source string, totalData int) []string {
	s := []string{
		"// Generated sourcecode, DO NOT EDIT!",
		"//",
		"// source: " + source,
		"// date: " + time.Now().String(),
		"// length of the string array: " + strconv.Itoa(totalData),
	}
	return s
}

func generateGo(data []string, source, pkgname, varname string) []byte {
	s := header(source, len(data))
	s = append(s, "")
	s = append(s, "package "+pkgname+"\n")
	s = append(s, "var "+varname+" = []string{")
	for i := 0; i < len(data); i++ {
		s = append(s, fmt.Sprintf("  %q,", data[i]))
	}
	s = append(s, "}")
	return []byte(strings.Join(s, "\n"))
}

// func generateGoTest(pkgUrl string, varNames []string) []byte {
// 	s := make([]byte, 0)
// 	s = append(s, []byte("package "+pkgUrl+"\n\n")...)
// 	s = append(s, []byte("import \"testing\"")...)
// 	s = append(s, []byte("func Test_y(t *testing.T) {\n")...)
// 	s = append(s, []byte("\t\n")...)
// 	s = append(s, []byte("\tt.Error(123)\n")...)
// 	s = append(s, []byte("}\n")...)
// 	return s
// }

func generateJs(data []string, source, varname string) []byte {
	s := header(source, len(data))
	s = append(s, "")
	s = append(s, "var "+varname+" = [")
	for i := 0; i < len(data); i++ {
		s = append(s, fmt.Sprintf("  %q,", data[i]))
	}
	s = append(s, "]\n")
	s = append(s, "module.exports = "+varname+"\n")
	return []byte(strings.Join(s, "\n"))
}

func generateJson(data []string, varname string) {
	// TODO:
}
