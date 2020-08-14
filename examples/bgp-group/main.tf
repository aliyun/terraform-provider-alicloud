resource "alicloud_bgp_group" "foo" {
    peer_asn = 2
    router_id = var.router_id
    description = "test-description11"
    name = "test-name"
    is_fake_asn = true
    auth_key= "dasdasda"
}
