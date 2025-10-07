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
	"regexp"
	"sort"
	"strings"
	"unicode"
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
		"plugins.Secret": "../secret.md",
		"Secret":         "secret.md",
		"plugins.TLS":    "../tls.md",
	}

	kubernetes_link_templates = map[string]string{
		"corev1.": "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#%s-v1-core",
		"metav1.": "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#%s-v1-meta",
		"rbacv1.": "https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#%s-v1-rbac-authorization-k8s-io",
	}

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
	pluginsLocations := []DocumentsLocation{
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

	crdsLocations := []DocumentsLocation{
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
				var flag = true
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

			types := ParseDocumentationFrom(src, dl.name, true)

			for _, t := range types {
				strukt := t[0]
				if len(t) > 1 {
					buffer.WriteString(fmt.Sprintf("# %s\n\n%s\n\n\n", strukt.Name, strukt.Doc))

					buffer.WriteString("| Field | Description | Scheme |\n")
					buffer.WriteString("| ----- | ----------- | ------ |\n")
					fields := t[1:]
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

			_, _ = f.WriteString(buffer.String())
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
			types = append(types, ParseDocumentationFrom(src, dl.name, false)...)
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
				fields := t[1:]
				for _, f := range fields {
					buffer.WriteString(fmt.Sprintf("| %s | %s | %s |\n", f.Name, f.Doc, f.Type))
				}
				buffer.WriteString("\n")
				buffer.WriteString("[Back to TOC](#table-of-contents)\n")
			}
		}

		f, _ := os.Create(fmt.Sprintf("./docs/%s.md", dl.name))
		_, _ = f.WriteString(fmt.Sprintf(firstParagraph, dl.name) + buffer.String())
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
	name = strings.ReplaceAll(name, " ", "-")
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

func ParseDocumentationFrom(src string, dl_name string, shouldSort bool) []KubeTypes {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	pkg, err := doc.NewFromFiles(fset, []*ast.File{f}, "")
	if err != nil {
		fmt.Println(err)
		return nil
	}

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

	docForTypes := make([]KubeTypes, 0, len(types))
	for _, kubType := range types {
		structType, _ := kubType.Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType)

		var ks KubeTypes
		ks = append(ks, Pair{kubType.Name, fmtRawDoc(kubType.Doc), ""})

		for _, field := range structType.Fields.List {
			typeString := fieldType(field.Type, dl_name)
			if n := fieldName(field); n != "-" {
				fieldDoc := fmtRawDoc(field.Doc.Text())
				ks = append(ks, Pair{n, fieldDoc, typeString})
			}
		}
		docForTypes = append(docForTypes, ks)
	}

	return docForTypes
}

func fmtRawDoc(rawDoc string) string {
	var buffer bytes.Buffer
	delPrevChar := func() {
		if buffer.Len() > 0 {
			buffer.Truncate(buffer.Len() - 1) // Delete the last " " or "\n"
		}
	}

	for line := range strings.SplitSeq(rawDoc, "\n") {
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
	postDoc = strings.ReplaceAll(postDoc, "\\\"", "\"") // replace user's \" to "
	postDoc = strings.ReplaceAll(postDoc, "\"", "\\\"") // Escape "
	postDoc = strings.ReplaceAll(postDoc, "\n", "\\n")
	postDoc = strings.ReplaceAll(postDoc, "\t", "\\t")
	postDoc = strings.ReplaceAll(postDoc, "|", "\\|")

	return postDoc
}

func tryKubernetesLink(typeName string) (string, bool) {
	for prefix, link_template := range kubernetes_link_templates {
		if strings.HasPrefix(typeName, prefix) {
			typeName = strings.ToLower(strings.TrimPrefix(typeName, prefix))
			return fmt.Sprintf(link_template, typeName), true
		}
	}
	return "", false
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func toLink(typeName string, documentLocationName string) string {
	primitiveTypes := map[string]bool{
		"string":            true,
		"bool":              true,
		"int32":             true,
		"int64":             true,
		"map[string]string": true,
	}

	if _, ok := primitiveTypes[typeName]; ok {
		return typeName
	}

	link, hasLink := links[typeName]
	if hasLink {
		return wrapInLink(typeName, link)
	}

	if strings.Contains(typeName, "input.") ||
		strings.Contains(typeName, "output.") ||
		strings.Contains(typeName, "filter.") ||
		strings.Contains(typeName, "parser.") ||
		strings.Contains(typeName, "custom.") {
		// Eg. *output.Elasticsearch => ../plugins/output/elasticsearch.md
		pluginType, pluginName, _ := strings.Cut(typeName, ".")
		link := fmt.Sprintf("plugins/%s/%s/%s.md", documentLocationName, pluginType, ToSnakeCase(pluginName))
		return wrapInLink(typeName, link)
	}

	k8sLink, hasK8sLink := tryKubernetesLink(typeName)
	if hasK8sLink {
		return wrapInLink(typeName, k8sLink)
	}

	if !strings.Contains(typeName, ".") && len(typeName) > 0 && unicode.IsUpper([]rune(typeName)[0]) {
		link := fmt.Sprintf("#%s", strings.ToLower(typeName))
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
		jsonTag = reflect.StructTag(strings.Trim(field.Tag.Value, "`")).Get("json")
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

func fieldType(typ ast.Expr, dl_name string) string {
	switch typ := typ.(type) {
	case *ast.Ident:
		return toLink(typ.Name, dl_name)
	case *ast.StarExpr:
		return "*" + fieldType(typ.X, dl_name)
	case *ast.SelectorExpr:
		pkg := typ.X.(*ast.Ident)
		t := typ.Sel
		return toLink(pkg.Name+"."+t.Name, dl_name)
	case *ast.ArrayType:
		return "[]" + fieldType(typ.Elt, dl_name)
	case *ast.MapType:
		return "map[" + toLink(fieldType(typ.Key, dl_name), dl_name) + "]" + toLink(fieldType(typ.Value, dl_name), dl_name)
	default:
		return ""
	}
}
