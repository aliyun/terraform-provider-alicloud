---
subcategory: "Anti-DDoS Pro (DdosBgp)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddos_bgp_policy"
description: |-
  Provides a Alicloud Ddos Bgp Policy resource.
---

# alicloud_ddos_bgp_policy

Provides a Ddos Bgp Policy resource.

Ddos protection policy.

For information about Ddos Bgp Policy and how to use it, see [What is Policy](https://www.alibabacloud.com/help/en/anti-ddos/anti-ddos-origin/developer-reference/api-ddosbgp-2018-07-20-createpolicy).

-> **NOTE:** Available since v1.226.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ddos_bgp_policy&exampleId=b4414d58-5e25-2ea6-d88e-5be5bc4125d9ae989720&activeTab=example&spm=docs.r.ddos_bgp_policy.0.b4414d585e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf_exampleacc_bgp32594"
}

variable "policy_name" {
  default = "example_l4_policy"
}

resource "alicloud_ddos_bgp_policy" "default" {
  content {
    enable_defense = "false"
    layer4_rule_list {
      method  = "hex"
      match   = "1"
      action  = "1"
      limited = "0"
      condition_list {
        arg      = "3C"
        position = "1"
        depth    = "2"
      }
      name     = "11"
      priority = "10"
    }
  }

  type        = "l4"
  policy_name = "tf_exampleacc_bgp32594"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ddos_bgp_policy&spm=docs.r.ddos_bgp_policy.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `content` - (Optional) Configuration Content See [`content`](#content) below.
* `policy_name` - (Required) The name of the resource
* `type` - (Required, ForceNew) Type

### `content`

The content supports the following:
* `black_ip_list_expire_at` - (Optional) Blacklist and whitelist timeout.
* `enable_defense` - (Optional) Whether to enable L4 protection.
* `enable_drop_icmp` - (Optional) Switch to discard ICMP.
* `enable_intelligence` - (Optional) Whether the intelligent switch is on.
* `finger_print_rule_list` - (Optional) Fingerprint Rules. See [`finger_print_rule_list`](#content-finger_print_rule_list) below.
* `intelligence_level` - (Optional) Smart mode. Valid values: weak, hard, and default.
* `layer4_rule_list` - (Optional) L4 protection rules. See [`layer4_rule_list`](#content-layer4_rule_list) below.
* `port_rule_list` - (Optional) Port Rule List. See [`port_rule_list`](#content-port_rule_list) below.
* `reflect_block_udp_port_list` - (Optional) Reflective port filtering.
* `region_block_country_list` - (Optional) List of Regional Banned Countries.
* `region_block_province_list` - (Optional) List of Prohibited Provinces by Region.
* `source_block_list` - (Optional) Source pull Black. See [`source_block_list`](#content-source_block_list) below.
* `source_limit` - (Optional) Do not fill in when the source speed limit is deleted. See [`source_limit`](#content-source_limit) below.
* `whiten_gfbr_nets` - (Optional) Add white high protection back to source network segment switch.

### `content-finger_print_rule_list`

The content-finger_print_rule_list supports the following:
* `dst_port_end` - (Required) End of destination port 0-65535.
* `dst_port_start` - (Required) Destination Port start 0-65535.
* `finger_print_rule_id` - (Optional) The UUID of the rule is required to be deleted and modified, and it is not required to be created.
* `match_action` - (Required) Actions. Currently, drop, accept, session_rate, and ip_rate are supported.
* `max_pkt_len` - (Required) Maximum bag length.
* `min_pkt_len` - (Required) Minimum package length.
* `offset` - (Optional) Offset.
* `payload_bytes` - (Optional) Load match, hexadecimal string; Similar to 'abcd'.
* `protocol` - (Required) Protocol, tcp or udp.
* `rate_value` - (Optional) Speed limit value 1-100000.
* `seq_no` - (Required) Serial number 1-100 ● Affects the order issued by the bottom layer ● The larger the number, the lower it is.
* `src_port_end` - (Required) Source Port end 0-65535.
* `src_port_start` - (Required) Source port start 0-65535.

### `content-layer4_rule_list`

The content-layer4_rule_list supports the following:
* `action` - (Required) 1 for observation 2 for blocking.
* `condition_list` - (Required) Matching Condition. See [`condition_list`](#content-layer4_rule_list-condition_list) below.
* `limited` - (Required) .
* `match` - (Required) 0 indicates that the condition is not met 1 indicates that the condition is met.
* `method` - (Required) Char indicates a string match hex match.
* `name` - (Required) Rule Name.
* `priority` - (Required) 1-100, priority, the lower the number, the higher the priority.

### `content-port_rule_list`

The content-port_rule_list supports the following:
* `dst_port_end` - (Required) End of destination port 0-65535.
* `dst_port_start` - (Required) Destination Port start 0-65535.
* `match_action` - (Required, ForceNew) Action. Currently, only drop is supported.
* `port_rule_id` - (Optional) Rule UUID is required to be deleted and modified, and is not required to be created.
* `protocol` - (Required) Protocol, tcp or udp.
* `seq_no` - (Required) Serial number 1-100 ● Affects the order issued by the bottom layer ● The larger the number, the lower it is.
* `src_port_end` - (Required) Source Port end 0-65535.
* `src_port_start` - (Required) Source port start 0-65535.

### `content-source_block_list`

The content-source_block_list supports the following:
* `block_expire_seconds` - (Required) Statistical cycle range 60-1200.
* `every_seconds` - (Required) The time (unit second) for automatically releasing the black after triggering the speed limit is 60~2592000.
* `exceed_limit_times` - (Required) The number of times the speed limit is exceeded in a statistical period ranges from 1 to 1200.
* `type` - (Required) The type of black is optional source PPS speed limit Black: 3 source BPS speed limit Black: 4 SYNPPS speed limit Black: 5 SYNBPS speed limit Black: 6.

### `content-source_limit`

The content-source_limit supports the following:
* `bps` - (Optional) bps range 1024~268435456.
* `pps` - (Optional) Pps range 32~500000.
* `syn_bps` - (Optional) SynBps range 1024~268435456.
* `syn_pps` - (Optional) SynPps range 1~100000.

### `content-layer4_rule_list-condition_list`

The content-layer4_rule_list-condition_list supports the following:
* `arg` - (Required) Matching target character.
* `depth` - (Required) Depth of Matching.
* `position` - (Required) Position to start matching, starting from 0.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Policy.
* `update` - (Defaults to 5 mins) Used when update the Policy.

## Import

Ddos Bgp Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_ddos_bgp_policy.example <id>
```