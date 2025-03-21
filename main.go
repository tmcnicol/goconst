package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func run() error {
	var typeNames string
	var out string
	var cmd = &cobra.Command{
		Use:   "goconst",
		Short: "Generate types from consants",
		RunE: func(cmd *cobra.Command, args []string) error {
			return parseConst(cmd, typeNames, out, args)
		},
	}

	cmd.Flags().StringVar(
		&typeNames, "type", "", "comma separated list of type names")
	err := cmd.MarkFlagRequired("type")
	if err != nil {
		return err
	}
	cmd.Flags().StringVar(
		&out, "out", "stdout", "output target can be a filepath or stdout, defaults to stdout")

	return cmd.Execute()
}

func parseConst(_ *cobra.Command, typeNames string, out string, patterns []string) error {
	types := strings.Split(typeNames, ",")
	packages := loadPackages(patterns)
	for _, pkg := range packages {
		for _, typ := range types {
			values := findValues(typ, pkg)
			renderValues(values, out, typ)
		}
	}
	return nil
}

func loadPackages(patterns []string) []*Package {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedFiles,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) == 0 {
		log.Fatalf("error: no packages matching %v", strings.Join(patterns, " "))
	}

	out := make([]*Package, len(pkgs))
	for i, pkg := range pkgs {
		p := &Package{
			name:  pkg.Name,
			files: make([]*File, len(pkg.Syntax)),
		}

		for j, file := range pkg.Syntax {

			p.files[j] = &File{
				file: file,
				pkg:  p,
			}
		}

		out[i] = p
	}
	return out
}

func renderValues(values []Value, out string, typ string) {
	writer, cleanup, err := loadTarget(out)
	if err != nil {
		log.Fatalf("creating output: %v", err)
	}
	defer cleanup()

	renderer := NewRenderer(writer)

	// Capitalise first letter
	r, l := utf8.DecodeRuneInString(typ)
	typeName := fmt.Sprintf("%c%s", unicode.ToTitle(r), typ[l:])

	fields := make([]Field, len(values))
	for i, v := range values {
		fields[i] = Field{
			Name: v.name,
			Doc:  v.doc,
		}
	}

	data := Data{
		UnionName: fmt.Sprintf("%ss", typ),
		TypeName:  typeName,
		Fields:    fields,
	}

	err = renderer.Render(data)
	if err != nil {
		log.Fatalf("writing output: %v", err)
	}
}

func loadTarget(out string) (io.Writer, func(), error) {
	switch out {
	case "stdout":
		return os.Stdout, func() {}, nil
	default:
		return fileTarget(out)
	}
}

func fileTarget(path string) (io.Writer, func(), error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create file %s: %w", path, err)
	}

	// Return the writer and a cleanup function
	cleanup := func() {
		file.Close()
	}

	return file, cleanup, nil
}

func findValues(typeName string, pkg *Package) []Value {
	values := []Value{}
	for _, file := range pkg.files {
		// Set the state for this run of the walker.
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.findSymbols)
			values = append(values, file.values...)
		}
	}
	return values
}

type Package struct {
	name  string
	files []*File
}

type File struct {
	pkg      *Package  // Package to which this file belongs.
	file     *ast.File // Parsed AST.
	comments ast.CommentMap

	// These fields are reset for each type being generated.
	typeName string  // Name of the constant type.
	values   []Value // Accumulator for names of the consts.
}

type Value struct {
	name string
	doc  string
}

// Finds all of the const with the variable name
func (f *File) findSymbols(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.CONST {
		// We only care about const declarations.
		return true
	}

	for _, spec := range decl.Specs {
		vspec := spec.(*ast.ValueSpec) // Guaranteed to succeed as this is CONST.
		ident, ok := vspec.Type.(*ast.Ident)
		if !ok || ident.Name != f.typeName {
			continue
		}

		if len(vspec.Names) <= 0 {
			return true
		}
		name := vspec.Names[0].Name

		lines := []string{}
		for _, c := range vspec.Doc.List {
			text := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
			lines = append(lines, text)
		}
		comment := strings.Join(lines, "\n")

		value := Value{
			name: name,
			doc:  comment,
		}

		f.values = append(f.values, value)
	}
	return false
}
