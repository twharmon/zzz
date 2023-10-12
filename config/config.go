package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Request struct {
	Method  string            `yaml:"method"`
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
	Body    map[string]any    `yaml:"body"`
}

type Config struct {
	Request *Request `yaml:"request"`
}

func Parse(file string) (*Config, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) Exec() error {
	u, err := url.Parse(c.Request.URL)
	if err != nil {
		return err
	}
	req := &http.Request{
		Method: c.Request.Method,
		URL:    u,
		Header: make(http.Header),
	}
	for k, v := range c.Request.Headers {
		req.Header.Set(k, v)
	}
	if c.Request.Body != nil {
		b, err := json.Marshal(c.Request.Body)
		if err != nil {
			return err
		}
		reader := bytes.NewReader(b)
		readerCloser := io.NopCloser(reader)
		req.Body = readerCloser
		req.Header.Set("Content-Type", "application/json")
	}
	fmt.Println("Making http request...")
	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	elapsed := time.Since(start).Round(time.Microsecond * 100)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Elapsed: %s\n", elapsed)
	if strings.Contains(resp.Header.Get("content-type"), "application/json") {
		var buf bytes.Buffer
		if err := json.Indent(&buf, b, "", "\t"); err != nil {
			return err
		}
		fmt.Println(buf.String())
	} else {
		fmt.Println(string(b))
	}
	return nil
}
