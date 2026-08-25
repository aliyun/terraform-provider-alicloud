package alicloud

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func alicloudApigLifecycleBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "lifecycle" {
  network_access_config {
    type = "Intranet"
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  gateway_type = "API"
  payment_type = "PayAsYouGo"
  gateway_name = format("%%s-gateway", var.name)
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = false
    }
  }
}

resource "alicloud_apig_http_api" "lifecycle" {
  http_api_name = format("%%s-api", var.name)
  protocols     = ["HTTP"]
  type          = "Rest"
  base_path     = format("/%%s", var.name)
}

resource "alicloud_apig_route" "lifecycle" {
  route_name  = format("%%s-route", var.name)
  http_api_id = alicloud_apig_http_api.lifecycle.id
  environment_info {
    environment_id = alicloud_apig_gateway.lifecycle.environments.0.environment_id
  }
  match {
    path {
      type  = "Exact"
      value = format("/%%s", var.name)
    }
    methods = ["GET"]
  }
  backend {
    scene = "Mock"
  }
}

resource "alicloud_apig_consumer" "lifecycle" {
  consumer_name            = format("%%s-consumer", var.name)
  description              = "Terraform acceptance lifecycle consumer"
  enable                   = true
  gateway_type             = "API"
  credential_generate_mode = "System"
}
`, name)
}

func alicloudApigAuthorizationBasicDependence(name string) string {
	return alicloudApigLifecycleBasicDependence(name) + `
resource "alicloud_apig_http_api_deployment" "authorization" {
  http_api_id    = alicloud_apig_http_api.lifecycle.id
  route_id       = alicloud_apig_route.lifecycle.route_id
  environment_id = alicloud_apig_gateway.lifecycle.environments.0.environment_id
  gateway_id     = alicloud_apig_gateway.lifecycle.id
}
`
}

func alicloudApigPolicyBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "policy_primary" {
  network_access_config {
    type = "Intranet"
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  gateway_type = "API"
  payment_type = "PayAsYouGo"
  gateway_name = format("%%s-gateway-primary", var.name)
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = false
    }
  }
}

resource "alicloud_apig_gateway" "policy_secondary" {
  network_access_config {
    type = "Intranet"
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  gateway_type = "API"
  payment_type = "PayAsYouGo"
  gateway_name = format("%%s-gateway-secondary", var.name)
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = false
    }
  }
}

resource "alicloud_apig_http_api" "policy" {
  http_api_name = format("%%s-api", var.name)
  protocols     = ["HTTP"]
  type          = "Rest"
  base_path     = format("/%%s", var.name)
}

resource "alicloud_apig_route" "policy_primary" {
  route_name  = format("%%s-route-primary", var.name)
  http_api_id = alicloud_apig_http_api.policy.id
  environment_info {
    environment_id = alicloud_apig_gateway.policy_primary.environments.0.environment_id
  }
  match {
    path {
      type  = "Exact"
      value = format("/%%s-primary", var.name)
    }
    methods = ["GET"]
  }
  backend {
    scene = "Mock"
  }
}

resource "alicloud_apig_route" "policy_secondary" {
  route_name  = format("%%s-route-secondary", var.name)
  http_api_id = alicloud_apig_http_api.policy.id
  environment_info {
    environment_id = alicloud_apig_gateway.policy_secondary.environments.0.environment_id
  }
  match {
    path {
      type  = "Exact"
      value = format("/%%s-secondary", var.name)
    }
    methods = ["GET"]
  }
  backend {
    scene = "Mock"
  }
}
`, name)
}

func testAccApigPolicyImportID(state *terraform.State) (string, error) {
	resourceState, ok := state.RootModule().Resources["alicloud_apig_policy.default"]
	if !ok || resourceState.Primary == nil {
		return "", fmt.Errorf("alicloud_apig_policy.default was not found in state")
	}
	attachmentID := resourceState.Primary.Attributes["policy_attachment_id"]
	if attachmentID == "" {
		return "", fmt.Errorf("alicloud_apig_policy.default has no policy_attachment_id")
	}
	return resourceState.Primary.ID + ":" + attachmentID, nil
}

func testAccApigAuthorizationRulesImportID(state *terraform.State) (string, error) {
	resourceState, ok := state.RootModule().Resources["alicloud_apig_consumer_authorization_rules.default"]
	if !ok || resourceState.Primary == nil {
		return "", fmt.Errorf("alicloud_apig_consumer_authorization_rules.default was not found in state")
	}
	resourceIDs := make([]string, 0)
	for key, value := range resourceState.Primary.Attributes {
		if strings.HasPrefix(key, "resource_ids.") && key != "resource_ids.#" && value != "" {
			resourceIDs = append(resourceIDs, value)
		}
	}
	if len(resourceIDs) == 0 {
		return "", fmt.Errorf("alicloud_apig_consumer_authorization_rules.default has no resource_ids")
	}
	sort.Strings(resourceIDs)
	return resourceState.Primary.ID + ":" + strings.Join(resourceIDs, ","), nil
}
