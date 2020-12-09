package alicloud

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// # Generate a CA cert pair.
//echo '{"CN":"CA","key":{"algo":"rsa","size":2048}}' | cfssl gencert -initca - | cfssljson -bare ca -
//echo '{"signing":{"default":{"expiry":"438000h","usages":["signing","key encipherment","server auth","client auth"]}}}' > ca-config.json

//# Use CA cert to sign a client cert pair.
//export ADDRESS=
//export NAME=kubernetes-admin
//echo '{"CN":"'$NAME'","O": "system:masters","hosts":[""],"key":{"algo":"rsa","size":2048}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare $NAME

const caCert = `-----BEGIN CERTIFICATE-----
MIIC6jCCAdKgAwIBAgIUfSCqJB17C26e20n1wp3QZ0ypbsQwDQYJKoZIhvcNAQEL
BQAwDTELMAkGA1UEAxMCQ0EwHhcNMTkwMTA5MDcwNTAwWhcNMjQwMTA4MDcwNTAw
WjANMQswCQYDVQQDEwJDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AMwEC6JM3KxwWcg8l+bpehMBGAE/gUzglYaSwZt0SxlTYyLjx16HIUZcv2JQX5kP
D0jSLqTrd6C2lc3EaLGyU9SLAdjB6uV5jDKeCJhbAXJEtdHiaI5SuJPd1f/RwUym
7aUcG9puLN18203zvfp+Ot4uaoKlUd/sq+VREiojEz5oGbRrrHIVMD4VyqZidmfL
bOG2Zfz3XSKwcJEs2EuI7nlXYLEtWm2YDZQlC2goLfDbj/QkChjgpooyrlo9+Pp3
JTpydgrE3aTecrpVRRzioUKuUJ4RsXHqLfdDFVpN0GB0JdDfGdjGaptfUVD5Dn7I
Zfe69kTXmH9qwEALEmF1pl8CAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgEGMA8GA1Ud
EwEB/wQFMAMBAf8wHQYDVR0OBBYEFM/M1RKZToaXmUK434EthsoNsNunMA0GCSqG
SIb3DQEBCwUAA4IBAQA2cuOyLAb3mkYrfJsv8PHDuZ/c6TUNRDdHUpq6ItQRKuu3
a5fhmAcJD5MZp67n1gVzVZsQ95qrsduwnSCnwDBSZJP21vcqdeIaG+mjlg/zXXnw
b3qCqbtk27Yuypw91+3Jza834vzEAUvHQiWgVOiiHzFO5jImQhAosTMV838ae/kd
ws6mhM65UuWFg5HLbdM2J/zrjWrhuAJZgR1Kx2eReleUyDg97Bc0SPTBth28tGvH
UjY7X0eHM5vuv6NUOyElHVteY8oQ1f0f06K5K4lB7lJ1SB/9PdxYv/AawQwqJIQr
iPn9wR9zlLX6d0Zge293YJ/HGOsm7UzI65DxfDfp
-----END CERTIFICATE-----`

const clientCert = `-----BEGIN CERTIFICATE-----
MIIDXzCCAkegAwIBAgIUCaHVJZr6XAsKp+DCO2J7M4d2J6kwDQYJKoZIhvcNAQEL
BQAwDTELMAkGA1UEAxMCQ0EwIBcNMTkwMTA5MDczNzAwWhgPMjA2ODEyMjcwNzM3
MDBaMDQxFzAVBgNVBAoTDnN5c3RlbTptYXN0ZXJzMRkwFwYDVQQDExBrdWJlcm5l
dGVzLWFkbWluMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArvxZY4hy
rbAGwjYPx1naFDQi/JDdZkypZ11MGGevSf5516tRjJSgGPm0hCR4kKvoY29QLmk9
c67X0Gsih7xnr2ockl4lL9Mf+RrBDO3w55mnxfZ4PBTsuptBvj+MGu+y7jUySmDW
EH+PAJGG+UUj+CGN2yZrAGHlBSyhtVZ3CzX9WOeNj89cXv0cqdNtTyW+Nx3ny2rj
RI5zUeXfCELiwfNGPUbKMOHro302vqaYeFNzFxaf4aYqcPnG1Nk/LSHLSWXnKlsi
qACgSDZZmYEHtjNAMOC/sH0J3HneiIrmyQ2H1il5ai2ZXQvgRM9T5YPAu92O9ipl
aOdGEfO0q+DVzwIDAQABo4GNMIGKMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAU
BggrBgEFBQcDAQYIKwYBBQUHAwIwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUeC6+
JAmWJLIz3bF3wFgW0y9nUB4wHwYDVR0jBBgwFoAUz8zVEplOhpeZQrjfgS2Gyg2w
26cwCwYDVR0RBAQwAoIAMA0GCSqGSIb3DQEBCwUAA4IBAQAMiu/BpYXwtlmA+cp8
rZlRaBALAhFjEHSh4eJALcYXEEI1b1moWwqPJ8Cs3TGneK9MYhVOAPzsZcqNAzJq
UyrF9LKJPPymWvJztubYvqNBLHorAXGouruyc7Tf3HcHEf7j1bwtiNJN0SqHVYDT
qsUXxParkcOhV52CHGShwtLEmT7Tp1IuOJe+PGUrEkeePCDF+I3e5u3c8zUQ10eV
NOB1xGm2n5rdl8Qyl2aTr4Psmgzhx/OKQpDO4z4vK9ABbqJXxcPiVUrlbflsFniG
wzs74ptWLax2EQHY/lv6mdMFDRt24zyBJUQJ3rlgwRF64w6LgtX7cl9Qo4L0a5mt
SNex
-----END CERTIFICATE-----`

const clientCertKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEArvxZY4hyrbAGwjYPx1naFDQi/JDdZkypZ11MGGevSf5516tR
jJSgGPm0hCR4kKvoY29QLmk9c67X0Gsih7xnr2ockl4lL9Mf+RrBDO3w55mnxfZ4
PBTsuptBvj+MGu+y7jUySmDWEH+PAJGG+UUj+CGN2yZrAGHlBSyhtVZ3CzX9WOeN
j89cXv0cqdNtTyW+Nx3ny2rjRI5zUeXfCELiwfNGPUbKMOHro302vqaYeFNzFxaf
4aYqcPnG1Nk/LSHLSWXnKlsiqACgSDZZmYEHtjNAMOC/sH0J3HneiIrmyQ2H1il5
ai2ZXQvgRM9T5YPAu92O9iplaOdGEfO0q+DVzwIDAQABAoIBAEcAXcTlOKMBKbzj
8sMQ0kwgW5HftfYsZRBr6tR8PcPoXcgE27IPHGcF6xkzMziAiGrQX9h1G1o1N5x5
3Cj3aZrjk3RQfwZIxFBvaqW0ZmuTDWBmAaNfWi7dkG+BmXfUiuXc7+r+H93R5FR2
uC2swEuOUeD6VbByCFtxIKxTyTniooNhHaTRSZ0KLycCHgEQci1EabqnGcYCRnHg
wd8GeU0i70JxFDAdL4qnjOq6nwqpinU9DIp+4jhVdJ/ymvg+o2gwqVHt9qiwHm4x
ns/7P/Qrq+9lgAESVGf5q0r6uDkVTxrZ+KH6S76EpQog929w9GjrA97jd+Eq+xC3
cChEPQECgYEA0N9vNynoIN+T07GMU+c/ivB8foVkJlAYrWah9YJnwAP/hrvW4vCe
G0xzIq1ZRaTNQdq0AxzJ97PICRIV8RjbogzLS+cM8/g7iaqin7Gx0V0QOrUiSkrb
YIgm3yRi+jgNfukyCnPUz4oZmr+w5wAbRDXjEODAzSAEmBVvNzfSlzsCgYEA1neW
Na+xzHUN+ywsjr8/+CZf4W1L2wTTTeN+j+D2Ua/RR54VbFKU1mz/YJlkpGl30XCu
/qkGR8WXldCOXf52kZqtia2rwC6gXEEA9pQS1Sj3lWVlR72PlkjKVyTx+3iIbYGs
YvyiSIHEpOCOPK4vJlMnWLZ1Y96Tizd8x7MDGn0CgYAbDIRTiXrFHw7+wCRjDTRe
YsxMeivBBmhbtEnPCGc1J49kvFiUpQJkmJ7kY7yG11O5boAXUxgYmtCR1CTBRy3S
K4P8PVyhD4luR4mt0o4rhbi/UYuyQUVtl9Qo24ZxzuZ4g+x2DBAIHGM6dg6Lq6jc
SXoxSlnNdpMBuuzfIryD1wKBgExBPPFdxPQTcqMp87XVnmMXEeRPPjdjodYB21BB
BpPI1bqHJMrdGfqbyrmIENa8gVPAoxf89TSzttAX1WbqQTJIMwfO7lBow6/JlRQX
VhLgfBdsc/RsHA+tVfRiOH/XPXriLm8LsI/jRA3zod9Fd5JC4qySQ279BqzrT7yZ
k7LpAoGBAMBVPRGuDhLnVggOUdZ+nzJ61AggsJZd1pEI6/yt0Q/KcQiaNM4yZmMX
KdCIHS7bM3yDVqLYB1XoQ1fXiHpbojf+kkdu8f7WhTiMXLgPYfWpn74U0WqD6rcH
KKXMFi8uMDIellyOaUOsiPhAWnm4GqERwQuc4U4i3Z5k/eiKReWT
-----END RSA PRIVATE KEY-----`

func TestAccAlicloudCSKubernetes_basic(t *testing.T) {
	var v *cs.KubernetesClusterDetail

	resourceId := "alicloud_cs_kubernetes.default"
	ra := resourceAttrInit(resourceId, csKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKubernetes-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSKubernetesConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  name,
					"master_vswitch_ids":    []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default.id}", "${alicloud_vswitch.default.id}"},
					"worker_vswitch_ids":    []string{"${alicloud_vswitch.default.id}"},
					"master_instance_types": []string{"${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.0.id}"},
					"worker_instance_types": []string{"${data.alicloud_instance_types.default1.instance_types.0.id}"},
					"worker_number":         "1",
					"master_disk_category":  "cloud_ssd",
					"worker_disk_size":      "50",
					"password":              "Yourpassword1234",
					"pod_cidr":              "192.168.1.0/24",
					"service_cidr":          "192.168.2.0/24",
					"enable_ssh":            "true",
					"install_cloud_monitor": "true",
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"deletion_protection":   "true",
					"timezone":              "Asia/Shanghai",
					"os_type":               "Linux",
					"platform":              "CentOS",
					"node_port_range":       "30000-32767",
					"cluster_domain":        "cluster.local",
					"custom_san":            "www.terraform.io",
					"rds_instances":         []string{"${alicloud_db_instance.default.id}"},
					"taints":                []map[string]string{{"key": "tf-key1", "value": "tf-value1", "effect": "NoSchedule"}},
					"runtime":               map[string]interface{}{"Name": "docker", "Version": "19.03.5"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                  name,
						"worker_number":         "1",
						"master_disk_category":  "cloud_ssd",
						"worker_disk_size":      "50",
						"password":              "Yourpassword1234",
						"pod_cidr":              "192.168.1.0/24",
						"service_cidr":          "192.168.2.0/24",
						"enable_ssh":            "true",
						"install_cloud_monitor": "true",
						"resource_group_id":     CHECKSET,
						"deletion_protection":   "true",
						"timezone":              "Asia/Shanghai",
						"os_type":               "Linux",
						"platform":              "CentOS",
						"node_port_range":       "30000-32767",
						"cluster_domain":        "cluster.local",
						"custom_san":            "www.terraform.io",
						"rds_instances.#":       "1",
						"taints.#":              "1",
						"taints.0.key":          "tf-key1",
						"taints.0.value":        "tf-value1",
						"taints.0.effect":       "NoSchedule",
						"runtime.Name":          "docker",
						"runtime.Version":       "19.03.5",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"name", "new_nat_gateway", "pod_cidr",
					"service_cidr", "enable_ssh", "password", "install_cloud_monitor", "user_ca", "force_update",
					"master_disk_category", "master_disk_size", "master_instance_charge_type", "master_instance_types",
					"node_cidr_mask", "slb_internet_enabled", "vswitch_ids", "worker_disk_category", "worker_disk_size",
					"worker_instance_charge_type", "worker_instance_types", "log_config",
					"worker_data_disk_category", "worker_data_disk_size", "master_vswitch_ids", "worker_vswitch_ids", "exclude_autoscaler_nodes", "cpu_policy", "proxy_mode", "cluster_domain", "custom_san", "node_port_range", "os_type", "platform", "timezone", "runtime", "taints", "rds_instances"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"new_nat_gateway": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"new_nat_gateway": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"install_cloud_monitor": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"install_cloud_monitor": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "tf-dedicated-k8s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-dedicated-k8s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"new_nat_gateway":       REMOVEKEY,
					"worker_number":         "3",
					"install_cloud_monitor": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"new_nat_gateway":       "true",
						"worker_number":         "3",
						"install_cloud_monitor": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"worker_number": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"worker_number": "1",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCSKubernetes_ca(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "tf-acc-alicloud-cs-kubernetes-userca")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	err = ioutil.WriteFile(tmpFile.Name(), []byte(caCert), 0644)
	if err != nil {
		t.Fatal(err)
	}

	var v *cs.KubernetesClusterDetail

	resourceId := "alicloud_cs_kubernetes.default"
	ra := resourceAttrInit(resourceId, csKubernetesBasicMap)

	serviceFunc := func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKubernetes-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSKubernetesConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.KubernetesSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  name,
					"master_vswitch_ids":    []string{"${alicloud_vswitch.default.id}", "${alicloud_vswitch.default.id}", "${alicloud_vswitch.default.id}"},
					"worker_vswitch_ids":    []string{"${alicloud_vswitch.default.id}"},
					"master_instance_types": []string{"${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.0.id}", "${data.alicloud_instance_types.default.instance_types.0.id}"},
					"worker_instance_types": []string{"${data.alicloud_instance_types.default1.instance_types.0.id}"},
					"worker_number":         "1",
					"master_disk_category":  "cloud_ssd",
					"worker_disk_size":      "50",
					"password":              "Yourpassword1234",
					"pod_cidr":              "192.168.1.0/24",
					"service_cidr":          "192.168.2.0/24",
					"enable_ssh":            "true",
					"install_cloud_monitor": "true",
					"user_ca":               tmpFile.Name(),
					"resource_group_id":     "${alicloud_resource_manager_resource_group.default.id}",
					"deletion_protection":   "false",
					"timezone":              "Asia/Shanghai",
					"os_type":               "Linux",
					"platform":              "CentOS",
					"node_port_range":       "30000-32767",
					"cluster_domain":        "cluster.local",
					"custom_san":            "www.terraform.io",
					"rds_instances":         []string{"${alicloud_db_instance.default.id}"},
					"taints":                []map[string]string{{"key": "tf-key1", "value": "tf-value1", "effect": "NoSchedule"}},
					"runtime":               map[string]interface{}{"Name": "docker", "Version": "19.03.5"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserCA(resourceId, v),
					testAccCheck(map[string]string{
						"name":                  name,
						"worker_number":         "1",
						"master_disk_category":  "cloud_ssd",
						"worker_disk_size":      "50",
						"password":              "Yourpassword1234",
						"pod_cidr":              "192.168.1.0/24",
						"service_cidr":          "192.168.2.0/24",
						"enable_ssh":            "true",
						"install_cloud_monitor": "true",
						"resource_group_id":     CHECKSET,
						"deletion_protection":   "false",
						"timezone":              "Asia/Shanghai",
						"os_type":               "Linux",
						"platform":              "CentOS",
						"node_port_range":       "30000-32767",
						"cluster_domain":        "cluster.local",
						"custom_san":            "www.terraform.io",
						"rds_instances.#":       "1",
						"taints.#":              "1",
						"taints.0.key":          "tf-key1",
						"taints.0.value":        "tf-value1",
						"taints.0.effect":       "NoSchedule",
						"runtime.Name":          "docker",
						"runtime.Version":       "19.03.5",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"name", "new_nat_gateway", "pod_cidr",
					"service_cidr", "enable_ssh", "password", "install_cloud_monitor", "user_ca", "force_update",
					"master_disk_category", "master_disk_size", "master_instance_charge_type", "master_instance_types",
					"node_cidr_mask", "slb_internet_enabled", "vswitch_ids", "worker_disk_category", "worker_disk_size",
					"worker_instance_charge_type", "worker_instance_types", "log_config",
					"worker_data_disk_category", "worker_data_disk_size", "master_vswitch_ids", "worker_vswitch_ids", "exclude_autoscaler_nodes", "cpu_policy", "proxy_mode", "cluster_domain", "custom_san", "node_port_range", "os_type", "platform", "timezone", "runtime", "taints", "rds_instances"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"new_nat_gateway": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"new_nat_gateway": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"install_cloud_monitor": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"install_cloud_monitor": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"new_nat_gateway":       REMOVEKEY,
					"worker_number":         "3",
					"install_cloud_monitor": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"new_nat_gateway":       "true",
						"worker_number":         "3",
						"install_cloud_monitor": "true",
					}),
				),
			},
		},
	})
}

func resourceCSKubernetesConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	data "alicloud_zones" default {
	  available_resource_creation = "VSwitch"
	}
	
	data "alicloud_instance_types" "default" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		cpu_core_count = 2
		memory_size = 4
		kubernetes_node_role = "Master"
	}
	
	data "alicloud_instance_types" "default1" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		cpu_core_count = 2
		memory_size = 4
		kubernetes_node_role = "Worker"
	}

	data "alicloud_resource_manager_resource_groups" "default" {}

	resource "alicloud_vpc" "default" {
	  name = "${var.name}"
	  cidr_block = "10.1.0.0/21"
	}
	
	resource "alicloud_vswitch" "default" {
	  name = "${var.name}"
	  vpc_id = "${alicloud_vpc.default.id}"
	  cidr_block = "10.1.1.0/24"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	
	resource "alicloud_db_instance" "default" {
	  engine               = "MySQL"
	  engine_version       = "5.6"
	  instance_type        = "rds.mysql.s2.large"
	  instance_storage     = "30"
	  instance_charge_type = "Postpaid"
	  instance_name        = "${var.name}"
	  vswitch_id           = alicloud_vswitch.default.id
	  monitoring_period    = "60"
	}

	`, name)
}

func resourceCSKubernetesConfigDependence_multiAZ(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_zones" default {
	  available_resource_creation = "VSwitch"
	}

	data "alicloud_instance_types" "default_m1" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		cpu_core_count = 2
		memory_size = 4
		kubernetes_node_role = "Master"
	}
	data "alicloud_instance_types" "default_m2" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[length(data.alicloud_zones.default.zones)-1], "id")}"
		cpu_core_count = 2
		memory_size = 4
		kubernetes_node_role = "Master"
	}
	data "alicloud_instance_types" "default_m3" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[length(data.alicloud_zones.default.zones)-2], "id")}"
		cpu_core_count = 2
		memory_size = 4
		kubernetes_node_role = "Master"
	}

	data "alicloud_instance_types" "default_w1" {
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		cpu_core_count = 2
		memory_size = 4
		kubernetes_node_role = "Worker"
	}
	data "alicloud_instance_types" "default_w2" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[length(data.alicloud_zones.default.zones)-1], "id")}"
		cpu_core_count = 2
		memory_size = 4
		kubernetes_node_role = "Worker"
	}
	data "alicloud_instance_types" "default_w3" {
		availability_zone = "${lookup(data.alicloud_zones.default.zones[length(data.alicloud_zones.default.zones)-2], "id")}"
		cpu_core_count = 2
		memory_size = 4
		kubernetes_node_role = "Worker"
	}
	resource "alicloud_vpc" "default" {
	  name = "${var.name}"
	  cidr_block = "10.1.0.0/21"
	}

	resource "alicloud_vswitch" "default1" {
	  name = "${var.name}"
	  vpc_id = "${alicloud_vpc.default.id}"
	  cidr_block = "10.1.1.0/24"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}

	resource "alicloud_vswitch" "default2" {
	  name = "${var.name}"
	  vpc_id = "${alicloud_vpc.default.id}"
	  cidr_block = "10.1.2.0/24"
	  availability_zone = "${lookup(data.alicloud_zones.default.zones[length(data.alicloud_zones.default.zones)-1], "id")}"
	}

	resource "alicloud_vswitch" "default3" {
	  name = "${var.name}"
	  vpc_id = "${alicloud_vpc.default.id}"
	  cidr_block = "10.1.3.0/24"
	  availability_zone = "${lookup(data.alicloud_zones.default.zones[length(data.alicloud_zones.default.zones)-2], "id")}"
	}

	resource "alicloud_nat_gateway" "default" {
	  name = "${var.name}"
	  vpc_id = "${alicloud_vpc.default.id}"
	  specification   = "Small"
	}

	resource "alicloud_snat_entry" "default1" {
	  snat_table_id     = "${alicloud_nat_gateway.default.snat_table_ids}"
	  source_vswitch_id = "${alicloud_vswitch.default1.id}"
	  snat_ip           = "${alicloud_eip.default.ip_address}"
	}

	resource "alicloud_snat_entry" "default2" {
	  snat_table_id     = "${alicloud_nat_gateway.default.snat_table_ids}"
	  source_vswitch_id = "${alicloud_vswitch.default2.id}"
	  snat_ip           = "${alicloud_eip.default.ip_address}"
	}

	resource "alicloud_snat_entry" "default3" {
	  snat_table_id     = "${alicloud_nat_gateway.default.snat_table_ids}"
	  source_vswitch_id = "${alicloud_vswitch.default3.id}"
	  snat_ip           = "${alicloud_eip.default.ip_address}"
	}

	resource "alicloud_eip" "default" {
	  name = "${var.name}"
	  bandwidth = "100"
	}

	resource "alicloud_eip_association" "default" {
	  allocation_id = "${alicloud_eip.default.id}"
	  instance_id   = "${alicloud_nat_gateway.default.id}"
	}

	resource "alicloud_log_project" "default" {
	  name       = "${var.name}"
	}
	`, name)
}

var csKubernetesBasicMap = map[string]string{
	"name": CHECKSET,
}

func testAccCheckUserCA(n string, d *cs.KubernetesClusterDetail) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cluster, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cluster not found:%s", n)
		}
		if endpoint, ok := cluster.Primary.Attributes["connections.api_server_internet"]; ok {
			clientCertPair, err := tls.X509KeyPair([]byte(clientCert), []byte(clientCertKey))
			if err != nil {
				return fmt.Errorf("error loading client cert %s", err)
			}
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					Certificates:       []tls.Certificate{clientCertPair},
				},
			}
			client := &http.Client{Transport: tr}
			resp, err := client.Get(endpoint)
			if resp.StatusCode != 200 {
				return fmt.Errorf("accessing endpoint with client cert failed, http code %d", resp.StatusCode)
			}
			return nil
		} else {
			return fmt.Errorf("connections.api_server_internet not found in cluster %s", cluster.Primary.Attributes)
		}
	}
}
