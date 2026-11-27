package config

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitCommonLoadEndpointFromLocal(t *testing.T) {
	content := `<?xml version="1.0" encoding="UTF-8"?>
<Endpoints>
  <Endpoint name="cn-beijing">
    <RegionIds>
      <RegionId>cn-beijing</RegionId>
    </RegionIds>
    <Products>
      <Product>
        <ProductName>ECS</ProductName>
        <DomainName>ecs.cn-beijing.aliyuncs.com</DomainName>
      </Product>
    </Products>
  </Endpoint>
</Endpoints>`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "endpoints.xml")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	assert.NoError(t, err)

	t.Setenv("TF_ENDPOINT_PATH", tmpFile)
	defer os.Unsetenv("TF_ENDPOINT_PATH")

	config := &Config{
		Endpoints: new(sync.Map),
		RegionId:  "cn-beijing",
	}

	err = config.loadEndpointFromLocal()
	assert.NoError(t, err)

	val, ok := config.Endpoints.Load("ecs")
	assert.True(t, ok)
	assert.Equal(t, "ecs.cn-beijing.aliyuncs.com", val)
}

func TestUnitEfloControllerEndpoint(t *testing.T) {
	// The location service returns InvalidRegionId for efloctrl in some
	// regions (e.g. us-southeast-1). regularProductEndpoint acts as the
	// backup endpoint so eflo_controller resolves to
	// eflo-controller.<region>.aliyuncs.com instead of failing.
	endpointFmt, ok := regularProductEndpoint["eflo_controller"]
	assert.True(t, ok)
	assert.Equal(t, "eflo-controller.%s.aliyuncs.com", endpointFmt)

	for _, region := range []string{"us-southeast-1", "cn-hangzhou", "cn-beijing"} {
		assert.Equal(t, fmt.Sprintf("eflo-controller.%s.aliyuncs.com", region), fmt.Sprintf(endpointFmt, region))
	}
}

func newSilentListener(t *testing.T) (string, func(), <-chan struct{}) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	var mu sync.Mutex
	connections := make([]net.Conn, 0)
	accepted := make(chan struct{})
	var acceptOnce sync.Once
	go func() {
		for {
			connection, err := listener.Accept()
			if err != nil {
				return
			}
			acceptOnce.Do(func() { close(accepted) })
			mu.Lock()
			connections = append(connections, connection)
			mu.Unlock()
		}
	}()
	stop := func() {
		_ = listener.Close()
		mu.Lock()
		defer mu.Unlock()
		for _, connection := range connections {
			_ = connection.Close()
		}
	}
	return listener.Addr().String(), stop, accepted
}

func endpointProbeConfig(t *testing.T, locationEndpoint string) *Config {
	t.Helper()
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-key").
		SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("credential: %v", err)
	}
	t.Setenv("NO_PROXY", locationEndpoint)
	t.Setenv("no_proxy", locationEndpoint)
	return &Config{
		AccessKey:            "test-key",
		SecretKey:            "test-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "https",
		ClientReadTimeout:    60000,
		ClientConnectTimeout: 60000,
		LocationEndpoint:     locationEndpoint,
		Endpoints:            new(sync.Map),
		SignVersion:          new(sync.Map),
		SkipRegionValidation: true,
	}
}

func TestUnitLoadAPIEndpointWithContextStopsABlockingProbe(t *testing.T) {
	host, stop, accepted := newSilentListener(t)
	defer stop()

	config := endpointProbeConfig(t, host)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		select {
		case <-accepted:
		case <-time.After(2 * time.Second):
		}
		cancel()
	}()

	started := time.Now()
	endpoint, err := config.LoadAPIEndpointWithContext(ctx, "ecs", "cn-hangzhou", func() (string, error) { return "", nil })
	require.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, endpoint)
	assert.Less(t, time.Since(started), 5*time.Second)
	_, cached := config.Endpoints.Load("ecs")
	assert.False(t, cached, "a cancelled probe must not cache the regular fallback")
}

func TestUnitLoadAPIEndpointWithContextSkipsACancelledProbe(t *testing.T) {
	host, stop, accepted := newSilentListener(t)
	defer stop()

	config := endpointProbeConfig(t, host)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	started := time.Now()
	endpoint, err := config.LoadAPIEndpointWithContext(ctx, "ecs", "cn-hangzhou", func() (string, error) { return "", nil })
	require.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, endpoint)
	assert.Less(t, time.Since(started), 5*time.Second)
	_, cached := config.Endpoints.Load("ecs")
	assert.False(t, cached)
	select {
	case <-accepted:
		t.Fatal("a pre-cancelled read started a Location probe")
	default:
	}
}

func TestUnitLoadAPIEndpointWithContextServesLocation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"RequestId":"mock","Success":true,"Endpoints":{"Endpoint":[{"Type":"openAPI","Endpoint":"ecs.context.example"}]}}`))
	}))
	defer server.Close()
	host := strings.TrimPrefix(server.URL, "http://")
	config := endpointProbeConfig(t, host)
	config.Protocol = "http"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint, err := config.LoadAPIEndpointWithContext(ctx, "ecs", "cn-hangzhou", func() (string, error) { return "", nil })
	require.NoError(t, err)
	assert.Equal(t, "ecs.context.example", endpoint)
	cached, ok := config.Endpoints.Load("ecs")
	assert.True(t, ok)
	assert.Equal(t, "ecs.context.example", cached)
}

func TestUnitLoadAPIEndpointBackgroundServesLocation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"RequestId":"mock","Success":true,"Endpoints":{"Endpoint":[{"Type":"openAPI","Endpoint":"ecs.location.example"}]}}`))
	}))
	defer server.Close()
	host := strings.TrimPrefix(server.URL, "http://")
	config := endpointProbeConfig(t, host)
	config.Protocol = "http"

	endpoint, err := config.LoadAPIEndpoint("ecs", "cn-hangzhou", func() (string, error) { return "", nil })
	require.NoError(t, err)
	assert.Equal(t, "ecs.location.example", endpoint)
	cached, ok := config.Endpoints.Load("ecs")
	assert.True(t, ok)
	assert.Equal(t, "ecs.location.example", cached)
}

func TestUnitLoadAPIEndpointBackgroundKeepsRegularFallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"RequestId":"mock","Code":"InvalidRequest","Message":"rejected"}`))
	}))
	defer server.Close()
	host := strings.TrimPrefix(server.URL, "http://")
	config := endpointProbeConfig(t, host)
	config.Protocol = "http"

	endpoint, err := config.LoadAPIEndpoint("ecs", "cn-hangzhou", func() (string, error) { return "", nil })
	require.NoError(t, err)
	assert.Equal(t, "ecs.cn-hangzhou.aliyuncs.com", endpoint)
	cached, ok := config.Endpoints.Load("ecs")
	assert.True(t, ok)
	assert.Equal(t, "ecs.cn-hangzhou.aliyuncs.com", cached)
}

func TestUnitLoadAPIEndpointBackgroundRetriesLocation(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			time.Sleep(200 * time.Millisecond)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"RequestId":"mock","Success":true,"Endpoints":{"Endpoint":[{"Type":"openAPI","Endpoint":"ecs.retry.example"}]}}`))
	}))
	defer server.Close()
	host := strings.TrimPrefix(server.URL, "http://")
	config := endpointProbeConfig(t, host)
	config.Protocol = "http"
	config.ClientReadTimeout = 50
	config.ClientConnectTimeout = 50

	started := time.Now()
	endpoint, err := config.LoadAPIEndpoint("ecs", "cn-hangzhou", func() (string, error) { return "", nil })
	require.NoError(t, err)
	assert.Equal(t, "ecs.retry.example", endpoint)
	assert.GreaterOrEqual(t, atomic.LoadInt32(&calls), int32(2))
	assert.Less(t, time.Since(started), 15*time.Second)
}

// TestUnitLoadAPIEndpointForRegionBypassesTheWarmSharedCache pins the
// cross-region cache contract: the shared per-product entry only holds the
// endpoint of whichever region resolved first — a same-region call usually
// warms it during apply — so a request for another region must not read it
// (it would carry one region's signature against another region's host),
// must probe Location for its own region, and must cache under a
// region-scoped key without touching the shared entry.
func TestUnitLoadAPIEndpointForRegionBypassesTheWarmSharedCache(t *testing.T) {
	config := endpointProbeConfig(t, "location.invalid.example")
	// A same-region call during apply warmed the shared per-product entry.
	config.Endpoints.Store("vpc", "vpc.cn-hangzhou.example")

	var probed []string
	discover := func(ctx context.Context, productCode string, regionID string) (string, error) {
		probed = append(probed, productCode+"@"+regionID)
		return "vpc.cn-beijing.example", nil
	}
	endpoint, err := config.LoadAPIEndpointForRegion(context.Background(), "vpc", "cn-beijing", func() (string, error) { return "", nil }, discover)
	require.NoError(t, err)
	assert.Equal(t, "vpc.cn-beijing.example", endpoint)
	assert.Equal(t, []string{"vpc@cn-beijing"}, probed, "Location must describe the request's region")

	shared, ok := config.Endpoints.Load("vpc")
	require.True(t, ok)
	assert.Equal(t, "vpc.cn-hangzhou.example", shared, "the warm shared entry must stay untouched")

	regional, ok := config.Endpoints.Load(RegionalEndpointCacheKey("vpc", "cn-beijing"))
	require.True(t, ok)
	assert.Equal(t, "vpc.cn-beijing.example", regional)

	// The region-scoped entry serves the second cross-region call without
	// probing Location again.
	endpoint, err = config.LoadAPIEndpointForRegion(context.Background(), "vpc", "cn-beijing", func() (string, error) { return "", nil },
		func(ctx context.Context, productCode string, regionID string) (string, error) {
			probed = append(probed, productCode+"@"+regionID)
			return "vpc.somewhere-else.example", nil
		})
	require.NoError(t, err)
	assert.Equal(t, "vpc.cn-beijing.example", endpoint)
	assert.Len(t, probed, 1, "the region-scoped cache must serve the warm path")
}

// TestUnitLoadAPIEndpointForRegionKeepsTheEnvironmentOverride pins the
// override precedence on the cross-region path: the environment endpoint
// wins over the region rules without a Location probe, exactly as it does
// for the provider's own region.
func TestUnitLoadAPIEndpointForRegionKeepsTheEnvironmentOverride(t *testing.T) {
	config := endpointProbeConfig(t, "location.invalid.example")
	t.Setenv("ALIBABA_CLOUD_ENDPOINT_VPC", "vpc.env.example")

	endpoint, err := config.LoadAPIEndpointForRegion(context.Background(), "vpc", "cn-beijing", func() (string, error) { return "", nil },
		func(ctx context.Context, productCode string, regionID string) (string, error) {
			t.Fatal("the environment override must win without a Location probe")
			return "", nil
		})
	require.NoError(t, err)
	assert.Equal(t, "vpc.env.example", endpoint)
}
