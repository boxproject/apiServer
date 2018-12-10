// Copyright 2018. bolaxy.org authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		 http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"os"
	"os/exec"
	"github.com/boxproject/apiServer/config"
	"github.com/boxproject/apiServer/routers"
	"github.com/boxproject/apiServer/rpc"
	"github.com/boxproject/apiServer/timer"
	"github.com/gin-gonic/gin"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := newApp()
	app.Run(os.Args)
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "apiServer"
	app.Usage = "Serving BOX mobile APP"
	app.Author = "bolaxy.org"
	app.Copyright = "Copyright 2017-2018 The BOX Authors"
	app.Email = "devlop@bolaxy.org"
	app.Description = "The Server for BoxVault APP"
	app.Commands = []cli.Command{
		// 启动
		{
			Name:   "start",
			Usage:  "start the server",
			Action: StartCmd,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config,c",
					Usage: "Path of the config file",
					Value: "",
				},
				cli.StringFlag{
					Name:  "block-file,b",
					Usage: "Check point block number",
					Value: "",
				},
			},
		},
	}

	return app
}

func StopCmd(_ *cli.Context) error {
	_, err := exec.Command("sh", "-c", "pkill -SIGINT agent").Output()
	return err
}

func StartCmd(c *cli.Context) error {
	// Read the Config File
	cfg, err := config.LoadConfig()
	// init logger
	config.InitLogger()

	if err != nil {
		return err
	}
	config.Conf = cfg
	//init gprc
	go rpc.InitGrpc()
	go timer.TransferTimer()

	// init router
	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()
	// 获取证书
	r.StaticFile("server.cer", "./config/cer/server.cer")
	go r.Run(":10000")
	router := routers.InitRouter()
	router.RunTLS(cfg.Server.Port, "./config/cer/server.crt", "./config/cer/server.key")
	//router.Run(cfg.Server.Port)
	return nil
}
