package gomodule

import (
    //"fmt"
	"bytes"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"strings"
	"testing"
)

var fileSystems = []map[string][]byte{
	{
		"Blueprints": []byte(`
			go_binary {
			  name: "package-out",
			  pkg: ".",
              testPkg: ".",
			  srcs: [ "main_test.go", "main.go",],
			}
		`),
		"main.go":      nil,
		"main_test.go": nil,
	},
	{
		"Blueprints": []byte(`
			go_binary {
			  name: "package-out",
			  pkg: ".",
			  srcs: [ "main_test.go", "main.go",],
			}
		`),
		"main.go":      nil,
		"main_test.go": nil,
	},
}

var expectedOutput = [][]string{
	{
		"out/bin/package-out:",
		"g.gomodule.binaryBuild",
		"out/bin/test.txt",
		"g.gomodule.test | main_test.go main.go",
	},
	{
		"out/bin/package-out:",
		"g.gomodule.binaryBuild ",
	},
}

func TestSimpleBinFactory(t *testing.T) {
	for index, fileSystem := range fileSystems {
		t.Run(string(index), func(t *testing.T) {
			ctx := blueprint.NewContext()

			ctx.MockFileSystem(fileSystem)

			ctx.RegisterModuleType("go_binary", SimpleBinFactory)

			cfg := bood.NewConfig()

			_, errs := ctx.ParseBlueprintsFiles(".", cfg)
			if len(errs) != 0 {
				t.Fatalf("Syntax errors in the test blueprint file: %s", errs)
			}

			_, errs = ctx.PrepareBuildActions(cfg)
			if len(errs) != 0 {
				t.Errorf("Unexpected errors while preparing build actions: %s", errs)
			}
			buffer := new(bytes.Buffer)
			if err := ctx.WriteBuildFile(buffer); err != nil {
				t.Errorf("Error writing ninja file: %s", err)
			} else {
				text := buffer.String()
                 //fmt.Println(text)
				for _, expectedStr := range expectedOutput[index] {
					if strings.Contains(text, expectedStr) != true {
						t.Errorf("Generated ninja file does not have expected string `%s`", expectedStr)
					}
				}

			}

		})
	}
}