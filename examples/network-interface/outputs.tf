output "primary_private_ip" {
    value = "${alicloud_network_interface.eni.private_ip}"
}

output "vswitch_id" {
    value = "${alicloud_network_interface.eni.vswitch_id}"
}

output "network_interface_id" {
    value = "${alicloud_network_interface_attachment.at.network_interface_id}"
}

output "instance_id" {
    value = "${alicloud_network_interface_attachment.at.instance_id}"
}
