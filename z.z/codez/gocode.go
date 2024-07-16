package codez

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"ezpkg.io/colorz"
	"ezpkg.io/errorz"
	"ezpkg.io/logz"
)

func parseSearch(code string) {
	panic("TODO")
}

func parseExpr(x string) (ast.Expr, error) {
	return parser.ParseExpr(x)
}

func parseStmts(log logz.Logger, x string) ([]ast.Stmt, error) {
	fset := token.NewFileSet()
	src := fmt.Sprintf(`
package x
func x() {
	%v
}`, x)
	file, err := parseSrc(fset, src)
	if err != nil {
		return nil, err
	}
	if log.Enabled(logz.LevelDebug) {
		printAst("parseStmts", fset, file)
	}
	bodyStmt := file.Decls[0].(*ast.FuncDecl).Body
	return bodyStmt.List, nil
}

func parseStmt(log logz.Logger, x string) (ast.Stmt, error) {
	stmts, err := parseStmts(log, x)
	if err != nil {
		return nil, err
	}
	return stmts[0], err
}

func parseDecl(x string) (ast.Decl, error) {
	panic("TODO")
}

func parseSrc(fset *token.FileSet, src string) (*ast.File, error) {
	fset.AddFile("", fset.Base(), len(src))
	return parser.ParseFile(fset, "", []byte(src), 0)
}

func printAst(msg string, fset *token.FileSet, x any) {
	if fset == nil {
		fset = token.NewFileSet()
	}
	fmt.Printf("\n%v--- %v ---%v\n", colorz.Yellow, msg, colorz.Reset)
	errorz.MustZ(ast.Print(fset, x))
}
