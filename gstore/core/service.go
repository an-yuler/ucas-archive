package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

const (
	gMethod = "POST"
	gAPIUrl = "http://cloud.gstore.cn/api"
	gAction = "queryDB"
	gDBName = "jinrong"
	gAccessKeyId = "b91da25baa694d74a19524381d878451"
	gAccessSecret = "82B110F8122AF05BECABF2F1CB228885"
)

type QueryService struct {
	client *http.Client
	req *http.Request
}

func NewQueryService() (*QueryService, error) {
	w := &QueryService{}
	w.client = &http.Client{}

	var err error
	w.req, err = http.NewRequest(gMethod, gAPIUrl, nil)
	if err !=nil {
		return nil, err
	}

	q := w.req.URL.Query()
	q.Add("action", gAction)
	q.Add("dbName", gDBName)
	q.Add("access_secret", gAccessSecret)
	q.Add("accesskeyid", gAccessKeyId)
	q.Add("sparql", "")
	w.req.URL.RawQuery = q.Encode()

	return w, nil
}

// main loop
// call gstore api, then read body and decode json
func (w *QueryService) query(stat string, respvo *ResponseVO)(interface{}, error) {
	// change sparql
	q := w.req.URL.Query()
	q.Set("sparql", stat)
	w.req.URL.RawQuery = q.Encode()

	resp, err := w.client.Do(w.req)
	if err != nil {
		respvo.Code = StatusGStoreClientErr
		respvo.Status = StatusText(StatusGStoreClientErr)
		log.Println("client do error")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respvo.Code = StatusGStoreReadBodyErr
		respvo.Status = StatusText(StatusGStoreReadBodyErr)
		log.Println("ioutil read error")
		return nil, err
	}

	var f interface{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		respvo.Code = StatusGStoreJsonDecodeErr
		respvo.Status = StatusText(StatusGStoreJsonDecodeErr)
		return nil, err
	}

	return f, nil
}

// 查询两个公司之间的关联路径（2-hop）
func (w *QueryService) Query1(comp1, comp2 string, respvo *ResponseVO){
	// construct the sparsql
	stat := fmt.Sprintf("SELECT ?a WHERE { %s ?x ?a . ?a ?y %s . }", fullName(comp1), fullName(comp2))
	respvo.Query = append(respvo.Query, JsonSql{Id: 0, Sparql: stat})

	log.Println(stat)
	var f interface{}
	f, err := w.query(stat, respvo)
	if err != nil {
		log.Println("query error")
		return
	}

	// extract data
	chains := extract(&f)

	// output
	log.Println("Query1 Statement:")
	log.Println(stat)
	log.Println("Result:")

	respvo.Code = StatusOK
	respvo.Status = StatusText(StatusOK)
	for i, chain := range chains{
		buf := append(append([]string{comp1}, chain...), comp2)
		log.Println(strings.Join(buf, "->"))
		ch := JsonChain{Id: i, Chain: buf}
		respvo.Results = append(respvo.Results, ch)
	}
}

// 多层股权的穿透式查询，可以根据指定层数获得对应层级的股东
func (w *QueryService)Query2(comp string, hop int, respvo *ResponseVO) {
	// construct statement
	stat := constructStat2(comp, hop)

	// main loop
	var f interface{}
	f, err := w.query(stat, respvo)
	if err != nil {
		log.Println("query error")
		return
	}

	// extract data
	chains := extract(&f)

	// construct ResponseVO
	log.Println("Query2 Statement:")
	log.Println(stat)
	log.Println("Result:")

	respvo.Code = StatusOK
	respvo.Status = StatusText(StatusOK)
	respvo.Query = append(respvo.Query, JsonSql{Id: 0, Sparql: stat})
	for i, chain := range chains{
		buf := append([]string{comp}, chain...)
		log.Println(strings.Join(buf, "->"))
		ch := JsonChain{Id: i, Chain: buf}
		respvo.Results = append(respvo.Results, ch)
	}
}



// construct statement for query2
// eg.
// for hop = 3:
// SELECT ?a0 ?a1 ?a2  WHERE {  { $comp ?b0 ?a0.  }
//									UNION
// 								{ $comp ?b0 ?a0. ?a0 ?b1 ?a1. }
//									UNION
//								{ $comp ?b0 ?a0. ?a0 ?b1 ?a1. ?a1 ?b2 ?a2. }  }
func constructStat2(comp string, hop int) string {
	// generate the nth params: ?a$number
	param := func(number int) string{return fmt.Sprintf("?a%d", number)}
	// generate the nth operation: ?b$number
	op := func(number int) string{return fmt.Sprintf("?b%d", number)}

	var obuf, ibuf, ubuf bytes.Buffer

	// {%s} UNION {%s} UNION {%S}
	for i := 0; i < hop; i++{
		x, y, z := param(i-1), op(i), param(i)
		if i == 0 {
			x = fullName(comp)
		}else{
			ubuf.WriteString("UNION")
		}
		obuf.WriteString(fmt.Sprintf("%s ", z))
		ibuf.WriteString(fmt.Sprintf("%s %s %s. ", x, y, z))
		ubuf.WriteString(fmt.Sprintf(" { %s } ", ibuf.String()))
	}
	return fmt.Sprintf("SELECT %s WHERE { %s } ", obuf.String(), ubuf.String())
}

// 环形持股查询，判断两家公司是否存在环形持股现象，环形持股是指两家公司彼此持有对方的股份
func (w *QueryService)Query3(comp1, comp2 string, respvo *ResponseVO){

	// construct statement
	const _MAX_HOP = 5
	stat1, stat2 := constructStat3(comp1, _MAX_HOP), constructStat3(comp2, _MAX_HOP)

	f1, err := w.query(stat1, respvo)
	if err != nil {
		return
	}

	f2, err := w.query(stat2, respvo)
	if err != nil {
		return
	}

	respvo.Code = StatusOK
	respvo.Status = StatusText(StatusOK)
	respvo.Query =  []JsonSql{{0, stat1}, {1, stat2}}

	// extract data
	chains1 := extract(&f1)
	chains2 := extract(&f2)

	// output
	log.Println("Query3 Statement:")
	log.Println(stat1)
	log.Println(stat2)
	log.Println("Result:")

	if len(chains1) == 0 || len(chains2) == 0 {
		log.Println("No Cycle")
		respvo.Code = StatusNoCycle
		respvo.Status = StatusText(StatusNoCycle)
		return
	}

	// TODO: Use map to improve performance
	var ok1, ok2 bool
	var c1, c2 []string
	for _, chain := range chains1{
		for i, el := range chain {
			if el == comp2 {
				c1 = chain[:i]
				ok1 = true
				break
			}
		}
		if ok1 {
			break
		}
	}

	for _, chain := range chains2{
		for i, el := range chain {
			if el == comp1 {
				c2 = chain[:i]
				ok2 = true
				break
			}
		}
		if ok2 {
			break
		}
	}

	if !ok1 || !ok2 {
		log.Println("No Cycle")
		respvo.Code = StatusNoCycle
		respvo.Status = StatusText(StatusNoCycle)
		return
	}
	c1 = append(append([]string{comp1}, c1...), comp2)
	c2 = append(append([]string{comp2}, c2...), comp1)

	log.Println("Cycle Exist!")
	log.Println(strings.Join(c1, "->"))
	log.Println(strings.Join(c2, "->"))
	respvo.Results = []JsonChain{{0, c1}, {1, c2} }
}

// construct statement for query3
// eg.
// 	for hop = 3: SELECT ?a0 ?a1 ?a2 WHERE { $comp ?b0 ?a0 . ?a0 ?b1 ?a1 . ?a1 ?b2 ?a2 . }
func constructStat3(comp string, hop int) string {
	// generate the nth params: ?a$number
	param := func(number int) string{return fmt.Sprintf("?a%d", number)}
	// generate the nth operation: ?b$number
	op := func(number int) string{return fmt.Sprintf("?b%d", number)}

	var obuf, ibuf bytes.Buffer
	for i := 0; i < hop; i++{
		x, y, z := param(i-1), op(i), param(i)
		if i == 0 {
			x = fullName(comp)
		}
		obuf.WriteString(fmt.Sprintf("%s ", z))
		ibuf.WriteString(fmt.Sprintf("%s %s %s. ", x, y, z))
	}
	return fmt.Sprintf("SELECT %s WHERE { %s } ", obuf.String(), ibuf.String())
}

// extract the data that represented by vars from response body
// eg.
// 	for query1 : [[a0], [a0], [a0]...]
//  for query2 : [[a0, a1, a2], [a0, a1, a2]...]
func extract(body interface{}) [][]string{
	f := reflect.ValueOf(body).Elem().Interface()

	d := f.(map[string]interface{})["data"]
	vars := d.(map[string]interface{})["head"].(map[string]interface{})["vars"].([]interface{})
	chains := d.(map[string]interface{})["results"].(map[string]interface{})["bindings"].([]interface{})

	m, _ := len(chains), len(vars)
	ans := make([][]string, m)
	for i, chain := range chains {
		for _, v := range vars {
			if _, ok := chain.(map[string]interface{})[v.(string)]; !ok {
				continue
			}
			ans[i] = append(ans[i], shortName(chain.(map[string]interface{})[v.(string)].(map[string]interface{})["value"].(string)))
		}
	}
	return ans
}

// convert the short company name to the full name.
// eg.
//   "安邦财产保险股份有限公司" -> "<file:///F:/d2r-server-0.7/holder8.nt#holder_copy/安邦财产保险股份有限公司>"
func fullName(sname string) string{
	return fmt.Sprintf("<file:///F:/d2r-server-0.7/holder8.nt#holder_copy/%s>", sname)
}

// Contrary to fullName(), convert the full name to the short name.
// Note that there is no "<" and ">" here.
// eg.
//	"file:///F:/d2r-server-0.7/holder8.nt#holder_copy/安邦财产保险股份有限公司" -> "安邦财产保险股份有限公司"
func shortName(fname string) string {
	return fname[strings.LastIndex(fname, "/")+1:]
}