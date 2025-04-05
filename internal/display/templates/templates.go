package templates

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/zimeg/emporia-time/internal/display"
	"github.com/zimeg/emporia-time/internal/errors"
)

// templateBuilder creates a string using a template and variables
func templateBuilder(templateStr string, body any) (string, error) {
	funcs := template.FuncMap{
		"Bold": func(f string) string {
			return fmt.Sprintf("\x1b[1m%s\x1b[0m", f)
		},
		"CommandName": func() string {
			return os.Args[0]
		},
		"Percent": func(f float64, spacing int) string {
			return fmt.Sprintf("%*.1f", spacing, f*100)
		},
		"Time": func(f float64, spacing int) string {
			return fmt.Sprintf("%*.2f", spacing, f)
		},
		"TimeF": func(f float64, spacing int) string {
			return fmt.Sprintf("%*s", spacing, display.FormatSeconds(f))
		},
		"Value": func(f float64, spacing int) string {
			return fmt.Sprintf("%*.2f", spacing, f)
		},
	}
	tmpl, err := template.New("outputs").Funcs(funcs).Parse(templateStr)
	if err != nil {
		return "", errors.Wrap(errors.ErrTemplateParse, err)
	}
	var result strings.Builder
	if err := tmpl.Execute(&result, body); err != nil {
		return "", errors.Wrap(errors.ErrTemplateBuild, err)
	}
	return result.String(), nil
}
