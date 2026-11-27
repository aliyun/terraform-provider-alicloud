package client

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	sdkerrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
)

func TestUnitVpcSdkRpcLegacyRetries(t *testing.T) {
	var count atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count.Add(1)
		w.WriteHeader(500)
		io.WriteString(w, `{"Code":"InternalError","Message":"fixture"}`)
	}))
	defer server.Close()
	c := newLegacyWireClient(t, "vpc", server.Listener.Addr().String())
	response := vpc.CreateDescribeRegionsResponse()
	got, err := c.VpcSdkRpc(vpc.CreateDescribeRegionsRequest(), response)
	var remote *sdkerrors.ServerError
	if got != response || !errors.As(err, &remote) || count.Load() != 6 {
		t.Fatalf("attempts=%d, response=%T, err=%v", count.Load(), got, err)
	}
}
