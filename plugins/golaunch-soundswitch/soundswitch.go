package main

import (
	"encoding/json"
	"fmt"
	sdk "golaunch/sdk/go"
	"golaunch/sdk/go/plugin"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	aepFormat = `{"index":%d,"friendlyName":"%ws","state":%d,"default":%d,"description":"%ws","interfaceFriendlyName":"%ws","deviceID":"%ws"}`
)

type Device struct {
	Index                 int    `json:"index"`
	FriendlyName          string `json:"friendlyName"`
	State                 int    `json:"state"` //active or not
	Default               int    `json:"default"`
	Description           string `json:"description"`
	InterfaceFriendlyName string `json:"interfaceFriendlyName"`
	DeviceID              string `json:"deviceID"`
}

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
	if strings.HasPrefix(q, "soundswitch") {
		cmd := exec.Command(".\\vendor\\AudioEndPointController\\Release\\EndPointController.exe", "-f", aepFormat)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println(err)
			return
		}

		if err := cmd.Start(); err != nil {
			log.Println(err)
			return
		}

		decoder := json.NewDecoder(stdout)
		var results []sdk.QueryResult
		for {
			var device Device
			if err := decoder.Decode(&device); err == io.EOF {
				break
			} else if err != nil {
				log.Println(err)
			}

			result := sdk.QueryResult{
				Program: sdk.Program{
					Icon: p.metadata.Icon,
				},
				Title: device.FriendlyName,
				//Subtitle: fmt.Sprintf("Audio Device %d", device.Index),
				Query: q,
				ID:    p.metadata.ID,
				Score: -1,
				Data:  fmt.Sprintf("%d", device.Index),
			}

			if device.Default == 1 {
				result.Subtitle = "Default"
			}

			results = append(results, result)
		}

		if err := cmd.Wait(); err != nil {
			log.Println(err)
			return
		}

		p.client.Call("queryResults", results)
		return
	}

	p.client.Call("noQueryResults", nil)
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
