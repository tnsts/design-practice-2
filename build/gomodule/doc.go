package gomodule

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
)

var (
	// Ninja rule to execute godoc.
	goDocs = pctx.StaticRule("docs", blueprint.RuleParams{
		Command:     "cd $workDir && godoc -url $pkg > $outputPath",
		Description: "generate docs for $pkg",
	}, "workDir", "outputPath", "pkg")
)

type docModule struct {
	blueprint.SimpleName

	properties struct {
		Name string
		Pkg string
		Srcs []string
		SrcsExclude []string
	}
}

func (tb *docModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
		config := bood.ExtractConfig(ctx)
		config.Debug.Printf("Adding build actions for go binary module '%s'", name)

		outputPath := path.Join(config.BaseOutputDir, "docs", "my-docs.html")

		var inputs []string
		inputErors := false
		for _, src := range tb.properties.Srcs {
			if matches, err := ctx.GlobWithDeps(src, tb.properties.SrcsExclude); err == nil {
				inputs = append(inputs, matches...)
			} else {
				ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
				inputErors = true
			}
		}
		if inputErors {
			return
		}

		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Genaration of docs for %s package", name),
			Rule:        goDocs,
			Outputs:     []string{outputPath},
			Implicits:   inputs,
			Args: map[string]string{
				"outputPath": outputPath,
				"workDir":    ctx.ModuleDir(),
				"pkg":        tb.properties.Pkg,
			},
		})
}

func DocFactory() (blueprint.Module, []interface{}) {
	mType := &docModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
