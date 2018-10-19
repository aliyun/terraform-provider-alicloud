### SLB CA Certificate  Example

The example create SLB CA Certificate for SLB listener of the protocol "https".

### SLB CA Certificate parameter describe

The following arguments are supported:

* `name` - (Optional) Name of the Server Certificate.
* `ca_certificate` - (Required, ForceNew) the content of the ssl certificate. where `alicloud_certificate_id` is null, it is required, otherwise it is ignored.

### Attributes Reference

The following attributes are exported:

* `id` - The Id of CA Certificate (CA Certificate).


### Get up and running

* Planning phase
        terraform plan

* Apply phase
        terraform apply

* Destroy
        terraform destroy