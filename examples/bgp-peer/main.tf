resource "alicloud_bgp_group" "foo" {
    peer_asn = 213
    router_id = "vbr-xxxxxxxxxx"
    description = "test-description11"
    name = "test-name"
    is_fake_asn = true
    auth_key= "dasdasda"
}


resource "alicloud_bgp_peer" "foo" {
    depends_on = [alicloud_bgp_group.foo]
    bgp_group_id = alicloud_bgp_group.foo.id
    peer_ip_address = "192.168.0.10"
}
