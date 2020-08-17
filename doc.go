package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"reflect"
	"regexp"
	"strings"
)


var (
	gatherRegexp = regexp.MustCompile("([^A-Z]+|[A-Z][^A-Z]+|[A-Z]+)")
)

// walks over struct, extracts tags and comments and write it to writer
func walkStructPrefix(prefix string, currStruct *ast.StructType, fileDesc io.Writer) {
			logDebug("processing struct %v", currStruct)
			for _, field := range currStruct.Fields.List {
				switch  reflect.TypeOf(field.Type).String() {
				case "*ast.StructType":
					logDebug("its struct: %s, walk deeper",field.Names[0].Name)
					walkStructPrefix(prefix+strings.ToUpper(field.Names[0].Name), field.Type.(*ast.StructType), fileDesc)
				case "*ast.Ident":
					logDebug("its ident %s, %s \n", field.Names[0].Name, field.Type.(*ast.Ident).Name)
				case "*ast.SelectorExpr":
					logDebug("ast exr selector: var: %s",field.Names[0].Name)
				case "*ast.ArrayType":
					logDebug("its array type : %s",field.Names[0].Name)
				case "*ast.MapType":
					logDebug("its map type: %s",field.Names[0].Name)
				default:
					logDebug("error, unsupported kind of field %s", reflect.TypeOf(field.Type).String())
					continue
				}
				// we skip fields without tags
				if field.Tag == nil {
					continue
				}
				// extract tags
				tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])


				fieldDefaultValue, fieldDescription, fieldRequired := "-","-","false"

				if tag.Get("default") != "" {
					fieldDefaultValue = tag.Get("default")
					logDebug("found default value of %s with value %s |\n", prefix+strings.ToUpper(field.Names[0].Name), fieldDefaultValue)

				}
				if tag.Get("required") != "" {
					fieldRequired = tag.Get("required")
				}
				if tag.Get("description") != "" {
					fieldDescription = tag.Get("description")
				}

				if field.Doc != nil {
					var fieldComments strings.Builder
					logDebug("comments isn't nil %s",field.Doc.Text())
					for _,comment := range field.Doc.List{
						commentValue := comment.Text
						commentValue = strings.TrimLeft(commentValue,"//")
						commentValue = strings.TrimLeft(commentValue," ")
						// skip TODO and lines started with generation comment +
						if strings.Contains(commentValue,"TODO") || strings.HasPrefix(commentValue,"+"){
							continue
						}
						logDebug("comment: %v",commentValue)

						fieldComments.WriteString(commentValue)
					}
					// set description to comments, if it not set explicitly with tag
					if fieldDescription == "-"{
						logDebug("override fieldDescription")
						fieldDescription = fieldComments.String()
					}
				}

				if tag.Get("split_words") == "true" {
					words := gatherRegexp.FindAllStringSubmatch(field.Names[0].Name, -1)
					if len(words) > 0 {
						var name []string
						for _, words := range words {
							name = append(name, words[0])
						}

						field.Names[0].Name = strings.Join(name, "_")
					}
				}

				override := prefix  + strings.ToUpper(field.Names[0].Name)

				if tag.Get("envconfig") != "" {
					override = tag.Get("envconfig")

				}
				if *truncate {
					_, _ = fileDesc.Write([]byte(fmt.Sprintf("| %v | %.30v | %.10v | %.50v |\n", override, fieldDefaultValue, fieldRequired, fieldDescription)))

				} else {
					_, _ = fileDesc.Write([]byte(fmt.Sprintf("| %v | %v | %v | %v |\n", override, fieldDefaultValue, fieldRequired, fieldDescription)))

				}

				logDebug("processed struct: %v", currStruct)

			}



}

// list constants at file
// try to find prefix constant name
func extractConstPrefix(node []ast.Decl)string {
	var prefString string
	for _, v := range node {
		switch decl := v.(type) {
		case *ast.GenDecl:
			switch decl.Tok {
			case token.CONST:
			CONL:
				for _, spec := range decl.Specs {
					vspec := spec.(*ast.ValueSpec)
					if *prefixConstName != vspec.Names[0].Name {
						continue CONL
					}
					logDebug("constant found for prefix %v", vspec.Names[0].Name)

					for _, v := range vspec.Values {

						reflected := strings.Trim(reflect.ValueOf((v.(*ast.BasicLit)).Value).String(), "\"")
						logDebug("value of prefix, that we found  %v  ", reflected)
						prefString = reflected
						break
					}
				}
			}
		}

	}
	if prefString != ""{
		prefString += "_"
	}
	return prefString
}
func findStructsAndWalk(nodes []ast.Decl,prefix string,fileDesc io.Writer){

	for _, v := range nodes {
		g, ok := v.(*ast.GenDecl)
		if !ok {
			continue
		}
	SPECL:
		for _, d := range g.Specs {
			spec, ok := d.(*ast.TypeSpec)
			if !ok {
				continue
			}
			_, ok = spec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			if g.Doc == nil {
				continue
			}
			needGen := false
			for _, doc := range g.Doc.List {
				//we only parse struct with matched comment
				needGen = needGen || strings.HasPrefix(doc.Text, *structNeedGenComment)
			}
			if !needGen {
				continue SPECL
			}
			//now we are walking over struct and writing to file
			walkStructPrefix(prefix, spec.Type.(*ast.StructType), fileDesc)
		}

	}
}