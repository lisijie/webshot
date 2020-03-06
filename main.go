package main

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/netutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var (
	// 默认设置
	addr      = ":80"                   // http监听地址
	debug     = false                   // 是否debug
	remoteUrl = "http://127.0.0.1:9222" // headless-shell 地址
	connLimit = 20                      // 并发数限制
	timeout   = time.Second * 30        // 处理超时设置
	startTime = time.Now()
)

func main() {
	if v, ok := os.LookupEnv("HTTP_ADDR"); ok && v == "true" {
		addr = v
	}
	if v, ok := os.LookupEnv("DEBUG"); ok && v == "true" {
		debug = true
	}
	if v, ok := os.LookupEnv("REMOTE_URL"); ok {
		remoteUrl = v
	}
	if v, ok := os.LookupEnv("CONN_LIMIT"); ok {
		if n, _ := strconv.Atoi(v); n > 0 {
			connLimit = n
		}
	}
	if remoteUrl != "" {
		u, err := toIpUrl(remoteUrl)
		if err != nil {
			log.Fatalf("解析URL %s 失败: %s \n", remoteUrl, err)
		}
		log.Printf("解析URL %s => %s\n", remoteUrl, u)
		remoteUrl = u
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	l = netutil.LimitListener(l, connLimit)
	log.Println("启动HTTP服务 " + addr)
	if err := http.Serve(l, &Handler{}); err != nil {
		log.Fatalln(err)
	}
}

func toIpUrl(rawUrl string) (string, error) {
	if rawUrl == "" {
		return rawUrl, nil
	}
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	if u.Hostname() == "localhost" || net.ParseIP(u.Hostname()) != nil {
		return rawUrl, nil
	}
	addrs, err := net.LookupHost(u.Hostname())
	if err != nil {
		return "", err
	}

	if len(addrs) == 0 {
		return "", errors.Errorf("解析 %s 失败", u.Hostname())
	}
	u.Host = fmt.Sprintf("%s:%s", addrs[0], u.Port())
	return u.String(), nil
}
