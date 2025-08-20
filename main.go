package traefik_plugin_response_cache_control

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Config 設定結構
type Config struct {
	Value               string   `json:"value,omitempty"`
	Override            bool     `json:"override,omitempty"`
	ExcludedStatusCodes []string `json:"excludedStatusCodes,omitempty"`
}

// CreateConfig 建立預設設定
func CreateConfig() *Config {
	return &Config{
		Value:    "public, max-age=3600",
		Override: true,
	}
}

type ResponseCacheControl struct {
	next     http.Handler
	name     string
	config   *Config
	excluded [][2]int
}

// New 建立 middleware
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	plugin := &ResponseCacheControl{
		next:   next,
		name:   name,
		config: config,
	}

	// 解析排除的狀態碼範圍
	for _, v := range config.ExcludedStatusCodes {
		parts := strings.Split(v, "-")
		if len(parts) == 1 {
			code, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid status code: %s", v)
			}
			plugin.excluded = append(plugin.excluded, [2]int{code, code})
		} else if len(parts) == 2 {
			min, err1 := strconv.Atoi(parts[0])
			max, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid status code range: %s", v)
			}
			plugin.excluded = append(plugin.excluded, [2]int{min, max})
		}
	}

	return plugin, nil
}

func (p *ResponseCacheControl) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// 包裝 ResponseWriter
	wrapped := &responseWriter{
		ResponseWriter: rw,
		header:         make(http.Header),
		plugin:         p,
	}

	// 複製原始標頭到包裝器
	for k, v := range rw.Header() {
		wrapped.header[k] = v
	}

	p.next.ServeHTTP(wrapped, req)
}

// responseWriter 包裝 status code 和標頭
type responseWriter struct {
	http.ResponseWriter
	status      int
	header      http.Header
	wroteHeader bool
	plugin      *ResponseCacheControl
}

func (w *responseWriter) Header() http.Header {
	return w.header
}

func (w *responseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.status = code
	w.wroteHeader = true

	// 在寫入標頭前應用 Cache-Control 規則
	w.applyCacheControlHeader()

	// 把 wrapped header 複製到實際的 ResponseWriter
	for k, v := range w.header {
		w.ResponseWriter.Header()[k] = v
	}

	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) applyCacheControlHeader() {
	// 判斷是否在排除清單
	status := w.status
	if status == 0 {
		status = http.StatusOK
	}

	for _, rng := range w.plugin.excluded {
		if status >= rng[0] && status <= rng[1] {
			return
		}
	}

	// 判斷是否要覆蓋
	if !w.plugin.config.Override {
		if w.header.Get("Cache-Control") != "" {
			return
		}
	}

	// 設定 Cache-Control
	w.header.Set("Cache-Control", w.plugin.config.Value)
}
