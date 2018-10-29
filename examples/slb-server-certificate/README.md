### SLB Server Certificate (ssl Certificate) Example

The example create SLB Server Certificate for SLB listener of the protocol "https".

### SLB Server Certificate parameter describe

The following arguments are supported:

* `name` - (Optional) Name of the Server Certificate.
* `server_certificate` - (Optional, ForceNew) the content of the ssl certificate. where `alicloud_certificate_id` is null, it is required, otherwise it is ignored.
* `private_key` - (Optional, ForceNew) the content of privat key of the ssl certificate specified by `server_certificate`. where `alicloud_certificate_id` is null, it is required, otherwise it is ignored.
* `alicloud_certificate_id` - (Optional) an id of server certificate ssued/proxied by alibaba cloud. but it is not supported on the international site  of alibaba cloud now.
* `alicloud_certificate_name`- (Optional) the name of the certificate specified by `alicloud_certificate_id`.but it is not supported on the international site  of alibaba cloud now.

### Attributes Reference

The following attributes are exported:

* `id` - The Id of Server Certificate (SSL Certificate).

### Get up and running

* Planning phase
        terraform plan

* Apply phase
        terraform apply

* Destroy
        terraform destroy