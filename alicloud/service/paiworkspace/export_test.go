package paiworkspace

import "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

// ConfigExists reports whether a PAI Workspace Config entry exists for the
// given (workspaceId, categoryName, configKey) triple. It forwards to the
// unexported getConfig, which returns (nil, nil) when the entry is absent —
// the same NotFound signal the Framework resource's Read relies on. Exposed
// via export_test.go so the external paiworkspace_test package can drive a real
// destroy check against the cloud API without re-implementing the read path.
//
// The wrapper is deliberately not named Test* even though it lives in a
// _test.go file: Go's testing framework claims the Test prefix and requires
// func(t *testing.T), which would not compile for this helper signature.
func ConfigExists(client *connectivity.AliyunClient, workspaceId, categoryName, configKey string) (bool, error) {
	object, err := getConfig(client, workspaceId, categoryName, configKey)
	if err != nil {
		return false, err
	}
	return object != nil, nil
}

// ParseConfigID splits a composite Config id
// "<workspace_id>:<category_name>:<config_key>" into its three parts. It
// forwards to the unexported parseConfigID, returning ok=false for ids that do
// not have exactly three colon-separated parts. Like ConfigExists it avoids
// the Test prefix so it can carry a non-(t *testing.T) signature in a _test.go
// file.
func ParseConfigID(id string) (workspaceId, categoryName, configKey string, ok bool) {
	return parseConfigID(id)
}
