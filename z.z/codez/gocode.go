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

func parseStmts(log logz.Logger, code string) ([]ast.Stmt, error) {
	fset := token.NewFileSet()
	src := fmt.Sprintf("package 𝐳\nfunc 𝐳() {\n%v\n}", code)
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

func parseStmt(log logz.Logger, code string) (ast.Stmt, error) {
	stmts, err := parseStmts(log, code)
	if err != nil {
		return nil, err
	}
	if len(stmts) == 0 {
		return nil, errorz.New("no statement")
	}
	return stmts[0], err
}

func parseDecls(log logz.Logger, code string) ([]ast.Decl, error) {
	fset := token.NewFileSet()
	src := fmt.Sprintf("package 𝐳\n%v\n", code)
	file, err := parseSrc(fset, src)
	if err != nil {
		return nil, err
	}
	if log.Enabled(logz.LevelDebug) {
		printAst("parseDecl", fset, file)
	}
	return file.Decls, nil
}

func parseDecl(log logz.Logger, code string) (ast.Decl, error) {
	decls, err := parseDecls(log, code)
	if err != nil {
		return nil, err
	}
	if len(decls) == 0 {
		return nil, errorz.New("no declaration")
	}
	return decls[0], nil
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
