package main

import (
	"bufio"
	"fmt"
	sdk "golaunch/sdk/go"
	"golaunch/sdk/go/plugin"
	"log"
	"os"
	"os/exec"
	"strings"
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
	if q == "soundswitch" {
		cmd := exec.Command(".\\vendor\\AudioEndPointController\\Release\\EndPointController.exe")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println(err)
			return
		}

		if err := cmd.Start(); err != nil {
			log.Println(err)
			return
		}

		var results []sdk.QueryResult
		scanner := bufio.NewScanner(stdout)
		index := 1
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, ": ", 2)
			if len(parts) != 2 {
				continue
			}

			results = append(results, sdk.QueryResult{
				Program: sdk.Program{
					Icon: p.metadata.Icon,
				},
				Title:    parts[1],
				Subtitle: parts[0],
				Query:    q,
				ID:       p.metadata.ID,
				Score:    -1,
				Data:     fmt.Sprintf("%d", index),
			})

			index += 1
		}

		if err := cmd.Wait(); err != nil {
			log.Println(err)
			return
		}

		p.client.QueryResults(results)
	}
}

func (p *Plugin) Action(a sdk.Action) {
	cmd := exec.Command(".\\vendor\\AudioEndPointController\\Release\\EndPointController.exe", a.QueryResult.Data.(string))
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stderr)

	s := plugin.NewServer()
	s.Register(NewPlugin())
	s.Serve()
}
