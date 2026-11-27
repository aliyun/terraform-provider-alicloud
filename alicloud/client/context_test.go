package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/client/config"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-log/tflogtest"
)

func TestProviderContextPreservesParent(t *testing.T) {
	type key struct{}
	var output bytes.Buffer
	parent, cancel := context.WithDeadline(context.WithValue(context.Background(), key{}, "kept"), time.Now().Add(time.Minute))
	defer cancel()
	parent = tflogtest.RootLogger(parent, &output)
	client := &AliyunClient{config: &config.Config{TerraformTraceId: "provider-1", RegionId: "cn-hangzhou"}}
	ctx := client.ProviderContext(parent)
	deadline, ok := ctx.Deadline()
	want, _ := parent.Deadline()
	if !ok || !deadline.Equal(want) || ctx.Value(key{}) != "kept" || ctx.Done() != parent.Done() {
		t.Fatal("Provider 公共上下文改变了父 ctx")
	}
	tflog.Info(ctx, "child")
	tflog.Info(parent, "parent")
	entries, err := tflogtest.MultilineJSONDecode(&output)
	if err != nil || len(entries) != 2 || entries[0]["alicloud_trace_id"] != "provider-1" || entries[0]["alicloud_region"] != "cn-hangzhou" {
		t.Fatalf("公共字段缺失：%v %v", entries, err)
	}
	if _, exists := entries[1]["alicloud_trace_id"]; exists {
		t.Fatal("公共字段污染父 ctx")
	}
	cancel()
	if !errors.Is(ctx.Err(), context.Canceled) {
		t.Fatal("Provider 公共上下文未传播取消")
	}
	var empty *AliyunClient
	if empty.ProviderContext(parent) != parent || (&AliyunClient{}).ProviderContext(parent) != parent {
		t.Fatal("空 client 应保留 incoming ctx")
	}
}

func TestNewClientContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	client, err := NewClientWithValidRegionsContext(ctx, &config.Config{}, nil)
	if client != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("取消必须先于配置校验返回：client=%p err=%v", client, err)
	}
}

func TestNewClientContextCancelsAccountProbe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stopped := make(chan bool, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		cancel()
		select {
		case <-r.Context().Done():
			stopped <- true
		case <-time.After(2 * time.Second):
			stopped <- false
			http.Error(w, "未传播取消", http.StatusBadRequest)
		}
	}))
	defer server.Close()
	cfg := &config.Config{
		AccessKey: "test-ak", SecretKey: "test-sk", RegionId: "cn-hangzhou", Region: config.Region("cn-hangzhou"),
		Protocol: "HTTP", SkipRegionValidation: true, Endpoints: &sync.Map{}, SignVersion: &sync.Map{},
		Credential: &mockCredential{accessKeyId: "test-ak", accessKeySecret: "test-sk"},
	}
	cfg.Endpoints.Store("bssopenapi", strings.TrimPrefix(server.URL, "http://"))
	client, err := NewClientWithValidRegionsContext(ctx, cfg, nil)
	if client != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("账号探测取消未终止构造：client=%p err=%v", client, err)
	}
	if !<-stopped {
		t.Fatal("账号探测 HTTP 请求未响应取消")
	}
}

func TestNewClientContextAccountProbeResults(t *testing.T) {
	for _, tc := range []struct {
		name        string
		status      int
		body        string
		failures    int32
		wantHits    int32
		wantAccount string
	}{
		{"success", 200, `{"Code":"Success"}`, 0, 1, "Domestic"},
		{"not-applicable", 400, `{"Code":"NotApplicable","Message":"not applicable"}`, 1, 1, "International"},
		{"ordinary-error", 403, `{"Code":"Forbidden","Message":"denied"}`, 1, 1, ""},
		{"http-retry", 503, `{"Code":"ServiceUnavailable","Message":"busy"}`, 3, 4, "Domestic"},
		{"decode-retry", 502, `<html>gateway</html>`, 3, 4, "Domestic"},
		{"decode-retry-exhausted", 502, `<html>gateway</html>`, 4, 4, ""},
		{"business-retry", 200, `{"Code":"Throttling","Message":"busy"}`, 1, 2, "Domestic"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var hits atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil || r.FormValue("Action") != "QueryAvailableInstances" || r.FormValue("ProductType") != "vipcloudfw" {
					t.Errorf("账号探测请求发生变化：%v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				if hits.Add(1) <= tc.failures {
					w.WriteHeader(tc.status)
					fmt.Fprint(w, tc.body)
					return
				}
				fmt.Fprint(w, `{"Code":"Success"}`)
			}))
			defer server.Close()
			cfg := &config.Config{
				AccessKey: "test-ak", SecretKey: "test-sk", RegionId: "cn-hangzhou", Region: config.Region("cn-hangzhou"),
				Protocol: "HTTP", SkipRegionValidation: true, Endpoints: &sync.Map{}, SignVersion: &sync.Map{},
				Credential: &mockCredential{accessKeyId: "test-ak", accessKeySecret: "test-sk"},
			}
			endpoint := strings.TrimPrefix(server.URL, "http://")
			cfg.Endpoints.Store("bssopenapi", endpoint)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			client, err := NewClientWithValidRegionsContext(ctx, cfg, nil)
			if err != nil || client == nil || cfg.AccountType != tc.wantAccount || hits.Load() != tc.wantHits {
				t.Fatalf("账号探测结果错误：account=%q hits=%d err=%v", cfg.AccountType, hits.Load(), err)
			}
			if tc.wantAccount == "International" {
				endpoint = config.BssOpenAPIEndpointInternational
			}
			if got, _ := cfg.Endpoints.Load("bssopenapi"); got != endpoint {
				t.Fatalf("账号探测 endpoint 错误：%v", got)
			}
		})
	}
}

func TestRequestContextPreservesParent(t *testing.T) {
	type key struct{}
	var output bytes.Buffer
	parent, cancel := context.WithDeadline(context.WithValue(context.Background(), key{}, "parent"), time.Now().Add(time.Minute))
	defer cancel()
	parent = tflogtest.RootLogger(parent, &output)
	parent = tflog.SetField(parent, "tf_req_id", "request-1")

	ctx := RequestContext(parent, "Vpc", "Read")
	if ctx.Value(key{}) != "parent" {
		t.Fatal("父 context 的值丢失")
	}
	want, _ := parent.Deadline()
	if got, ok := ctx.Deadline(); !ok || !got.Equal(want) {
		t.Fatalf("deadline = %v, %v，期望 %v", got, ok, want)
	}
	cancel()
	if ctx.Err() != context.Canceled {
		t.Fatalf("取消未传播：%v", ctx.Err())
	}

	tflog.Info(ctx, "child")
	tflog.Info(parent, "parent")
	entries, err := tflogtest.MultilineJSONDecode(&output)
	if err != nil || len(entries) != 2 {
		t.Fatalf("日志条目 = %v，错误 = %v", entries, err)
	}
	if entries[0]["tf_req_id"] != "request-1" || entries[0]["alicloud_service"] != "Vpc" || entries[0]["alicloud_operation"] != "Read" {
		t.Fatalf("请求日志字段不完整：%v", entries[0])
	}
	for _, field := range []string{"alicloud_service", "alicloud_operation"} {
		if _, ok := entries[1][field]; ok {
			t.Fatalf("派生字段污染父 context：%v", entries[1])
		}
	}
}

func TestRequestContextDoesNotInventServiceOrDeadline(t *testing.T) {
	var output bytes.Buffer
	ctx := RequestContext(tflogtest.RootLogger(context.Background(), &output), "", "Create")
	if _, ok := ctx.Deadline(); ok {
		t.Fatal("公共补充不应新增 deadline")
	}
	tflog.Info(ctx, "request")
	entries, err := tflogtest.MultilineJSONDecode(&output)
	if err != nil || len(entries) != 1 {
		t.Fatalf("日志条目 = %v，错误 = %v", entries, err)
	}
	if _, ok := entries[0]["alicloud_service"]; ok {
		t.Fatalf("没有产品声明时不应补充 service：%v", entries[0])
	}
	if entries[0]["alicloud_operation"] != "Create" {
		t.Fatalf("操作字段缺失：%v", entries[0])
	}
}

func TestRequestContextWithoutLogger(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	ctx := RequestContext(parent, "Vpc", "Read")
	cancel()
	if ctx.Err() != context.Canceled {
		t.Fatal("没有 logger 时也必须保留取消链")
	}
}

func TestRequestContextConcurrentIsolation(t *testing.T) {
	type key struct{}
	parent := context.WithValue(context.Background(), key{}, "shared-parent")
	for _, service := range []string{"Vpc", "Ecs", "Ims"} {
		t.Run(service, func(t *testing.T) {
			t.Parallel()
			var output bytes.Buffer
			base := tflogtest.RootLogger(parent, &output)
			first := RequestContext(base, service, "Create")
			second := RequestContext(base, service, "Delete")
			tflog.Info(first, "first")
			tflog.Info(second, "second")
			entries, err := tflogtest.MultilineJSONDecode(&output)
			if err != nil || len(entries) != 2 {
				t.Fatalf("日志条目 = %v，错误 = %v", entries, err)
			}
			for i, op := range []string{"Create", "Delete"} {
				if entries[i]["alicloud_service"] != service || entries[i]["alicloud_operation"] != op {
					t.Fatalf("请求字段串值：%v", entries[i])
				}
			}
			if first.Value(key{}) != "shared-parent" || second.Value(key{}) != "shared-parent" {
				t.Fatal("父 context 值未保留")
			}
		})
	}
}
