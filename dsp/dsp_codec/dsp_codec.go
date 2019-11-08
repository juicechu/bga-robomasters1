package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"git.bug-br.org.br/bga/robomasters1/dsp"
)

var (
	creator = flag.String("creator", "Anonymous", "program creator")
	title   = flag.String("title", "Untitled", "program title")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s filename\n",
			os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(),
			"filename must be end in .dsp or .py\n\n")
	}

	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(2)
	}

	fileName := flag.Arg(0)

	isDsp := strings.HasSuffix(strings.ToLower(fileName), ".dsp")
	isPy := strings.HasSuffix(strings.ToLower(fileName), ".py")

	if !isDsp && !isPy {
		flag.Usage()
		os.Exit(2)
	}

	extension := filepath.Ext(fileName)
	baseFileName := strings.TrimSuffix(fileName, extension)

	var err error
	if isDsp {
		err = writePythonFile(baseFileName, extension)
	} else {
		err = writeDspFile(baseFileName, extension)
	}

	if err != nil {
		panic(err)
	}
}

func writePythonFile(baseFileName, extension string) error {
	f, err := dsp.Load(baseFileName + extension)
	if err != nil {
		return err
	}

	f.Dump()

	fd, err := os.OpenFile(baseFileName+".py", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.WriteString(f.PythonCode())
	if err != nil {
		return err
	}

	return nil
}

func writeDspFile(baseFileName, extension string) error {
	fd, err := os.Open(baseFileName + extension)
	if err != nil {
		return err
	}
	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return err
	}

	f, err := dsp.NewWithPythonCode(*creator, *title, string(data))
	if err != nil {
		return err
	}

	err = f.Save(baseFileName + ".dsp")
	if err != nil {
		return err
	}

	return nil
}
