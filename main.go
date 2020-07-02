package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"agent/cron"
	"agent/g"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func prepare() {
	//设置线程
	runtime.GOMAXPROCS(runtime.NumCPU())
	//设置日志
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// init 方法是包的初始化方法,在main方法之前执行
func init() {
	// 准备函数
	prepare()

	// 加载配置文件
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	// 命令行参数接口
	flag.Parse()

	handleVersion(*version)
	handleHelp(*help)
	handleConfig(*cfg)

	// 初始化客户端
	// backend.InitClients(g.Config.Web.Addrs)
	g.Init()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}
func main() {
	// go cron.Push()
	// cron.StartCheck()
	http.HandleFunc("/test", indexHandler)
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		// sc.Lock()
		// conf := sc.C
		// sc.Unlock()
		probeHandler(w, r) //, conf, logger, rh)
	})
	// http.ListenAndServe(g.Config.ListenPort, nil)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html>
    <head><title>Blackbox Exporter</title></head>
    <body>
    <h1>Blackbox Exporter</h1>
    <p><a href="probe?target=prometheus.io&module=http_2xx">Probe prometheus.io for http_2xx</a></p>
    <p><a href="probe?target=prometheus.io&module=http_2xx&debug=true">Debug probe prometheus.io for http_2xx</a></p>
    <p><a href="metrics">Metrics</a></p>
    <p><a href="config">Configuration</a></p>
    <h2>Recent Probes</h2>
    <table border='1'><tr><th>Module</th><th>Target</th><th>Result</th><th>Debug</th>`))

		w.Write([]byte(`</table></body>
    </html>`))
	})

	srv := http.Server{Addr: g.Config.ListenPort}
	err := srv.ListenAndServe()
	log.Println(err)
}

func handleVersion(displayVersion bool) {
	if displayVersion {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
}

func handleHelp(displayHelp bool) {
	if displayHelp {
		flag.Usage()
		os.Exit(0)
	}
}

func handleConfig(configFile string) {
	err := g.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}

// func metrics(w http.ResponseWriter, r *http.Request) {
// 	cron.StartCheck(w, r)
// }

func probeHandler(w http.ResponseWriter, r *http.Request) {
	//, c *config.Config, logger log.Logger, rh *resultHistory) {

	start := time.Now()

	/////////////////////////
	cron.StartCheck(w, r) //prober(ctx, target, module, registry, sl)
	duration := time.Since(start).Seconds()
	log.Println("msg", "Probe succeeded", "duration_seconds", duration)
	// promhttp.Handler()
	// registry.MustRegister(promhttp.Handler())
	h := promhttp.HandlerFor(g.Registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
