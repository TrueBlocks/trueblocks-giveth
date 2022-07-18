package output

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"reflect"
	"strings"
)

func Header(data interface{}, w io.Writer, format string) {
	if format == "txt" || format == "csv" {
		tt := reflect.TypeOf(data)
		w.Write([]byte(GetHeader(&tt, format)))
		w.Write([]byte("\n"))
	} else if format == "json" || format == "api" {
		w.Write([]byte("{\n  \"data\": [\n    "))
	} else {
		// do nothing
	}
}

func Footer(data interface{}, w io.Writer, format string) {
	if format == "txt" || format == "csv" {
		// do nothing
	} else if format == "json" || format == "api" {
		w.Write([]byte("  ]\n}\n"))
	} else {
		// do nothing
	}
}

func Line(data interface{}, w io.Writer, format string, first bool) {
	var outputBytes []byte

	preceeds := ""
	switch format {
	case "api":
		fallthrough
	case "json":
		if !first {
			fmt.Fprintf(w, "%s", ",")
		}
		fmt.Fprintf(w, "%s", data)
	case "csv":
		fallthrough
	case "txt":
		tt := reflect.TypeOf(data)
		rowTemplate, _ := GetRowTemplate(&tt, format)
		rowTemplate.Execute(w, data)
		return
	default:
		tt := reflect.TypeOf(data)
		rowTemplate, _ := GetRowTemplate(&tt, format)
		rowTemplate.Execute(w, data)
		return
	}

	w.Write([]byte(preceeds))
	w.Write(outputBytes)
}

func GetHeader(t *reflect.Type, format string) string {
	fields, sep, quote := GetFields(t, format, true)
	var sb strings.Builder
	for i, field := range fields {
		if i > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(quote + field + quote)
	}
	return sb.String()
}

func GetRowTemplate(t *reflect.Type, format string) (*template.Template, error) {
	fields, sep, quote := GetFields(t, format, false)
	var sb strings.Builder
	for i, field := range fields {
		if i > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(quote + "{{." + field + "}}" + quote)
	}
	tt, err := template.New("").Parse(sb.String() + "\n")
	return tt, err
}

func GetFields(t *reflect.Type, format string, header bool) (fields []string, sep string, quote string) {
	sep = "\t"
	quote = ""
	if format == "csv" || strings.Contains(format, ",") {
		sep = ","
	}

	if format == "csv" || strings.Contains(format, "\"") {
		quote = "\""
	}

	if strings.Contains(format, "\t") || strings.Contains(format, ",") {
		custom := strings.Replace(format, "\t", ",", -1)
		custom = strings.Replace(custom, "\"", ",", -1)
		fields = strings.Split(custom, ",")

	} else {
		if (*t).Kind() != reflect.Struct {
			fields = append(fields, "")
		} else {
			for i := 0; i < (*t).NumField(); i++ {
				fn := (*t).Field(i).Name
				if header {
					fields = append(fields, MakeFirstLowerCase(fn))
				} else {
					fields = append(fields, fn)
				}
			}
		}
	}

	return fields, sep, quote
}

func MakeFirstLowerCase(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}
	bts := []byte(s)
	lc := bytes.ToLower([]byte{bts[0]})
	rest := bts[1:]
	return string(bytes.Join([][]byte{lc, rest}, nil))
}
