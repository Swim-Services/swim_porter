package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/swedeachu/swim_porter/port"
)

func main() {
	input := flag.String("i", "", "input pack")
	output := flag.String("o", "", "output path")
	flag.Parse()
	if *input == "" || *output == "" {
		flag.Usage()
		os.Exit(-1)
	}
	dat, err := os.ReadFile(*input)
	if err != nil {
		log.Fatalln(err)
	}
	name := filepath.Base(*input)
	nameNoExt := name[:strings.LastIndex(name, path.Ext(name))]

	out, err := port.Port(dat, nameNoExt)
	if err != nil {
		log.Fatalln(err)
	}
	outFile := (*output) + "/" + nameNoExt + ".mcpack"
	if err := os.WriteFile(outFile, out, 0644); err != nil {
		log.Fatalln(err)
	}
}
