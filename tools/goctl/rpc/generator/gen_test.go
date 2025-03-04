package generator

import (
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/l306287405/go-zero/core/logx"
	"github.com/l306287405/go-zero/core/stringx"
	conf "github.com/l306287405/go-zero/tools/goctl/config"
	"github.com/l306287405/go-zero/tools/goctl/rpc/execx"
	"github.com/l306287405/go-zero/tools/goctl/util/pathx"
)

var cfg = &conf.Config{
	NamingFormat: "gozero",
}

func TestRpcGenerate(t *testing.T) {
	_ = Clean()
	dispatcher := NewDefaultGenerator()
	err := dispatcher.Prepare()
	if err != nil {
		logx.Error(err)
		return
	}
	projectName := stringx.Rand()
	g := NewRPCGenerator(dispatcher, cfg)

	src := filepath.Join(build.Default.GOPATH, "src")
	_, err = os.Stat(src)
	if err != nil {
		return
	}

	projectDir := filepath.Join(src, projectName)
	srcDir := projectDir
	defer func() {
		_ = os.RemoveAll(srcDir)
	}()
	common, err := filepath.Abs(".")
	assert.Nil(t, err)

	// case go path
	t.Run("GOPATH", func(t *testing.T) {
		err = g.Generate("./test.proto", projectDir, []string{common}, "Mbase/common.proto=./base")
		assert.Nil(t, err)
		_, err = execx.Run("go test "+projectName, projectDir)
		if err != nil {
			assert.True(t, func() bool {
				return strings.Contains(err.Error(), "not in GOROOT") || strings.Contains(err.Error(), "cannot find package")
			}())
		}
	})

	// case go mod
	t.Run("GOMOD", func(t *testing.T) {
		workDir := pathx.MustTempDir()
		name := filepath.Base(workDir)
		_, err = execx.Run("go mod init "+name, workDir)
		if err != nil {
			logx.Error(err)
			return
		}

		projectDir = filepath.Join(workDir, projectName)
		err = g.Generate("./test.proto", projectDir, []string{common}, "Mbase/common.proto=./base")
		assert.Nil(t, err)
		_, err = execx.Run("go test "+projectName, projectDir)
		if err != nil {
			assert.True(t, func() bool {
				return strings.Contains(err.Error(), "not in GOROOT") || strings.Contains(err.Error(), "cannot find package")
			}())
		}
	})

	// case not in go mod and go path
	t.Run("OTHER", func(t *testing.T) {
		err = g.Generate("./test.proto", projectDir, []string{common, src}, "Mbase/common.proto=./base")
		assert.Nil(t, err)
		_, err = execx.Run("go test "+projectName, projectDir)
		if err != nil {
			assert.True(t, func() bool {
				return strings.Contains(err.Error(), "not in GOROOT") || strings.Contains(err.Error(), "cannot find package")
			}())
		}
	})
}
