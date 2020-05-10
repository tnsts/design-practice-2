package gomodule

import (
	"fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
)

var (
	// Package context used to define Ninja build rules.
	pctx = blueprint.NewPackageContext("github.com/tnsts/design-practice-2/build/gomodule")

	// Ninja rule to execute go build.
	goBuild = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
		Command:     "cd $workDir && go build -o $outputPath $pkg",
		Description: "build go command $pkg",
	}, "workDir", "outputPath", "pkg")

	// Ninja rule to execute go mod vendor.
	goVendor = pctx.StaticRule("vendor", blueprint.RuleParams{
		Command:     "cd $workDir && go mod vendor",
		Description: "vendor dependencies of $name",
	}, "workDir", "name")

	goTest = pctx.StaticRule("test", blueprint.RuleParams{
		Command:     "cd ${workDir} && go test -v ${pkg} > ${outputPath}",
		Description: "test ${pkg}",
	}, "workDir", "outputPath", "pkg")


)

type testedBinaryModule struct {
	blueprint.SimpleName

	properties struct {
		Name string
		Pkg string
		TestPkg string
		Srcs []string
		SrcsExclude []string
		VendorFirst bool
		TestSrcs []string
		TestSrcsExclude []string
		TestsResFile string
	}
}

func (tb *testedBinaryModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
		config := bood.ExtractConfig(ctx)
		config.Debug.Printf("Adding build actions for go binary module '%s'", name)

		outputPath := path.Join(config.BaseOutputDir, "bin", name)
		testOutputPath := path.Join(config.BaseOutputDir, "test-results", "test-res.txt")
		if len(tb.properties.TestsResFile) > 0{
			testOutputPath = path.Join(config.BaseOutputDir, "test-results", tb.properties.TestsResFile)
		}

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

		var testInputs []string
		for _, src := range tb.properties.TestSrcs {
			if matches, err := ctx.GlobWithDeps(src, tb.properties.TestSrcsExclude); err == nil {
				 testInputs = append( testInputs, matches...)
			} else {
				ctx.PropertyErrorf("testSrcs", "Cannot resolve files that match pattern %s", src)
				inputErors = true
			}
		}
		if inputErors {
			return
		}
		testInputs = append( testInputs, inputs...)

		if tb.properties.VendorFirst {
			vendorDirPath := path.Join(ctx.ModuleDir(), "vendor")
			ctx.Build(pctx, blueprint.BuildParams{
				Description: fmt.Sprintf("Vendor dependencies of %s", name),
				Rule:        goVendor,
				Outputs:     []string{vendorDirPath},
				Implicits:   []string{path.Join(ctx.ModuleDir(), "go.mod")},
				Optional:    true,
				Args: map[string]string{
					"workDir": ctx.ModuleDir(),
					"name":    name,
				},
			})
			inputs = append(inputs, vendorDirPath)
		}

		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Build %s as Go binary", name),
			Rule:        goBuild,
			Outputs:     []string{outputPath},
			Implicits:   inputs,
			Args: map[string]string{
				"outputPath": outputPath,
				"workDir":    ctx.ModuleDir(),
				"pkg":        tb.properties.Pkg,
			},
		})

	  if len(tb.properties.TestPkg) > 0 {
			ctx.Build(pctx, blueprint.BuildParams{
				Description: fmt.Sprintf("%s tests to Go binary", name),
	  		Rule:        goTest,
				Outputs:     []string{testOutputPath},
				Implicits:   testInputs,
				Args: map[string]string{
					"outputPath": testOutputPath,
					"workDir":    ctx.ModuleDir(),
					"pkg":        tb.properties.TestPkg,
				},
			})
		}
}

func TestBinFactory() (blueprint.Module, []interface{}) {
	mType := &testedBinaryModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
