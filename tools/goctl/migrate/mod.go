package migrate

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/l306287405/go-zero/core/stringx"
	"github.com/l306287405/go-zero/tools/goctl/rpc/execx"
	"github.com/l306287405/go-zero/tools/goctl/util/console"
	"github.com/l306287405/go-zero/tools/goctl/util/ctx"
)

const deprecatedGoZeroMod = "github.com/tal-tech/go-zero"
const deprecatedBuilderx = "github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
const replacementBuilderx = "github.com/l306287405/go-zero/core/stores/builder"
const goZeroMod = "github.com/l306287405/go-zero"

var errInvalidGoMod = errors.New("it's only working for go module")

func editMod(version string, verbose bool) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	isGoMod, _ := ctx.IsGoMod(wd)
	if !isGoMod {
		return nil
	}

	latest, err := getLatest(goZeroMod, verbose)
	if err != nil {
		return err
	}

	if !stringx.Contains(latest, version) {
		return fmt.Errorf("release version %q is not found", version)
	}
	mod := fmt.Sprintf("%s@%s", goZeroMod, version)
	err = removeRequire(deprecatedGoZeroMod, verbose)
	if err != nil {
		return err
	}
	return addRequire(mod, verbose)
}

func addRequire(mod string, verbose bool) error {
	if verbose {
		console.Info("adding require %s ...", mod)
		time.Sleep(200 * time.Millisecond)
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	isGoMod, _ := ctx.IsGoMod(wd)
	if !isGoMod {
		return errInvalidGoMod
	}

	_, err = execx.Run("go mod edit -require "+mod, wd)
	return err
}

func removeRequire(mod string, verbose bool) error {
	if verbose {
		console.Info("remove require %s ...", mod)
		time.Sleep(200 * time.Millisecond)
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = execx.Run("go mod edit -droprequire "+mod, wd)
	return err
}

func tidy(verbose bool) error {
	if verbose {
		console.Info("go mod tidy ...")
		time.Sleep(200 * time.Millisecond)
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	isGoMod, _ := ctx.IsGoMod(wd)
	if !isGoMod {
		return nil
	}

	_, err = execx.Run("go mod tidy", wd)
	return err
}
