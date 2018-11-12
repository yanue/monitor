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
	"log"
	"net/http"
	"os"
)

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

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	// write pid file
	cmd.WritePidFile()

	err = gracehttp.ListenAndServe(cmdr.ListenAddr, router)
	if err != nil {
		log.Fatalf("listener error: %v", err)
	}
}
