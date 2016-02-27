package main

import (
	sdk "golaunch/sdk/go"
	"golaunch/sdk/go/plugin"
	"log"
	"os"
)

type Plugin struct {
	metadata sdk.Metadata
	client   *plugin.Client
}

func NewPlugin() *Plugin {
	return &Plugin{
		client: plugin.NewClient(),
	}
}

func (p *Plugin) Init(m sdk.Metadata) {
	p.metadata = m
}

func (p *Plugin) Query(q string) {
	if q == "helloworld" {
		p.client.Call("queryResults", []sdk.QueryResult{{
			Title:    "Hello back to you!",
			Subtitle: "From the helloworld plugin",
			Query:    q,
			ID:       p.metadata.ID,
			Score:    -1,
		}})
		return
	}

	p.client.Call("noQueryResults", nil)
}

func (p *Plugin) Action(a sdk.Action) {
	//log.Print("action for helloworld pressed")
}

func main() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stderr)

	s := plugin.NewServer()
	s.Register(NewPlugin())
	s.Serve()
}
