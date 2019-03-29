---
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_rule"
sidebar_current: "docs-alicloud-resource-slb-rule"
description: |-
  Provides a Load Banlancer Forwarding Rule Resource and add it to one Listener.
---

# alicloud\_slb\_rule

A forwarding rule is configured in `HTTP`/`HTTPS` listener and it used to listen a list of backend servers which in one specified virtual backend server group.
You can add forwarding rules to a listener to forward requests based on the domain names or the URL in the request.

-> **NOTE:** One virtual backend server group can be attached in multiple forwarding rules.

-> **NOTE:** At least one "Domain" or "Url" must be specified when creating a new rule.

-> **NOTE:** Having the same 'Domain' and 'Url' rule can not be created repeatedly in the one listener.

-> **NOTE:** Rule only be created in the `HTTP` or `HTTPS` listener.

-> **NOTE:** Only rule's virtual server group can be modified.

## Example Usage

```
# Create a new load balancer and virtual rule

resource "alicloud_slb" "instance" {
  name = "new-slb"
  vswitch_id = "<one vswitch id>"
}

resource "alicloud_slb_listener" "listener" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  protocol = "http"
  ...
}

resource "alicloud_slb_server_group" "group" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  ...
}

resource "alicloud_slb_rule" "rule" {
  count = 2
  load_balancer_id = "${alicloud_slb.instance.id}"
  frontend_port = "${alicloud_slb_listener.listener.frontend_port}"
  name = "from-tf"
  domain = "*.test.com"
  url = "/image/${count.index}"
  server_group_id = "${alicloud_slb_server_group.group.id}"
}

```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The Load Balancer ID which is used to launch the new forwarding rule.
* `name` - (Optional, ForceNew) Name of the forwarding rule. Our plugin provides a default name: "tf-slb-rule".
* `frontend_port` - (Required, ForceNew) The listener frontend port which is used to launch the new forwarding rule. Valid range: [1-65535].
* `domain` - (Optional, ForceNew) Domain name of the forwarding rule. It can contain letters a-z, numbers 0-9, hyphens (-), and periods (.),
and wildcard characters. The following two domain name formats are supported:
   - Standard domain name: www.test.com
   - Wildcard domain name: *.test.com. wildcard (*) must be the first character in the format of (*.)
* `url` - (Optional, ForceNew) Domain of the forwarding rule. It must be 2-80 characters in length. Only letters a-z, numbers 0-9,
and characters '-' '/' '?' '%' '#' and '&' are allowed. URLs must be started with the character '/', but cannot be '/' alone.
* `server_group_id` - (Required) ID of a virtual server group that will be forwarded.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the forwarding rule.
* `load_balancer_id` - The Load Balancer ID in which forwarding rule belongs.
* `name` - The name of the forwarding rule.
* `forntend_port` - The listener port in which forwarding rule belongs.
* `domain` - The domain name of the forwarding rule.
* `url` - The url of the forwarding rule.
* `server_group_id` - The Id of the virtual server group.

## Import

Load balancer forwarding rule can be imported using the id, e.g.

```
$ terraform import alicloud_slb_rule.example rule-abc123456
```
