package new

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/urfave/cli"
	"github.com/l306287405/go-zero/tools/goctl/api/gogen"
	conf "github.com/l306287405/go-zero/tools/goctl/config"
	"github.com/l306287405/go-zero/tools/goctl/util"
	"github.com/l306287405/go-zero/tools/goctl/util/pathx"
)

const apiTemplate = `
type Request {
  Name string ` + "`" + `path:"name,options=you|me"` + "`" + ` 
}

type Response {
  Message string ` + "`" + `json:"message"` + "`" + `
}

service {{.name}}-api {
  @handler {{.handler}}Handler
  get /from/:name(Request) returns (Response)
}
`

// CreateServiceCommand fast create service
func CreateServiceCommand(c *cli.Context) error {
	args := c.Args()
	dirName := args.First()
	if len(dirName) == 0 {
		dirName = "greet"
	}

	dirStyle := c.String("style")
	if len(dirStyle) == 0 {
		dirStyle = conf.DefaultFormat
	}
	if strings.Contains(dirName, "-") {
		return errors.New("api new command service name not support strikethrough, because this will used by function name")
	}

	abs, err := filepath.Abs(dirName)
	if err != nil {
		return err
	}

	err = pathx.MkdirIfNotExist(abs)
	if err != nil {
		return err
	}

	dirName = filepath.Base(filepath.Clean(abs))
	filename := dirName + ".api"
	apiFilePath := filepath.Join(abs, filename)
	fp, err := os.Create(apiFilePath)
	if err != nil {
		return err
	}

	defer fp.Close()

	home := c.String("home")
	remote := c.String("remote")
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote)
		if len(repo) > 0 {
			home = repo
		}
	}

	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	text, err := pathx.LoadTemplate(category, apiTemplateFile, apiTemplate)
	if err != nil {
		return err
	}

	t := template.Must(template.New("template").Parse(text))
	if err := t.Execute(fp, map[string]string{
		"name":    dirName,
		"handler": strings.Title(dirName),
	}); err != nil {
		return err
	}

	err = gogen.DoGenProject(apiFilePath, abs, dirStyle)
	return err
}
