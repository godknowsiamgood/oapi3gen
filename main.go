package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func templateMap(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func log(format string, vars ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", vars...)
}

func main() {
	if len(os.Args) == 1 {
		log("yml s file should be provided")
		return
	}

	serverFlag := flag.String("server", "", "server implementation")
	outputFlag := flag.String("output", "", "output file")

	flag.Parse()

	specFileName := os.Args[len(os.Args)-1]

	specFileData, err := ioutil.ReadFile(specFileName)
	if err != nil {
		log("fiel not found: %v", err)
		return
	}

	out, err := generate(specFileData, *serverFlag)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	if *outputFlag != "" {
		genDirOutput := filepath.Dir(*outputFlag)
		if err := os.MkdirAll(genDirOutput, 0750); err != nil {
			log("saving generated code failed: %v", err)
		}
		if err := ioutil.WriteFile(*outputFlag, out, 0755); err != nil {
			log("saving generated code failed: %v", err)
			return
		}
	} else {
		if _, err := fmt.Fprintf(os.Stdout, "%s", out); err != nil {
			log("%v", err)
		}
	}
}
