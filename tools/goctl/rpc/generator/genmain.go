package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	conf "github.com/l306287405/go-zero/tools/goctl/config"
	"github.com/l306287405/go-zero/tools/goctl/rpc/parser"
	"github.com/l306287405/go-zero/tools/goctl/util"
	"github.com/l306287405/go-zero/tools/goctl/util/format"
	"github.com/l306287405/go-zero/tools/goctl/util/pathx"
	"github.com/l306287405/go-zero/tools/goctl/util/stringx"
)

const mainTemplate = `package main

import (
	"flag"
	"fmt"

	{{.imports}}

	"github.com/l306287405/go-zero/core/conf"
	"github.com/l306287405/go-zero/core/service"
	"github.com/l306287405/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.New{{.serviceNew}}Server(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		{{.pkg}}.Register{{.service}}Server(grpcServer, srv)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
`

// GenMain generates the main file of the rpc service, which is an rpc service program call entry
func (g *DefaultGenerator) GenMain(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
	mainFilename, err := format.FileNamingFormat(cfg.NamingFormat, ctx.GetServiceName().Source())
	if err != nil {
		return err
	}

	fileName := filepath.Join(ctx.GetMain().Filename, fmt.Sprintf("%v.go", mainFilename))
	imports := make([]string, 0)
	pbImport := fmt.Sprintf(`"%v"`, ctx.GetPb().Package)
	svcImport := fmt.Sprintf(`"%v"`, ctx.GetSvc().Package)
	remoteImport := fmt.Sprintf(`"%v"`, ctx.GetServer().Package)
	configImport := fmt.Sprintf(`"%v"`, ctx.GetConfig().Package)
	imports = append(imports, configImport, pbImport, remoteImport, svcImport)
	text, err := pathx.LoadTemplate(category, mainTemplateFile, mainTemplate)
	if err != nil {
		return err
	}

	etcFileName, err := format.FileNamingFormat(cfg.NamingFormat, ctx.GetServiceName().Source())
	if err != nil {
		return err
	}

	return util.With("main").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"serviceName": etcFileName,
		"imports":     strings.Join(imports, pathx.NL),
		"pkg":         proto.PbPackage,
		"serviceNew":  stringx.From(proto.Service.Name).ToCamel(),
		"service":     parser.CamelCase(proto.Service.Name),
	}, fileName, false)
}
