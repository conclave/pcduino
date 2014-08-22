package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const template = `package main

import (
  . "github.com/conclave/pcduino/core"
)

func init() {
  Init()
  setup()
}

func main() {
  for {
    loop()
  }
}

// impl your own setup() and loop()
func setup() {

}

func loop() {

}
`

func main() {
	flag.Parse()
	if flag.NArg() == 0 || flag.Arg(0) == "" {
		fmt.Println("No package name provided.")
		os.Exit(1)
	}
	package_name := flag.Arg(0)
	file_name := "main.go"
	if flag.NArg() > 1 || flag.Arg(1) != "" {
		file_name = flag.Arg(1)
		if !strings.HasSuffix(file_name, ".go") {
			file_name += ".go"
		}
	}
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
	}()
	os.Mkdir(package_name, 0755)
	if err = os.Chdir(package_name); err != nil {
		return
	}
	file, err := os.Create(file_name)
	if err != nil {
		return
	}
	_, err = io.WriteString(file, template)
	file.Close()
	if file, err = os.Create(".gitignore"); err != nil {
		return
	}
	_, err = io.WriteString(file, package_name+"\n")
	return
}
