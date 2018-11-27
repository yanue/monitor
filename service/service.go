/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : service.go
 Time    : 2018/11/12 13:48
 Author  : yanue
 
 - 
 
------------------------------- go ---------------------------------*/

package service

import (
	"github.com/yanue/monitor/cmd"
	"sync"
)

type Client struct {
	id     string
	ip     string
	name   string
	status int
}

type Hub struct {
	Clients []*Client
	mux     *sync.Mutex
}

func NewService(cmdr *cmd.Cmd) *Hub {
	mux := &sync.Mutex{}
	clients := make([]*Client, 0)
	hub := &Hub{
		mux:     mux,
		Clients: clients,
	}
	return hub
}

func (this *Hub) AddClient(name, ip string, status int) {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.Clients = append(this.Clients, &Client{
		ip:     ip,
		name:   name,
		status: status,
	})
}

func (this *Hub) DelClient() {
	this.mux.Lock()
	defer this.mux.Unlock()
}

func (this *Hub) get() {

}
