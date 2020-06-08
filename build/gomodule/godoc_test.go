package gomodule

import (
  "github.com/google/blueprint"
  "github.com/roman-mazur/bood"
  "testing"
  "bytes"
  "strings"
)

func TestDocFactory(t *testing.T) {
  substrings := []struct {
    str string
    err string
  }{
    {"example.go", "Wrong source file in build.ninja"},
    {"build out/docs/my-docs.html:", "Wrong result file in build.ninja"},
  }
  ctx := blueprint.NewContext()

  ctx.MockFileSystem(map[string][]byte{
    "Blueprints": []byte(`
      go_doc {
        name: "package",
        srcs:["example.go"],
        pkg: ".",
        }
        `),
      "example.go": nil,
  })

  ctx.RegisterModuleType("go_doc", DocFactory)

  cfg := bood.NewConfig()

  _, errs1 := ctx.ParseBlueprintsFiles(".", cfg)
  if len(errs1) != 0 {
    t.Errorf("Parsing errors %s", errs1)
  }
  _, errs2 := ctx.PrepareBuildActions(cfg)
  if len(errs2) != 0 {
    t.Errorf("Preparing errors %s", errs2)
  }
  buffer := new(bytes.Buffer)

  if err := ctx.WriteBuildFile(buffer); err != nil {
    t.Errorf("Writing error %s", err)
  } else {
    text := buffer.String()
		for _, substring := range substrings{
      if !strings.Contains(text, substring.str){
        t.Errorf(substring.err)
      }
    }
  }
}