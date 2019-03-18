package configgenerator

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

// Yaml parsing helpers.

// GetString casts the given interface (expected string) as string.
// An array of length 1 is also considered a single string.
func GetString(s interface{}) string {
	if _, ok := s.([]interface{}); ok {
		values := GetStringArray(s)
		if len(values) == 1 {
			return values[0]
		}
		log.Fatalf("Entry %v is not a string or string array of size 1", s)
	}
	if str, ok := s.(string); ok {
		return str
	}
	log.Fatalf("Entry %v is not a string", s)
	return ""
}

// GetInt casts the given interface (expected int) as int.
func GetInt(s interface{}) int {
	if value, ok := s.(int); ok {
		return value
	}
	log.Fatalf("Entry %v is not an integer", s)
	return 0
}

// GetBool casts the given interface (expected bool) as bool.
func GetBool(s interface{}) bool {
	if value, ok := s.(bool); ok {
		return value
	}
	log.Fatalf("Entry %v is not a boolean", s)
	return false
}

// GetInterfaceArray casts the given interface (expected interface array) as interface array.
func GetInterfaceArray(s interface{}) []interface{} {
	if interfaceArray, ok := s.([]interface{}); ok {
		return interfaceArray
	}
	log.Fatalf("Entry %v is not an interface array", s)
	return nil
}

// GetStringArray casts the given interface (expected string array) as string array.
func GetStringArray(s interface{}) []string {
	interfaceArray := GetInterfaceArray(s)
	strArray := make([]string, len(interfaceArray))
	for i := range interfaceArray {
		strArray[i] = GetString(interfaceArray[i])
	}
	return strArray
}

// GetMapSlice casts the given interface (expected MapSlice) as MapSlice.
func GetMapSlice(m interface{}) yaml.MapSlice {
	if mm, ok := m.(yaml.MapSlice); ok {
		return mm
	}
	log.Fatalf("Entry %v is not a yaml.MapSlice", m)
	return nil
}

// general helper functions for processing the template text

// Quote returns the given string quoted if it's not a key/value pair or already quoted.
func Quote(s string) string {
	if strings.Contains(s, "\"") || strings.Contains(s, ": ") || strings.HasSuffix(s, ":") {
		return s
	}
	return "\"" + s + "\""
}

// IndentBase is a helper function which returns the given array indented.
func IndentBase(indentation int, prefix string, indentFirstLine bool, quote bool, array []string) string {
	s := ""
	if len(array) == 0 {
		return s
	}
	indent := strings.Repeat(" ", indentation)
	for i := 0; i < len(array); i++ {
		if i > 0 || indentFirstLine {
			s += indent
		}
		if quote {
			s += prefix + Quote(array[i]) + "\n"
		} else {
			s += prefix + array[i] + "\n"
		}
	}
	return s
}

// IndentArray returns the given array indented, prefixed by "-".
func IndentArray(indentation int, array []string) string {
	return IndentBase(indentation, "- ", false, true, array)
}

// IndentArrayWithoutQuote returns the given array indented, prefixed by "-" and the content unquoted
func IndentArrayWithoutQuote(indentation int, array []string) string {
	return IndentBase(indentation, "- ", false, false, array)
}

// IndentKeys returns the given array of key/value pairs indented.
func IndentKeys(indentation int, array []string) string {
	return IndentBase(indentation, "", false, true, array)
}

// IndentSectionBase is a helper function which returns the given array of key/value pairs indented inside a section.
func IndentSectionBase(indentation int, title string, prefix string, array []string) string {
	keys := IndentBase(indentation, prefix, true, true, array)
	if keys == "" {
		return keys
	}
	return title + ":\n" + keys
}

// IndentArraySection returns the given array indented inside a section.
func IndentArraySection(indentation int, title string, array []string) string {
	return IndentSectionBase(indentation, title, "- ", array)
}

// IndentSection returns the given array of key/value pairs indented inside a section.
func IndentSection(indentation int, title string, array []string) string {
	return IndentSectionBase(indentation, title, "", array)
}

// IndentMap returns the given map indented, with each key/value separated by ": "
func IndentMap(indentation int, mp map[string]string) string {
	arr := make([]string, len(mp))
	i := 0
	for k, v := range mp {
		arr[i] = k + ": " + v
		i++
	}
	return IndentBase(indentation, "", false, true, arr)
}

// OutputConfig outputs the given line, if not empty, to stdout.
func OutputConfig(line string) {
	s := strings.TrimSpace(line)
	if s != "" {
		fmt.Println(line)
	}
}

// StrExists checks if the given string exists in the array
func StrExists(arr []string, str string) bool {
	for _, s := range arr {
		if str == s {
			return true
		}
	}
	return false
}
