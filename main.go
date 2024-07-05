package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	defaultPrefixName = "prefixVar"
)

var (
	prefixConstName      = flag.String("prefixConst", defaultPrefixName, "prefix constant name at conf.go")
	prefixOverride       = flag.String("prefix", "", "override prefix const value")
	structNeedGenComment = flag.String("matchComment", "//genvars:true", "struct comment line, that must be added to struct for match")
	errLog               = log.New(os.Stderr, "", 0)
	debug                = flag.String("debugs", "", "enables debug mode if not empty, debug will be written to stderr")
	outputFile           = flag.String("output", "", "out put file for variables")
	inputFile            = flag.String("input", "", "input go file with config struct, by default conf.go")
	truncate             = flag.Bool("truncate", true, "truncate variable  value longer then 30 symbols")
	sumRegex             = regexp.MustCompile("\\[envconfig-sum\\]\\:[\\s]*(?P<sum>[^\\s]+)")
)

func main() {
	flag.Parse()

	fset := token.NewFileSet()
	var data interface{}
	if *inputFile == "" {
		fmt.Printf("you must provide filename for doc generation, use --help for more information \n")
		os.Exit(1)
	}
	logDebug("starting generation readme for file %v, writing output to %v", *inputFile, *outputFile)
	node, err := parser.ParseFile(fset, *inputFile, data, parser.ParseComments)
	if err != nil {
		logDebug("error parsing file %v", err)
		os.Exit(1)
	}
	fileDesc := &strings.Builder{}

	packName := node.Name.Name
	_, err = fileDesc.Write([]byte(fmt.Sprintf("# Auto Generated vars for package %v \n", packName)))
	if err != nil {
		panic(err)
	}

	fileDesc.Write([]byte(" updated at %v\n\n\n"))

	fileDesc.Write([]byte(fmt.Sprintf("| variable name | variable default value | variable required | variable description |\n")))
	fileDesc.Write([]byte(fmt.Sprintf("| --- | --- | --- | --- |\n")))

	logDebug("starting parsing file")
	logDebug("start searching prefix for our configuration")

	var prefix string
	if *prefixOverride != "" {
		logDebug("overriding prefix with value: %s", *prefixOverride)
		prefix = *prefixOverride
	} else {
		prefix = extractConstPrefix(node.Decls)

	}
	findStructsAndWalk(node.Decls, prefix, fileDesc)
	hasher := md5.New()
	hasher.Write([]byte(fileDesc.String()))
	newSum := hex.EncodeToString(hasher.Sum(nil))
	if *outputFile != "" {
		outputFileContent, err := os.ReadFile(*outputFile)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				fmt.Printf("cannot read output file: %s, error: %v", *outputFile, err)
				os.Exit(1)
			}
		} else {
			match := sumRegex.FindSubmatch(outputFileContent)
			var foundSum string
			for i, name := range sumRegex.SubexpNames() {
				if i != 0 && name == "sum" {
					foundSum = string(match[i])
				}
			}
			if foundSum == newSum {
				logDebug("file is already up to date!")
				os.Exit(0)
			}
		}
	}

	fileDesc.Write([]byte(fmt.Sprintf("[envconfig-sum]: %s", newSum)))

	logDebug("finished program!")
	generatedTime := time.Now().UTC().Format(time.UnixDate)
	output := fmt.Sprintf(fileDesc.String(), generatedTime)
	if *outputFile != "" {
		err := os.WriteFile(*outputFile, []byte(output), 0o600)
		if err != nil {
			fmt.Printf("cannot write to file: %s, error: %v", *outputFile, err)
			os.Exit(1)
		}
	}
	_, _ = os.Stdout.Write([]byte(output))
}

func logDebug(line string, exprs ...interface{}) {
	if *debug != "" {
		errLog.Printf(line, exprs...)
	}
}
