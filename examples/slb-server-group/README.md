### SLB Server Group Example

The example create SLB Server Group

### SLB Server Group parameter describe

The following arguments are supported:

* `load_balancer_id` - (Required, ForceNew) The Load Balancer ID which is used to launch a new virtual server group.
* `name` - (Optional) Name of the virtual server group. Our plugin provides a default name: "tf-server-group".
* `servers` - (Required) A list of ECS instances to be added. At most 20 ECS instances can be supported in one resource. It contains three sub-fields as `Block server` follows.

### Block servers

The servers mapping supports the following:

* `server_ids` - (Deprecated) Field 'server_ids' has been deprecated from provider version 1.93.0. Use 'server_id' replaces it.
* `server_id` - (Required, Available in 1.93.0+) Backend server ID (ECS instance ID).
* `port` - (Required) The port used by the backend server. Valid value range: [1-65535].
* `weight` - (Optional) Weight of the backend server. Valid value range: [0-100]. Default to 100.

### Attributes Reference

The following attributes are exported:

* `id` - The ID of the virtual server group.
* `load_balancer_id` - The Load Balancer ID which is used to launch a new virtual server group.
* `name` - The name of the virtual server group.
* `servers` - A list of ECS instances that have be added.


### Get up and running

* Planning phase

        terraform plan

* Apply phase

        terraform apply

* Destroy

        terraform destroy