package main

import (
	"fmt"
	"testing"
)


func TestConstructStat2(t *testing.T) {
	stat := constructStat2("上海天发投资有限公司", 3)
	fmt.Println(stat)
}

func TestConstructStat3(t *testing.T) {
	stat := constructStat3("上海天发投资有限公司", 3)
	fmt.Println(stat)
}


type QueryCase struct {
	c1, c2 string
	hop int
}

func TestQuery(t *testing.T) {
	service, _ := NewQueryService()

	c := QueryCase{c1: "招商局轮船股份有限公司", c2: "招商银行股份有限公司", hop: 3}

	respvo := new(ResponseVO)
	t.Run("query1", func(t *testing.T) {
		service.Query1(c.c1, c.c2, respvo)
	})

	t.Run("query2", func(t *testing.T) {
		service.Query2(c.c1, c.hop, respvo)
	})

	t.Run("query3", func(t *testing.T) {
		service.Query3("A", "C", respvo)
	})
}

