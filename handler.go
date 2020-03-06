package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Handler struct {
	TotalCount int64
	ErrorCount int64
	TotalTime  time.Duration
	MaxTime    time.Duration
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		h.process(w, r)
	case "/status":
		h.status(w, r)
	default:
		w.WriteHeader(404)
		w.Write([]byte("404 not found"))
	}
}

func (h *Handler) process(w http.ResponseWriter, r *http.Request) {
	urlStr := r.URL.Query().Get("url")
	dimensions := r.URL.Query().Get("size") // 图片缩放尺寸，格式：800x600
	if urlStr == "" {
		w.WriteHeader(500)
		w.Write([]byte("url不能为空"))
		return
	}
	var (
		bs  []byte
		err error
	)
	t1 := time.Now()
	defer func() {
		atomic.AddInt64(&h.TotalCount, 1)
		ut := time.Now().Sub(t1)
		if err != nil {
			atomic.AddInt64(&h.ErrorCount, 1)
			log.Printf("[error] 截图失败，URL: %s, 错误: %s, 耗时: %s", urlStr, err, ut)
		} else {
			log.Printf("[success] 截图成功，URL: %s, 大小: %s, 耗时: %s", urlStr, SizeFormat(float64(len(bs))), ut)
		}
		h.TotalTime = h.TotalTime + ut
		if ut > h.MaxTime {
			h.MaxTime = ut
		}
	}()
	bs, err = screenshot(urlStr)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	if dimensions != "" {
		ss := strings.Split(dimensions, "x")
		w, _ := strconv.Atoi(ss[0])
		h, _ := strconv.Atoi(ss[1])
		if b, err := ImageResize(bytes.NewReader(bs), uint(w), uint(h)); err == nil {
			bs = b
		}
	}

	w.Header().Set("content-type", "image/png")
	w.Write(bs)
}

func (h *Handler) status(w http.ResponseWriter, r *http.Request) {
	avgTime := time.Duration(0)
	if h.TotalCount > 0 {
		avgTime = h.TotalTime / time.Duration(h.TotalCount)
	}
	hostname, _ := os.Hostname()
	data := map[string]interface{}{
		"start_time":  startTime,
		"total_count": h.TotalCount,
		"error_count": h.ErrorCount,
		"max_time":    h.MaxTime.String(),
		"avg_time":    avgTime.String(),
		"conn_limit":  connLimit,
		"hostname":    hostname,
	}
	b, _ := json.Marshal(data)
	w.Write(b)
}
