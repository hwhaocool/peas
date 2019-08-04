package task

import (
	"net/http"
)

//TaskReq 任务请求，或者说 任务体
type TaskReq struct {
	N       int         `json:"n"`
	C       int         `json:"c"`
	Timeout int         `json:"timeout"`
	Method  string      `json:"method"`
	URL     string      `json:"url"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body"`
}

type ResponseModel struct {
	Error         string           `json:"error,omitempty"`
	Status        string           `json:"status"`
	StatusCode    int              `json:"statusCode"`
	ContentLength int64            `json:"contentLength"`
	Proto         string           `json:"proto"`
	Headers       http.Header      `json:"headers"`
	Cookies       []*http.Cookie   `json:"cookies"`
	Body          string           `json:"body"`
	Code          string           `json:"code"`
	Duration      ResponseDuration `json:"duration"`
}

type ResponseDuration struct {
	DNS    string `json:"dns"`
	Conn   string `json:"conn"`
	Req    string `json:"req"`
	Res    string `json:"res"`
	Delay  string `json:"delay"`
	Finish string `json:"finish"`
}
