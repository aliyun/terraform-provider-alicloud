output "name" {
  description = "Two instances on a single account in the same region cannot have the same name and the number of instances in the same region cannot exceed 8. The length must be 3 to 64 characters. Chinese characters, English letters digits and hyphen are allowed."
  value       = "${alicloud_ons_instance.instance.name}"
}

output "instance_id" {
  description = "The ID used to identify ons.instance resource.Instance_id is unique. Any two instances will not have the same instance_id, even if they are in the same region."
  value       = "${alicloud_ons_instance.instance.id}"
}

output "remark" {
  description = "This attribute is a concise description of instance."
  value       = "${alicloud_ons_instance.instance.remark}"
}

