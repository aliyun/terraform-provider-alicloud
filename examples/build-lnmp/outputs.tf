output "nginx_url" {
  value = "${element(alicloud_eip.default.*.ip_address, 1)}:80/test.php"
}