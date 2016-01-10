package main

import (
	"fmt"
	sdk "golaunch/sdk/go"
	"golaunch/sdk/go/fuzzy"
	"golaunch/sdk/go/plugin"
	"golaunch/sdk/go/system"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mozilla/masche/process"
	"github.com/spf13/cast"
)

type Plugin struct {
	metadata sdk.Metadata
	client   *plugin.Client
	system   *system.System
}

func NewPlugin() *Plugin {
	return &Plugin{
		client: plugin.NewClient(),
		// FIXME: how to defer system.Close()??
		system: system.NewSystem(),
	}
}

func (p *Plugin) Init(m sdk.Metadata) {
	p.metadata = m
}

func (p *Plugin) Query(q string) {
	// if q == "helloworld" {
	// 	p.client.Call("queryresults", sdk.QueryResult{
	// 		Title:    "Hello back to you!",
	// 		Subtitle: "From the helloworld plugin",
	// 		Query:    q,
	// 		ID:       p.metadata.ID,
	// 		Score:    -1,
	// 	})
	// }

	if strings.HasPrefix(q, "kill ") {
		processQuery := q[5:]
		var results []sdk.QueryResult

		procs, herr, _ := process.OpenAll()
		if herr != nil {
			log.Fatal(herr)
		}

		var allPids []uint
		for _, v := range procs {
			name, _, _ := v.Name()
			shortname := filepath.Base(name[:len(name)-len(filepath.Ext(name))])
			mr := fuzzy.Match(processQuery, shortname)
			if mr.Success || len(processQuery) == 0 {
				icon, _ := p.system.EmbeddedAppIcon(name)
				results = append(results, sdk.QueryResult{
					Program: sdk.Program{
						Icon: icon,
					},
					Title:    fmt.Sprintf("%s - %d", shortname, v.Pid()),
					Subtitle: name,
					Query:    q,
					ID:       p.metadata.ID,
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
						Icon: p.metadata.Icon,
					},
					Title:    "Kill all \"" + processQuery + "\" processes",
					Subtitle: "Make sure it matches what you want!",
					Query:    q,
					ID:       p.metadata.ID,
					Score:    -1,
					Data:     allPids,
				}}, results...)
			}

			p.client.Call("queryresults", results)
		}
	}
}

func (p *Plugin) Action(a sdk.Action) {
	for _, pid := range cast.ToIntSlice(a.QueryResult.Data) {
		proc, _ := os.FindProcess(pid)
		proc.Kill()
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stderr)

	s := plugin.NewServer()
	s.Register(NewPlugin())
	s.Serve()
}
