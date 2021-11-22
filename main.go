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

func logError(format string, vars ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", vars...)
}

func log(format string, vars ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, format+"\n", vars...)
}

var isVerbose = false

func main() {
	if len(os.Args) == 1 {
		logError("yml file should be provided")
		return
	}

	serverFlag := flag.String("server", "", "server implementation")
	outputFlag := flag.String("output", "", "output file")
	verboseFlag := flag.Bool("verbose", false, "show additional info")

	flag.Parse()

	isVerbose = *verboseFlag
	specFileName := os.Args[len(os.Args)-1]

	if isVerbose {
		log("Reading spec file %v...", specFileName)
	}

	specFileData, err := ioutil.ReadFile(specFileName)
	if err != nil {
		logError("file not found: %v", err)
		return
	}

	out, err := generate(specFileData, *serverFlag)
	if err != nil {
		logError("$v", err)
		return
	}

	if isVerbose {
		log("Output code...")
	}

	if *outputFlag != "" {
		genDirOutput := filepath.Dir(*outputFlag)
		if err := os.MkdirAll(genDirOutput, 0750); err != nil {
			logError("saving generated code failed: %v", err)
		}
		if err := ioutil.WriteFile(*outputFlag, out, 0755); err != nil {
			logError("saving generated code failed: %v", err)
			return
		}
	} else {
		if _, err := fmt.Fprintf(os.Stdout, "%s", out); err != nil {
			logError("%v", err)
		}
	}
}
