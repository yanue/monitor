/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : base.go
 Time    : 2018/11/27 14:14
 Author  : yanue
 
 - 
 
------------------------------- go ---------------------------------*/

package service

import (
	"github.com/yanue/monitor/cmd"
	"sync"
)

type JsonResult struct {
	Result int    `json:"code"`
	Msg    string `json:"msg"`
	Data   string `json:"data"`
}

type Client struct {
	Id     int    `json:"id"`
	Ip     string `json:"ip"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type ServerHub struct {
	Clients []*Client
	mux     *sync.Mutex
	lastId  int
	cmd     *cmd.Cmd
}

type ClientHub struct {
	Client *Client
	cmd    *cmd.Cmd
}

func NewService(cmdr *cmd.Cmd) *ServerHub {
	mux := &sync.Mutex{}
	clients := make([]*Client, 0)
	hub := &ServerHub{
		mux:     mux,
		Clients: clients,
		cmd:     cmdr,
	}
	hub.loadConfig()

	go hub.checkState()

	return hub
}

func NewClient(cmdr *cmd.Cmd) *ClientHub {
	return &ClientHub{
		Client: &Client{},
		cmd:    cmdr,
	}
}
