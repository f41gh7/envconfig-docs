package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
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
	outPutFile           = flag.String("output", "", "out put file for variables")
	inputFile            = flag.String("input", "", "input go file with config struct, by default conf.go")
	truncate             = flag.Bool("truncate", true, "truncate variable  value longer then 30 symbols")
)

func main() {
	flag.Parse()

	fset := token.NewFileSet()
	var data interface{}
	if *inputFile == "" {
		fmt.Printf("you must provide filename for doc generation, use --help for more information \n")
		os.Exit(1)
	}
	logDebug("starting generation readme for file %v, writing output to %v", *inputFile, *outPutFile)
	node, err := parser.ParseFile(fset, *inputFile, data, parser.ParseComments)
	if err != nil {
		logDebug("error parsing file %v", err)
		os.Exit(1)
	}
	fileDesc := &strings.Builder{}
	packName := node.Name.Name
	generatedTime := time.Now().UTC().Format(time.UnixDate)
	_, err = fileDesc.Write([]byte(fmt.Sprintf("# Auto Generated vars for package %v \n", packName)))
	if err != nil {
		panic(err)
	}

	fileDesc.Write([]byte(fmt.Sprintf(" updated at %v \n\n\n", generatedTime)))

	fileDesc.Write([]byte(fmt.Sprintf("| varible name | variable default value | variable required | variable description |\n")))
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

	logDebug("finished program!")
	if *outPutFile != "" {
		err := ioutil.WriteFile(*outPutFile, []byte(fileDesc.String()), os.ModeExclusive)
		if err != nil {
			fmt.Printf("cannot write to file: %s, error: %v", *outPutFile, err)
			os.Exit(1)
		}
	}
	_, _ = os.Stdout.Write([]byte(fileDesc.String()))
}

func logDebug(line string, exprs ...interface{}) {
	if *debug != "" {
		errLog.Printf(line, exprs...)
	}
}
