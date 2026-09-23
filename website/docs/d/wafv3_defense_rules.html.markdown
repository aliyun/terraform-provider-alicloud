---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_wafv3_defense_rules"
sidebar_current: "docs-alicloud-datasource-wafv3-defense-rules"
description: |-
  Provides a list of Wafv3 Defense Rule owned by an Alibaba Cloud account.
---

# alicloud_wafv3_defense_rules

This data source provides Wafv3 Defense Rule available to the user.[What is Defense Rule](https://next.api.alibabacloud.com/document/waf-openapi/2021-10-01/CreateDefenseRule)

-> **NOTE:** Available since v1.283.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "region_id" {
  default = "cn-hangzhou"
}

resource "alicloud_wafv3_instance" "defaultnxb04D" {
}

resource "alicloud_wafv3_defense_template" "defaultfIoHt5" {
  status                = "1"
  description           = "testCreate"
  instance_id           = alicloud_wafv3_instance.defaultnxb04D.id
  defense_template_name = "1782219650"
  template_origin       = "custom"
  defense_scene         = "custom_acl"
  template_type         = "user_custom"
}

resource "alicloud_wafv3_address_book" "default9dtEmt" {
  description       = "test"
  instance_id       = alicloud_wafv3_instance.defaultnxb04D.id
  address_book_name = "1782219650"
  address_list      = ["100.100.100.100/32", "101.101.101.101/32", "102.102.102.102/32"]
  address_book_type = "ip"
}

resource "alicloud_wafv3_address_book" "defaultSB0uHV" {
  description       = "test"
  instance_id       = alicloud_wafv3_instance.defaultnxb04D.id
  address_book_name = "1782219650"
  address_list      = ["100.100.100.100/32", "101.101.101.101/32", "102.102.102.102/32"]
  address_book_type = "ip"
}


resource "alicloud_wafv3_defense_rule" "default" {
  defense_origin = "custom"
  instance_id    = alicloud_wafv3_instance.defaultnxb04D.id
  config {
    rule_action = "block"
    conditions {
      op_value = "contain"
      values   = "abc"
      key      = "URL"
    }
    conditions {
      op_value = "contain"
      values   = "abc"
      key      = "URLPath"
    }
    conditions {
      op_value = "contain"
      values   = "1.1.1.2"
      key      = "IP"
    }
    conditions {
      key      = "IP"
      op_value = "in-list"
      values   = alicloud_wafv3_address_book.default9dtEmt.address_book_id
    }
    cc_status = "0"
    cc_effect = "service"
    rate_limit {
      target    = "remote_addr"
      interval  = "16"
      threshold = "204"
      ttl       = "68"
      status {
        code  = "414"
        count = "333"
      }
      sub_key = "testky1"
    }
    gray_status = "1"
    gray_config {
      gray_target = "remote_addr"
      gray_rate   = "80"
    }
    time_config {
      time_scope = "period"
      time_zone  = "8"
      time_periods {
        start = "1760174804000"
        end   = "1760175804000"
      }
      time_periods {
        start = "1760171804000"
        end   = "1760172804000"
      }
      time_periods {
        start = "1760176804000"
        end   = "1760177804000"
      }
      time_periods {
        start = "1760178804000"
        end   = "1760179804000"
      }
      time_periods {
        start = "1760170804000"
        end   = "1760171804000"
      }
    }
  }
  defense_scene = "custom_acl"
  rule_status   = "1"
  defense_type  = "template"
  template_id   = alicloud_wafv3_defense_template.defaultfIoHt5.defense_template_id
  rule_name     = "custom_acl-create"
}

data "alicloud_wafv3_defense_rules" "default" {
  ids          = ["${alicloud_wafv3_defense_rule.default.id}"]
  defense_type = "template"
  instance_id  = alicloud_wafv3_instance.defaultnxb04D.id
}

output "alicloud_wafv3_defense_rule_example_id" {
  value = data.alicloud_wafv3_defense_rules.default.rules.0.id
}
```

## Argument Reference

The following arguments are supported:
* `defense_type` - (Optional) The protection rule type. Value:
  - `template` (default): indicates the template protection rule.
  - `resource`: indicates the rule of the protected object dimension.
* `instance_id` - (Required) The ID of the Web Application Firewall (WAF) instance.
* `query` - (Optional) The query condition is converted into a string in JSON format constructed with a series of parameters.

-> **NOTE:**  The results of the protection rules vary according to different query conditions. For more information, see **Query parameter description**.

* `rule_type` - (Optional) The protection rule type. Value:
  - `whitelist`: whitelist rules
  - `defense` (default): protection rules
* `ids` - (Optional, Computed) A list of Defense Rule IDs. The value is formulated as `<instance_id>:<defense_type>:<rule_id>`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `enable_details` - (Optional) Default to `false`. Set it to `true` to call the per-rule `DescribeDefenseRule` API and fetch the full `config` block for each rule.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Defense Rule IDs.
* `rules` - A list of Defense Rule Entries. Each element contains the following attributes:
    * `config` - Rule configuration content, in JSON format, constructed with a series of parameters.
        * `abroad_regions` - The regions outside China from which you want to block requests.
        * `account_identifiers` - The policies for account extraction.
            * `decode_type` - The authentication mode.
            * `key` - The field from which you want to extract account information.
            * `position` - The field that stores the decoded account information.
            * `priority` - The priority of the current extraction configuration.
            * `sub_key` - The child match field.
        * `auto_update` - Whether the new Web core protection rules are automatically updated.
        * `bypass_regular_rules` - The list of regular rule IDs that are not detected.
        * `bypass_regular_types` - The regular rule type is not detected.
        * `bypass_tags` - The modules to which the whitelist applies.
        * `cc_effect` - Set the effective range of the speed limit.
        * `cc_status` - Whether to open the speed limit.
        * `cn_regions` - The regions in China from which you want to block requests.
        * `codec_list` - The type to enable decoding.
        * `conditions` - The traffic characteristics of ACL, which are described in JSON format.
            * `key` - Match field.
            * `op_value` - Logical character.
            * `sub_key` - Custom child match fields.
            * `values` - Match the content and fill in the corresponding content as needed.
        * `gray_config` - The canary release configuration for the rule.
            * `gray_rate` - The percentage of traffic for which the canary release takes effect.
            * `gray_sub_key` - The sub-feature of the statistical object.
            * `gray_target` - The type of the canary release object.
        * `gray_status` - Specifies whether to enable canary release for the rule.
        * `mode` - The HTTP flood protection mode.
        * `protocol` - The protocol type of the cached page address.
        * `rate_limit` - The detailed speed limit configuration, which is described in the JSON string format.
            * `interval` - The statistical period, in seconds.
            * `status` - Response code frequency setting.
                * `code` - Required.
                * `count` - The threshold for the number of occurrences.
                * `ratio` - The threshold for the proportion of occurrences (percentage).
            * `sub_key` - The characteristics of the statistical object.
            * `target` - The type of the statistical object.
            * `threshold` - The maximum number of requests that can be sent from a statistical object.
            * `ttl` - The period of time during which you want the specified action to be valid.
        * `remote_addr` - The IP addresses that you want to add to the blacklist.
        * `rule_action` - Protection rule action.
        * `throttle_threhold` - The throttling threshold.
        * `throttle_type` - The throttling method.
        * `time_config` - The scheduled rule configuration.
            * `time_periods` - The time period during which the rule is effective.
                * `end` - The end time of the rule.
                * `start` - The start time of the rule.
            * `time_scope` - The effective period of the rule.
            * `time_zone` - The time zone in which the rule is effective.
            * `week_time_periods` - The periodic time period during which the rule is effective.
                * `day` - The time period of each day when the rule is effective.
                * `day_periods` - The time period of each day when the rule is effective.
                    * `end` - The end time of each day when the rule is effective.
                    * `start` - The start time of each day when the rule is effective.
        * `ua` - The User-Agent string that is allowed for access to the address.
        * `url` - The address of the cached page.
        * `waf_base_config` - The configuration of the Web core protection rules to be modified.
            * `rule_batch_operation_config` - The batch operation on rules.
            * `rule_detail` - The configuration of the Web core protection rules to be modified.
                * `rule_action` - Web core protection rule action.
                * `rule_id` - The ID of the core Web protection rule.
                * `rule_status` - The status of Web core protection rules.
            * `rule_type` - The type of the rule.
    * `defense_origin` - Sources of protection.
    * `defense_scene` - The WAF protection scenario to be created.
    * `defense_type` - The protection rule type.
    * `gmt_modified` - The modification time of the protection rule.
    * `resource` - The protection object corresponding to the rule to be queried.
    * `rule_id` - The protection rule ID.
    * `rule_name` - The rule name.
    * `rule_status` - Protection rule status.
    * `template_id` - The protection template ID of the protection rule to be created.
    * `id` - The ID of the resource supplied above.
