package ctx

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/l306287405/go-zero/core/jsonx"
	"github.com/l306287405/go-zero/tools/goctl/rpc/execx"
	"github.com/l306287405/go-zero/tools/goctl/util/pathx"
)

// Module contains the relative data of go module,
// which is the result of the command go list
type Module struct {
	Path      string
	Main      bool
	Dir       string
	GoMod     string
	GoVersion string
}

// projectFromGoMod is used to find the go module and project file path
// the workDir flag specifies which folder we need to detect based on
// only valid for go mod project
func projectFromGoMod(workDir string) (*ProjectContext, error) {
	if len(workDir) == 0 {
		return nil, errors.New("the work directory is not found")
	}
	if _, err := os.Stat(workDir); err != nil {
		return nil, err
	}

	workDir, err := pathx.ReadLink(workDir)
	if err != nil {
		return nil, err
	}

	data, err := execx.Run("go list -json -m", workDir)
	if err != nil {
		return nil, err
	}

	var m Module
	err = jsonx.Unmarshal([]byte(data), &m)
	if err != nil {
		return nil, err
	}
	var ret ProjectContext
	ret.WorkDir = workDir
	ret.Name = filepath.Base(m.Dir)
	dir, err := pathx.ReadLink(m.Dir)
	if err != nil {
		return nil, err
	}

	ret.Dir = dir
	ret.Path = m.Path
	return &ret, nil
}
