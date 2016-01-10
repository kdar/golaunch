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

func (c *Client) QueryResults(results []sdk.QueryResult) {
	msg := sdk.Response{
		Result: results,
	}
	c.enc.Encode(msg)
}

func (c *Client) Call(method string, params ...string) {
	data, _ := json.Marshal(params)
	msg := sdk.Request{
		Method: method,
		Params: data,
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
