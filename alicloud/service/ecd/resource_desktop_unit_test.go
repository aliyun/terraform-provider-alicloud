package ecd

import "testing"

// A configured-but-empty desktop_timers block can surface as a nil element from
// GetOk. The request builder must skip it, as ValidateDesktopTimers does at plan
// time, instead of panicking at apply.
func TestUnitAliCloudEcdDesktopTimersNilElement(t *testing.T) {
	request := map[string]interface{}{}
	addDesktopTimersToRequest(request, []interface{}{nil})
	if len(request) != 0 {
		t.Fatalf("nil timer element produced request parameters: %#v", request)
	}

	request = map[string]interface{}{}
	addDesktopTimersToRequest(request, []interface{}{
		nil,
		map[string]interface{}{
			"timer_type":           "TimerShutdown",
			"cron_expression":      "0 0 16 ? * 1,2,3,4,5,6,7",
			"interval":             10,
			"enforce":              false,
			"reset_type":           "RESET_TYPE_BOTH",
			"allow_client_setting": true,
			"operation_type":       "",
		},
	})
	if got := request["DesktopTimers.1.TimerType"]; got != "TimerShutdown" {
		t.Fatalf("DesktopTimers.1.TimerType = %v, want TimerShutdown", got)
	}
	if _, ok := request["DesktopTimers.1.OperationType"]; ok {
		t.Fatalf("empty operation_type must be omitted, got %v", request["DesktopTimers.1.OperationType"])
	}
	if _, ok := request["DesktopTimers.2.TimerType"]; ok {
		t.Fatalf("skipped nil element must not reserve a request slot: %#v", request)
	}
}
