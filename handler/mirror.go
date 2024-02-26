package handler

import (
	"fmt"
	"go-relay/common"
	"io"
	"net/http"
	"golang.org/x/net/proxy"
	"os"
)

func MirrorHandler(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(r.Method, fmt.Sprintf("%s%s", common.MirrorWebsite, r.URL.String()), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header = r.Header.Clone()

	// 创建一个代理地址的拨号器
	dialer, err := proxy.SOCKS5("tcp", "172.17.0.5:1080", nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}

	// 设置 http.Transport 结构体
	httpTransport := &http.Transport{}

	// 设置连接的代理
	httpTransport.Dial = dialer.Dial

	// 创建连接
	client := &http.Client{Transport: httpTransport}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
