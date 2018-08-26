package tfx

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
	"text/template"
	"text/template/parse"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Factory interface {
	Load(tb testing.TB, name string, v interface{}, opts ...LoadOption)
}

var DefaultFactory Factory = New()

func Load(tb testing.TB, name string, v interface{}, opts ...LoadOption) {
	tb.Helper()
	DefaultFactory.Load(tb, name, v, opts...)
}

func New() Factory {
	return &factoryImpl{}
}

type factoryImpl struct {
}

func (f *factoryImpl) Load(tb testing.TB, name string, v interface{}, opts ...LoadOption) {
	tb.Helper()
	var err error

	createTemplate := func() *template.Template {
		return template.New("tfx").Funcs(template.FuncMap{
			"Seq": func(key string) int { return 1 },
			"Now": time.Now,
		})
	}

	tmpl := createTemplate()

	path := filepath.Join("testdata", name+".yaml")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		tb.Fatalf("failed to read %s: %v", path, err)
	}

	tmpl, err = tmpl.Parse(string(data))
	if err != nil {
		tb.Fatalf("failed to load %s: %v", path, err)
	}

	params := createParams(opts)

	for _, fn := range collectFields(tmpl.Tree.Root) {
		if len(fn.Ident) > 1 {
			continue
		}
		key := fn.Ident[0]
		if _, ok := params[key]; !ok {
			params[key] = "{{." + key + "}}"
		}
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, params)

	fmt.Println(buf.String())

	data = buf.Bytes()

	m := yaml.MapSlice{}
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		tb.Fatalf("failed to parse %s: %v", path, err)
	}

	for i := 0; i < len(m); i++ {
		item := m[i]

		switch v := item.Value.(type) {
		case yaml.MapSlice:
			// no-op
			continue
		case yaml.MapItem:
			// no-op
			continue
		case string:
			tmpl, err := createTemplate().Parse(v)
			if err != nil {
				fmt.Println(err)
				continue
			}
			buf := new(bytes.Buffer)
			err = tmpl.Execute(buf, params)
			if err != nil {
				fmt.Println(err)
				continue
			}
			item.Value = buf.String()
			m[i] = item
		}
		if k, ok := item.Key.(string); ok {
			params[k] = item.Value
		}
	}

	for _, item := range m {
		fmt.Printf("%#v\n", item)
	}

	// TODO
}

func collectFields(n parse.Node) (nodes []*parse.FieldNode) {
	fmt.Printf("%#v\n", n)
	switch n := n.(type) {
	case *parse.ActionNode:
		nodes = append(nodes, collectFields(n.Pipe)...)
	case *parse.BoolNode:
		// no-op
	case *parse.BranchNode:
		nodes = append(nodes, collectFields(n.ElseList)...)
		nodes = append(nodes, collectFields(n.List)...)
		nodes = append(nodes, collectFields(n.Pipe)...)
	case *parse.ChainNode:
		nodes = append(nodes, collectFields(n.Node)...)
		fmt.Println(n.Field)
	case *parse.CommandNode:
		for _, arg := range n.Args {
			nodes = append(nodes, collectFields(arg)...)
		}
	case *parse.DotNode:
		// no-op
	case *parse.FieldNode:
		nodes = append(nodes, n)
	case *parse.IdentifierNode:
		// no-op
	case *parse.IfNode:
		nodes = append(nodes, collectFields(&n.BranchNode)...)
	case *parse.ListNode:
		if n == nil {
			return
		}
		for _, item := range n.Nodes {
			nodes = append(nodes, collectFields(item)...)
		}
	case *parse.NilNode:
		// no-op
	case *parse.NumberNode:
		// no-op
	case *parse.PipeNode:
		for _, cn := range n.Cmds {
			nodes = append(nodes, collectFields(cn)...)
		}
		for _, dn := range n.Decl {
			nodes = append(nodes, collectFields(dn)...)
		}
	case *parse.RangeNode:
		nodes = append(nodes, collectFields(&n.BranchNode)...)
	case *parse.StringNode:
		// no-op
	case *parse.TemplateNode:
		nodes = append(nodes, collectFields(n.Pipe)...)
	case *parse.TextNode:
		// no-op
	case *parse.VariableNode:
		// no-op
	case *parse.WithNode:
		nodes = append(nodes, collectFields(&n.BranchNode)...)
	default:
		// no-op
	}

	return
}
