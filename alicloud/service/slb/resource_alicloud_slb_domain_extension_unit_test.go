package slb

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// listener_port and frontend_port are mutually exclusive and at least one is
// required; both rules are enforced by schema ExactlyOneOf at plan time, so a
// config-level unit test locks them without creating cloud resources.
func TestUnitAlicloudSLBDomainExtensionPortExactlyOneOf(t *testing.T) {
	baseConfig := map[string]interface{}{
		"load_balancer_id":      "lb-bp1unitcheck",
		"domain":                "www.test.com",
		"server_certificate_id": "cert-bp1unitcheck",
	}

	configWith := func(ports map[string]interface{}) map[string]interface{} {
		config := make(map[string]interface{}, len(baseConfig)+len(ports))
		for key, value := range baseConfig {
			config[key] = value
		}
		for key, value := range ports {
			config[key] = value
		}
		return config
	}

	errorDetails := func(config map[string]interface{}) []string {
		var details []string
		for _, diagnostic := range ResourceDomainExtension().Validate(terraform.NewResourceConfigRaw(config)) {
			if diagnostic.Severity == diag.Error {
				details = append(details, diagnostic.Detail)
			}
		}
		return details
	}

	t.Run("both ports configured", func(t *testing.T) {
		details := errorDetails(configWith(map[string]interface{}{
			"listener_port": 443,
			"frontend_port": 443,
		}))
		want := "only one of `frontend_port,listener_port` can be specified, but `frontend_port,listener_port` were specified."
		if !strings.Contains(strings.Join(details, "\n"), want) {
			t.Fatalf("want rejection reason %q, got %v", want, details)
		}
	})

	t.Run("neither port configured", func(t *testing.T) {
		details := errorDetails(configWith(nil))
		want := "one of `frontend_port,listener_port` must be specified"
		if !strings.Contains(strings.Join(details, "\n"), want) {
			t.Fatalf("want rejection reason %q, got %v", want, details)
		}
	})

	t.Run("single port configured", func(t *testing.T) {
		if details := errorDetails(configWith(map[string]interface{}{"listener_port": 443})); len(details) != 0 {
			t.Fatalf("want no errors, got %v", details)
		}
	})
}
