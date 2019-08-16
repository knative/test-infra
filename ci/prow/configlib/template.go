package configlib

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"knative.dev/test-infra/shared/common"
)

// TemplateEvaluator provides a set of default template functions used to evaluate the prow config templates
type TemplateEvaluator struct {
	FuncMap template.FuncMap
	OutFile *os.File
}

// NewTemplateEvaluator creates a new TemplateEvaluator
func NewTemplateEvaluator(f *os.File, repoFn interface{}) *TemplateEvaluator {
	return &TemplateEvaluator{
		FuncMap: template.FuncMap{
			"indent_section":       indentSection,
			"indent_array_section": indentArraySection,
			"indent_array":         indentArray,
			"indent_keys":          indentKeys,
			"indent_map":           indentMap,
			"repo":                 repoFn,
		},
		OutFile: f,
	}
}

// ExecuteTemplate outputs the given template with the given data.
func (te *TemplateEvaluator) ExecuteTemplate(name, templ string, data interface{}) {
	var res bytes.Buffer
	t := template.Must(template.New(name).Funcs(te.FuncMap).Delims("[[", "]]").Parse(templ))
	if err := t.Execute(&res, data); err != nil {
		log.Fatalf("Error in template %s: %v", name, err)
	}
	for _, line := range strings.Split(res.String(), "\n") {
		if strings.TrimSpace(line) != "" {
			fmt.Fprintln(te.OutFile, line)
		}
	}
}

// indentBase is a helper function which returns the given array indented.
func indentBase(indentation int, prefix string, indentFirstLine bool, array []string) string {
	s := ""
	if len(array) == 0 {
		return s
	}
	indent := strings.Repeat(" ", indentation)
	for i := 0; i < len(array); i++ {
		if i > 0 || indentFirstLine {
			s += indent
		}
		s += prefix + common.Quote(array[i]) + "\n"
	}
	return s
}

// indentArray returns the given array indented, prefixed by "-".
func indentArray(indentation int, array []string) string {
	return indentBase(indentation, "- ", false, array)
}

// indentKeys returns the given array of key/value pairs indented.
func indentKeys(indentation int, array []string) string {
	return indentBase(indentation, "", false, array)
}

// indentSectionBase is a helper function which returns the given array of key/value pairs indented inside a section.
func indentSectionBase(indentation int, title string, prefix string, array []string) string {
	keys := indentBase(indentation, prefix, true, array)
	if keys == "" {
		return keys
	}
	return title + ":\n" + keys
}

// indentArraySection returns the given array indented inside a section.
func indentArraySection(indentation int, title string, array []string) string {
	return indentSectionBase(indentation, title, "- ", array)
}

// indentSection returns the given array of key/value pairs indented inside a section.
func indentSection(indentation int, title string, array []string) string {
	return indentSectionBase(indentation, title, "", array)
}

// indentMap returns the given map indented, with each key/value separated by ": "
func indentMap(indentation int, mp map[string]string) string {
	// Extract map keys to keep order consistent.
	keys := make([]string, 0, len(mp))
	for key := range mp {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	arr := make([]string, len(mp))
	for i := 0; i < len(mp); i++ {
		arr[i] = keys[i] + ": " + common.Quote(mp[keys[i]])
	}
	return indentBase(indentation, "", false, arr)
}
