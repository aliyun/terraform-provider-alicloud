package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA List. >>> Resource test cases, automatically generated.
// Case resource_List_test_ip
func TestAccAliCloudESAListresource_List_test_ip(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_list.default"
	ra := resourceAttrInit(resourceId, AliCloudESAListresource_List_test_ipMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaList")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAList%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAListresource_List_test_ipBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test ip list",
					"kind":        "ip",
					"items": []string{
						"10.1.1.1",
						"10.1.1.2",
						"10.1.1.3",
					},
					"name": "resource_test_ip_list",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test ip list",
					"items": []string{
						"10.1.1.1",
						"10.1.1.2",
						"10.1.1.3",
					},
					"name": "resource_test_ip_list_modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test ip list-modify",
					"items": []string{
						"10.1.1.1",
						"10.1.1.2",
						"10.1.1.3",
					},
					"name": "resource_test_ip_list_modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test ip list-modify",
					"items": []string{
						"10.1.1.4",
						"10.1.1.5",
					},
					"name": "resource_test_ip_list_modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test ip list modify2",
					"items": []string{
						"10.1.1.6",
						"10.1.1.7",
					},
					"name": "resource_test_ip_list_modify2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test ip list modify2",
					"items":       []string{},
					"name":        "resource_test_ip_list_modify2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudESAListresource_List_test_ipMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAListresource_List_test_ipBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_List_test_asn
func TestAccAliCloudESAListresource_List_test_asn(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_list.default"
	ra := resourceAttrInit(resourceId, AliCloudESAListresource_List_test_asnMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaList")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAList%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAListresource_List_test_asnBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test asn list",
					"kind":        "asn",
					"items": []string{
						"652312",
						"652313",
					},
					"name": "resource_test_asn_list",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test asn list",
					"items": []string{
						"652312",
						"652313",
					},
					"name": "resource_test_asn_list_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test asn list update",
					"items": []string{
						"652312",
						"652313",
					},
					"name": "resource_test_asn_list_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test asn list update",
					"items": []string{
						"652312",
						"652313",
						"652314",
					},
					"name": "resource_test_asn_list_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test asn list update2",
					"items": []string{
						"652326",
						"65231234",
					},
					"name": "resource_test_asn_list_update2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test asn list update2",
					"items":       []string{},
					"name":        "resource_test_asn_list_update2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudESAListresource_List_test_asnMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAListresource_List_test_asnBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_List_test_host
func TestAccAliCloudESAListresource_List_test_host(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_list.default"
	ra := resourceAttrInit(resourceId, AliCloudESAListresource_List_test_hostMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaList")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAList%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAListresource_List_test_hostBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test host list",
					"kind":        "host",
					"items": []string{
						"api.example.com",
						"test.example.com",
					},
					"name": "resource_test_host_list",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test host list",
					"items": []string{
						"api.example.com",
						"test.example.com",
					},
					"name": "resource_test_host_list_modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test host list modfy",
					"items": []string{
						"api.example.com",
						"test.example.com",
					},
					"name": "resource_test_host_list_modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test host list modify",
					"items": []string{
						"api.example.com",
						"test.example.com",
						"test1.example.com",
					},
					"name": "resource_test_host_list_modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test host list modify2",
					"items": []string{
						"test3.example.com",
						"test4.example.com",
					},
					"name": "resource_test_host_list_modify2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "resource test host list modify2",
					"items":       []string{},
					"name":        "resource_test_host_list_modify2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudESAListresource_List_test_hostMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAListresource_List_test_hostBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ESA List. <<< Resource test cases, automatically generated.
