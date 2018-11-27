/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : service.go
 Time    : 2018/11/12 13:48
 Author  : yanue
 
 - 
 
------------------------------- go ---------------------------------*/

package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path"
	"time"
)

func (this *ServerHub) AddClient(name, ip, key string, status int) {
	this.mux.Lock()
	defer this.mux.Unlock()

	var index = -1
	for idx, item := range this.Clients {
		if item.Ip == ip {
			index = idx
		}
	}

	if index > -1 {
		this.Clients[index].Name = name
		this.Clients[index].Status = status
	} else {
		this.lastId += 1

		cli := &Client{
			Id:     this.lastId,
			Ip:     ip,
			Name:   name,
			Status: status,
		}

		this.Clients = append(this.Clients, cli)
	}

	this.saveConfig()
}

func (this *ServerHub) saveConfig() {
	file := path.Join(this.cmd.WorkDir, "server.json")

	clis := this.Clients

	js, err := json.Marshal(clis)
	if err != nil {
		fmt.Println("save file error", err.Error())
		return
	}

	ioutil.WriteFile(file, js, 0755)
}

func (this *ServerHub) loadConfig() {
	file := path.Join(this.cmd.WorkDir, "server.json")
	body, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("load file error", err.Error())
		return
	}

	err = json.Unmarshal(body, &this.Clients)
	if err != nil {
		fmt.Println("save file error", err.Error())
		return
	}

	// get last id
	for _, val := range this.Clients {
		if val.Id > this.lastId {
			this.lastId = val.Id
		}
	}
}

func (this *ServerHub) DelClient() {
	this.mux.Lock()
	defer this.mux.Unlock()
}

func (this *ServerHub) checkState() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		this.checkAllState()
	}
}

func (this *ServerHub) checkAllState() {
	this.mux.Lock()
	defer this.mux.Unlock()

	for _, item := range this.Clients {
		this.getState(item)
	}
}

func (this *ServerHub) getState(item *Client) {
	// todo
	fmt.Println("get state", item.Name)
}

func (this *ServerHub) JoinServer(c *gin.Context) {
	ip := c.DefaultQuery("ip", "")
	name := c.DefaultQuery("name", "")
	key := c.DefaultQuery("key", "")

	this.AddClient(name, ip, key, 1)

	c.JSON(http.StatusOK, &JsonResult{Result: 1, Msg: "", Data: ""})
}

func (this *ServerHub) State(c *gin.Context) {
	c.JSON(http.StatusOK, &JsonResult{Result: 1, Msg: "", Data: ""})
}
