package slb

import (
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// TestUnitAlicloudSLBTlsCipherPoliciesDataSource was moved into the internal
// package slb test scope because it patches the unexported data source Read
// func (dataSourceTlsCipherPoliciesRead) and the AliyunClient via gomonkey;
// the external slb_test package cannot reach unexported symbols or import
// alicloud/acctest without a cycle. It constructs a ListTLSCipherPolicies
// response whose RelateListeners carry a listener "Protocol" field and
// verifies the output "protocol" attribute is populated from that response
// map key. This guards the migration fidelity fixed in PR #11398: the
// listener map key m1["Protocol"] is a production response field and must
// not be rewritten to a package-qualified token such as m1["helper.Protocol"].

// lintignore: R001
func TestUnitAlicloudSLBTlsCipherPoliciesDataSource(t *testing.T) {
	p := DataSourceTlsCipherPolicies()
	d, _ := schema.InternalMap(p.Schema).Data(nil, nil)

	// A mock-only AliyunClient (dummy creds, region validation and account-type
	// probing skipped, SLB endpoint pre-populated) keeps this offline: gomonkey
	// intercepts (*client.Client).DoRequest, so no real cloud call or credential
	// is needed and the test no longer SKIPs when ALICLOUD_* env is unset.
	rawClient, err := unitTestMockClient()
	if err != nil {
		t.Fatalf("failed to build mock client: %s", err)
	}

	// ListTLSCipherPolicies response with a relate listener carrying a real
	// "Protocol" field; the data source Read maps it to output "protocol".
	readMockResponse := map[string]interface{}{
		"TLSCipherPolicies": []interface{}{
			map[string]interface{}{
				"Ciphers":     []interface{}{"AES256-SHA256"},
				"Status":      "normal",
				"TLSVersions": []interface{}{"TLSv1.1", "TLSv1.2"},
				"Name":        "tls_cipher_policy_name",
				"InstanceId":  "MockTLSCipherPolicyId",
				"RelateListeners": []interface{}{
					map[string]interface{}{
						"LoadBalancerId": "lb-mock",
						"Port":           float64(443),
						"Protocol":       "https",
					},
				},
			},
		},
		"NextToken": "",
	}

	t.Run("ReadNormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			return readMockResponse, nil
		})
		defer patches.Reset()
		err := dataSourceTlsCipherPoliciesRead(d, rawClient)
		assert.Nil(t, err)

		policies := d.Get("policies").([]interface{})
		assert.Equal(t, 1, len(policies))
		policy := policies[0].(map[string]interface{})
		relateListeners := policy["relate_listeners"].([]interface{})
		assert.Equal(t, 1, len(relateListeners))
		listener := relateListeners[0].(map[string]interface{})
		// The "Protocol" response map key must reach the output "protocol"
		// attribute; a migration-misnamed key would leave this empty.
		assert.Equal(t, "https", listener["protocol"])
		assert.Equal(t, "lb-mock", listener["load_balancer_id"])
		assert.Equal(t, 443, listener["port"])
		assert.Equal(t, "MockTLSCipherPolicyId", policy["id"])
	})
}
