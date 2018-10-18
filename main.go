package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/gommon/color"
	"github.com/urfave/cli"
	"github.com/xuender/go-utils"

	"./rich"
)

//go:generate go run $GOROOT/src/crypto/tls/generate_cert.go --host=localhost
//go:generate mv cert.pem key.pem keys
//go:generate go-bindata -nomemcopy -pkg keys -o ./keys/bindata.go keys/cert.pem keys/key.pem
//go:generate go-bindata -nomemcopy -pkg rich -o ./rich/bindata.go www/...

func main() {
	app := cli.NewApp()
	app.Name = "Go Rich"
	app.Usage = "致力服务小商家"
	app.Version = "v0.1.12"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "address,a",
			Value: "6181",
			Usage: "访问地址的端口号",
		},
		cli.StringFlag{
			Name:  "protocol,p",
			Value: "http",
			Usage: "访问协议 http, TLS, AutoTLS",
		},
		cli.StringFlag{
			Name:  "cert,c",
			Usage: "TLS证书文件, TLS必须",
		},
		cli.StringFlag{
			Name:  "key,k",
			Usage: "TLS密钥文件, TLS必须",
		},
		cli.StringFlag{
			Name:  "db,b",
			Value: "db",
			Usage: "数据库目录",
		},
		cli.StringFlag{
			Name:  "temp,t",
			Value: "temp",
			Usage: "临时文件目录",
		},
		cli.StringFlag{
			Name:  "logfile,l",
			Value: "rich.log",
			Usage: "日志输出文件",
		},
		cli.BoolFlag{
			Name:  "zip-backup,z",
			Usage: "数据库备份",
		},
		cli.BoolFlag{
			Name:  "open,o",
			Usage: "启动浏览器显示QR码",
		},
		cli.BoolFlag{
			Name:  "develop-mode,d",
			Usage: "开发模式: 静态资源读自www目录,支持跨域访问,访问日志显示",
		},
		cli.BoolFlag{
			Name:  "upgrade,u",
			Usage: "自动升级",
		},
		cli.BoolFlag{
			Name:  "reset,r",
			Usage: "帐目索引重置",
		},
	}
	app.Action = runAction
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
func runAction(c *cli.Context) error {
	// 备份
	if c.Bool("z") {
		bak := time.Now().Format(rich.DayFormat) + ".zip"
		if err := rich.Zip(c.String("b"), bak); err != nil {
			log.Println(err.Error())
		} else {
			log.Println("备份:", bak)
		}
	}

	address := c.String("a")
	if !strings.HasPrefix(address, ":") {
		address = ":" + address
	}
	web := rich.Web{
		Temp:    c.String("t"),
		Db:      c.String("b"),
		Dev:     c.Bool("d"),
		LogFile: c.String("l"),
		App:     rich.NewApp(c.App),
	}
	if url, err := rich.GetURL(address, !strings.EqualFold(c.String("p"), "http")); err == nil {
		web.URL = url
	} else {
		return err
	}
	// 扫码提示
	color.Println(" ", color.Blue(web.URL, color.U))
	color.Println("  使用", color.Green("微信"), "扫码，即刻使用")
	color.Println(" ", color.Red("《  "+c.App.Name+"  》", color.B, color.WhtBg), color.Cyan(c.App.Version, color.In))
	color.Println("      ", color.Red(c.App.Usage, color.U))
	rich.PrintQR(web.URL)
	color.Println("请保持本窗口处于打开状态，", color.Red("CTRL+C", color.B), "退出")
	// 初始化
	if err := web.Init(); err != nil {
		return err
	}
	if c.Bool("r") {
		if err := web.Reset(); err != nil {
			return err
		}
	}
	// 升级
	if c.Bool("u") {
		go func() {
			if url, tag, err := rich.LastURL("xuender", "go-rich", c.App.Version); err == nil && tag != c.App.Version {
				log.Println("更新版本:", tag, "网址:", url)
				rich.Download(url, c.String("t"))
				os.Rename(filepath.Join(c.String("t"), "go-rich"), "go-rich")
				log.Println("更新完毕")
			}
		}()
	}
	// 退出
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	// 运行
	switch strings.ToLower(c.String("p")) {
	case "tls":
		go web.StartTLS(address, c.String("c"), c.String("k"))
	case "autotls":
		go web.StartAutoTLS(address)
	default:
		go web.Start(address)
	}
	// 打开浏览器
	if c.Bool("o") {
		utils.Open(web.URL + "/qr")
	}
	fmt.Println(<-quitChan)
	web.Close()
	return nil
}
