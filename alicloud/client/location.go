package client

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/location"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/client/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// endpointRetryBackoff returns the pause before the next discovery attempt.
func endpointRetryBackoff(attempt int) time.Duration {
	if attempt <= 1 {
		return 3 * time.Second
	}
	return 5 * time.Second
}

// endpointRetryWait keeps retry backoff within the caller's operation budget.
func endpointRetryWait(ctx context.Context, waitTime time.Duration) error {
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return context.DeadlineExceeded
		}
		if remaining <= waitTime {
			<-ctx.Done()
			return ctx.Err()
		}
	}
	if waitTime <= 0 {
		return nil
	}
	timer := time.NewTimer(waitTime)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

type locationContextTransport struct {
	ctx       context.Context
	transport *http.Transport
}

func (t *locationContextTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if err := t.ctx.Err(); err != nil {
		return nil, err
	}
	return t.transport.RoundTrip(request.WithContext(t.ctx))
}

func (client *AliyunClient) loadEndpoint(ctx context.Context, productCode string) error {
	return client.config.LoadEndpointWithContext(ctx, productCode, client.RegionId, func(ctx context.Context, productCode, regionId string) (string, error) {
		return client.describeEndpointForService(ctx, productCode, regionId)
	})
}

var describeLocationEndpoints = func(ctx context.Context, locationClient *location.Client, args *location.DescribeEndpointsRequest, readTimeout, connectTimeout time.Duration, transport *http.Transport, proxy *url.URL) (*location.DescribeEndpointsResponse, error) {
	attemptCtx := ctx
	if readTimeout > 0 {
		var cancel context.CancelFunc
		attemptCtx, cancel = context.WithTimeout(ctx, readTimeout)
		defer cancel()
	}
	if transport == nil {
		transport = &http.Transport{}
	}
	attemptTransport := transport.Clone()
	defer attemptTransport.CloseIdleConnections()
	attemptTransport.Proxy = nil
	if proxy != nil {
		attemptTransport.Proxy = http.ProxyURL(proxy)
	}
	dialer := &net.Dialer{Timeout: connectTimeout, DualStack: true}
	attemptTransport.DialContext = dialer.DialContext
	locationClient.SetReadTimeout(readTimeout)
	locationClient.SetTransport(&locationContextTransport{ctx: attemptCtx, transport: attemptTransport})
	return locationClient.DescribeEndpoints(args)
}

// describeEndpointForService probes Location for the endpoints of the region
// the request targets. The provider's client region only applies when the
// request itself carries no region; a cross-region request must describe its
// own region, or Location would answer with the client region's host.
func (client *AliyunClient) describeEndpointForService(ctx context.Context, productCode, regionId string) (string, error) {
	args := location.CreateDescribeEndpointsRequest()
	args.ServiceCode = config.LocationCodeForProduct(productCode)
	args.Id = regionId
	args.Domain = client.config.LocationEndpoint
	if args.Domain == "" {
		args.Domain = config.LoadEndpoint(regionId, config.LOCATIONCode)
	}
	if args.Domain == "" {
		args.Domain = "location.aliyuncs.com"
	}

	sdkConfig := client.getSdkConfig()
	locationClient, err := location.NewClientWithOptions(regionId, sdkConfig, client.config.AuthCredential(true))
	if err != nil {
		return "", fmt.Errorf("unable to initialize the location client: %w", err)
	}
	locationClient.SourceIp = client.config.SourceIp
	locationClient.SecureTransport = client.config.SecureTransport
	defer locationClient.Shutdown()
	var endpointResult string
	attempt := 0
	err = retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		attempt++
		if ctxErr := ctx.Err(); ctxErr != nil {
			return retry.NonRetryableError(ctxErr)
		}
		readTimeout, connectTimeout := client.endpointDiscoveryTimeouts(ctx)
		proxy, err := client.getHttpProxy()
		if err != nil {
			return retry.NonRetryableError(err)
		}
		skip, err := client.skipProxy(args.Domain)
		if err != nil {
			return retry.NonRetryableError(err)
		}
		if proxy != nil && skip {
			proxy = nil
		}
		endpointsResponse, err := describeLocationEndpoints(ctx, locationClient, args, readTimeout, connectTimeout, sdkConfig.HttpTransport, proxy)
		if ctxErr := ctx.Err(); ctxErr != nil {
			return retry.NonRetryableError(ctxErr)
		}
		if err != nil {
			if e, ok := err.(*errors.ClientError); ok && e.ErrorCode() == errors.TimeoutErrorCode {
				if waitErr := endpointRetryWait(ctx, endpointRetryBackoff(attempt)); waitErr != nil {
					return retry.NonRetryableError(waitErr)
				}
				args.Domain = "location-readonly.aliyuncs.com"
				return retry.RetryableError(err)
			}
			re := regexp.MustCompile("^Post [\"]*https://.*")
			if err.Error() != "" && re.MatchString(err.Error()) {
				if waitErr := endpointRetryWait(ctx, endpointRetryBackoff(attempt)); waitErr != nil {
					return retry.NonRetryableError(waitErr)
				}
				args.Domain = "location-readonly.aliyuncs.com"
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if endpointsResponse != nil && len(endpointsResponse.Endpoints.Endpoint) > 0 {
			for _, e := range endpointsResponse.Endpoints.Endpoint {
				if e.Type == "openAPI" {
					endpointResult = e.Endpoint
					return nil
				}
			}
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("Describe %s endpoint using region: %v got an error: %w.", productCode, regionId, err)
	}
	if endpointResult == "" {
		return "", fmt.Errorf("There is no any available endpoint for %s in region %s.", productCode, regionId)
	}
	return endpointResult, nil
}

// endpointDiscoveryTimeouts clamps each Location attempt to the remaining operation budget.
func (client *AliyunClient) endpointDiscoveryTimeouts(ctx context.Context) (time.Duration, time.Duration) {
	readTimeout := time.Duration(client.config.ClientReadTimeout) * time.Millisecond
	connectTimeout := time.Duration(client.config.ClientConnectTimeout) * time.Millisecond
	deadline, ok := ctx.Deadline()
	if !ok {
		return readTimeout, connectTimeout
	}
	remaining := time.Until(deadline)
	if remaining <= 0 {
		return time.Millisecond, time.Millisecond
	}
	if remaining < time.Millisecond {
		remaining = time.Millisecond
	}
	if readTimeout <= 0 || remaining < readTimeout {
		readTimeout = remaining
	}
	if connectTimeout <= 0 || remaining < connectTimeout {
		connectTimeout = remaining
	}
	return readTimeout, connectTimeout
}
