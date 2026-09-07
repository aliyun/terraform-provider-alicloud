### SLB ACL Example

The example create SLB ACL

### SLB Acl parameter describe

The following arguments are supported:

* `name` - (Required) Name of the access control list.
* `ip_version` - (Optional, ForceNew) The IP Version of access control list is the type of its entry (IP addresses or CIDR blocks). It values ipv4/ipv6. Our plugin provides a default ip_version: "ipv4".
* `entry_list` - (Optional) A list of entry (IP addresses or CIDR blocks) to be added. At most 50 etnry can be supported in one resource. It contains two sub-fields as `Entry Block` follows.

#### Entry Block

The entry mapping supports the following:

* `entry` - (Required) An IP addresses or CIDR blocks.
* `comment` - (Optional) the comment of the entry.

### Attributes Reference

The following attributes are exported:

* `id` - The Id of the access control list.


### Get up and running

* Planning phase
        terraform plan

* Apply phase
        terraform apply

* Destroy
        terraform destroy