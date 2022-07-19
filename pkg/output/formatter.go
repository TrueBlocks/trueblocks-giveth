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
		sb.WriteString(quote + field.Name + quote)
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
		if field.Kind == reflect.String || field.IsDate {
			sb.WriteString(quote + "{{." + field.Name + "}}" + quote)
		} else {
			sb.WriteString("{{." + field.Name + "}}")
		}
	}
	tt, err := template.New("").Parse(sb.String() + "\n")
	return tt, err
}

type Field struct {
	Name   string
	Kind   reflect.Kind
	IsDate bool
}

func GetFields(t *reflect.Type, format string, header bool) (fields []Field, sep string, quote string) {
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
		names := strings.Split(custom, ",")
		for _, name := range names {
			fields = append(fields, Field{
				Name:   name,
				Kind:   reflect.String,
				IsDate: strings.Contains(strings.ToLower(name), "date"),
			})
		}

	} else {
		if (*t).Kind() != reflect.Struct {
			fields = append(fields, Field{})
		} else {
			for i := 0; i < (*t).NumField(); i++ {
				name := (*t).Field(i).Name
				if header {
					fields = append(fields, Field{
						Name:   MakeFirstLowerCase(name),
						Kind:   (*t).Field(i).Type.Kind(),
						IsDate: false,
					})
				} else {
					fields = append(fields, Field{
						Name:   name,
						Kind:   (*t).Field(i).Type.Kind(),
						IsDate: strings.Contains(strings.ToLower(name), "date"),
					})
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
