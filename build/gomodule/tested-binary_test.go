package gomodule

import (
  "github.com/google/blueprint"
  "github.com/roman-mazur/bood"
  "testing"
  "bytes"
  "strings"
)

func TestTestBinFactory(t *testing.T) {
  substrings := []struct {
    str string
    err string
  }{
    {"out/bin/test-binary:", "Wrong binary name in build.ninja"},
    {"main.go", "Wrong source file in build.ninja"},
    {"out/test-results/test-res.txt", "Wrong result file in build.ninja"},
  }
  ctx := blueprint.NewContext()

  ctx.MockFileSystem(map[string][]byte{
    "Blueprints": []byte(`
      go_binary {
        name: "test-binary",
        srcs:["main.go"],
        pkg: ".",
        testPkg: ".",
        }
        `),
      "main.go": nil,
  })

  ctx.RegisterModuleType("go_binary", TestBinFactory)

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
    t.Logf("build.ninja:   %s", text)
		for _, substring := range substrings{
      if !strings.Contains(text, substring.str){
        t.Errorf(substring.err)
      }
    }
  }
}
