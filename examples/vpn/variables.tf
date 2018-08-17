variable "vpc" {
  default = "vpc-2zek0m2zqtj98bulj4t6m"
}

variable "bandwidth" {
  default = 10
}

variable "long_name" {
  default = "alicloud_vpn"
}

variable "auto_pay" {
   default = true
}

variable "instance_charge_type" {
   default = "postpaid"
}

variable "cgw_name" {
    default = "tf-cgw_test_01"
}

variable "cgw_ip_address" {
    default = "39.104.22.228"
}

variable "cgw_description" {
    default = "tf-test_to_ap_southeast_5"
}

variable "vco_name" {
    default = "tf-vco_test1"
}

variable "local_subnet" {
    default = "192.168.0.0/16"
}

variable "remote_subnet" {
    default = "10.0.0.0/8"
}

variable "effect_immediately" {
    default = true
}

variable "ike_config" {
    default = "{\"IkeAuthAlg\":\"sha1\",\"IkeEncAlg\":\"aes\",\"IkeVersion\":\"ikev2\",\"IkeMode\":\"aggressive\",\"IkeLifetime\":86400,\"Psk\":\"tf-testvpn1\",\"IkePfs\":\"group2\"}"
}

variable "ipsec_config" {
    default = "{\"IpsecPfs\":\"group2\",\"IpsecEncAlg\":\"aes\",\"IpsecAuthAlg\":\"sha1\",\"IpsecLifetime\":86400}"
}

variable "ssl_server_name" {
    default = "ssl_vpn_server_2"
}

variable "ssl_vpn_client_pool_1" {
    default = "172.16.10.0/24"
}

variable "ssl_vpn_local_subnet_1" {
    default = "192.168.8.0/21"
}

variable "ssl_vpn_proto" {
    default = "UDP"
}

variable "ssl_vpn_cipher" {
    default = "AES-192-CBC"
}

variable "ssl_vpn_port" {
    default = 1194
}

variable "ssl_vpn_compress" {
    default = false
}

variable "ssl_vpn_client_cert_1" {
    default = "tf-ssl_vpn_client_cert_1"
}