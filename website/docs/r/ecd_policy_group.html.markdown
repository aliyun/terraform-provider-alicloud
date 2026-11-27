---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_policy_group"
sidebar_current: "docs-alicloud-resource-ecd-policy-group"
description: |-
  Provides a Alicloud Elastic Desktop Service (ECD) Policy Group resource.
---

# alicloud_ecd_policy_group

Provides a Elastic Desktop Service (ECD) Policy Group resource.

For information about Elastic Desktop Service (ECD) Policy Group and how to use it, see [What is Policy Group](https://www.alibabacloud.com/help/en/wuying-workspace/developer-reference/api-ecd-2020-09-30-createpolicygroup).

-> **NOTE:** Available since v1.130.0.

-> **NOTE:** Importing a policy group records the domain access control the server returns: `domain_resolve_rule_type` comes from the service, and the rules are recorded as `domain_list` when they carry no descriptions or as `domain_resolve_rule` blocks when they do. Dropping the recorded fields without managing them is a visible removal diff, and applying it disables the control. To take over an imported configuration, manage the same form with matching values; to clear it, set `domain_resolve_rule_type = "off"`. Refresh records the switch the server reports, so a close that failed or was reverted out of band reappears as a plan diff against the configured value and the next apply re-sends it.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecd_policy_group&exampleId=8b6500b3-4da7-1d11-22dc-ef58e8a22b9db4492c0a&activeTab=example&spm=docs.r.ecd_policy_group.0.8b6500b34d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = "terraform-example"
  clipboard         = "read"
  local_drive       = "read"
  usb_redirect      = "off"
  watermark         = "off"

  authorize_access_policy_rules {
    description = "terraform-example"
    cidr_ip     = "1.2.3.45/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = "terraform-example"
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "1.2.3.4/24"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecd_policy_group&spm=docs.r.ecd_policy_group.example&intl_lang=EN_US)

Domain Access Control Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name        = "terraform-example"
  domain_resolve_rule_type = "on"

  domain_resolve_rule {
    domain      = "*.trycloudflare.com"
    policy      = "block"
    description = "block wildcard domain"
  }
}
```

To remove the domain access control configuration, set `domain_resolve_rule_type` to `off` (and drop the `domain_resolve_rule` blocks):

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name        = "terraform-example"
  domain_resolve_rule_type = "off"
}
```

## Argument Reference

The following arguments are supported:

* `authorize_access_policy_rules` - (Optional) The rule of authorize access rule. See [`authorize_access_policy_rules`](#authorize_access_policy_rules) below.
* `authorize_security_policy_rules` - (Optional) The policy rule. See [`authorize_security_policy_rules`](#authorize_security_policy_rules) below.
* `clipboard` - (Optional, Computed) The clipboard policy. Valid values: `off`, `read`, `readwrite`.
* `domain_list` - (Optional) The legacy domain access control list. The value must start with the rule type `[black:]` or `[white:]` followed by comma-separated domains that support the `*` wildcard, e.g. `[black:],*.example.com,www.example.com`.

  -> **NOTE:** `DomainList` is rejected by the service with `InvalidDomainRuleType` when the `[black:]`/`[white:]` prefix is missing. Removing `domain_list` that this configuration managed disables domain access control by sending `DomainResolveRuleType=OFF`, since `ModifyPolicyGroup` keeps omitted parameters and an empty `DomainList` alone cannot clear an existing configuration. Domain access control that was never recorded in state (for example changed out of band on a group this configuration never managed) is left untouched by unrelated updates; domain access control recorded in state — managed here or imported — plans a visible removal when the configuration drops it, and applying that removal sends the OFF. If the server has not converged to `off` yet, refresh keeps mirroring the rules it still reports, so the removal stays a visible diff and the next apply re-sends the OFF. Prefer `domain_resolve_rule_type` together with `domain_resolve_rule` for new configurations. Switching between `domain_list` and `domain_resolve_rule` in a single apply is rejected: disable domain access control first (`domain_resolve_rule_type = "off"` and drop the configured form), apply, then configure the other form.
* `domain_resolve_rule_type` - (Optional, Computed, Available in 2.3.0+) The switch of the domain resolve policy. Valid values: `on`, `off`. Setting it to `off` disables and clears the domain access control configuration; drop the configured `domain_list`/`domain_resolve_rule` fields in the same apply, otherwise the plan is rejected. The switch is refreshed from the value the server reports, so drift against a configured value is detected and re-sent.
* `domain_resolve_rule` - (Optional, Available in 2.3.0+) The list of domain resolve rules. See [`domain_resolve_rule`](#domain_resolve_rule) below.
* `html_access` - (Optional, Computed) The access of html5. Valid values: `off`, `on`.
* `html_file_transfer` - (Optional, Computed) The html5 file transfer. Valid values: `all`, `download`, `off`, `upload`.
* `local_drive` - (Optional, Computed) Local drive redirect policy. Valid values: `readwrite`, `off`, `read`.
* `policy_group_name` - (Optional, Computed) The name of policy group. When omitted, `CreatePolicyGroup` generates the name and the generated value is stored in state.
* `usb_redirect` - (Optional, Computed) The usb redirect policy. Valid values: `off`, `on`.
* `visual_quality` - (Optional, Computed) The quality of visual. Valid values: `high`, `lossless`, `low`, `medium`.
* `watermark` - (Optional, Computed) The watermark policy. Valid values: `off`, `on`.
* `watermark_transparency` - (Optional, Computed) The watermark transparency. Valid values: `DARK`, `LIGHT`, `MIDDLE`.
* `watermark_type` - (Optional) The type of watemark. Valid values: `EndUserId`, `HostName`.
* `recording` - (Optional, Computed, Available in 1.171.0+) Whether to enable screen recording. Valid values: `off`, `all-time`, `period`.
* `recording_start_time` - (Optional, Available in 1.171.0+) The start time of recording, value: `HH:MM:SS`. This return value is meaningful only when the value of `recording` is `period`.
* `recording_end_time` - (Optional, Available in 1.171.0+) The end time of recording, value: `HH:MM:SS`. This return value is meaningful only when the value of `recording` is `period`.
* `recording_fps` - (Optional, Computed, Available in 1.171.0+) The fps of recording. Valid values: `2`, `5`, `10`, `15`.
* `camera_redirect` - (Optional, Computed, Available in 1.171.0+) Whether to enable local camera redirection. Valid values: `on`, `off`.
* `recording_expires` - (Optional, Available in 1.186.0+) The screen recording video retention. Valid values between 30 and 180. This return value is meaningful only when the value of `recording` is `period` or `all-time`.
* `internet_communication_protocol` - (Optional, Computed, Available in 2.2.0+) The network communication protocol. Valid values: `BOTH`, `TCP`.
* `recording_duration` - (Optional, Computed, Available in 2.2.0+) The recording duration. Valid values: `10`, `20`, `30`, `60`.
* `client_types` - (Optional, Available in 2.2.0+) The list of client type access rules used to control which clients can access cloud desktops. See [`client_types`](#client_types) below.
* `end_user_apply_admin_coordinate` - (Optional, Computed, Available in 2.2.0+) Whether to allow end users to apply for administrator assistance. Valid values: `off`, `on`.
* `end_user_group_coordinate` - (Optional, Computed, Available in 2.2.0+) Whether to allow end users to join user groups. Valid values: `off`, `on`.

---

### `client_types`

The client_types supports the following:

* `client_type` - (Optional) The client type. Valid values: `android`, `html5`, `ios`, `macos`, `windows`.
* `status` - (Optional) Whether to allow the client type to access. Valid values: `off`, `on`.

  -> **NOTE:** `client_types` manages only the client types listed in the configuration; types that are not listed are not sent to the API and are left to the service. `DescribePolicyGroups` returns every client type, but unconfigured entries are not written back to state, and entries that omit `status` are not read back either, so out-of-band changes to those values are not detected. `terraform import` does not populate `client_types`; apply the configuration again to manage the client types. Removing all `client_types` blocks calls `ModifyPolicyGroup` without client-type parameters, which the API treats as allowing every client type.

---

### `domain_resolve_rule`

The domain_resolve_rule supports the following:

* `domain` - (Optional) The domain. Must be non-empty when a rule block is configured. The domain supports the `*` wildcard, e.g. `*.example.com`.
* `policy` - (Optional) The resolve policy. Valid values: `allow`, `block`.
* `description` - (Optional) The description of the rule.

  -> **NOTE:** `domain_resolve_rule` is effective only when `domain_resolve_rule_type` is `on`, and at least one rule block is required while the switch is `on`. An `on` switch without rule blocks, rule blocks together with `off`, rule blocks without the switch, and a rule block with an empty `domain` are all rejected at plan time. To remove the domain access control configuration, set `domain_resolve_rule_type` to `off` and drop the rule blocks in the same apply; `ModifyPolicyGroup` keeps omitted parameters, so dropping the blocks without disabling the switch does not clear the rules already stored on the server.

---

### `authorize_security_policy_rules`

The authorize_security_policy_rules supports the following:

* `cidr_ip` - (Optional) The cidrip of security rules.
* `description` - (Optional) The description of security rules.
* `ip_protocol` - (Optional) The ip protocol of security rules.
* `policy` - (Optional) The policy of security rules.
* `port_range` - (Optional) The port range of security rules.
* `priority` - (Optional) The priority of security rules.
* `type` - (Optional) The type of security rules.

---

### `authorize_access_policy_rules`

The authorize_access_policy_rules supports the following:

* `cidr_ip` - (Optional) The cidrip of authorize access rule.
* `description` - (Optional) The description of authorize access rule.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Policy Group.
* `status` - The status of policy.

## Import

Elastic Desktop Service (ECD) Policy Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecd_policy_group.example <id>
```
