package connectivity

import (
	"net/http"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/recording"
)

func TestRecordingClientRejectsUnsafeCredentialSources(t *testing.T) {
	t.Setenv(recording.ModeEnv, "replay")
	cases := []Config{
		{AccountType: "Domestic", AccessKey: "real-key", SecretKey: "real-secret"},
		{AccountType: "Domestic", AccessKey: "recording-access-key", SecretKey: "recording-secret-key", SecurityToken: "real-token"},
		{AccountType: "Domestic", AccessKey: "recording-access-key", SecretKey: "recording-secret-key", RamRoleArn: "role"},
		{AccountType: "Domestic", AccessKey: "recording-access-key", SecretKey: "recording-secret-key", EcsRoleName: "metadata-role"},
	}
	for i, c := range cases {
		if _, err := c.Client(); err == nil {
			t.Fatalf("case %d accepted unsafe replay credentials", i)
		}
	}
}
func TestRecordingClientRequiresExplicitAccountType(t *testing.T) {
	t.Setenv(recording.ModeEnv, "record")
	c := &Config{AccessKey: "test-key", SecretKey: "test-secret"}
	if _, err := c.Client(); err == nil || !strings.Contains(err.Error(), "account_type") {
		t.Fatalf("unexpected account discovery result: %v", err)
	}
}
func TestRecordingUnsupportedClientsFailBeforeRequests(t *testing.T) {
	t.Setenv(recording.ModeEnv, "record")
	c := &AliyunClient{}
	if _, err := c.RpcPost("Product", "2020-01-01", "Action", nil, nil, true); err == nil {
		t.Fatal("unsupported Tea RPC accepted")
	}
	if _, err := c.RoaGet("Product", "2020-01-01", "/", nil, nil, nil); err == nil {
		t.Fatal("unsupported Tea ROA accepted")
	}
	if _, err := c.NewTeaCommonClient("https://example.invalid"); err == nil {
		t.Fatal("unsupported direct Tea factory accepted")
	}
}

func TestRecordingTransportPreservesConfiguredProxy(t *testing.T) {
	t.Setenv(recording.ModeEnv, "record")
	t.Setenv("HTTPS_PROXY", "http://127.0.0.1:1234")
	t.Setenv("NO_PROXY", "")
	t.Setenv("no_proxy", "")
	client := &AliyunClient{config: &Config{Protocol: "HTTPS"}}
	config := client.getSdkConfig(0)
	request, _ := http.NewRequest("GET", "https://example.invalid/", nil)
	proxy, err := config.HttpTransport.Proxy(request)
	if err != nil || proxy == nil || proxy.String() != "http://127.0.0.1:1234" {
		t.Fatalf("recording lost proxy configuration: %v %v", proxy, err)
	}
}
