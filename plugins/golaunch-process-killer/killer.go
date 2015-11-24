package main

import (
	"encoding/json"
	"fmt"
	sdk "golaunch/sdk/go"
	"golaunch/sdk/go/fuzzy"
	"golaunch/sdk/go/system"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mozilla/masche/process"
	"github.com/spf13/cast"
)

func main() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stderr)

	// empty for now, but will be filled by "init"
	var metadata sdk.Metadata

	system := system.NewSystem()
	defer system.Close()

	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)

	for {
		var v sdk.Request
		if err := dec.Decode(&v); err != nil {
			continue
		}

		switch v.Method {
		case "init":
			json.Unmarshal(v.Params, &metadata)
		case "query":
			var query string
			json.Unmarshal(v.Params, &query)

			if strings.HasPrefix(query, "kill ") {
				processQuery := query[5:]
				var results []sdk.QueryResult

				p, herr, _ := process.OpenAll()
				if herr != nil {
					log.Fatal(herr)
				}

				var allPids []uint
				for _, v := range p {
					name, _, _ := v.Name()
					shortname := filepath.Base(name[:len(name)-len(filepath.Ext(name))])
					mr := fuzzy.Match(processQuery, shortname)
					if mr.Success || len(processQuery) == 0 {
						image, _ := system.EmbeddedAppIcon(name)
						results = append(results, sdk.QueryResult{
							Program: sdk.Program{
								Image: image,
							},
							Title:    fmt.Sprintf("%s - %d", shortname, v.Pid()),
							Subtitle: name,
							Query:    query,
							ID:       metadata.ID,
							Score:    -1,
							Data:     []uint{v.Pid()},
						})
						allPids = append(allPids, v.Pid())
					}
				}

				if len(results) > 0 {
					if len(processQuery) != 0 {
						results = append([]sdk.QueryResult{{
							Program: sdk.Program{
								Image: metadata.Icon,
							},
							Title:    "Kill all \"" + processQuery + "\" processes",
							Subtitle: "Make sure it matches what you want!",
							Query:    query,
							ID:       metadata.ID,
							Score:    -1,
							Data:     allPids,
						}}, results...)
					}

					msg := sdk.Response{
						Result: results,
					}

					enc.Encode(msg)
				}
			}

		case "action":
			var param sdk.Action
			json.Unmarshal(v.Params, &param)

			for _, pid := range cast.ToIntSlice(param.QueryResult.Data) {
				p, _ := os.FindProcess(pid)
				p.Kill()
			}
		}
	}
}
