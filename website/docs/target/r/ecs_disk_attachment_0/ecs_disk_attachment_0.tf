# Create a new ECS disk-attachment and use it attach one disk to a new instance.
resource "alicloud_security_group" "ecs_sg" {
  name        = "terraform-test-group"
  description = "New security group"
}

resource "alicloud_ecs_disk" "ecs_disk" {
  availability_zone = "cn-beijing-a"
  size              = "50"
  tags = {
    Name = "TerraformTest-disk"
  }
}

resource "alicloud_instance" "ecs_instance" {
  image_id             = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
  instance_type        = "ecs.n4.small"
  availability_zone    = "cn-beijing-a"
  security_groups      = [alicloud_security_group.ecs_sg.id]
  instance_name        = "Hello"
  internet_charge_type = "PayByBandwidth"
  tags = {
    Name = "TerraformTest-instance"
  }
}

resource "alicloud_ecs_disk_attachment" "ecs_disk_att" {
  disk_id     = alicloud_ecs_disk.ecs_disk.id
  instance_id = alicloud_instance.ecs_instance.id
}
