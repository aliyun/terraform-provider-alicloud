resource "alicloud_ots_instance" "foo" {
	name = "${var.name}"
	description = "${var.name}"
	accessed_by = "ConsoleOrVpc"
	instance_type = "Capacity"
	tags {
		Updated = "TF"
		For = "acceptance test"
		From = "TF"
	}
}
variable "name" {
	default = "tf-testAcc12345"
}