package output

import (
	"bytes"
	"errors"
	"reflect"
	"text/template"
)

func formatTemplate(buf *bytes.Buffer, v interface{}, config *Config) error {
	if config.Template == "" {
		return errors.New("template must not be empty")
	}

	tpl, err := template.New("template").
		Option(config.TemplateConfig.Options...).
		Funcs(config.TemplateConfig.Funcs).
		Parse(config.Template)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(v)

	if config.TemplateItems && rv.Kind() == reflect.Slice {
		n := rv.Len()

		for i := 0; i < n; i++ {
			v := rv.Index(i).Interface()

			if err := tpl.Execute(buf, v); err != nil {
				return err
			}

			if i+1 < n {
				buf.WriteByte('\n')
			}
		}

		return nil
	}

	return tpl.Execute(buf, v)
}
