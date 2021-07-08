package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type RequestVO struct {
	Params []string `json:"params""`
}

type JsonSql struct {
	Id int `json:"id"`
	Sparql string `json:"sparql"`
}

type JsonChain struct {
	Id int `json:"id"`
	Chain []string `json:"chain"`
}

type ResponseVO struct {
	Code int            `json:"code"`
	Status  string      `json:"http""`
	Query   []JsonSql   `json:"query"`
	Results []JsonChain `json:"results""`
}


type GServer struct {
	qs *QueryService
	http.Handler
}

func decodeRequest(request *http.Request, req *RequestVO, resp *ResponseVO) error{
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		resp.Code = StatusReadBodyErr
		resp.Status = StatusText(StatusReadBodyErr)
		return err
	}

	err = json.Unmarshal(body, req)
	if err != nil {
		resp.Code = StatusJsonDecodeErr
		resp.Status = StatusText(StatusJsonDecodeErr)
		return err
	}
	return nil
}

//// write error response
//func writeError(w http.ResponseWriter, err error) {
//	log.Println(err)
//	json.NewEncoder(w).Encode(ErrResponse)
//}

// write success response
func writeResponse(w http.ResponseWriter, resp *ResponseVO) {
	json.NewEncoder(w).Encode(resp)
}

func (g *GServer) Query1Handle(w http.ResponseWriter, request *http.Request) {
	reqvo, respvo := new(RequestVO), new(ResponseVO)

	err := decodeRequest(request, reqvo, respvo)
	if err != nil {
		writeResponse(w, respvo)
		return
	}

	if len(reqvo.Params) != 2 {
		respvo.Code = StatusParamsErr
		respvo.Status = StatusText(StatusParamsErr)
		writeResponse(w, respvo)
		return
	}

	g.qs.Query1(reqvo.Params[0], reqvo.Params[1], respvo)

	writeResponse(w, respvo)
}


func (g *GServer) Query2Handle(w http.ResponseWriter, request *http.Request) {

	reqvo, respvo := new(RequestVO), new(ResponseVO)
	err := decodeRequest(request, reqvo, respvo)
	if err != nil {
		writeResponse(w, respvo)
		return
	}

	if len(reqvo.Params) != 2 {
		respvo.Code = StatusParamsErr
		respvo.Status = StatusText(StatusParamsErr)
		writeResponse(w, respvo)
		return
	}

	hop, err := strconv.Atoi(reqvo.Params[1])
	if err != nil {
		respvo.Code = StatusParamsErr
		respvo.Status = StatusText(StatusParamsErr)
		writeResponse(w, respvo)
		return
	}

	g.qs.Query2(reqvo.Params[0], hop, respvo)
	writeResponse(w, respvo)
}

func (g *GServer) Query3Handle(w http.ResponseWriter, request *http.Request) {
	reqvo, respvo := new(RequestVO), new(ResponseVO)
	err := decodeRequest(request, reqvo, respvo)
	if err != nil {
		writeResponse(w, respvo)
		return
	}

	if len(reqvo.Params) != 2 {
		respvo.Code = StatusParamsErr
		respvo.Status = StatusText(StatusParamsErr)
		writeResponse(w, respvo)
		return
	}

	g.qs.Query3(reqvo.Params[0], reqvo.Params[1], respvo)
	writeResponse(w, respvo)
}

// 中间件
// 1. 拦截option请求
// 2. 设置跨域头
// 3. 设置content-type
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//cros
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers","Content-Type,AccessToken,X-CSRF-Token,Authorization,Token")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials","true")
		method := r.Method
		if method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// contype-type
		w.Header().Set("Content-Type", "application/json") //返回数据格式是json

		//end
		handler.ServeHTTP(w, r)
	})
}

func NewGServer() *GServer {
	g := new(GServer)

	g.qs, _ = NewQueryService()

	router := http.NewServeMux()
	router.HandleFunc("/query1", g.Query1Handle)
	router.HandleFunc("/query2", g.Query2Handle)
	router.HandleFunc("/query3", g.Query3Handle)

	g.Handler = setupGlobalMiddleware(router)
	return g
}