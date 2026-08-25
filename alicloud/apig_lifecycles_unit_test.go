package alicloud

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestApigLifecycleResourceSchemas(t *testing.T) {
	t.Parallel()

	resources := map[string]*schema.Resource{
		"consumer":            resourceAliCloudApigConsumer(),
		"policy":              resourceAliCloudApigPolicy(),
		"authorization_rules": resourceAliCloudApigConsumerAuthorizationRules(),
		"deployment":          resourceAliCloudApigHttpApiDeployment(),
	}
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("provider schema is invalid: %v", err)
	}
	if !resources["policy"].Schema["config"].Sensitive {
		t.Fatal("policy config must be sensitive")
	}
	for _, forbidden := range []string{"ak", "sk", "access_key", "secret_key"} {
		if _, exists := resources["consumer"].Schema[forbidden]; exists {
			t.Fatalf("consumer schema must not expose %q", forbidden)
		}
	}
}

func TestApigConsumerCreateRequestUsesSystemGeneratedCredentials(t *testing.T) {
	t.Parallel()

	resource := resourceAliCloudApigConsumer()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"consumer_name": "terraform-test",
		"description":   "test consumer",
	})
	request := apigConsumerCreateRequest(d)
	expected := []map[string]interface{}{{"type": "AkSk", "generateMode": "System"}}
	if !reflect.DeepEqual(request["akSkIdentityConfigs"], expected) {
		t.Fatalf("unexpected identity request: %#v", request["akSkIdentityConfigs"])
	}
	if _, exists := request["ak"]; exists {
		t.Fatal("consumer request must not contain ak")
	}
	if _, exists := request["sk"]; exists {
		t.Fatal("consumer request must not contain sk")
	}
}

func TestApigAuthorizationRulesRequestIsExactAndSorted(t *testing.T) {
	t.Parallel()

	resource := resourceAliCloudApigConsumerAuthorizationRules()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"consumer_id":        "cs-test",
		"environment_id":     "env-test",
		"parent_resource_id": "api-test",
		"resource_ids":       []interface{}{"hr-b", "hr-a"},
	})
	request := apigAuthorizationRulesRequest(d)
	rules, ok := request["authorizationRules"].([]map[string]interface{})
	if !ok || len(rules) != 2 {
		t.Fatalf("unexpected authorization request: %#v", request)
	}
	for index, resourceID := range []string{"hr-a", "hr-b"} {
		rule := rules[index]
		if rule["consumerId"] != "cs-test" || rule["principalType"] != "Consumer" || rule["resourceType"] != "HttpApiRoute" || rule["expireMode"] != "LongTerm" {
			t.Fatalf("unexpected authorization rule: %#v", rule)
		}
		identifier := rule["resourceIdentifier"].(map[string]interface{})
		if identifier["resourceId"] != resourceID || identifier["parentResourceId"] != "api-test" || identifier["environmentId"] != "env-test" {
			t.Fatalf("unexpected resource identifier: %#v", identifier)
		}
	}
}

func TestApigHttpApiDeploymentRequestUsesCurrentRouteContract(t *testing.T) {
	t.Parallel()

	resource := resourceAliCloudApigHttpApiDeployment()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"http_api_id":    "api-test",
		"route_id":       "hr-test",
		"environment_id": "env-test",
		"gateway_id":     "gw-test",
	})
	request := apigHttpApiDeploymentRequest(d)
	expected := map[string]interface{}{"routeId": "hr-test", "isInternalRoute": false}
	if !reflect.DeepEqual(request, expected) {
		t.Fatalf("unexpected deployment request: %#v", request)
	}
	if _, deprecated := request["httpApiConfig"]; deprecated {
		t.Fatal("deployment request must not use deprecated httpApiConfig")
	}
}

func TestApigRouteMatchesDeployment(t *testing.T) {
	t.Parallel()

	route := map[string]interface{}{
		"environmentInfo": map[string]interface{}{
			"environmentId": "env-test",
			"gatewayInfo": map[string]interface{}{
				"gatewayId": "gw-test",
			},
		},
	}
	if !apigRouteMatchesDeployment(route, "env-test", "gw-test") {
		t.Fatal("expected deployment target to match")
	}
	if apigRouteMatchesDeployment(route, "env-other", "gw-test") {
		t.Fatal("different environment must not match")
	}
}

func TestApigSelectPolicyAttachmentPreservesBatchResourceIDs(t *testing.T) {
	t.Parallel()

	attachment := apigSelectPolicyAttachment([]map[string]interface{}{
		{
			"policyId": "plc-other",
			"attachments": []interface{}{
				map[string]interface{}{"policyAttachmentId": "pr-test"},
			},
		},
		{
			"policyId": "plc-test",
			"attachments": []interface{}{
				map[string]interface{}{
					"policyAttachmentId": "pr-test",
					"attachResourceIds":  []interface{}{"hr-b", "hr-a"},
				},
			},
		},
	}, "plc-test", "pr-test")
	if attachment == nil {
		t.Fatal("expected policy attachment")
	}
	if got := apigStringSlice(attachment["attachResourceIds"]); !reflect.DeepEqual(got, []string{"hr-b", "hr-a"}) {
		t.Fatalf("unexpected attachment resource IDs: %#v", got)
	}
}

func TestApigCompositeImportValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		resource *schema.Resource
		id       string
		wantID   string
		wantErr  bool
	}{
		{name: "policy", resource: resourceAliCloudApigPolicy(), id: "plc-test:pr-test", wantID: "plc-test"},
		{name: "policy invalid", resource: resourceAliCloudApigPolicy(), id: "plc-test", wantErr: true},
		{name: "deployment", resource: resourceAliCloudApigHttpApiDeployment(), id: "api-test:hr-test:env-test:gw-test", wantID: "api-test:hr-test:env-test:gw-test"},
		{name: "authorization", resource: resourceAliCloudApigConsumerAuthorizationRules(), id: "cs-test:env-test:api-test:HttpApiRoute:hr-b,hr-a", wantID: "cs-test:env-test:api-test:HttpApiRoute"},
		{name: "authorization duplicate", resource: resourceAliCloudApigConsumerAuthorizationRules(), id: "cs-test:env-test:api-test:HttpApiRoute:hr-a,hr-a", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, test.resource.Schema, nil)
			d.SetId(test.id)
			var err error
			switch test.name {
			case "policy", "policy invalid":
				_, err = resourceAliCloudApigPolicyImport(d, nil)
			case "deployment":
				_, err = resourceAliCloudApigHttpApiDeploymentImport(d, nil)
			default:
				_, err = resourceAliCloudApigConsumerAuthorizationRulesImport(d, nil)
			}
			if test.wantErr {
				if err == nil {
					t.Fatal("expected import error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected import error: %v", err)
			}
			if d.Id() != test.wantID {
				t.Fatalf("unexpected imported ID %q", d.Id())
			}
		})
	}
}

func TestApigAcceptanceConfigSyntax(t *testing.T) {
	t.Parallel()

	configs := map[string]string{
		"consumer": resourceTestAccConfigFunc("alicloud_apig_consumer.default", "tfacc", func(name string) string {
			return fmt.Sprintf("variable \"name\" { default = %q }\n", name)
		})(map[string]interface{}{
			"consumer_name":            "tfacc-consumer",
			"description":              "syntax check",
			"enable":                   true,
			"gateway_type":             "API",
			"credential_generate_mode": "System",
		}),
		"policy": resourceTestAccConfigFunc("alicloud_apig_policy.default", "tfacc", alicloudApigPolicyBasicDependence)(map[string]interface{}{
			"policy_name":          "tfacc-policy",
			"class_name":           "RateLimit",
			"config":               "{\\\"enable\\\":true,\\\"threshold\\\":10}",
			"description":          "syntax check",
			"attach_resource_ids":  []string{"${alicloud_apig_route.policy_primary.route_id}"},
			"attach_resource_type": "GatewayRoute",
			"environment_id":       "${alicloud_apig_gateway.policy_primary.environments.0.environment_id}",
			"gateway_id":           "${alicloud_apig_gateway.policy_primary.id}",
		}),
		"authorization": resourceTestAccConfigFunc("alicloud_apig_consumer_authorization_rules.default", "tfacc", alicloudApigAuthorizationBasicDependence)(map[string]interface{}{
			"consumer_id":        "${alicloud_apig_consumer.lifecycle.id}",
			"environment_id":     "${alicloud_apig_http_api_deployment.authorization.environment_id}",
			"parent_resource_id": "${alicloud_apig_http_api.lifecycle.id}",
			"resource_type":      "HttpApiRoute",
			"resource_ids":       []string{"${alicloud_apig_route.lifecycle.route_id}"},
			"principal_type":     "Consumer",
			"expire_mode":        "LongTerm",
		}),
		"deployment": resourceTestAccConfigFunc("alicloud_apig_http_api_deployment.default", "tfacc", alicloudApigLifecycleBasicDependence)(map[string]interface{}{
			"http_api_id":    "${alicloud_apig_http_api.lifecycle.id}",
			"route_id":       "${alicloud_apig_route.lifecycle.route_id}",
			"environment_id": "${alicloud_apig_gateway.lifecycle.environments.0.environment_id}",
			"gateway_id":     "${alicloud_apig_gateway.lifecycle.id}",
		}),
	}
	for name, config := range configs {
		name, config := name, config
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, diagnostics := hclsyntax.ParseConfig([]byte(config), name+".tf", hcl.Pos{Line: 1, Column: 1})
			if diagnostics.HasErrors() {
				t.Fatalf("generated acceptance config is invalid: %s", diagnostics.Error())
			}
		})
	}
}
