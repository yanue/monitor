/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : monitor.go
 Time    : 2018/11/12 11:38
 Author  : yanue
 
 - 
 
------------------------------- go ---------------------------------*/

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tabalt/gracehttp"
	"github.com/yanue/monitor/cmd"
	"github.com/yanue/monitor/service"
	"log"
	"net/http"
	"os"
)

var hub *service.ServerHub

func main() {
	// reload config
	if len(os.Args) > 1 && os.Args[1] == "reload" {
		cmd.ReloadConfig()
		return
	}

	cmdr, err := cmd.ParseCmd()

	if err != nil {
		if len(err.Error()) > 0 {
			fmt.Println("error:", err)
		}
		return
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	if cmdr.RunMode == cmd.RUN_MODE_SERVER {
		hub = service.NewService(cmdr)
		router.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Welcome Gin Server")
		})

		router.GET("/join", hub.JoinServer)
		router.GET("/state", hub.State)
	} else {
		cli := service.NewClient(cmdr)

		// join server
		cli.JoinServer()

		router.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		router.GET("/join", cli.JoinCustom)
		router.GET("/install", cli.InstallMonit)
		router.GET("/config", cli.ConfigMonit)
		router.GET("/state", cli.GetState)
		router.GET("/run", cli.RunCmd)
	}

	// write pid file
	cmd.WritePidFile()

	err = gracehttp.ListenAndServe(cmdr.ListenAddr, router)
	if err != nil {
		log.Fatalf("listener error: %v", err)
	}
}
