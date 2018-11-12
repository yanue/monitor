/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : cmd.go
 Time    : 2018/11/12 11:43
 Author  : yanue
 
 - 
 
------------------------------- go ---------------------------------*/

package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const RUN_MODE_SERVER = "server"
const RUN_MODE_CLIENT = "client"

const keyConfig = "config"
const keyListen = "listen"
const keyMode = "mode"
const keyJoin = "join"
const keyPath = "path"

const defaultConfigFile = "/etc/monitor.conf"

type Cmd struct {
	ConfigFile string
	ListenAddr string
	RunMode    string
	JoinAddr   string
	WorkDir    string
}

func usage() {
	fmt.Printf(`Usage:
  %s [-config]
  %s [-listen] [-mode] [-join] [-path]

Options:
`, os.Args[0], os.Args[0])
	flag.PrintDefaults()
	fmt.Println("")
	fmt.Printf("  --help - usage help.\n")
	fmt.Printf("  reload - reload config.\n")
	fmt.Println("")
}

func ParseCmd() (cmd *Cmd, err error) {
	cmd = new(Cmd)

	flag.StringVar(&cmd.ConfigFile, keyConfig, "", "ini config file path, default is: "+defaultConfigFile)
	flag.StringVar(&cmd.ListenAddr, keyListen, ":9090", "listen address.")
	flag.StringVar(&cmd.RunMode, keyMode, "", "run mode: server or client.")
	flag.StringVar(&cmd.JoinAddr, keyJoin, "", "server address to join. (when run mode is client)")
	flag.StringVar(&cmd.WorkDir, keyPath, "", "work path.")

	flag.Usage = usage
	flag.Parse()

	// or has default config file
	var configFile = ""

	// if user does not set flags, check default config file
	if flag.NFlag() == 0 {
		// default config exist
		if _, e := os.Stat(defaultConfigFile); !os.IsNotExist(e) {
			configFile = defaultConfigFile
		} else {
			err = errors.New("missing args")
			usage()
			return
		}
	} else {
		// if config file has set with -config
		if len(cmd.ConfigFile) > 0 {
			if _, e := os.Stat(cmd.ConfigFile); os.IsNotExist(e) {
				err = errors.New("config file is not exist")
				return
			}
			configFile = cmd.ConfigFile
		}
	}

	if len(configFile) > 0 {
		// parse ini file
		fmt.Println("load config:", configFile)
		cfg, e := ini.Load(configFile)
		if e != nil {
			err = errors.New("parse config file error:" + e.Error())
			return
		}

		cmd.ListenAddr = cfg.Section("").Key(keyListen).String()
		cmd.RunMode = cfg.Section("").Key(keyMode).String()
		cmd.JoinAddr = cfg.Section("").Key(keyJoin).String()
		cmd.WorkDir = cfg.Section("").Key(keyPath).String()
	}

	if cmd.RunMode != RUN_MODE_SERVER && cmd.RunMode != RUN_MODE_CLIENT {
		err = errors.New("missing run mode: server or client")
		return
	}

	if len(cmd.WorkDir) == 0 {
		err = errors.New("missing work dir")
		return
	}

	if cmd.RunMode == RUN_MODE_CLIENT && len(cmd.JoinAddr) == 0 {
		err = errors.New("missing json addr when run at client mode")
		return
	}

	return
}

func (cmd *Cmd) String() string {
	return fmt.Sprintf("[ConfigFile -> '%s',ListenAddr -> '%s', RunMode -> '%s', JoinAddr -> '%s', WorkDir -> '%s']",
		cmd.ConfigFile,
		cmd.ListenAddr,
		cmd.RunMode,
		cmd.JoinAddr,
		cmd.WorkDir)
}

func ReloadConfig() {
	fmt.Println("reload config")
	pidStr, err := getPid()
	if err != nil {
		fmt.Println("get pid error:", err.Error())
		return
	}
	pidNum, err := strconv.Atoi(pidStr)
	if err != nil {
		fmt.Println("get pid error:", err.Error())
		return
	}

	if pidNum < 1024 {
		fmt.Println("run at a unsafe port, please restart manual")
		return
	}

	fmt.Println("cmd", "kill", "-USR2", fmt.Sprintf("%d", pidNum))
	cmd := exec.Command("kill", "-USR2", fmt.Sprintf("%d", pidNum))
	cmd.Stdout = os.Stdout
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

// get pid file
func getPid() (string, error) {
	pidFile := getPidFile()

	// re-open file
	dat, err := ioutil.ReadFile(pidFile)
	return string(dat), err
}

// get pid file
func getPidFile() string {
	runName := strings.Split(os.Args[0], "/")
	pidFile := "/tmp/" + runName[len(runName)-1] + ".pid"
	return pidFile
}

// Write a pid file
func WritePidFile() {
	pidFile := getPidFile()
	if err := ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0664); err != nil {
		fmt.Println("WritePidFile err", err)
	}
}

// delete pid file
func DeletePidFile() {
	// delete file
	os.Remove(getPidFile())
}
