package main

import (
	"flag"
	"fmt"
	"github.com/ChainSecurity/scilla_static_checker/pkg/ast"
	"github.com/ChainSecurity/scilla_static_checker/pkg/ir"
	"github.com/ChainSecurity/scilla_static_checker/pkg/souffle"
	"io/ioutil"
	"os"
	"path"
)

func main() {

	// Handling CLI flags

	var analysisDir string
	flag.StringVar(&analysisDir, "analysis_dir", "./souffle_analysis", "folder where facts_in and facts_out will be created")

	flag.Parse()

	// Reading json ast
	astJsonPath := flag.Arg(0)
	jsonFile, err := os.Open(astJsonPath)
	fmt.Println(astJsonPath, analysisDir)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	// Parsing ast
	byteValue, _ := ioutil.ReadAll(jsonFile)
	cm, err := ast.Parse_mod(byteValue)
	if err != nil {
		panic(err)
	}

	// Building IR CFG
	b := ir.BuildCFG(cm)

	// Creating souffle output files

	factsOutFolder := path.Join(analysisDir, "facts_out")
	factsInFolder := path.Join(analysisDir, "facts_in")
	fmt.Printf("The facts folders: %s %s\n", factsOutFolder, factsInFolder)
	err = souffle.MakeCleanFolder(factsOutFolder)
	if err != nil {
		panic(err)
	}

	err = souffle.MakeCleanFolder(factsInFolder)
	if err != nil {
		panic(err)
	}

	ir.DumpFacts(b, factsInFolder)

	// Running souffle
	fmt.Println("Starting Souffle analysis")
	souffle.RunSouffle("souffle_analysis/analysis.dl", factsInFolder, factsOutFolder)
	fmt.Println("Souffle analysis finished")

	// Results output
	fmt.Println("======RESULTS======")

	outputFiles := []string{"patternMatch.csv", "patternMatchInfo.csv"}
	for _, f := range outputFiles {
		result, err := souffle.ReadOutput(path.Join(factsOutFolder, f))
		if err != nil {
			panic(err)
		}
		fmt.Println()
		fmt.Println(f)
		for _, fact := range result {
			fmt.Println(fact)
		}
		fmt.Println()
	}
}
