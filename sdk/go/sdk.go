package sdk

import (
	"encoding/json"
	//"golaunch/flatapi"

	//	flatbuffers "github.com/google/flatbuffers/go"
)

type Request struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type Response struct {
	Result []QueryResult `json:"result"`
}

type Program struct {
	Path  string `json:"path"`
	Image string `json:"image"`
	Usage int    `json:"usage"`
}

type QueryResult struct {
	Program
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Query    string `json:"query"`
	// ID of plugin returning this query results
	ID      string `json:"id"`
	Score   int    `json:"score"`
	LowName string `json:"-"`
	// extra data for plugins
	Data interface{} `json:"data"`
}

type Metadata struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Version     string `json:"version"`
	Type        string `json:"type"`
	Main        string `json:"main"`
	Icon        string `json:"_icon"`
	// where to store your app's data and settings
	AppData string `json:"_appdata"`
}

// type Api struct {
// 	builder *flatbuffers.Builder
// }
//
// func NewApi() *Api {
// 	return &Api{builder: flatbuffers.NewBuilder(0)}
// }
//
// func (a *Api) Reset() {
// 	a.builder.Reset()
// }
//
// func (a *Api) CreateQueryResult(title, subtitle, image, query string, score int) flatbuffers.UOffsetT {
// 	imagev := a.builder.CreateString(image)
// 	titlev := a.builder.CreateString(title)
// 	subtitlev := a.builder.CreateString(subtitle)
// 	queryv := a.builder.CreateString(query)
//
// 	flatapi.QueryResultStart(a.builder)
// 	flatapi.QueryResultAddImage(a.builder, imagev)
// 	flatapi.QueryResultAddTitle(a.builder, titlev)
// 	flatapi.QueryResultAddSubtitle(a.builder, subtitlev)
// 	flatapi.QueryResultAddScore(a.builder, int32(score))
// 	flatapi.QueryResultAddQuery(a.builder, queryv)
// 	queryresult := flatapi.QueryResultEnd(a.builder)
//
// 	flatapi.ResultStart(a.builder)
// 	flatapi.ResultAddResultType(a.builder, flatapi.AnyResultQueryResult)
// 	flatapi.ResultAddResult(a.builder, queryresult)
// 	return flatapi.ResultEnd(a.builder)
// }
//
// func (a *Api) CreateResponse(id string, results []flatbuffers.UOffsetT) []byte {
// 	idv := a.builder.CreateString(id)
//
// 	flatapi.ResponseStartResultVector(a.builder, len(results))
// 	for x := 0; x < len(results); x++ {
// 		a.builder.PrependUOffsetT(results[x])
// 	}
// 	rv := a.builder.EndVector(len(results))
//
// 	flatapi.ResponseStart(a.builder)
// 	flatapi.ResponseAddId(a.builder, idv)
// 	flatapi.ResponseAddResult(a.builder, rv)
// 	responseOffset := flatapi.ResponseEnd(a.builder)
//
// 	a.builder.Finish(responseOffset)
//
// 	return a.builder.Bytes[a.builder.Head():]
// }
//
// // buf, _ := ioutil.ReadFile("fb.bin")
// // response := flatapi.GetRootAsResponse(buf, 0)
// // fmt.Println(string(response.Id()))
// // var result flatapi.Result
// // response.Result(&result, 0)
// // var qresult flatapi.QueryResult
// // var qtable flatbuffers.Table
// // result.Result(&qtable)
// // qresult.Init(buf, qtable.Pos)
// // fmt.Println(string(qresult.Title()))
// //
// // return
