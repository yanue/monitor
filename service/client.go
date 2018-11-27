/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : client.go
 Time    : 2018/11/27 14:13
 Author  : yanue
 
 - 
 
------------------------------- go ---------------------------------*/

package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
)

func (this *ClientHub) JoinServer() {
	ip := this.getIntranetIp()
	if ip == "" {
		panic("get ip error:")
	}

	name, err := os.Hostname()
	if err != nil {
		panic("get hostname error:" + err.Error())
	}

	val := url.Values{}
	val.Set("ip", ip)
	val.Set("name", name)

	uri := this.cmd.JoinAddr + "?" + val.Encode()
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Println("join server error:", err.Error())
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}

	js := &JsonResult{}
	err = json.Unmarshal(body, js)
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}

	if js.Result == 0 {
		println("error:", js.Msg)
	} else {
		println("done", js.Data)
	}
}

func (this *ClientHub) JoinCustom(c *gin.Context) {
	res := &JsonResult{}
	c.JSON(200, res)
}

func (this *ClientHub) RunCmd(c *gin.Context) {
	res := &JsonResult{}
	c.JSON(200, res)
}

func (this *ClientHub) InstallMonit(c *gin.Context) {
	file := path.Join(this.cmd.WorkDir, "bin", "monit.sh")

	cmd := exec.Command("sh", "-c", file)

	res := &JsonResult{}
	msg, err := cmd.CombinedOutput()
	res.Msg = string(msg)

	if err != nil {
		res.Result = 0
	} else {
		res.Result = 1
	}

	c.JSON(200, res)
}

func (this *ClientHub) ConfigMonit(c *gin.Context) {
	res := &JsonResult{}
	c.JSON(200, res)
}

func (this *ClientHub) GetState(c *gin.Context) {
	res := &JsonResult{}
	c.JSON(200, res)
}

func (this *ClientHub) getIntranetIp() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
