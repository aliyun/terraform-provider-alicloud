resource "alicloud_ros_template" "example" {
  template_name = "example_value"
  template_body = <<EOF
    {
    	"ROSTemplateFormatVersion": "2015-09-01"
    }
    EOF
}

