resource "alicloud_ros_stack_group" "example" {
  stack_group_name = "example_value"
  template_body    = <<EOF
    {
    	"ROSTemplateFormatVersion": "2015-09-01"
    }
    EOF
}

