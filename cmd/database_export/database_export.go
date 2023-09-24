package main

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/vikpe/qw-demobot/internal/pkg/demo/collection"
	"github.com/vikpe/qw-demobot/internal/pkg/demo/export"
	"os"
)

func main() {
	demoCollection := collection.New("/home/vikpe/games/demoquake/qw/demos")
	itemCallback := func(demoFilepath string, demoExport export.Export) export.Export {
		return demoExport
	}
	demoExports := demoCollection.Export(itemCallback)
	jsonData, err := json.MarshalIndent(demoExports, "", "  ")
	err = os.WriteFile("demos.json", jsonData, 0644)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(fmt.Sprintf("Done (%d demos)", len(demoExports)))
	}

	foo := make(map[string]string)

	for _, demoExport := range demoExports {
		foo[demoExport.Sha256] = demoExport.Filepath
	}

	jsonData, err = json.MarshalIndent(foo, "", "  ")
	err = os.WriteFile("demos_map.json", jsonData, 0644)

	if err != nil {
		fmt.Println(err)
	}
}
