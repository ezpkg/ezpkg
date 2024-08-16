package codez

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"

	"ezpkg.io/colorz"
	"ezpkg.io/errorz"
	"ezpkg.io/logz"
)

type parseOutput struct {
	file  *ast.File
	ident *ast.Ident
	expr  ast.Expr
	stmt  ast.Stmt
	decl  ast.Decl
	stmts []ast.Stmt
	decls []ast.Decl
}

func (p parseOutput) IsEmpty() bool {
	return p.ident == nil && p.expr == nil &&
		p.stmt == nil && p.decl == nil &&
		len(p.stmts) == 0 && len(p.decls) == 0
}

func parseSearch(log logz.Logger, code string) (output parseOutput, _ error) {
	maybeKind := detectCode(code)
	switch maybeKind {
	case zExpr:
		expr, err := parseExpr(log, code)
		if err != nil {
			return output, err
		}
		output.expr = expr
		if ident, ok := expr.(*ast.Ident); ok {
			output.ident = ident
		}
		return output, nil

	case zStmt:
		stmts, err := parseStmts(log, code)
		if err != nil {
			return output, err
		}
		output.stmts = stmts
		if len(stmts) == 1 {
			output.stmt = stmts[0]
			if decl, ok := output.stmt.(*ast.DeclStmt); ok {
				output.decl = decl.Decl
			}
		}
		return output, nil

	case zDecl:
		decls, err := parseDecls(log, code)
		if err != nil {
			return output, err
		}
		output.decls = decls
		if len(decls) == 1 {
			output.decl = decls[0]
		}
		return output, nil

	case zFile:
		file, err := parseSrc(token.NewFileSet(), code)
		if err != nil {
			return output, err
		}
		output.file = file
		return output, nil

	default:
		return output, nil // empty
	}
}

func detectCode(code string) (maybe zKind) {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(code))
	sc := scanner.Scanner{}
	sc.Init(file, []byte(code), nil, 0)

	nextToken := func() token.Token {
		for {
			_, tok, _ := sc.Scan()
			switch tok {
			case token.EOF:
				return token.EOF
			case token.COMMENT:
				continue
			case token.ILLEGAL:
				return token.ILLEGAL
			default:
				return tok
			}
		}
	}
	tok := nextToken()
	switch {
	case tok == token.EOF:
		return 0
	case tok == token.PACKAGE:
		return zFile
	case stmtStart[tok]:
		return zStmt
	case declStart[tok]:
		return zDecl
	default:
		return zExpr
	}
}

func parseExpr(log logz.Logger, code string) (ast.Expr, error) {
	expr, err := parser.ParseExpr(code)
	if err != nil {
		return nil, err
	}
	if log.Enabled(context.Background(), logz.LevelDebug) {
		printAst("parseExpr", nil, expr)
	}
	return expr, nil
}

func parseStmts(log logz.Logger, code string) ([]ast.Stmt, error) {
	fset := token.NewFileSet()
	src := fmt.Sprintf("package 𝐳\nfunc 𝐳() {\n%v\n}", code)
	file, err := parseSrc(fset, src)
	if err != nil {
		return nil, err
	}
	if log.Enabled(context.Background(), logz.LevelDebug) {
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
	if log.Enabled(context.Background(), logz.LevelDebug) {
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
	if fset == nil {
		fset = token.NewFileSet()
	}
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
