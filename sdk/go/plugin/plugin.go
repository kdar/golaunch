package plugin

import (
	"encoding/json"
	"golaunch/sdk/go"
	"os"
)

type Plugin interface {
	Init(sdk.Metadata)
	Query(string)
	Action(sdk.Action)
}

type Client struct {
	enc *json.Encoder
}

func (c *Client) Call(method string, params interface{}) {
	msg := sdk.Response{
		Method: method,
		Params: params,
	}
	c.enc.Encode(msg)
}

func NewClient() *Client {
	return &Client{
		enc: json.NewEncoder(os.Stdout),
	}
}

type Server struct {
	dec *json.Decoder
	p   Plugin
}

func NewServer() *Server {
	return &Server{
		dec: json.NewDecoder(os.Stdin),
	}
}

func (s *Server) Register(p Plugin) {
	s.p = p
}

func (s *Server) Serve() {
	for {
		var v sdk.Request
		if err := s.dec.Decode(&v); err != nil {
			continue
		}

		switch v.Method {
		case "init":
			var metadata sdk.Metadata
			json.Unmarshal(v.Params, &metadata)
			s.p.Init(metadata)
		case "query":
			var query string
			json.Unmarshal(v.Params, &query)
			s.p.Query(query)
		case "action":
			var action sdk.Action
			json.Unmarshal(v.Params, &action)
			s.p.Action(action)
		}
	}
}
