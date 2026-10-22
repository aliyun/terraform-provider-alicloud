package redis

import (
	"strings"

	tea "github.com/alibabacloud-go/tea/tea"
)

// isSocketTimeoutException reports whether err is the server-side
// SocketTimeout.Exception that R-kvstore returns (as a *tea.SDKError whose Code
// or Message carries "SocketTimeout") while CreateTairInstance is still being
// processed. This is distinct from a transport Client.Timeout (handled by
// helper.NeedRetry): it signals the create may still be in flight server-side,
// so the caller retries with the SAME ClientToken (idempotency) after a fixed
// 3-minute backoff rather than failing fast.
//
// R-kvstore surfaces business errors through connectivity.formatError as a
// *tea.SDKError whose Code is the raw API error code (e.g.
// "SocketTimeout.Exception"); the err.Error() fallback covers any wrapping that
// loses the typed SDKError so the predicate stays robust.
func isSocketTimeoutException(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok {
		if e.Code != nil && strings.Contains(*e.Code, "SocketTimeout") {
			return true
		}
		if e.Message != nil && strings.Contains(*e.Message, "SocketTimeout") {
			return true
		}
	}
	return strings.Contains(err.Error(), "SocketTimeout")
}

// classifyCreateError reports whether err is a SocketTimeout.Exception that
// CreateTairInstance should retry with a fixed socketTimeoutRetryInterval backoff,
// before the socketTimeoutRetryCount budget (maxSocketTimeoutRetries) is exhausted. It
// is the single SocketTimeout retry-decision point shared by the production create
// closure (resourceTairInstanceCreate) and TestUnitRedisTairInstanceCreateRetryLogic,
// so the tested and shipped SocketTimeout branches cannot diverge.
//
// The pre-existing InternalError/NeedRetry retry decision is intentionally NOT folded
// in here: it stays inline in the create closure under the CreateTairInstance action,
// which both keeps SocketTimeout on its own counter/backoff (independent of that path)
// and leaves the InternalError retry error codes where the breaking-change checker can
// observe them per-action. classifyCreateError therefore owns only the SocketTimeout
// contract; a drift in the inline InternalError branch is still caught by the unit test,
// which mirrors that branch against the same real helper.NeedRetry /
// errs.IsExpectedErrors calls.
func classifyCreateError(err error, socketTimeoutRetryCount int) bool {
	return isSocketTimeoutException(err) && socketTimeoutRetryCount < maxSocketTimeoutRetries
}
