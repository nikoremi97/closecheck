package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/packages"
)

var (
	// Analyzer defines the analyzer for check-close
	Analyzer = &analysis.Analyzer{
		Name:     "checkclose",
		Doc:      "check for any closable return value",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		//FactTypes: []analysis.Fact{new(closer)},
	}

	closerType *types.Interface
)

func run(pass *analysis.Pass) (interface{}, error) {
	visitor := &Visitor{
		pass: pass,
	}

	for _, f := range pass.Files {
		ast.Inspect(f, visitor.Do)
	}

	return nil, nil
}

// init finds the io.Closer interface
func init() {
	cfg := &packages.Config{Mode: packages.NeedTypes, Tests: false}

	pkgs, err := packages.Load(cfg, "io")
	if err != nil {
		panic(err)
	}

	if len(pkgs) != 1 {
		panic("couldn't load io package")
	}

	closerType = pkgs[0].Types.Scope().Lookup("Closer").Type().Underlying().(*types.Interface)
	if closerType == nil {
		panic("io.Closer not found")
	}
}
