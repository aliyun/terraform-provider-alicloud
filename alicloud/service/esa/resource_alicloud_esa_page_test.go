package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ESA Page. >>> Resource test cases, automatically generated.
// Case resource_Page_test_html
func TestAccAliCloudESAPageresource_Page_test_html(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_page.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESAPageresource_Page_test_htmlMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaPage")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAPage%d", acctest.DefaultRegionToTest, rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESAPageresource_Page_test_htmlBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "test resource html page",
					"content_type": "text/html",
					"content":      "PCFET0NUWVBFIGh0bWw+CjxodG1sIGxhbmc9InpoLUNOIj4KICA8aGVhZD4KICAgIDx0aXRsZT40MDMgRm9yYmlkZGVuPC90aXRsZT4KICA8L2hlYWQ+CiAgPGJvZHk+CiAgICA8aDE+NDAzIEZvcmJpZGRlbjwvaDE+CiAgPC9ib2R5Pgo8L2h0bWw+",
					"page_name":    "resource_test_html_page",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test modify resource html page",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "test resource json page",
					"content_type": "application/json",
					"content":      "ewogICJNZXNzYWdlIjogIlNvcnJ5LCB5b3UgaGF2ZSBiZWVuIGJsb2NrZWQiLAogICJSZXF1ZXN0SUQiOiB7e1JFUVVFU1RfSUR9fSwKICAiUnVsZUlEIjp7e1JVTEVfSUR9fQp9",
					"page_name":    "resource_test_json_page_1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "test resource html page modify description",
					"content_type": "text/html",
					"content":      "PCFET0NUWVBFIGh0bWw+CjxodG1sIGxhbmc9InpoLUNOIj4KICA8aGVhZD4KICAgIDx0aXRsZT40MDMgRm9yYmlkZGVuPC90aXRsZT4KICA8L2hlYWQ+CiAgPGJvZHk+CiAgICA8aDE+NDAzIEZvcmJpZGRlbjwvaDE+CiAgPC9ib2R5Pgo8L2h0bWw+",
					"page_name":    "resource_test_html_page_1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "test resource html page modify description",
					"content_type": "text/html",
					"content":      "PCFET0NUWVBFIGh0bWw+CjxodG1sIGxhbmc9InpoLUNOIj4KICA8aGVhZD4KICAgIDx0aXRsZT40MDQgTm90IGZvdW5kPC90aXRsZT4KICA8L2hlYWQ+CiAgPGJvZHk+CiAgICA8aDE+NDA0IE5vdCBmb3VuZDwvaDE+CiAgPC9ib2R5Pgo8L2h0bWw+",
					"page_name":    "resource_test_html_page_1",
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

var AliCloudESAPageresource_Page_test_htmlMap = map[string]string{
	"id": acctest.CHECKSET,
}

func AliCloudESAPageresource_Page_test_htmlBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_Page_test_json
func TestAccAliCloudESAPageresource_Page_test_json(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_page.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESAPageresource_Page_test_jsonMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaPage")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAPage%d", acctest.DefaultRegionToTest, rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESAPageresource_Page_test_jsonBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "test resource json page",
					"content_type": "application/json",
					"content":      "ewogICJNZXNzYWdlIjogIlNvcnJ5LCB5b3UgaGF2ZSBiZWVuIGJsb2NrZWQiLAogICJSZXF1ZXN0SUQiOiB7e1JFUVVFU1RfSUR9fSwKICAiUnVsZUlEIjp7e1JVTEVfSUR9fQp9",
					"page_name":    "resource_test_json_page",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "test resource json page",
					"content_type": "application/json",
					"content":      "ewogICJNZXNzYWdlIjogIlNvcnJ5LCB5b3UgaGF2ZSBiZWVuIGJsb2NrZWQiLAogICJSZXF1ZXN0SUQiOiB7e1JFUVVFU1RfSUR9fSwKICAiUnVsZUlEIjp7e1JVTEVfSUR9fQp9",
					"page_name":    "resource_test_json_page_1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "test resource json page modify description",
					"content_type": "application/json",
					"content":      "ewogICJNZXNzYWdlIjogIlNvcnJ5LCB5b3UgaGF2ZSBiZWVuIGJsb2NrZWQiLAogICJSZXF1ZXN0SUQiOiB7e1JFUVVFU1RfSUR9fSwKICAiUnVsZUlEIjp7e1JVTEVfSUR9fQp9",
					"page_name":    "resource_test_json_page_1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "test resource json page modify description",
					"content_type": "application/json",
					"content":      "ewogICJNZXNzYWdlIjogIm5vdCBmb3VuZCIsCiAgIlJlcXVlc3RJRCI6IHt7UkVRVUVTVF9JRH19LAogICJSdWxlSUQiOnt7UlVMRV9JRH19Cn0=",
					"page_name":    "resource_test_json_page_1",
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

var AliCloudESAPageresource_Page_test_jsonMap = map[string]string{
	"id": acctest.CHECKSET,
}

func AliCloudESAPageresource_Page_test_jsonBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ESA Page. <<< Resource test cases, automatically generated.
