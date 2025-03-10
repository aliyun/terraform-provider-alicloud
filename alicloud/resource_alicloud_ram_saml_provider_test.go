package alicloud

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudRAMSamlProvider_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_saml_provider.default"
	ra := resourceAttrInit(resourceId, AliCloudRAMSamlProviderMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ImsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamSamlProvider")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudRamSamlProvider%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRAMSamlProviderBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"saml_provider_name":            name,
					"encodedsaml_metadata_document": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48bWQ6RW50aXR5RGVzY3JpcHRvciBlbnRpdHlJRD0iaHR0cDovL2V4YW1wbGUuYWxpeXVuLmNvbS9leGFtcGxlLWlkcCIgeG1sbnM6bWQ9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDptZXRhZGF0YSI+PG1kOklEUFNTT0Rlc2NyaXB0b3IgV2FudEF1dGhuUmVxdWVzdHNTaWduZWQ9ImZhbHNlIiBwcm90b2NvbFN1cHBvcnRFbnVtZXJhdGlvbj0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOnByb3RvY29sIj48bWQ6S2V5RGVzY3JpcHRvciB1c2U9InNpZ25pbmciPjxkczpLZXlJbmZvIHhtbG5zOmRzPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwLzA5L3htbGRzaWcjIj48ZHM6WDUwOURhdGE+PGRzOlg1MDlDZXJ0aWZpY2F0ZT5NSUlEL3pDQ0F1ZWdBd0lCQWdJRU1yb0tjakFOQmdrcWhraUc5dzBCQVFzRkFEQ0JnREVuTUNVR0ExVUVBeE1lClFXeHBlWFZ1SUZKQlRTQkZlR0Z0Y0d4bElFTmxjblJwWm1sallYUmxNUkF3RGdZRFZRUUxFd2RCYkdsaVlXSmgKTVJBd0RnWURWUVFLRXdkQmJHbGlZV0poTVJFd0R3WURWUVFIRXdoSVlXNW5XbWh2ZFRFUk1BOEdBMVVFQ0JNSQpXbWhsU21saGJtY3hDekFKQmdOVkJBWVRBa05PTUNBWERUSXpNVEl3TkRFeU1EY3dNRm9ZRHpJd05URXdOREl4Ck1USXdOekF3V2pDQmdERW5NQ1VHQTFVRUF4TWVRV3hwZVhWdUlGSkJUU0JGZUdGdGNHeGxJRU5sY25ScFptbGoKWVhSbE1SQXdEZ1lEVlFRTEV3ZEJiR2xpWVdKaE1SQXdEZ1lEVlFRS0V3ZEJiR2xpWVdKaE1SRXdEd1lEVlFRSApFd2hJWVc1bldtaHZkVEVSTUE4R0ExVUVDQk1JV21obFNtbGhibWN4Q3pBSkJnTlZCQVlUQWtOT01JSUJJakFOCkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQW9KVGVndWc0eXFaalNzQzFQUWpzbGxreUxWZHEKcXR0UGFqNmNOYldQdVFRNThSMkF4ZHNYeng4c05lOElLYUFYbW84azdDTFhDcXFLVzNjNEpzRWtTOUcva3B2NApJWFpBOGFVcDBCeXdQUFBocmFjUXd4cmJ5dkhja0dqVUpkNEZrOEVjbVVjNjRrSE5LbjBCaVJpL0NEZlM3MXBuCjh5T3dDNUZPSUlYWXhWMGtTTnNQMnozV2tBbFBXWm1sVkZSd1dxeHhGS2xCTjVpdVhaVHA4dk5rU2htVndBTW8KcjVpb2VBaFdXd0N1L0pvdUhLa3lnbVJnSDNhRjlSRlkrOGZ4NkMzR2hjZktISUszRTFBbVVtWlpjR3NDUCtxNApXeTBuSFp4QStaZEhTeE1OYUJPMm5JbkxJSHVDWHgza096eWpGV3dUaTVGSTlwdE5vNktBay9wRThRSURBUUFCCm8zMHdlekFQQmdOVkhSTUVDREFHQVFIL0FnRURNQjBHQTFVZEpRUVdNQlFHQ0NzR0FRVUZCd01CQmdnckJnRUYKQlFjREFqQUxCZ05WSFE4RUJBTUNBb1F3SFFZRFZSMFJBUUgvQkJNd0VZSUpiRzlqWVd4b2IzTjBod1IvQUFBQgpNQjBHQTFVZERnUVdCQlQ2TXluMnJjMVhEQTZqQkZYWVBOYitGaldMVmpBTkJna3Foa2lHOXcwQkFRc0ZBQU9DCkFRRUFoWHpUUzJJaHZjY3hzSVNzVVNFcldNNDJiQlZESHhTa05EemhPRmd0eGNtNUxuNHdjWXJvdkM3NHZxS1oKUWdQWmpGcWw3YUJTb1ZyNFdseWFaZlVBdHdNL3pZZytJbklUSVpBQ0dhM1VNK3h5V0NLSVhRNGpJVldnNG9QWQpxTStjNWllLzJFVlE0YmhObEQyL0lYZUVEZFd2TXMzdmFyRTFCUE5PQXJZZ2tZTmNER3lDSnA0ZmQ3d3ladWxhCllEZFFIWDdpdUJ1R0JOZFRBajlCUW5xaTJRcTc5RndMVTBrQkFvdUpVVVBPUjBpMGtwZ24vc2dSbHhvaHY1bHgKVTFwYVhtMEZRWHpUUDEvdjV5Y24wM3NVckFUekg2VkRpVlQ2N0NRQjR4MXJpOTFvUVRkWERXN1RvRkVhOGIrOApPdE8wZERMdDlnbCtNMkxYRzJTWnBZTkJoZz09PC9kczpYNTA5Q2VydGlmaWNhdGU+PC9kczpYNTA5RGF0YT48L2RzOktleUluZm8+PC9tZDpLZXlEZXNjcmlwdG9yPjxtZDpOYW1lSURGb3JtYXQ+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6MS4xOm5hbWVpZC1mb3JtYXQ6ZW1haWxBZGRyZXNzPC9tZDpOYW1lSURGb3JtYXQ+PG1kOk5hbWVJREZvcm1hdD51cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6bmFtZWlkLWZvcm1hdDpwZXJzaXN0ZW50PC9tZDpOYW1lSURGb3JtYXQ+PG1kOlNpbmdsZVNpZ25PblNlcnZpY2UgQmluZGluZz0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmJpbmRpbmdzOkhUVFAtUE9TVCIgTG9jYXRpb249Imh0dHA6Ly9leGFtcGxlLmFsaXl1bi5jb20vZXhhbXBsZS1pZHAvc3NvL3NhbWwiLz48bWQ6U2luZ2xlU2lnbk9uU2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHA6Ly9leGFtcGxlLmFsaXl1bi5jb20vZXhhbXBsZS1pZHAvc3NvL3NhbWwiLz48L21kOklEUFNTT0Rlc2NyaXB0b3I+PC9tZDpFbnRpdHlEZXNjcmlwdG9yPg==",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"saml_provider_name":            name,
						"encodedsaml_metadata_document": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encodedsaml_metadata_document": "PD95bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48bWQ6RW50aXR5RGVzY3JpcHRvciBlbnRpdHlJRD0iaHR0cDovL2V4YW1wbGUuYWxpeXVuLmNvbS9leGFtcGxlLWlkcCIgeG1sbnM6bWQ9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDptZXRhZGF0YSI+PG1kOklEUFNTT0Rlc2NyaXB0b3IgV2FudEF1dGhuUmVxdWVzdHNTaWduZWQ9ImZhbHNlIiBwcm90b2NvbFN1cHBvcnRFbnVtZXJhdGlvbj0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOnByb3RvY29sIj48bWQ6S2V5RGVzY3JpcHRvciB1c2U9InNpZ25pbmciPjxkczpLZXlJbmZvIHhtbG5zOmRzPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwLzA5L3htbGRzaWcjIj48ZHM6WDUwOURhdGE+PGRzOlg1MDlDZXJ0aWZpY2F0ZT5NSUlEL3pDQ0F1ZWdBd0lCQWdJRU1yb0tjakFOQmdrcWhraUc5dzBCQVFzRkFEQ0JnREVuTUNVR0ExVUVBeE1lClFXeHBlWFZ1SUZKQlRTQkZlR0Z0Y0d4bElFTmxjblJwWm1sallYUmxNUkF3RGdZRFZRUUxFd2RCYkdsaVlXSmgKTVJBd0RnWURWUVFLRXdkQmJHbGlZV0poTVJFd0R3WURWUVFIRXdoSVlXNW5XbWh2ZFRFUk1BOEdBMVVFQ0JNSQpXbWhsU21saGJtY3hDekFKQmdOVkJBWVRBa05PTUNBWERUSXpNVEl3TkRFeU1EY3dNRm9ZRHpJd05URXdOREl4Ck1USXdOekF3V2pDQmdERW5NQ1VHQTFVRUF4TWVRV3hwZVhWdUlGSkJUU0JGZUdGdGNHeGxJRU5sY25ScFptbGoKWVhSbE1SQXdEZ1lEVlFRTEV3ZEJiR2xpWVdKaE1SQXdEZ1lEVlFRS0V3ZEJiR2xpWVdKaE1SRXdEd1lEVlFRSApFd2hJWVc1bldtaHZkVEVSTUE4R0ExVUVDQk1JV21obFNtbGhibWN4Q3pBSkJnTlZCQVlUQWtOT01JSUJJakFOCkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQW9KVGVndWc0eXFaalNzQzFQUWpzbGxreUxWZHEKcXR0UGFqNmNOYldQdVFRNThSMkF4ZHNYeng4c05lOElLYUFYbW84azdDTFhDcXFLVzNjNEpzRWtTOUcva3B2NApJWFpBOGFVcDBCeXdQUFBocmFjUXd4cmJ5dkhja0dqVUpkNEZrOEVjbVVjNjRrSE5LbjBCaVJpL0NEZlM3MXBuCjh5T3dDNUZPSUlYWXhWMGtTTnNQMnozV2tBbFBXWm1sVkZSd1dxeHhGS2xCTjVpdVhaVHA4dk5rU2htVndBTW8KcjVpb2VBaFdXd0N1L0pvdUhLa3lnbVJnSDNhRjlSRlkrOGZ4NkMzR2hjZktISUszRTFBbVVtWlpjR3NDUCtxNApXeTBuSFp4QStaZEhTeE1OYUJPMm5JbkxJSHVDWHgza096eWpGV3dUaTVGSTlwdE5vNktBay9wRThRSURBUUFCCm8zMHdlekFQQmdOVkhSTUVDREFHQVFIL0FnRURNQjBHQTFVZEpRUVdNQlFHQ0NzR0FRVUZCd01CQmdnckJnRUYKQlFjREFqQUxCZ05WSFE4RUJBTUNBb1F3SFFZRFZSMFJBUUgvQkJNd0VZSUpiRzlqWVd4b2IzTjBod1IvQUFBQgpNQjBHQTFVZERnUVdCQlQ2TXluMnJjMVhEQTZqQkZYWVBOYitGaldMVmpBTkJna3Foa2lHOXcwQkFRc0ZBQU9DCkFRRUFoWHpUUzJJaHZjY3hzSVNzVVNFcldNNDJiQlZESHhTa05EemhPRmd0eGNtNUxuNHdjWXJvdkM3NHZxS1oKUWdQWmpGcWw3YUJTb1ZyNFdseWFaZlVBdHdNL3pZZytJbklUSVpBQ0dhM1VNK3h5V0NLSVhRNGpJVldnNG9QWQpxTStjNWllLzJFVlE0YmhObEQyL0lYZUVEZFd2TXMzdmFyRTFCUE5PQXJZZ2tZTmNER3lDSnA0ZmQ3d3ladWxhCllEZFFIWDdpdUJ1R0JOZFRBajlCUW5xaTJRcTc5RndMVTBrQkFvdUpVVVBPUjBpMGtwZ24vc2dSbHhvaHY1bHgKVTFwYVhtMEZRWHpUUDEvdjV5Y24wM3NVckFUekg2VkRpVlQ2N0NRQjR4MXJpOTFvUVRkWERXN1RvRkVhOGIrOApPdE8wZERMdDlnbCtNMkxYRzJTWnBZTkJoZz09PC9kczpYNTA5Q2VydGlmaWNhdGU+PC9kczpYNTA5RGF0YT48L2RzOktleUluZm8+PC9tZDpLZXlEZXNjcmlwdG9yPjxtZDpOYW1lSURGb3JtYXQ+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6MS4xOm5hbWVpZC1mb3JtYXQ6ZW1haWxBZGRyZXNzPC9tZDpOYW1lSURGb3JtYXQ+PG1kOk5hbWVJREZvcm1hdD51cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6bmFtZWlkLWZvcm1hdDpwZXJzaXN0ZW50PC9tZDpOYW1lSURGb3JtYXQ+PG1kOlNpbmdsZVNpZ25PblNlcnZpY2UgQmluZGluZz0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmJpbmRpbmdzOkhUVFAtUE9TVCIgTG9jYXRpb249Imh0dHA6Ly9leGFtcGxlLmFsaXl1bi5jb20vZXhhbXBsZS1pZHAvc3NvL3NhbWwiLz48bWQ6U2luZ2xlU2lnbk9uU2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHA6Ly9leGFtcGxlLmFsaXl1bi5jb20vZXhhbXBsZS1pZHAvc3NvL3NhbWwiLz48L21kOklEUFNTT0Rlc2NyaXB0b3I+PC9tZDpFbnRpdHlEZXNjcmlwdG9yPg==",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encodedsaml_metadata_document": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudRAMSamlProvider_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_saml_provider.default"
	ra := resourceAttrInit(resourceId, AliCloudRAMSamlProviderMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ImsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamSamlProvider")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAliCloudRamSamlProvider%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRAMSamlProviderBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"saml_provider_name":            name,
					"encodedsaml_metadata_document": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48bWQ6RW50aXR5RGVzY3JpcHRvciBlbnRpdHlJRD0iaHR0cDovL2V4YW1wbGUuYWxpeXVuLmNvbS9leGFtcGxlLWlkcCIgeG1sbnM6bWQ9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDptZXRhZGF0YSI+PG1kOklEUFNTT0Rlc2NyaXB0b3IgV2FudEF1dGhuUmVxdWVzdHNTaWduZWQ9ImZhbHNlIiBwcm90b2NvbFN1cHBvcnRFbnVtZXJhdGlvbj0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOnByb3RvY29sIj48bWQ6S2V5RGVzY3JpcHRvciB1c2U9InNpZ25pbmciPjxkczpLZXlJbmZvIHhtbG5zOmRzPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwLzA5L3htbGRzaWcjIj48ZHM6WDUwOURhdGE+PGRzOlg1MDlDZXJ0aWZpY2F0ZT5NSUlEL3pDQ0F1ZWdBd0lCQWdJRU1yb0tjakFOQmdrcWhraUc5dzBCQVFzRkFEQ0JnREVuTUNVR0ExVUVBeE1lClFXeHBlWFZ1SUZKQlRTQkZlR0Z0Y0d4bElFTmxjblJwWm1sallYUmxNUkF3RGdZRFZRUUxFd2RCYkdsaVlXSmgKTVJBd0RnWURWUVFLRXdkQmJHbGlZV0poTVJFd0R3WURWUVFIRXdoSVlXNW5XbWh2ZFRFUk1BOEdBMVVFQ0JNSQpXbWhsU21saGJtY3hDekFKQmdOVkJBWVRBa05PTUNBWERUSXpNVEl3TkRFeU1EY3dNRm9ZRHpJd05URXdOREl4Ck1USXdOekF3V2pDQmdERW5NQ1VHQTFVRUF4TWVRV3hwZVhWdUlGSkJUU0JGZUdGdGNHeGxJRU5sY25ScFptbGoKWVhSbE1SQXdEZ1lEVlFRTEV3ZEJiR2xpWVdKaE1SQXdEZ1lEVlFRS0V3ZEJiR2xpWVdKaE1SRXdEd1lEVlFRSApFd2hJWVc1bldtaHZkVEVSTUE4R0ExVUVDQk1JV21obFNtbGhibWN4Q3pBSkJnTlZCQVlUQWtOT01JSUJJakFOCkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQW9KVGVndWc0eXFaalNzQzFQUWpzbGxreUxWZHEKcXR0UGFqNmNOYldQdVFRNThSMkF4ZHNYeng4c05lOElLYUFYbW84azdDTFhDcXFLVzNjNEpzRWtTOUcva3B2NApJWFpBOGFVcDBCeXdQUFBocmFjUXd4cmJ5dkhja0dqVUpkNEZrOEVjbVVjNjRrSE5LbjBCaVJpL0NEZlM3MXBuCjh5T3dDNUZPSUlYWXhWMGtTTnNQMnozV2tBbFBXWm1sVkZSd1dxeHhGS2xCTjVpdVhaVHA4dk5rU2htVndBTW8KcjVpb2VBaFdXd0N1L0pvdUhLa3lnbVJnSDNhRjlSRlkrOGZ4NkMzR2hjZktISUszRTFBbVVtWlpjR3NDUCtxNApXeTBuSFp4QStaZEhTeE1OYUJPMm5JbkxJSHVDWHgza096eWpGV3dUaTVGSTlwdE5vNktBay9wRThRSURBUUFCCm8zMHdlekFQQmdOVkhSTUVDREFHQVFIL0FnRURNQjBHQTFVZEpRUVdNQlFHQ0NzR0FRVUZCd01CQmdnckJnRUYKQlFjREFqQUxCZ05WSFE4RUJBTUNBb1F3SFFZRFZSMFJBUUgvQkJNd0VZSUpiRzlqWVd4b2IzTjBod1IvQUFBQgpNQjBHQTFVZERnUVdCQlQ2TXluMnJjMVhEQTZqQkZYWVBOYitGaldMVmpBTkJna3Foa2lHOXcwQkFRc0ZBQU9DCkFRRUFoWHpUUzJJaHZjY3hzSVNzVVNFcldNNDJiQlZESHhTa05EemhPRmd0eGNtNUxuNHdjWXJvdkM3NHZxS1oKUWdQWmpGcWw3YUJTb1ZyNFdseWFaZlVBdHdNL3pZZytJbklUSVpBQ0dhM1VNK3h5V0NLSVhRNGpJVldnNG9QWQpxTStjNWllLzJFVlE0YmhObEQyL0lYZUVEZFd2TXMzdmFyRTFCUE5PQXJZZ2tZTmNER3lDSnA0ZmQ3d3ladWxhCllEZFFIWDdpdUJ1R0JOZFRBajlCUW5xaTJRcTc5RndMVTBrQkFvdUpVVVBPUjBpMGtwZ24vc2dSbHhvaHY1bHgKVTFwYVhtMEZRWHpUUDEvdjV5Y24wM3NVckFUekg2VkRpVlQ2N0NRQjR4MXJpOTFvUVRkWERXN1RvRkVhOGIrOApPdE8wZERMdDlnbCtNMkxYRzJTWnBZTkJoZz09PC9kczpYNTA5Q2VydGlmaWNhdGU+PC9kczpYNTA5RGF0YT48L2RzOktleUluZm8+PC9tZDpLZXlEZXNjcmlwdG9yPjxtZDpOYW1lSURGb3JtYXQ+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6MS4xOm5hbWVpZC1mb3JtYXQ6ZW1haWxBZGRyZXNzPC9tZDpOYW1lSURGb3JtYXQ+PG1kOk5hbWVJREZvcm1hdD51cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6bmFtZWlkLWZvcm1hdDpwZXJzaXN0ZW50PC9tZDpOYW1lSURGb3JtYXQ+PG1kOlNpbmdsZVNpZ25PblNlcnZpY2UgQmluZGluZz0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmJpbmRpbmdzOkhUVFAtUE9TVCIgTG9jYXRpb249Imh0dHA6Ly9leGFtcGxlLmFsaXl1bi5jb20vZXhhbXBsZS1pZHAvc3NvL3NhbWwiLz48bWQ6U2luZ2xlU2lnbk9uU2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHA6Ly9leGFtcGxlLmFsaXl1bi5jb20vZXhhbXBsZS1pZHAvc3NvL3NhbWwiLz48L21kOklEUFNTT0Rlc2NyaXB0b3I+PC9tZDpFbnRpdHlEZXNjcmlwdG9yPg==",
					"description":                   name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"saml_provider_name":            name,
						"encodedsaml_metadata_document": CHECKSET,
						"description":                   name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudRAMSamlProviderMap0 = map[string]string{
	"arn":         CHECKSET,
	"update_date": CHECKSET,
}

func AliCloudRAMSamlProviderBasicDependence0(name string) string {
	return ""
}

func TestUnitAliCloudRAMSamlProvider(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_ram_saml_provider"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_ram_saml_provider"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"saml_provider_name":            "saml_provider_name",
		"encodedsaml_metadata_document": "encodedsaml_metadata_document",
		"description":                   "description",
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"SAMLProvider": map[string]interface{}{
			"SAMLProviderName":            "MockSAMLProviderName",
			"Arn":                         "arn",
			"Description":                 "description",
			"EncodedSAMLMetadataDocument": "encodedsaml_metadata_document",
			"UpdateDate":                  "update_date",
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_ram_saml_provider", "MockSAMLProviderName"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewImsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudRamSamlProviderCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudRamSamlProviderCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudRamSamlProviderCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("MockSAMLProviderName")
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewImsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAliCloudRamSamlProviderUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateUpdateSAMLProviderAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"description", "encodedsaml_metadata_document"} {
			switch p["alicloud_ram_saml_provider"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ram_saml_provider"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudRamSamlProviderUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("UpdateUpdateSAMLProviderNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"description", "encodedsaml_metadata_document"} {
			switch p["alicloud_ram_saml_provider"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_ram_saml_provider"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudRamSamlProviderUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewImsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudRamSamlProviderDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudRamSamlProviderDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudRamSamlProviderDelete(d, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})
	t.Run("DeleteNonRetryableError", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudRamSamlProviderDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	//Read
	t.Run("ReadGetSAMLProviderNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudRamSamlProviderRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})
	t.Run("ReadGetSAMLProviderAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudRamSamlProviderRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

// Test Ram SamlProvider. >>> Resource test cases, automatically generated.
// Case SamlProvider测试 9203
func TestAccAliCloudRamSamlProvider_basic9203(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_saml_provider.default"
	ra := resourceAttrInit(resourceId, AlicloudRamSamlProviderMap9203)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamSamlProvider")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRamSamlProviderBasicDependence9203)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"saml_provider_name":            name,
					"description":                   "for test",
					"encodedsaml_metadata_document": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48bWQ6RW50aXR5RGVzY3JpcHRvciB4bWxuczptZD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOm1ldGFkYXRhIiBlbnRpdHlJRD0iaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tL28vc2FtbDI/aWRwaWQ9QzAza3p3Z2oxIiB2YWxpZFVudGlsPSIyMDI2LTEyLTIzVDA5OjMwOjM4LjAwMFoiPgogIDxtZDpJRFBTU09EZXNjcmlwdG9yIFdhbnRBdXRoblJlcXVlc3RzU2lnbmVkPSJmYWxzZSIgcHJvdG9jb2xTdXBwb3J0RW51bWVyYXRpb249InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpwcm90b2NvbCI+CiAgICA8bWQ6S2V5RGVzY3JpcHRvciB1c2U9InNpZ25pbmciPgogICAgICA8ZHM6S2V5SW5mbyB4bWxuczpkcz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC8wOS94bWxkc2lnIyI+CiAgICAgICAgPGRzOlg1MDlEYXRhPgogICAgICAgICAgPGRzOlg1MDlDZXJ0aWZpY2F0ZT5NSUlEZERDQ0FseWdBd0lCQWdJR0FYM3J4eXpoTUEwR0NTcUdTSWIzRFFFQkN3VUFNSHN4RkRBU0JnTlZCQW9UQzBkdmIyZHNaU0JKCmJtTXVNUll3RkFZRFZRUUhFdzFOYjNWdWRHRnBiaUJXYVdWM01ROHdEUVlEVlFRREV3WkhiMjluYkdVeEdEQVdCZ05WQkFzVEQwZHYKYjJkc1pTQkdiM0lnVjI5eWF6RUxNQWtHQTFVRUJoTUNWVk14RXpBUkJnTlZCQWdUQ2tOaGJHbG1iM0p1YVdFd0hoY05NakV4TWpJMApNRGt6TURNNFdoY05Nall4TWpJek1Ea3pNRE00V2pCN01SUXdFZ1lEVlFRS0V3dEhiMjluYkdVZ1NXNWpMakVXTUJRR0ExVUVCeE1OClRXOTFiblJoYVc0Z1ZtbGxkekVQTUEwR0ExVUVBeE1HUjI5dloyeGxNUmd3RmdZRFZRUUxFdzlIYjI5bmJHVWdSbTl5SUZkdmNtc3gKQ3pBSkJnTlZCQVlUQWxWVE1STXdFUVlEVlFRSUV3cERZV3hwWm05eWJtbGhNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QQpNSUlCQ2dLQ0FRRUF1T0RRUFhKL3BKczdQZWhyK3RZdWhVUHJGNHZvSStHSW1wT1M2bkEwWUM2SHBnUDFyeXo3WVBNaEF1dGQrT2NVCjRDUWF0ZDFrOERtV0dBS09DZEtNQ0JlaG8wY2lsREVZQTk0Mk1GNEhTVDNFNlBGcFhhb3htdjNTRkptMXFHNkhrdklRbzBzRTg1YkUKaGMwN3N2Y0FUQ0JQUTNRNjV1dndPeVhQRXlaaVVzL1BLeEx2NHYzMXBiYjNSdmJoRldMUkw1enVnUGh6M3JqTXpoY05GdUhnUTlKNgoySkw3dHY4azdDTDdLOXFZZWZKdXlYZy9wZXpYU2xXbFIvcFlCV3QxY1l0M2Y5azdHRDlydGkrb3pZN0RxVGZqeGU5QVpjaEM0cHVGClBlSkRSVzdIWDFrT2ROWW1BbWVhcitBQzJaOUhZdmNJTmVlelhCaEd2MXR6cHpqM09RSURBUUFCTUEwR0NTcUdTSWIzRFFFQkN3VUEKQTRJQkFRQ1g4WlY0eGQ4TGFOQUZ3YS9QN2tOeFVadG1jWE5KWFphdW5xM2hxL3o2YWhXc0VlWDlEbk4rUnZJSkRXZUw3azg5KzluZAo2cUh5UWxYS1hmbVNVeGVJeDdwckpuMDg5dzN2UnMzdE1nNThOK0NQOFFEYjFNdWFkdHlYSDduYlBhSVVvU3lhTlphQjUwV2R2R0pzCnd3THJOWHQxZGg0Um1zVG8rZVJHUysvY08zZkE0M1d2WHlHNURsTnJNd3FDbUE4d3VKVkNFaWpsbldDQm5zYm4ybDZUMjgwSVdoMTAKWXROVVpSSXhVdWRTRy9QY0FFYjJMRnY1c1B5eGRPdEFHbytFakp0bEplaFdWbHpTMnNHaVpscEhER1kxellBeDJRcVhVM1RPRTU2TQpHYjJWNEc0YUs4bDEvT0NxTVJOd3VlNzhudFlwS29PZXI2RnZaTE9WaVN4MzwvZHM6WDUwOUNlcnRpZmljYXRlPgogICAgICAgIDwvZHM6WDUwOURhdGE+CiAgICAgIDwvZHM6S2V5SW5mbz4KICAgIDwvbWQ6S2V5RGVzY3JpcHRvcj4KICAgIDxtZDpOYW1lSURGb3JtYXQ+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6MS4xOm5hbWVpZC1mb3JtYXQ6ZW1haWxBZGRyZXNzPC9tZDpOYW1lSURGb3JtYXQ+CiAgICA8bWQ6U2luZ2xlU2lnbk9uU2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL3NhbWwyL2lkcD9pZHBpZD1DMDNrendnajEiLz4KICAgIDxtZDpTaW5nbGVTaWduT25TZXJ2aWNlIEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLVBPU1QiIExvY2F0aW9uPSJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20vby9zYW1sMi9pZHA/aWRwaWQ9QzAza3p3Z2oxIi8+CiAgPC9tZDpJRFBTU09EZXNjcmlwdG9yPgo8L21kOkVudGl0eURlc2NyaXB0b3I+",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"saml_provider_name":            name,
						"description":                   "for test",
						"encodedsaml_metadata_document": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48bWQ6RW50aXR5RGVzY3JpcHRvciB4bWxuczptZD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOm1ldGFkYXRhIiBlbnRpdHlJRD0iaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tL28vc2FtbDI/aWRwaWQ9QzAza3p3Z2oxIiB2YWxpZFVudGlsPSIyMDI2LTEyLTIzVDA5OjMwOjM4LjAwMFoiPgogIDxtZDpJRFBTU09EZXNjcmlwdG9yIFdhbnRBdXRoblJlcXVlc3RzU2lnbmVkPSJmYWxzZSIgcHJvdG9jb2xTdXBwb3J0RW51bWVyYXRpb249InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpwcm90b2NvbCI+CiAgICA8bWQ6S2V5RGVzY3JpcHRvciB1c2U9InNpZ25pbmciPgogICAgICA8ZHM6S2V5SW5mbyB4bWxuczpkcz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC8wOS94bWxkc2lnIyI+CiAgICAgICAgPGRzOlg1MDlEYXRhPgogICAgICAgICAgPGRzOlg1MDlDZXJ0aWZpY2F0ZT5NSUlEZERDQ0FseWdBd0lCQWdJR0FYM3J4eXpoTUEwR0NTcUdTSWIzRFFFQkN3VUFNSHN4RkRBU0JnTlZCQW9UQzBkdmIyZHNaU0JKCmJtTXVNUll3RkFZRFZRUUhFdzFOYjNWdWRHRnBiaUJXYVdWM01ROHdEUVlEVlFRREV3WkhiMjluYkdVeEdEQVdCZ05WQkFzVEQwZHYKYjJkc1pTQkdiM0lnVjI5eWF6RUxNQWtHQTFVRUJoTUNWVk14RXpBUkJnTlZCQWdUQ2tOaGJHbG1iM0p1YVdFd0hoY05NakV4TWpJMApNRGt6TURNNFdoY05Nall4TWpJek1Ea3pNRE00V2pCN01SUXdFZ1lEVlFRS0V3dEhiMjluYkdVZ1NXNWpMakVXTUJRR0ExVUVCeE1OClRXOTFiblJoYVc0Z1ZtbGxkekVQTUEwR0ExVUVBeE1HUjI5dloyeGxNUmd3RmdZRFZRUUxFdzlIYjI5bmJHVWdSbTl5SUZkdmNtc3gKQ3pBSkJnTlZCQVlUQWxWVE1STXdFUVlEVlFRSUV3cERZV3hwWm05eWJtbGhNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QQpNSUlCQ2dLQ0FRRUF1T0RRUFhKL3BKczdQZWhyK3RZdWhVUHJGNHZvSStHSW1wT1M2bkEwWUM2SHBnUDFyeXo3WVBNaEF1dGQrT2NVCjRDUWF0ZDFrOERtV0dBS09DZEtNQ0JlaG8wY2lsREVZQTk0Mk1GNEhTVDNFNlBGcFhhb3htdjNTRkptMXFHNkhrdklRbzBzRTg1YkUKaGMwN3N2Y0FUQ0JQUTNRNjV1dndPeVhQRXlaaVVzL1BLeEx2NHYzMXBiYjNSdmJoRldMUkw1enVnUGh6M3JqTXpoY05GdUhnUTlKNgoySkw3dHY4azdDTDdLOXFZZWZKdXlYZy9wZXpYU2xXbFIvcFlCV3QxY1l0M2Y5azdHRDlydGkrb3pZN0RxVGZqeGU5QVpjaEM0cHVGClBlSkRSVzdIWDFrT2ROWW1BbWVhcitBQzJaOUhZdmNJTmVlelhCaEd2MXR6cHpqM09RSURBUUFCTUEwR0NTcUdTSWIzRFFFQkN3VUEKQTRJQkFRQ1g4WlY0eGQ4TGFOQUZ3YS9QN2tOeFVadG1jWE5KWFphdW5xM2hxL3o2YWhXc0VlWDlEbk4rUnZJSkRXZUw3azg5KzluZAo2cUh5UWxYS1hmbVNVeGVJeDdwckpuMDg5dzN2UnMzdE1nNThOK0NQOFFEYjFNdWFkdHlYSDduYlBhSVVvU3lhTlphQjUwV2R2R0pzCnd3THJOWHQxZGg0Um1zVG8rZVJHUysvY08zZkE0M1d2WHlHNURsTnJNd3FDbUE4d3VKVkNFaWpsbldDQm5zYm4ybDZUMjgwSVdoMTAKWXROVVpSSXhVdWRTRy9QY0FFYjJMRnY1c1B5eGRPdEFHbytFakp0bEplaFdWbHpTMnNHaVpscEhER1kxellBeDJRcVhVM1RPRTU2TQpHYjJWNEc0YUs4bDEvT0NxTVJOd3VlNzhudFlwS29PZXI2RnZaTE9WaVN4MzwvZHM6WDUwOUNlcnRpZmljYXRlPgogICAgICAgIDwvZHM6WDUwOURhdGE+CiAgICAgIDwvZHM6S2V5SW5mbz4KICAgIDwvbWQ6S2V5RGVzY3JpcHRvcj4KICAgIDxtZDpOYW1lSURGb3JtYXQ+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6MS4xOm5hbWVpZC1mb3JtYXQ6ZW1haWxBZGRyZXNzPC9tZDpOYW1lSURGb3JtYXQ+CiAgICA8bWQ6U2luZ2xlU2lnbk9uU2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL3NhbWwyL2lkcD9pZHBpZD1DMDNrendnajEiLz4KICAgIDxtZDpTaW5nbGVTaWduT25TZXJ2aWNlIEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLVBPU1QiIExvY2F0aW9uPSJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20vby9zYW1sMi9pZHA/aWRwaWQ9QzAza3p3Z2oxIi8+CiAgPC9tZDpJRFBTU09EZXNjcmlwdG9yPgo8L21kOkVudGl0eURlc2NyaXB0b3I+",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":                   "for test new",
					"encodedsaml_metadata_document": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48bWQ6RW50aXR5RGVzY3JpcHRvciB4bWxuczptZD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOm1ldGFkYXRhIiBlbnRpdHlJRD0iaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tL28vc2FtbDI/aWRwaWQ9QzAza3p3Z2oxIiB2YWxpZFVudGlsPSIyMDI2LTEyLTI0VDA5OjMwOjM4LjAwMFoiPgogIDxtZDpJRFBTU09EZXNjcmlwdG9yIFdhbnRBdXRoblJlcXVlc3RzU2lnbmVkPSJmYWxzZSIgcHJvdG9jb2xTdXBwb3J0RW51bWVyYXRpb249InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpwcm90b2NvbCI+CiAgICA8bWQ6S2V5RGVzY3JpcHRvciB1c2U9InNpZ25pbmciPgogICAgICA8ZHM6S2V5SW5mbyB4bWxuczpkcz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC8wOS94bWxkc2lnIyI+CiAgICAgICAgPGRzOlg1MDlEYXRhPgogICAgICAgICAgPGRzOlg1MDlDZXJ0aWZpY2F0ZT5NSUlEZERDQ0FseWdBd0lCQWdJR0FYM3J4eXpoTUEwR0NTcUdTSWIzRFFFQkN3VUFNSHN4RkRBU0JnTlZCQW9UQzBkdmIyZHNaU0JKCmJtTXVNUll3RkFZRFZRUUhFdzFOYjNWdWRHRnBiaUJXYVdWM01ROHdEUVlEVlFRREV3WkhiMjluYkdVeEdEQVdCZ05WQkFzVEQwZHYKYjJkc1pTQkdiM0lnVjI5eWF6RUxNQWtHQTFVRUJoTUNWVk14RXpBUkJnTlZCQWdUQ2tOaGJHbG1iM0p1YVdFd0hoY05NakV4TWpJMApNRGt6TURNNFdoY05Nall4TWpJek1Ea3pNRE00V2pCN01SUXdFZ1lEVlFRS0V3dEhiMjluYkdVZ1NXNWpMakVXTUJRR0ExVUVCeE1OClRXOTFiblJoYVc0Z1ZtbGxkekVQTUEwR0ExVUVBeE1HUjI5dloyeGxNUmd3RmdZRFZRUUxFdzlIYjI5bmJHVWdSbTl5SUZkdmNtc3gKQ3pBSkJnTlZCQVlUQWxWVE1STXdFUVlEVlFRSUV3cERZV3hwWm05eWJtbGhNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QQpNSUlCQ2dLQ0FRRUF1T0RRUFhKL3BKczdQZWhyK3RZdWhVUHJGNHZvSStHSW1wT1M2bkEwWUM2SHBnUDFyeXo3WVBNaEF1dGQrT2NVCjRDUWF0ZDFrOERtV0dBS09DZEtNQ0JlaG8wY2lsREVZQTk0Mk1GNEhTVDNFNlBGcFhhb3htdjNTRkptMXFHNkhrdklRbzBzRTg1YkUKaGMwN3N2Y0FUQ0JQUTNRNjV1dndPeVhQRXlaaVVzL1BLeEx2NHYzMXBiYjNSdmJoRldMUkw1enVnUGh6M3JqTXpoY05GdUhnUTlKNgoySkw3dHY4azdDTDdLOXFZZWZKdXlYZy9wZXpYU2xXbFIvcFlCV3QxY1l0M2Y5azdHRDlydGkrb3pZN0RxVGZqeGU5QVpjaEM0cHVGClBlSkRSVzdIWDFrT2ROWW1BbWVhcitBQzJaOUhZdmNJTmVlelhCaEd2MXR6cHpqM09RSURBUUFCTUEwR0NTcUdTSWIzRFFFQkN3VUEKQTRJQkFRQ1g4WlY0eGQ4TGFOQUZ3YS9QN2tOeFVadG1jWE5KWFphdW5xM2hxL3o2YWhXc0VlWDlEbk4rUnZJSkRXZUw3azg5KzluZAo2cUh5UWxYS1hmbVNVeGVJeDdwckpuMDg5dzN2UnMzdE1nNThOK0NQOFFEYjFNdWFkdHlYSDduYlBhSVVvU3lhTlphQjUwV2R2R0pzCnd3THJOWHQxZGg0Um1zVG8rZVJHUysvY08zZkE0M1d2WHlHNURsTnJNd3FDbUE4d3VKVkNFaWpsbldDQm5zYm4ybDZUMjgwSVdoMTAKWXROVVpSSXhVdWRTRy9QY0FFYjJMRnY1c1B5eGRPdEFHbytFakp0bEplaFdWbHpTMnNHaVpscEhER1kxellBeDJRcVhVM1RPRTU2TQpHYjJWNEc0YUs4bDEvT0NxTVJOd3VlNzhudFlwS29PZXI2RnZaTE9WaVN4MzwvZHM6WDUwOUNlcnRpZmljYXRlPgogICAgICAgIDwvZHM6WDUwOURhdGE+CiAgICAgIDwvZHM6S2V5SW5mbz4KICAgIDwvbWQ6S2V5RGVzY3JpcHRvcj4KICAgIDxtZDpOYW1lSURGb3JtYXQ+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6MS4xOm5hbWVpZC1mb3JtYXQ6ZW1haWxBZGRyZXNzPC9tZDpOYW1lSURGb3JtYXQ+CiAgICA8bWQ6U2luZ2xlU2lnbk9uU2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL3NhbWwyL2lkcD9pZHBpZD1DMDNrendnajEiLz4KICAgIDxtZDpTaW5nbGVTaWduT25TZXJ2aWNlIEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLVBPU1QiIExvY2F0aW9uPSJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20vby9zYW1sMi9pZHA/aWRwaWQ9QzAza3p3Z2oxIi8+CiAgPC9tZDpJRFBTU09EZXNjcmlwdG9yPgo8L21kOkVudGl0eURlc2NyaXB0b3I+",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                   "for test new",
						"encodedsaml_metadata_document": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48bWQ6RW50aXR5RGVzY3JpcHRvciB4bWxuczptZD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOm1ldGFkYXRhIiBlbnRpdHlJRD0iaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tL28vc2FtbDI/aWRwaWQ9QzAza3p3Z2oxIiB2YWxpZFVudGlsPSIyMDI2LTEyLTI0VDA5OjMwOjM4LjAwMFoiPgogIDxtZDpJRFBTU09EZXNjcmlwdG9yIFdhbnRBdXRoblJlcXVlc3RzU2lnbmVkPSJmYWxzZSIgcHJvdG9jb2xTdXBwb3J0RW51bWVyYXRpb249InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpwcm90b2NvbCI+CiAgICA8bWQ6S2V5RGVzY3JpcHRvciB1c2U9InNpZ25pbmciPgogICAgICA8ZHM6S2V5SW5mbyB4bWxuczpkcz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC8wOS94bWxkc2lnIyI+CiAgICAgICAgPGRzOlg1MDlEYXRhPgogICAgICAgICAgPGRzOlg1MDlDZXJ0aWZpY2F0ZT5NSUlEZERDQ0FseWdBd0lCQWdJR0FYM3J4eXpoTUEwR0NTcUdTSWIzRFFFQkN3VUFNSHN4RkRBU0JnTlZCQW9UQzBkdmIyZHNaU0JKCmJtTXVNUll3RkFZRFZRUUhFdzFOYjNWdWRHRnBiaUJXYVdWM01ROHdEUVlEVlFRREV3WkhiMjluYkdVeEdEQVdCZ05WQkFzVEQwZHYKYjJkc1pTQkdiM0lnVjI5eWF6RUxNQWtHQTFVRUJoTUNWVk14RXpBUkJnTlZCQWdUQ2tOaGJHbG1iM0p1YVdFd0hoY05NakV4TWpJMApNRGt6TURNNFdoY05Nall4TWpJek1Ea3pNRE00V2pCN01SUXdFZ1lEVlFRS0V3dEhiMjluYkdVZ1NXNWpMakVXTUJRR0ExVUVCeE1OClRXOTFiblJoYVc0Z1ZtbGxkekVQTUEwR0ExVUVBeE1HUjI5dloyeGxNUmd3RmdZRFZRUUxFdzlIYjI5bmJHVWdSbTl5SUZkdmNtc3gKQ3pBSkJnTlZCQVlUQWxWVE1STXdFUVlEVlFRSUV3cERZV3hwWm05eWJtbGhNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QQpNSUlCQ2dLQ0FRRUF1T0RRUFhKL3BKczdQZWhyK3RZdWhVUHJGNHZvSStHSW1wT1M2bkEwWUM2SHBnUDFyeXo3WVBNaEF1dGQrT2NVCjRDUWF0ZDFrOERtV0dBS09DZEtNQ0JlaG8wY2lsREVZQTk0Mk1GNEhTVDNFNlBGcFhhb3htdjNTRkptMXFHNkhrdklRbzBzRTg1YkUKaGMwN3N2Y0FUQ0JQUTNRNjV1dndPeVhQRXlaaVVzL1BLeEx2NHYzMXBiYjNSdmJoRldMUkw1enVnUGh6M3JqTXpoY05GdUhnUTlKNgoySkw3dHY4azdDTDdLOXFZZWZKdXlYZy9wZXpYU2xXbFIvcFlCV3QxY1l0M2Y5azdHRDlydGkrb3pZN0RxVGZqeGU5QVpjaEM0cHVGClBlSkRSVzdIWDFrT2ROWW1BbWVhcitBQzJaOUhZdmNJTmVlelhCaEd2MXR6cHpqM09RSURBUUFCTUEwR0NTcUdTSWIzRFFFQkN3VUEKQTRJQkFRQ1g4WlY0eGQ4TGFOQUZ3YS9QN2tOeFVadG1jWE5KWFphdW5xM2hxL3o2YWhXc0VlWDlEbk4rUnZJSkRXZUw3azg5KzluZAo2cUh5UWxYS1hmbVNVeGVJeDdwckpuMDg5dzN2UnMzdE1nNThOK0NQOFFEYjFNdWFkdHlYSDduYlBhSVVvU3lhTlphQjUwV2R2R0pzCnd3THJOWHQxZGg0Um1zVG8rZVJHUysvY08zZkE0M1d2WHlHNURsTnJNd3FDbUE4d3VKVkNFaWpsbldDQm5zYm4ybDZUMjgwSVdoMTAKWXROVVpSSXhVdWRTRy9QY0FFYjJMRnY1c1B5eGRPdEFHbytFakp0bEplaFdWbHpTMnNHaVpscEhER1kxellBeDJRcVhVM1RPRTU2TQpHYjJWNEc0YUs4bDEvT0NxTVJOd3VlNzhudFlwS29PZXI2RnZaTE9WaVN4MzwvZHM6WDUwOUNlcnRpZmljYXRlPgogICAgICAgIDwvZHM6WDUwOURhdGE+CiAgICAgIDwvZHM6S2V5SW5mbz4KICAgIDwvbWQ6S2V5RGVzY3JpcHRvcj4KICAgIDxtZDpOYW1lSURGb3JtYXQ+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6MS4xOm5hbWVpZC1mb3JtYXQ6ZW1haWxBZGRyZXNzPC9tZDpOYW1lSURGb3JtYXQ+CiAgICA8bWQ6U2luZ2xlU2lnbk9uU2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL3NhbWwyL2lkcD9pZHBpZD1DMDNrendnajEiLz4KICAgIDxtZDpTaW5nbGVTaWduT25TZXJ2aWNlIEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLVBPU1QiIExvY2F0aW9uPSJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20vby9zYW1sMi9pZHA/aWRwaWQ9QzAza3p3Z2oxIi8+CiAgPC9tZDpJRFBTU09EZXNjcmlwdG9yPgo8L21kOkVudGl0eURlc2NyaXB0b3I+",
					}),
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

var AlicloudRamSamlProviderMap9203 = map[string]string{
	"update_date": CHECKSET,
	"arn":         CHECKSET,
}

func AlicloudRamSamlProviderBasicDependence9203(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Ram SamlProvider. <<< Resource test cases, automatically generated.
