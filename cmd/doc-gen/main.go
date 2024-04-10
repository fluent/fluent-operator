package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

const (
	firstParagraph = `# API Docs
This Document documents the types introduced by the %s Operator.
> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.`

	fluentbitPluginPath = "apis/fluentbit/v1alpha2/plugins/"
	fluentdPluginPath   = "apis/fluentd/v1alpha1/plugins/"
	fluentbitCrdsPath   = "apis/fluentbit/v1alpha2/"
	fluentdCrdsPath     = "apis/fluentd/v1alpha1/"
)

var (
	links = map[string]string{
		"metav1.ObjectMeta":        "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta",
		"metav1.ListMeta":          "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta",
		"metav1.LabelSelector":     "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta",
		"corev1.SecretKeySelector": "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#secretkeyselector-v1-core",
		"corev1.Toleration":        "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#toleration-v1-core",
		"corev1.VolumeSource":      "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#volume-v1-core",
		"plugins.Secret":           "../secret.md",
		"Secret":                   "secret.md",
		"plugins.TLS":              "../tls.md",
	}

	selfLinks = map[string]string{}

	unincludedKeyWords = []string{
		"deepcopy.go",
		"params",
		"test",
		"interface",
	}
)

type DocumentsLocation struct {
	path string
	name string
}

// Inspired by coreos/prometheus-operator: https://github.com/coreos/prometheus-operator
func main() {
	var pluginsLocations = []DocumentsLocation{
		{
			path: fluentbitPluginPath,
			name: "fluentbit",
		},
		{
			path: fluentdPluginPath,
			name: "fluentd",
		},
	}
	plugins(pluginsLocations)

	var crdsLocations = []DocumentsLocation{
		{
			path: fluentbitCrdsPath,
			name: "fluentbit",
		},
		{
			path: fluentdCrdsPath,
			name: "fluentd",
		},
	}
	crds(crdsLocations)
}

func plugins(docsLocations []DocumentsLocation) {
	for _, dl := range docsLocations {
		var srcs []string
		err := filepath.Walk(dl.path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".go") {
				var flag bool = true
				for _, keyword := range unincludedKeyWords {
					flag = flag && !strings.Contains(path, keyword)
				}
				if flag {
					srcs = append(srcs, path)
				}
			}
			return nil
		})
		if err != nil {
			panic(err)
		}

		for _, src := range srcs {
			var buffer bytes.Buffer

			types := ParseDocumentationFrom(src, true)

			for _, t := range types {
				strukt := t[0]
				if len(t) > 1 {
					buffer.WriteString(fmt.Sprintf("# %s\n\n%s\n\n\n", strukt.Name, strukt.Doc))

					buffer.WriteString("| Field | Description | Scheme |\n")
					buffer.WriteString("| ----- | ----------- | ------ |\n")
					fields := t[1:(len(t))]
					for _, f := range fields {
						buffer.WriteString(fmt.Sprintf("| %s | %s | %s |\n", f.Name, f.Doc, f.Type))
					}
					buffer.WriteString("")
				}
			}

			src_name := strings.TrimPrefix(src, dl.path)
			if strings.HasSuffix(src_name, "_types.go") {
				src_name = strings.TrimSuffix(src_name, "_types.go")
			} else {
				src_name = strings.TrimSuffix(src_name, ".go")
			}

			dst := fmt.Sprintf("./docs/plugins/%s/%s.md", dl.name, src_name)

			if err := genDocDirs(dst); err != nil {
				fmt.Printf("Error while generating documentation directories: %s\n", err.Error())
			}

			f, err := os.Create(dst)
			if err != nil {
				fmt.Printf("Error while generating documentation: %s\n", err.Error())
			}

			f.WriteString(buffer.String())
		}
	}
}

func crds(docsLocations []DocumentsLocation) {
	for _, dl := range docsLocations {
		var srcs []string
		err := filepath.Walk(dl.path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !strings.Contains(path, "/plugins") && strings.HasSuffix(path, "_types.go") {
				srcs = append(srcs, path)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}

		var buffer bytes.Buffer
		var types []KubeTypes
		for _, src := range srcs {
			types = append(types, ParseDocumentationFrom(src, false)...)
		}

		sort.Slice(types, func(i, j int) bool {
			return interface{}(types[i]).(KubeTypes)[0].Name < interface{}(types[j]).(KubeTypes)[0].Name
		})

		buffer = printTOC(types)

		for _, t := range types {
			strukt := t[0]
			if len(t) > 1 {
				buffer.WriteString(fmt.Sprintf("# %s\n\n%s\n\n\n", strukt.Name, strukt.Doc))

				buffer.WriteString("| Field | Description | Scheme |\n")
				buffer.WriteString("| ----- | ----------- | ------ |\n")
				fields := t[1:(len(t))]
				for _, f := range fields {
					buffer.WriteString(fmt.Sprintf("| %s | %s | %s |\n", f.Name, f.Doc, f.Type))
				}
				buffer.WriteString("\n")
				buffer.WriteString("[Back to TOC](#table-of-contents)\n")
			}
		}
		
		f, _ := os.Create(fmt.Sprintf("./docs/%s.md", dl.name))
		f.WriteString(fmt.Sprintf(firstParagraph, dl.name) + buffer.String())
	}
}

func genDocDirs(docPath string) error {
	dirPath := filepath.Dir(docPath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, os.ModePerm)
	}
	return nil
}

func toSectionLink(name string) string {
	name = strings.ToLower(name)
	name = strings.Replace(name, " ", "-", -1)
	return name
}

func printTOC(types []KubeTypes) bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString("\n## Table of Contents\n")
	for _, t := range types {
		strukt := t[0]
		if len(t) > 1 {
			buffer.WriteString(fmt.Sprintf("* [%s](#%s)\n", strukt.Name, toSectionLink(strukt.Name)))
		}
	}
	return buffer
}

// Pair of strings. We need the name of fields and the doc
type Pair struct {
	Name, Doc, Type string
}

// KubeTypes is an array to represent all available types in a parsed file. [0] is for the type itself
type KubeTypes []Pair

func ParseDocumentationFrom(src string, shouldSort bool) []KubeTypes {
	var docForTypes []KubeTypes

	pkg := astFrom(src)

	var types []*doc.Type
	for _, kubType := range pkg.Types {
		if _, ok := kubType.Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType); ok {
			types = append(types, kubType)
		}
	}

	if shouldSort {
		sort.Slice(types, func(i, j int) bool {
			return types[i].Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Pos() <
				types[j].Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Pos()
		})
	}

	for _, kubType := range types {
		structType, _ := kubType.Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType)

		var ks KubeTypes
		ks = append(ks, Pair{kubType.Name, fmtRawDoc(kubType.Doc), ""})

		for _, field := range structType.Fields.List {
			typeString := fieldType(field.Type)
			if n := fieldName(field); n != "-" {
				fieldDoc := fmtRawDoc(field.Doc.Text())
				ks = append(ks, Pair{n, fieldDoc, typeString})
			}
		}
		docForTypes = append(docForTypes, ks)
	}

	return docForTypes
}

func astFrom(filePath string) *doc.Package {
	fset := token.NewFileSet()
	m := make(map[string]*ast.File)

	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	m[filePath] = f
	apkg, _ := ast.NewPackage(fset, m, nil, nil)

	return doc.New(apkg, "", 0)
}

func fmtRawDoc(rawDoc string) string {
	var buffer bytes.Buffer
	delPrevChar := func() {
		if buffer.Len() > 0 {
			buffer.Truncate(buffer.Len() - 1) // Delete the last " " or "\n"
		}
	}

	for _, line := range strings.Split(rawDoc, "\n") {
		line = strings.TrimRight(line, " ")
		leading := strings.TrimLeft(line, " ")
		switch {
		case len(line) == 0: // Keep paragraphs
			delPrevChar()
			buffer.WriteString("\n\n")
		case strings.HasPrefix(leading, "+"): // Ignore instructions to go2idl
		default:
			line += " "
			buffer.WriteString(line)
		}
	}

	postDoc := strings.TrimRight(buffer.String(), "\n")
	postDoc = strings.Replace(postDoc, "\\\"", "\"", -1) // replace user's \" to "
	postDoc = strings.Replace(postDoc, "\"", "\\\"", -1) // Escape "
	postDoc = strings.Replace(postDoc, "\n", "\\n", -1)
	postDoc = strings.Replace(postDoc, "\t", "\\t", -1)
	postDoc = strings.Replace(postDoc, "|", "\\|", -1)

	return postDoc
}

func toLink(typeName string) string {
	if strings.Contains(typeName, "input.") ||
		strings.Contains(typeName, "output.") ||
		strings.Contains(typeName, "filter.") ||
		strings.Contains(typeName, "parser.") {
		// Eg. *output.Elasticsearch => ../plugins/output/elasticsearch.md
		link := fmt.Sprintf("plugins/%s.md", strings.ReplaceAll(strings.ToLower(typeName), ".", "/"))
		return wrapInLink(typeName, link)
	}

	selfLink, hasSelfLink := selfLinks[typeName]
	if hasSelfLink {
		return wrapInLink(typeName, selfLink)
	}

	link, hasLink := links[typeName]
	if hasLink {
		return wrapInLink(typeName, link)
	}

	return typeName
}

func wrapInLink(text, link string) string {
	return fmt.Sprintf("[%s](%s)", text, link)
}

// fieldName returns the name of the field as it should appear in JSON format
// "-" indicates that this field is not part of the JSON representation
func fieldName(field *ast.Field) string {
	jsonTag := ""
	if field.Tag != nil {
		jsonTag = reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]).Get("json") // Delete first and last quotation
		if strings.Contains(jsonTag, "inline") {
			return "-"
		}
	}

	jsonTag = strings.Split(jsonTag, ",")[0] // This can return "-"
	if jsonTag == "" {
		if field.Names != nil {
			return field.Names[0].Name
		}
		return field.Type.(*ast.Ident).Name
	}
	return jsonTag
}

func fieldType(typ ast.Expr) string {
	switch typ.(type) {
	case *ast.Ident:
		return toLink(typ.(*ast.Ident).Name)
	case *ast.StarExpr:
		return "*" + fieldType(typ.(*ast.StarExpr).X)
	case *ast.SelectorExpr:
		e := typ.(*ast.SelectorExpr)
		pkg := e.X.(*ast.Ident)
		t := e.Sel
		return toLink(pkg.Name + "." + t.Name)
	case *ast.ArrayType:
		return "[]" + toLink(fieldType(typ.(*ast.ArrayType).Elt))
	case *ast.MapType:
		mapType := typ.(*ast.MapType)
		return "map[" + toLink(fieldType(mapType.Key)) + "]" + toLink(fieldType(mapType.Value))
	default:
		return ""
	}
}
