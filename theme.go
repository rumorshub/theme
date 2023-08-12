package theme

import (
	"context"
	"html/template"
	"strings"

	"github.com/Masterminds/sprig/v3"
	et "github.com/gowool/extends-template"
)

var _ Theme = (*Environment)(nil)

var FuncMap = sprig.FuncMap()

func init() {
	FuncMap["raw"] = func(s string) template.HTML {
		return template.HTML(s)
	}
}

type Theme interface {
	Debug(debug bool) *et.Environment
	Funcs(funcMap template.FuncMap) *et.Environment
	Global(global ...string) *et.Environment
	Load(ctx context.Context, name string) (*et.TemplateWrapper, error)
	HTML(ctx context.Context, name string, data any) (string, error)
}

type Environment struct {
	*et.Environment
}

func NewEnvironment(loaders []et.Loader) *Environment {
	env := et.NewEnvironment(et.NewChainLoader(loaders...))
	env.Funcs(FuncMap)

	return &Environment{Environment: env}
}

func (t *Environment) HTML(ctx context.Context, name string, data any) (string, error) {
	wrap, err := t.Load(ctx, name)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	if err = wrap.HTML.ExecuteTemplate(&sb, name, data); err != nil {
		return "", err
	}

	return sb.String(), nil
}
