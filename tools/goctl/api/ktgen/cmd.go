package ktgen

import (
	"errors"

	"github.com/urfave/cli"
	"github.com/l306287405/go-zero/tools/goctl/api/parser"
)

// KtCommand the generate kotlin code command entrance
func KtCommand(c *cli.Context) error {
	apiFile := c.String("api")
	if apiFile == "" {
		return errors.New("missing -api")
	}
	dir := c.String("dir")
	if dir == "" {
		return errors.New("missing -dir")
	}
	pkg := c.String("pkg")
	if pkg == "" {
		return errors.New("missing -pkg")
	}

	api, e := parser.Parse(apiFile)
	if e != nil {
		return e
	}

	api.Service = api.Service.JoinPrefix()
	e = genBase(dir, pkg, api)
	if e != nil {
		return e
	}
	e = genApi(dir, pkg, api)
	if e != nil {
		return e
	}
	return nil
}
