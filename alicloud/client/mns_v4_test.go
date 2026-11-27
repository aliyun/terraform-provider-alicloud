package client

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/valyala/fasthttp"
)

// The physical region used to sign the captured MNS requests. It deliberately
// differs from the classic endpoint regions embedded in the test URLs so a
// passing signature check proves the region is taken from the provider config.
const mnsV4TestRegion = "cn-hangzhou"

// newMnsV4WireClient wraps a legacy client whose URL keeps the 5-segment public
// cloud shape the SDK parser requires (localhost hosts make it panic) while the
// wrapper's transport targets the capture server.
func newMnsV4WireClient(t *testing.T, owner *AliyunClient, sendURL string) *mnsV4Client {
	t.Helper()
	legacy := ali_mns.NewAliMNSClientWithToken(
		fmt.Sprintf("https://123456789.mns.%s.aliyuncs.com", mnsV4TestRegion), "test-ak", "test-secret", "")
	v4conn, err := newMnsV4Client(legacy, sendURL, owner)
	if err != nil {
		t.Fatalf("building the v4 MNS client: %s", err)
	}
	return v4conn
}

// verifyMnsDataPlaneWire recomputes the ACS4 signature and checks the
// MNS-specific wire headers the legacy clients also sent. The body is read
// once and re-seeded so the signature recomputation and the digest check both
// observe the real payload.
func verifyMnsDataPlaneWire(t *testing.T, r *http.Request) {
	t.Helper()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	verifyWireV4(t, r, "mns", mnsV4TestRegion)
	sum := md5.Sum(body)
	if got := r.Header.Get("Content-MD5"); got != base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%x", sum))) {
		t.Errorf("Content-MD5 = %q, want the legacy digest of the body", got)
	}
	if got := r.Header.Get("X-Mns-Version"); got != "2015-06-06" {
		t.Errorf("x-mns-version = %q, want 2015-06-06", got)
	}
	if got := r.Header.Get("Content-Type"); got != "application/xml" {
		t.Errorf("Content-Type = %q, want application/xml", got)
	}
	if auth := r.Header.Get("Authorization"); !strings.HasPrefix(auth, "ACS4-HMAC-SHA256 Credential=test-ak/") || strings.HasPrefix(auth, "MNS ") {
		t.Errorf("Authorization = %q, want an ACS4 credential for test-ak", auth)
	}
}

// TestUnitMnsV4DataPlaneWireCapture drives the queue/topic managers on a v4-wrapped
// client against a local HTTP server: signature, path, query, XML body, legacy
// MNS headers, temporary token, error decoding and the single-attempt retry
// semantics are all asserted from the captured requests.
func TestUnitMnsV4DataPlaneWireCapture(t *testing.T) {
	clearProxyEnv(t)

	const queueName = "tf-mns-v4-queue"
	const topicName = "tf-mns-v4-topic"
	const receiptHandle = "1-ODM4OTkxMzI2NC4xMjM0NTEzODQwODYzJTNEJTNE"

	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
		verifyMnsDataPlaneWire(t, r)
		switch {
		case r.Method == "PUT" && r.URL.Path == "/queues/"+queueName && r.URL.RawQuery == "":
			for _, want := range []string{"<Queue>", "<DelaySeconds>10</DelaySeconds>", "<MaximumMessageSize>2048</MaximumMessageSize>"} {
				if !strings.Contains(string(body), want) {
					t.Errorf("create queue body %q misses %q", body, want)
				}
			}
			w.WriteHeader(http.StatusCreated)
		case r.Method == "PUT" && r.URL.Path == "/queues/"+queueName && r.URL.RawQuery == "metaoverride=true":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == "GET" && r.URL.Path == "/queues/"+queueName:
			io.WriteString(w, "<Queue><QueueName>"+queueName+"</QueueName><DelaySeconds>20</DelaySeconds><MaximumMessageSize>2048</MaximumMessageSize><MessageRetentionPeriod>345600</MessageRetentionPeriod><VisibilityTimeout>30</VisibilityTimeout><PollingWaitSeconds>0</PollingWaitSeconds></Queue>")
		case r.Method == "GET" && r.URL.Path == "/queues":
			if r.Header.Get("X-Mns-Marker") != "marker-1" || r.Header.Get("X-Mns-Ret-Number") != "10" || r.Header.Get("X-Mns-Prefix") != "tf-" {
				t.Errorf("list queue headers = %q / %q / %q, want the caller supplied x-mns-* values",
					r.Header.Get("X-Mns-Marker"), r.Header.Get("X-Mns-Ret-Number"), r.Header.Get("X-Mns-Prefix"))
			}
			io.WriteString(w, "<Queues><Queue><QueueURL>http://123456789.mns.cn-hangzhou.aliyuncs.com/queues/"+queueName+"</QueueURL></Queue></Queues>")
		case r.Method == "DELETE" && r.URL.Path == "/queues/"+queueName:
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "<Error><Code>QueueNotExist</Code><Message>the queue does not exist</Message></Error>")
		case r.Method == "PUT" && r.URL.Path == "/topics/"+topicName:
			if !strings.Contains(string(body), "<MaximumMessageSize>2048</MaximumMessageSize>") {
				t.Errorf("create topic body %q misses the message size", body)
			}
			w.WriteHeader(http.StatusCreated)
		case r.Method == "GET" && r.URL.Path == "/topics/"+topicName:
			io.WriteString(w, "<Topic><TopicName>"+topicName+"</TopicName><MaximumMessageSize>2048</MaximumMessageSize><MessageRetentionPeriod>86400</MessageRetentionPeriod></Topic>")
		case r.Method == "GET" && r.URL.Path == "/queues/"+queueName+"/messages":
			if got := r.URL.Query().Get("ReceiptHandle"); got != receiptHandle {
				t.Errorf("ReceiptHandle = %q, want the decoded %q", got, receiptHandle)
			}
			io.WriteString(w, "<Messages></Messages>")
		default:
			t.Errorf("unexpected request %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	c := newLocalRpcClient(t, strings.TrimPrefix(server.URL, "http://"))
	c.config.RegionId, c.RegionId = mnsV4TestRegion, mnsV4TestRegion
	cred, err := credentials.NewCredential(new(credentials.Config).SetType("sts").
		SetAccessKeyId("test-ak").SetAccessKeySecret("test-secret").SetSecurityToken("test-token"))
	if err != nil {
		t.Fatal(err)
	}
	c.config.Credential = cred
	v4conn := newMnsV4WireClient(t, c, server.URL)

	queueManager := ali_mns.NewMNSQueueManager(v4conn)
	if err := queueManager.CreateQueue(queueName, 10, 2048, 345600, 30, 0, 3); err != nil {
		t.Fatalf("CreateQueue: %s", err)
	}
	attr, err := queueManager.GetQueueAttributes(queueName)
	if err != nil || attr.QueueName != queueName || attr.DelaySeconds != 20 {
		t.Fatalf("GetQueueAttributes = %+v, %v", attr, err)
	}
	if err := queueManager.SetQueueAttributes(queueName, 20, 2048, 345600, 30, 0, 3); err != nil {
		t.Fatalf("SetQueueAttributes: %s", err)
	}
	queues, err := queueManager.ListQueue("marker-1", 10, "tf-")
	if err != nil || len(queues.Queues) != 1 || !strings.Contains(queues.Queues[0].QueueURL, queueName) {
		t.Fatalf("ListQueue = %+v, %v", queues, err)
	}
	// A 404 error response must decode into the MNS QueueNotExist error the
	// service layer's not-found helpers match on.
	err = queueManager.DeleteQueue(queueName)
	if err == nil || !strings.Contains(err.Error(), "QueueNotExist") {
		t.Fatalf("DeleteQueue error = %v, want QueueNotExist", err)
	}

	topicManager := ali_mns.NewMNSTopicManager(v4conn)
	if err := topicManager.CreateTopic(topicName, 2048, false); err != nil {
		t.Fatalf("CreateTopic: %s", err)
	}
	topicAttr, err := topicManager.GetTopicAttributes(topicName)
	if err != nil || topicAttr.TopicName != topicName {
		t.Fatalf("GetTopicAttributes = %+v, %v", topicAttr, err)
	}

	// A query value the caller embedded already escaped must reach the server
	// decoded identically, with the signature computed over the same value.
	resp, err := v4conn.Send(ali_mns.GET, nil, nil, fmt.Sprintf("queues/%s/messages?ReceiptHandle=%s", queueName, url.QueryEscape(receiptHandle)))
	if err != nil || resp.StatusCode() != http.StatusOK {
		t.Fatalf("direct Send = %v, %v", resp, err)
	}

	// The legacy client never retries at the SDK level: one request per call.
	if got := atomic.LoadInt32(&requests); got != 8 {
		t.Errorf("captured %d requests, want 8 (no SDK-level retries)", got)
	}
}

// TestUnitWithMnsClientSignVersionGating pins the WithMnsClient selection: the
// legacy client stays untouched unless v4 is selected, an explicit mns = v1
// overrides a v4 default, and the v4 wrapper resolves the data plane URL from
// the endpoint configuration while the signing region comes from the provider.
func TestUnitWithMnsClientSignVersionGating(t *testing.T) {
	clearProxyEnv(t)
	t.Setenv("MNS_ENDPOINT", "")
	const accountId = "123456789"

	type capture struct {
		isV4    bool
		sendURL string
	}
	assertClient := func(t *testing.T, c *AliyunClient) capture {
		t.Helper()
		var got capture
		_, err := c.WithMnsClient(func(conn *ali_mns.MNSClient) (interface{}, error) {
			if v4conn, ok := (*conn).(*mnsV4Client); ok {
				got.isV4, got.sendURL = true, v4conn.sendURL
			}
			return nil, nil
		})
		if err != nil {
			t.Fatalf("WithMnsClient: %s", err)
		}
		return got
	}

	t.Run("unconfigured keeps the legacy client", func(t *testing.T) {
		c := newLocalRpcClient(t, "unused.example")
		c.accountId = accountId
		if got := assertClient(t, c); got.isV4 {
			t.Fatalf("unconfigured sign_version produced a v4 client: %+v", got)
		}
	})

	t.Run("default v4 wraps and resolves the region endpoint", func(t *testing.T) {
		c := newLocalRpcClient(t, "unused.example")
		c.accountId = accountId
		c.config.RegionId, c.RegionId = mnsV4TestRegion, mnsV4TestRegion
		c.config.SignVersion.Store("default", "v4")
		got := assertClient(t, c)
		want := fmt.Sprintf("https://%s.mns.%s.aliyuncs.com", accountId, mnsV4TestRegion)
		if !got.isV4 || got.sendURL != want {
			t.Fatalf("got %+v, want the v4 wrapper sending to %q", got, want)
		}
	})

	t.Run("explicit station endpoint overrides the region endpoint", func(t *testing.T) {
		c := newLocalRpcClient(t, "unused.example")
		c.accountId = accountId
		c.config.RegionId, c.RegionId = mnsV4TestRegion, mnsV4TestRegion
		c.config.SignVersion.Store("default", "v4")
		c.config.MnsEndpoint = "eo-station.example.com"
		got := assertClient(t, c)
		want := fmt.Sprintf("https://%s.mns.%s", accountId, "eo-station.example.com")
		if !got.isV4 || got.sendURL != want {
			t.Fatalf("got %+v, want the v4 wrapper sending to %q", got, want)
		}
	})

	t.Run("mns v1 overrides a v4 default", func(t *testing.T) {
		c := newLocalRpcClient(t, "unused.example")
		c.accountId = accountId
		c.config.RegionId, c.RegionId = mnsV4TestRegion, mnsV4TestRegion
		c.config.SignVersion.Store("default", "v4")
		c.config.SignVersion.Store("mns", "v1")
		if got := assertClient(t, c); got.isV4 {
			t.Fatalf("sign_version.mns = v1 produced a v4 client: %+v", got)
		}
	})

	t.Run("v4 without a region fails", func(t *testing.T) {
		c := newLocalRpcClient(t, "unused.example")
		c.accountId = accountId
		c.config.RegionId = ""
		c.config.SignVersion.Store("default", "v4")
		_, err := c.WithMnsClient(func(conn *ali_mns.MNSClient) (interface{}, error) { return nil, nil })
		if err == nil || !strings.Contains(err.Error(), "requires a signing region") {
			t.Fatalf("got %v, want the missing signing region diagnostic", err)
		}
	})
}

// TestUnitMnsOpenManagementPlaneSignVersion pins the mnsopen (Mns-open) management
// plane carried by the shared RPC path: unconfigured it keeps the legacy
// HMAC-SHA1 RPC signature, default v4 signs it with ACS4.
func TestUnitMnsOpenManagementPlaneSignVersion(t *testing.T) {
	clearProxyEnv(t)
	for _, v4 := range []bool{false, true} {
		t.Run(fmt.Sprintf("v4=%t", v4), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if v4 {
					verifyWireV4(t, r, "mns-open", mnsV4TestRegion)
				} else {
					q := r.URL.Query()
					if q.Get("SignatureMethod") != "HMAC-SHA1" || q.Get("Signature") == "" {
						t.Errorf("expected the legacy RPC signature, got %v", q)
					}
				}
				w.Header().Set("Content-Type", "application/json")
				io.WriteString(w, `{"RequestId":"request-id"}`)
			}))
			defer server.Close()
			endpoint := strings.TrimPrefix(server.URL, "http://")
			c := newLocalRpcClient(t, endpoint)
			c.config.Endpoints.Store("mns_open", endpoint)
			c.config.RegionId, c.RegionId = mnsV4TestRegion, mnsV4TestRegion
			c.teaRpcOpenapiConfig.Protocol = tea.String("HTTP")
			cred, err := credentials.NewCredential(new(credentials.Config).SetType("sts").
				SetAccessKeyId("test-ak").SetAccessKeySecret("test-secret").SetSecurityToken("test-token"))
			if err != nil {
				t.Fatal(err)
			}
			c.config.Credential = cred
			if v4 {
				c.config.SignVersion.Store("default", "v4")
			}
			response, err := c.RpcPost("Mns-open", "2022-01-19", "ListTopics", map[string]interface{}{"PageSize": 10}, map[string]interface{}{"PageSize": 10}, false)
			if err != nil || response["RequestId"] != "request-id" {
				t.Fatalf("unexpected response: %v, %v", response, err)
			}
		})
	}
}

func TestUnitMnsV4ProxyRouting(t *testing.T) {
	clearProxyEnv(t)
	served := func(marker string) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.Header.Get("Authorization"), "ACS4-HMAC-SHA256 ") {
				t.Error("missing signature")
			}
			w.Header().Set("X-Route", marker)
		}))
	}
	origin, proxy := served("direct"), served("proxy")
	defer origin.Close()
	defer proxy.Close()
	t.Setenv("HTTP_PROXY", proxy.URL)
	t.Setenv("no_proxy", "")
	for _, tc := range []struct{ bypass, want string }{{"", "proxy"}, {"127.0.0.1", "direct"}} {
		t.Setenv("NO_PROXY", tc.bypass)
		c := newLocalRpcClient(t, strings.TrimPrefix(origin.URL, "http://"))
		client := newMnsV4WireClient(t, c, origin.URL)
		resp, err := client.Send(ali_mns.Method("GET"), nil, nil, "queues")
		if err != nil {
			t.Fatal(err)
		}
		if got := string(resp.Header.Peek("X-Route")); got != tc.want {
			t.Errorf("route %q, want %q", got, tc.want)
		}
		fasthttp.ReleaseResponse(resp)
	}
}

func TestUnitMnsV4HTTPSProxyConnect(t *testing.T) {
	clearProxyEnv(t)
	t.Setenv("NO_PROXY", "")
	t.Setenv("no_proxy", "")
	const host = "queues.eo.invalid"
	var connected atomic.Bool
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodConnect && r.Host == host+":443" {
			connected.Store(true)
		}
		if r.Header.Get("Authorization") != "" {
			t.Error("MNS signature sent to proxy during CONNECT")
		}
		w.WriteHeader(http.StatusProxyAuthRequired)
	}))
	defer proxy.Close()
	t.Setenv("HTTPS_PROXY", proxy.URL)
	c := newLocalRpcClient(t, host)
	c.config.Protocol = "HTTPS"
	client := newMnsV4WireClient(t, c, "https://"+host)
	if _, err := client.Send(ali_mns.Method("GET"), nil, nil, "queues"); err == nil || !connected.Load() {
		t.Fatalf("expected CONNECT via configured HTTPS proxy: %v", err)
	}
}
