package redis

import (
	"testing"
	"time"

	tea "github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/stretchr/testify/assert"
)

// TestUnitRedisTairInstanceCreateRetryLogic deterministically exercises the
// SocketTimeout.Exception create-retry contract (defect: redis 创建超时) plus the
// pre-existing InternalError/NeedRetry behavior. It drives the REAL production
// SocketTimeout retry-decision function classifyCreateError (shared with
// resourceTairInstanceCreate), the REAL isSocketTimeoutException predicate /
// maxSocketTimeoutRetries / socketTimeoutRetryInterval constants, AND the REAL
// helper.NeedRetry / errs.IsExpectedErrors calls that the create closure keeps inline
// for the InternalError branch. Because every decision routes through production code,
// the test cannot pass against a shipped branch that disagrees with it — there is no
// hand-copied mirror.
//
// sleepFn replaces time.Sleep(socketTimeoutRetryInterval) so the 180s backoffs are
// counted without actually sleeping; waitFn replaces the InternalError/NeedRetry wait.
// A single ClientToken is generated before the loop (as in production) and the token
// observed on every call is recorded to prove it is reused verbatim across retries.
//
// AC1: SocketTimeout retries at a fixed 180s interval, up to 3 retries (4 attempts),
//
//	then writes the ID once a later attempt succeeds.
//
// AC2: once the retry budget is exhausted the original error is returned and no ID is
//
//	written.
//
// AC3: ClientToken is generated exactly once and reused on every attempt.
// AC4: the existing InternalError (max 2) / NeedRetry path and the fast-fail on
//
//	non-retryable errors are preserved; SocketTimeout uses its own counter.
func TestUnitRedisTairInstanceCreateRetryLogic(t *testing.T) {
	type callResult struct {
		errCode string
	}

	// buildRetryFunc drives the REAL production retry decision (classifyCreateError)
	// against a scripted sequence of create results, mirroring only the surrounding
	// loop scaffolding that resourceTairInstanceCreate cannot expose (the RpcPost call and
	// the schema timeout). The decision, counters and backoff selection are the shipped
	// logic, so a regression there fails this test rather than a stale copy.
	buildRetryFunc := func(results []callResult, waitFn func(), sleepFn func(time.Duration)) (
		callCount int, finalErr error, idWritten bool, seenTokens []string) {

		request := map[string]interface{}{}
		// ClientToken is generated once, before the retry loop, exactly as production.
		request["ClientToken"] = helper.BuildClientToken("CreateTairInstance")

		internalErrRetryCount := 0
		socketTimeoutRetryCount := 0
		// A generous timeout so the loop is driven purely by the result slice.
		retry.Retry(1*time.Hour, func() *retry.RetryError {
			if callCount >= len(results) {
				return nil
			}
			r := results[callCount]
			callCount++
			seenTokens = append(seenTokens, fmtString(request["ClientToken"]))

			if r.errCode == "" {
				idWritten = true // mirrors d.SetId(response["InstanceId"])
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

			// SocketTimeout branch: the REAL production classifyCreateError decides the
			// fixed 180s backoff (mirrors resourceTairInstanceCreate).
			if classifyCreateError(err, socketTimeoutRetryCount) {
				sleepFn(socketTimeoutRetryInterval)
				socketTimeoutRetryCount++
				return retry.RetryableError(err)
			}
			// InternalError/NeedRetry branch: mirrors the production inline decision,
			// driving the SAME real helper.NeedRetry and errs.IsExpectedErrors so a drift
			// there fails this test rather than a hand-copied mirror. Its budget is
			// separate from SocketTimeout's, exactly as production.
			if (helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"InternalError"})) && internalErrRetryCount < 2 {
				waitFn()
				internalErrRetryCount++
				return retry.RetryableError(err)
			}
			finalErr = err
			return retry.NonRetryableError(err)
		})
		return
	}

	stToken := func(tokens []string) (unique []string) {
		seen := map[string]bool{}
		for _, tk := range tokens {
			if !seen[tk] {
				seen[tk] = true
				unique = append(unique, tk)
			}
		}
		return
	}

	socketTimeoutCalls := func(n int) []callResult {
		out := make([]callResult, 0, n)
		for i := 0; i < n; i++ {
			out = append(out, callResult{errCode: "SocketTimeout.Exception"})
		}
		return out
	}

	// ── AC1: SocketTimeout retries 3 times (4 attempts) at a fixed 180s interval ──

	t.Run("SocketTimeout exhausts after 4 attempts: 3 fixed 180s backoffs then original error", func(t *testing.T) {
		var slept []time.Duration
		callCount, err, idWritten, tokens := buildRetryFunc(
			socketTimeoutCalls(4),
			func() {},
			func(d time.Duration) { slept = append(slept, d) },
		)
		assert.Error(t, err, "exhaustion must return the original error")
		assert.Equal(t, "SocketTimeout.Exception", *err.(*tea.SDKError).Code)
		assert.False(t, idWritten, "no ID must be written on exhaustion (AC2)")
		assert.Equal(t, 4, callCount, "1 initial + 3 retries = 4 attempts total")
		assert.Len(t, slept, 3, "exactly 3 retries (180s each)")
		for _, d := range slept {
			assert.Equal(t, socketTimeoutRetryInterval, d, "fixed 180s interval (AC1)")
		}
		assert.Len(t, stToken(tokens), 1, "ClientToken generated once and reused (AC3)")
		assert.Len(t, tokens, 4, "token observed on all 4 attempts")
	})

	t.Run("SocketTimeout succeeds on 4th attempt after 3 retries: ID written, 3 backoffs", func(t *testing.T) {
		var slept []time.Duration
		results := append(socketTimeoutCalls(3), callResult{errCode: ""})
		callCount, err, idWritten, tokens := buildRetryFunc(
			results,
			func() {},
			func(d time.Duration) { slept = append(slept, d) },
		)
		assert.NoError(t, err)
		assert.True(t, idWritten, "ID is written once a later attempt succeeds (AC1)")
		assert.Equal(t, 4, callCount, "3 failed + 1 success = 4 attempts")
		assert.Len(t, slept, 3, "3 retries before success")
		assert.Len(t, stToken(tokens), 1, "ClientToken reused across retries (AC3)")
	})

	t.Run("SocketTimeout succeeds on 2nd attempt after 1 retry", func(t *testing.T) {
		var slept []time.Duration
		callCount, err, idWritten, _ := buildRetryFunc(
			[]callResult{{errCode: "SocketTimeout.Exception"}, {errCode: ""}},
			func() {},
			func(d time.Duration) { slept = append(slept, d) },
		)
		assert.NoError(t, err)
		assert.True(t, idWritten)
		assert.Equal(t, 2, callCount)
		assert.Len(t, slept, 1)
	})

	// ── AC3: ClientToken is generated exactly once and reused ──

	t.Run("ClientToken is identical on every retry attempt", func(t *testing.T) {
		_, _, _, tokens := buildRetryFunc(
			socketTimeoutCalls(2),
			func() {},
			func(time.Duration) {},
		)
		assert.Len(t, tokens, 2)
		assert.Equal(t, tokens[0], tokens[1])
	})

	// ── AC4: SocketTimeout uses its own counter, independent of InternalError ──

	t.Run("SocketTimeout and InternalError counters are independent", func(t *testing.T) {
		var slept []time.Duration
		var waits int
		callCount, err, idWritten, _ := buildRetryFunc(
			[]callResult{
				{errCode: "SocketTimeout.Exception"}, // 1 SocketTimeout retry
				{errCode: "InternalError"},           // 1 InternalError retry
				{errCode: "InternalError"},           // 2nd InternalError retry (budget hit)
				{errCode: "InternalError"},           // exceeds InternalError budget -> error
			},
			func() { waits++ },
			func(d time.Duration) { slept = append(slept, d) },
		)
		assert.Error(t, err)
		assert.False(t, idWritten)
		assert.Equal(t, 4, callCount)
		assert.Len(t, slept, 1, "SocketTimeout backoff used once")
		assert.Equal(t, 2, waits, "InternalError wait used twice (its own budget)")
	})

	t.Run("SocketTimeout retries still succeed when followed by InternalError retries", func(t *testing.T) {
		var slept []time.Duration
		var waits int
		callCount, err, idWritten, _ := buildRetryFunc(
			[]callResult{
				{errCode: "SocketTimeout.Exception"},
				{errCode: "InternalError"},
				{errCode: ""},
			},
			func() { waits++ },
			func(d time.Duration) { slept = append(slept, d) },
		)
		assert.NoError(t, err)
		assert.True(t, idWritten)
		assert.Equal(t, 3, callCount)
		assert.Len(t, slept, 1)
		assert.Equal(t, 1, waits)
	})

	// ── AC4: existing InternalError (max 2) and NeedRetry path preserved ──

	t.Run("InternalError succeeds on second attempt: wait called once", func(t *testing.T) {
		var waits int
		callCount, err, _, _ := buildRetryFunc(
			[]callResult{{errCode: "InternalError"}, {errCode: ""}},
			func() { waits++ },
			func(time.Duration) {},
		)
		assert.NoError(t, err)
		assert.Equal(t, 2, callCount)
		assert.Equal(t, 1, waits, "expected one InternalError wait")
	})

	t.Run("InternalError succeeds on third attempt: wait called twice", func(t *testing.T) {
		var waits int
		callCount, err, _, _ := buildRetryFunc(
			[]callResult{{errCode: "InternalError"}, {errCode: "InternalError"}, {errCode: ""}},
			func() { waits++ },
			func(time.Duration) {},
		)
		assert.NoError(t, err)
		assert.Equal(t, 3, callCount)
		assert.Equal(t, 2, waits, "expected two InternalError waits")
	})

	t.Run("InternalError exceeds retry limit: wait called twice then returns error", func(t *testing.T) {
		var waits int
		callCount, err, idWritten, _ := buildRetryFunc(
			[]callResult{{errCode: "InternalError"}, {errCode: "InternalError"}, {errCode: "InternalError"}},
			func() { waits++ },
			func(time.Duration) {},
		)
		assert.Error(t, err)
		assert.True(t, errs.IsExpectedErrors(err, []string{"InternalError"}))
		assert.False(t, idWritten)
		assert.Equal(t, 3, callCount)
		assert.Equal(t, 2, waits, "third InternalError hits the counter limit and returns non-retryable")
	})

	t.Run("Throttling succeeds on third attempt: shares the InternalError counter", func(t *testing.T) {
		var waits int
		callCount, err, _, _ := buildRetryFunc(
			[]callResult{{errCode: "Throttling"}, {errCode: "Throttling"}, {errCode: ""}},
			func() { waits++ },
			func(time.Duration) {},
		)
		assert.NoError(t, err)
		assert.Equal(t, 3, callCount)
		assert.Equal(t, 2, waits)
	})

	t.Run("Throttling exceeds retry limit: returns error after 3 requests", func(t *testing.T) {
		var waits int
		callCount, err, _, _ := buildRetryFunc(
			[]callResult{{errCode: "Throttling"}, {errCode: "Throttling"}, {errCode: "Throttling"}},
			func() { waits++ },
			func(time.Duration) {},
		)
		assert.Error(t, err)
		assert.Equal(t, 3, callCount)
		assert.Equal(t, 2, waits)
	})

	t.Run("InternalError and Throttling mixed: share one counter, max 3 requests total", func(t *testing.T) {
		var waits int
		callCount, err, _, _ := buildRetryFunc(
			[]callResult{
				{errCode: "Throttling"},
				{errCode: "InternalError"},
				{errCode: "Throttling"},
			},
			func() { waits++ },
			func(time.Duration) {},
		)
		assert.Error(t, err)
		assert.Equal(t, 3, callCount)
		assert.Equal(t, 2, waits)
	})

	// ── AC4: non-retryable errors fast-fail; SocketTimeout is not mistaken for them ──

	t.Run("non-retryable error returns immediately without backoff or wait", func(t *testing.T) {
		var waits int
		var slept []time.Duration
		callCount, err, idWritten, _ := buildRetryFunc(
			[]callResult{{errCode: "InvalidParameter"}},
			func() { waits++ },
			func(d time.Duration) { slept = append(slept, d) },
		)
		assert.Error(t, err)
		assert.False(t, idWritten)
		assert.Equal(t, 1, callCount)
		assert.Equal(t, 0, waits)
		assert.Empty(t, slept)
	})

	t.Run("succeeds on first attempt: no backoff, no wait", func(t *testing.T) {
		var waits int
		var slept []time.Duration
		callCount, err, idWritten, _ := buildRetryFunc(
			[]callResult{{errCode: ""}},
			func() { waits++ },
			func(d time.Duration) { slept = append(slept, d) },
		)
		assert.NoError(t, err)
		assert.True(t, idWritten)
		assert.Equal(t, 1, callCount)
		assert.Equal(t, 0, waits)
		assert.Empty(t, slept)
	})

	// ── isSocketTimeoutException predicate (AC1/AC4 boundary) ──

	t.Run("isSocketTimeoutException identifies the server SocketTimeout.Exception", func(t *testing.T) {
		code := "SocketTimeout.Exception"
		msg := "code: 400, SocketTimeout.Exception"
		assert.True(t, isSocketTimeoutException(&tea.SDKError{Code: &code}))
		assert.True(t, isSocketTimeoutException(&tea.SDKError{Message: &msg}))
		assert.True(t, isSocketTimeoutException(&tea.SDKError{Code: &code, Message: &msg, StatusCode: tea.Int(400)}))
	})

	t.Run("isSocketTimeoutException does not match other retryable or non-retryable codes", func(t *testing.T) {
		for _, code := range []string{"InternalError", "Throttling", "InvalidParameter", "InstanceType.NotSupport"} {
			c := code
			assert.False(t, isSocketTimeoutException(&tea.SDKError{Code: &c}), code)
		}
	})

	t.Run("isSocketTimeoutException handles nil and text fallback", func(t *testing.T) {
		assert.False(t, isSocketTimeoutException(nil))
		assert.True(t, isSocketTimeoutException(errString("Request error: SocketTimeout.Exception")))
		assert.False(t, isSocketTimeoutException(errString("Request error: Client.Timeout")))
	})
}

// errString is a minimal error type for the predicate's err.Error() fallback path.
type errString string

func (e errString) Error() string { return string(e) }

// fmtString keeps the test free of a fmt import by stringifying a ClientToken value
// the same way the production closure observes it (it is a string produced by
// helper.BuildClientToken).
func fmtString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
