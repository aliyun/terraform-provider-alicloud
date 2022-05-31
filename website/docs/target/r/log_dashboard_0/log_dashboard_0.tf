resource "alicloud_log_project" "default" {
  name        = "tf-project"
  description = "tf unit test"
}
resource "alicloud_log_store" "default" {
  project          = "tf-project"
  name             = "tf-logstore"
  retention_period = "3000"
  shard_count      = 1
}
resource "alicloud_log_dashboard" "example" {
  project_name   = "tf-project"
  dashboard_name = "tf-dashboard"
  char_list      = <<EOF
  [
    {
      "action": {},
      "title":"new_title",
      "type":"map",
      "search":{
        "logstore":"tf-logstore",
        "topic":"new_topic",
        "query":"* | SELECT COUNT(name) as ct_name, COUNT(product) as ct_product, name,product GROUP BY name,product",
        "start":"-86400s",
        "end":"now"
      },
      "display":{
        "xAxis":[
          "ct_name"
        ],
        "yAxis":[
          "ct_product"
        ],
        "xPos":0,
        "yPos":0,
        "width":10,
        "height":12,
        "displayName":"xixihaha911"
      }
    }
  ]
EOF
}
