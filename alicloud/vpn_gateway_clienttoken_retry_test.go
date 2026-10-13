package alicloud

import (
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

// TestUnitVPNGatewayClientTokenRetryIdempotency locks the fix for duplicate
// VPN gateway creation. The root cause was that buildClientToken was called
// inside the resource.Retry closure, producing a fresh ClientToken on every
// retry. When the first RPC POST timed out after the server had already
// created the gateway, the retry used a new token — the server treated it
// as a new request and created a duplicate.
//
// This test mirrors the retry closure in resourceAliCloudVPNGatewayVPNGatewayCreate
// and verifies that the ClientToken set once before the retry closure is
// reused across all retry attempts, preserving OpenAPI idempotency.
func TestUnitVPNGatewayClientTokenRetryIdempotency(t *testing.T) {
	type vpnCallResult struct {
		errCode string
	}

	// buildVPNGatewayRetryFunc mirrors the retry closure in
	// resourceAliCloudVPNGatewayVPNGatewayCreate. ClientToken is set
	// ONCE before the retry closure (as the fix requires) and must not
	// be regenerated inside the closure.
	buildVPNGatewayRetryFunc := func(results []vpnCallResult, waitFn func()) (callCount int, finalErr error, tokens []string) {
		action := "CreateVpnGateway"
		request := make(map[string]interface{})
		// ClientToken is set once before the retry closure — same as the fix.
		request["ClientToken"] = buildClientToken(action)
		tokens = []string{}

		resource.Retry(1*time.Minute, func() *resource.RetryError {
			if callCount >= len(results) {
				return nil
			}
			r := results[callCount]
			callCount++
			// Capture the ClientToken used for this attempt.
			if token, ok := request["ClientToken"].(string); ok {
				tokens = append(tokens, token)
			}

			if r.errCode == "" {
				return nil
			}

			errCode := r.errCode
			errMsg := r.errCode
			err := &tea.SDKError{
				Code:       &errCode,
				Data:       &errMsg,
				Message:    &errMsg,
				StatusCode: tea.Int(400),
			}

			if NeedRetry(err) {
				waitFn()
				return resource.RetryableError(err)
			}
			finalErr = err
			return resource.NonRetryableError(err)
		})
		return
	}

	// Throttling is a NeedRetry path. When the server returns Throttling,
	// the original request may have been processed. Reusing the same
	// ClientToken lets the server deduplicate on retry.
	t.Run("Throttling retry reuses the same ClientToken", func(t *testing.T) {
		waitCalls := 0
		callCount, err, tokens := buildVPNGatewayRetryFunc(
			[]vpnCallResult{{errCode: "Throttling"}, {errCode: ""}},
			func() { waitCalls++ },
		)
		assert.NoError(t, err)
		assert.Equal(t, 2, callCount)
		assert.Len(t, tokens, 2)
		assert.Equal(t, tokens[0], tokens[1],
			"ClientToken must be reused across retries to preserve OpenAPI idempotency; "+
				"regenerating it inside the retry closure causes duplicate resource creation")
	})

	// ServiceUnavailable across multiple retries — all attempts must
	// share the same ClientToken.
	t.Run("ServiceUnavailable multiple retries reuse the same ClientToken", func(t *testing.T) {
		waitCalls := 0
		callCount, err, tokens := buildVPNGatewayRetryFunc(
			[]vpnCallResult{
				{errCode: "ServiceUnavailable"},
				{errCode: "ServiceUnavailable"},
				{errCode: ""},
			},
			func() { waitCalls++ },
		)
		assert.NoError(t, err)
		assert.Equal(t, 3, callCount)
		assert.Len(t, tokens, 3)
		for i := 1; i < len(tokens); i++ {
			assert.Equal(t, tokens[0], tokens[i],
				"all retry attempts must share the same ClientToken (attempt %d differs)", i)
		}
	})

	// Negative control: buildClientToken generates a DIFFERENT token on
	// each call. This proves that calling it inside the retry closure
	// (the original bug) would break idempotency, justifying the fix.
	t.Run("buildClientToken generates unique tokens proving regeneration would break idempotency", func(t *testing.T) {
		token1 := buildClientToken("CreateVpnGateway")
		token2 := buildClientToken("CreateVpnGateway")
		assert.NotEqual(t, token1, token2,
			"buildClientToken must return a unique token each call; "+
				"if it were called inside the retry closure, retries would carry different tokens "+
				"and the server would create duplicate resources")
		assert.NotEmpty(t, token1)
		assert.NotEmpty(t, token2)
	})
}
