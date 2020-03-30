package gomodule

import (
    "fmt"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"path"
    "regexp"
)

var (
	// Package context used to define Ninja build rules.
	pctx = blueprint.NewPackageContext("github.com/dkalytovskyi/go-lab-2/build/gomodule")

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
        Pkg string
        TestPkg string
        Srcs []string
        SrcsExclude []string
        VendorFirst bool
		
	}
}

func (gb *testedBinaryModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	outputPath := path.Join(config.BaseOutputDir, "bin", name)
    testOutputPath := path.Join(config.BaseOutputDir, "bin", "test.txt")

	var testInputs []string
	inputErors := false
	for _, src := range gb.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, gb.properties.SrcsExclude); err == nil {
			testInputs = append(testInputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErors = true
		}
	}
	if inputErors {
		return
	}

    inputs := testInputs

    for i:=0; i<len(testInputs); i++ {
        if val, _ := regexp.Match(".*_test\\.go$", []byte(testInputs[i])); val == false {
          inputs = append(inputs, testInputs[i])
        }
    }

	if gb.properties.VendorFirst {
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
			"pkg":        gb.properties.Pkg,
		},
	})

if len(gb.properties.TestPkg) > 0 {
	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Tests to Go binary for module %s", name),
	  	Rule:        goTest,
		Outputs:     []string{testOutputPath},
		Implicits:   testInputs,
		Args: map[string]string{
			"outputPath": testOutputPath,
			"workDir":    ctx.ModuleDir(),
			"pkg":        gb.properties.TestPkg,
		},
	})
}

}

func SimpleBinFactory() (blueprint.Module, []interface{}) {
	mType := &testedBinaryModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
