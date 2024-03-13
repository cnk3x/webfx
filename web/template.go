package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/cnk3x/webfx/utils/bufpool"
)

func GetTemplateContext(r *http.Request, f http.FileSystem) TemplateContext {
	return TemplateContext{f: f, r: r}
}

// Context is the Context with which HTTP templates are executed.
type TemplateContext struct {
	f http.FileSystem    //
	r *http.Request      //
	t *template.Template //
	m any
}

// NewTemplate returns a new template intended to be evaluated with this context, as it is initialized with configuration from this context.
func (c *TemplateContext) NewTemplate(tplName string) *template.Template {
	c.t = template.New(tplName)
	return c.t.Funcs(funcMap).Funcs(contextFuncMap(c.r)).Funcs(c.funcMap())
}

func (c TemplateContext) Execute(templateFilename string, buf *bytes.Buffer, model any) (err error) {
	if err = c.readTemplateFile(templateFilename, buf); err != nil {
		return
	}
	if err = c.bufferExecute(templateFilename, buf, model); err != nil {
		return
	}
	return
}

func (c TemplateContext) bufferExecute(tplName string, buf *bytes.Buffer, model any) (err error) {
	if _, err = c.NewTemplate(tplName).Parse(buf.String()); err != nil {
		return
	}
	buf.Reset()
	c.m = model
	return c.t.Execute(buf, model)
}

// funcImport parses the filename into the current template stack.
// The imported file will be rendered within the current template by calling {{ block }} or {{ template }} from the standard template library.
// If the imported file has no {{ define }} blocks, the name of the import will be the path
func (c *TemplateContext) funcImport(filename string) (out string, err error) {
	bodyBuf := bufpool.Get()
	defer bufpool.Put(bodyBuf)
	if err = c.readTemplateFile(filename, bodyBuf); err != nil {
		return
	}
	if _, err = c.t.Parse(bodyBuf.String()); err != nil {
		return
	}
	return
}

// funcInclude returns the contents of filename relative to the site root and renders it in place.
// Note that included files are NOT escaped, so you should only include trusted files.
// If it is not trusted, be sure to use escaping functions in your template.
func (c TemplateContext) funcInclude(filename string, args ...any) (out string, err error) {
	bodyBuf := bufpool.Get()
	defer bufpool.Put(bodyBuf)

	if err = c.readTemplateFile(filename, bodyBuf); err != nil {
		return
	}

	var model any
	if len(args) > 1 {
		model = args
	} else if len(args) == 1 {
		model = args[0]
	} else {
		model = c.m
	}

	if err = c.bufferExecute(filename, bodyBuf, model); err != nil {
		return
	}

	return bodyBuf.String(), nil
}

// funcReadFile returns the contents of a filename relative to the site root.
// Note that included files are NOT escaped, so you should only include trusted files.
// If it is not trusted, be sure to use escaping functions in your template.
func (c TemplateContext) funcReadFile(filename string) (out string, err error) {
	bodyBuf := bufpool.Get()
	defer bufpool.Put(bodyBuf)
	if err = c.readTemplateFile(filename, bodyBuf); err != nil {
		return
	}
	out = bodyBuf.String()
	return
}

// ReadFileTo reads a file into a buffer
func (c TemplateContext) readTemplateFile(filename string, w io.Writer) (err error) {
	if c.f == nil {
		err = fmt.Errorf("root file system not specified")
		return
	}
	var file http.File
	file, err = c.f.Open(filename)
	if os.IsNotExist(err) && !strings.HasSuffix(filename, ".gohtml") {
		if f, e := c.f.Open(filename + ".gohtml"); e == nil {
			err = nil
			file = f
		}
	}
	if err != nil {
		return
	}
	defer file.Close()
	_, err = io.Copy(w, file)
	return
}

func (c TemplateContext) funcMap() template.FuncMap {
	return template.FuncMap{"include": c.funcInclude, "import": c.funcImport, "readFile": c.funcReadFile}
}

func contextFuncMap(r *http.Request) template.FuncMap {
	c := GetContext(r)
	return template.FuncMap{
		"cookie":  c.Cookie,
		"host":    c.Host,
		"port":    c.Port,
		"ip":      c.RealIp,
		"istls":   c.IsTLS,
		"scheme":  c.Scheme,
		"origin":  c.Origin,
		"path":    c.Path,
		"pathDir": c.PathDir,
		"getjson": c.GetJSON,
		"query":   c.Query,
		"querys":  c.Querys,
		"header":  c.HeaderGet,
	}
}
