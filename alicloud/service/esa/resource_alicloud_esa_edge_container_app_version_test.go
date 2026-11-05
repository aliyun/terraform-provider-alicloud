package esa_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const esaEdgeContainerAppVersionResourceId = "alicloud_esa_edge_container_app_version.default"

// TestAccAliCloudEsaEdgeContainerAppVersion_basic covers the full version
// lifecycle: create a parent edge container app, create a version with a
// container group, publish the version, and re-publish after a publish
// parameter change.
func TestAccAliCloudEsaEdgeContainerAppVersion_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-esa-ecav-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName:            esaEdgeContainerAppVersionResourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEsaEdgeContainerAppVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEsaEdgeContainerAppVersionBasic(name, 80),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(esaEdgeContainerAppVersionResourceId, "id"),
					resource.TestCheckResourceAttrSet(esaEdgeContainerAppVersionResourceId, "version_id"),
					resource.TestCheckResourceAttr(esaEdgeContainerAppVersionResourceId, "name", name),
					resource.TestCheckResourceAttr(esaEdgeContainerAppVersionResourceId, "publish_env", "staging"),
					resource.TestCheckResourceAttr(esaEdgeContainerAppVersionResourceId, "percentage", "80"),
					resource.TestCheckResourceAttr(esaEdgeContainerAppVersionResourceId, "containers.#", "1"),
					resource.TestCheckResourceAttr(esaEdgeContainerAppVersionResourceId, "containers.0.spec", "1C2G"),
					resource.TestCheckResourceAttr(esaEdgeContainerAppVersionResourceId, "containers.0.storage", "0.5G"),
					resource.TestCheckResourceAttr(esaEdgeContainerAppVersionResourceId, "containers.0.probe_type", "tcpSocket"),
					resource.TestCheckResourceAttrSet(esaEdgeContainerAppVersionResourceId, "containers.0.probe_content.#"),
				),
			},
			{
				// Changing percentage re-publishes the version without replacing it.
				Config: testAccEsaEdgeContainerAppVersionBasic(name, 100),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(esaEdgeContainerAppVersionResourceId, "percentage", "100"),
					resource.TestCheckResourceAttrSet(esaEdgeContainerAppVersionResourceId, "version_id"),
				),
			},
			{
				ResourceName:            esaEdgeContainerAppVersionResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"publish_env", "publish_type", "percentage", "regions", "full_release", "start_time"},
			},
		},
	})
}

func testAccEsaEdgeContainerAppVersionBasic(name string, percentage int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

resource "alicloud_esa_edge_container_app" "default" {
  edge_container_app_name = var.name
  service_port            = 8080
  target_port             = 80
}

resource "alicloud_esa_edge_container_app_version" "default" {
  app_id     = alicloud_esa_edge_container_app.default.id
  name       = var.name
  remarks    = "tf test version"
  publish_env  = "staging"
  percentage   = %[2]d
  containers {
    name          = "tf-test-container"
    image         = "nginx:1.25"
    is_acr_image  = false
    spec          = "1C2G"
    storage       = "0.5G"
    command       = "nginx -g 'daemon off;'"
    probe_type    = "tcpSocket"
    probe_content {
      initial_delay_seconds = 5
      timeout_seconds       = 3
      period_seconds        = 10
      failure_threshold     = 3
      success_threshold     = 1
      port                  = 80
    }
  }
}
`, name, percentage)
}

func testAccCheckEsaEdgeContainerAppVersionDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	return checkEsaEdgeContainerAppVersionDestroyed(s, func(id string) error {
		_, err := esa.NewService(client).DescribeEsaEdgeContainerAppVersion(id)
		return err
	})
}

func checkEsaEdgeContainerAppVersionDestroyed(s *terraform.State, describe func(string) error) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_esa_edge_container_app_version" {
			continue
		}
		err := describe(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("edge container app version %s still exists", rs.Primary.ID)
		}
		if errs.IsProviderNotFound(err) {
			continue
		}
		return fmt.Errorf("describing edge container app version %s failed: %s", rs.Primary.ID, err)
	}
	return nil
}

func init() {
	resource.AddTestSweepers("alicloud_esa_edge_container_app_version", &resource.Sweeper{
		Name: "alicloud_esa_edge_container_app_version",
		F:    testSweepEsaEdgeContainerAppVersion,
	})
}

func testSweepEsaEdgeContainerAppVersion(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{"tf-testacc-", "tf_testacc-", "terraform-example"}

	// Versions live under apps, so list apps first, then their versions.
	appAction := "ListEdgeContainerApps"
	pageNumber := 1
	for {
		query := map[string]interface{}{
			"RegionId":   client.RegionId,
			"PageNumber": pageNumber,
			"PageSize":   helper.PageSizeXLarge,
		}
		var resp map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err := retry.RetryContext(context.Background(), 2*time.Minute, func() *retry.RetryError {
			var callErr error
			resp, callErr = client.RpcGet("ESA", "2024-09-10", appAction, query, map[string]interface{}{})
			if callErr != nil {
				if helper.NeedRetry(callErr) {
					wait()
					return retry.RetryableError(callErr)
				}
				return retry.NonRetryableError(callErr)
			}
			return nil
		})
		logs.AddDebug(appAction, resp, query)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_esa_edge_container_app_version", appAction, errs.AlibabaCloudSdkGoERROR)
		}
		apps, ok := resp["Apps"].([]interface{})
		if !ok {
			break
		}
		for _, raw := range apps {
			app, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			appName := fmt.Sprint(app["Name"])
			skip := true
			for _, p := range prefixes {
				if len(appName) >= len(p) && appName[:len(p)] == p {
					skip = false
					break
				}
			}
			if skip {
				continue
			}
			appId := fmt.Sprint(app["AppId"])
			if appId == "" {
				continue
			}
			if sweepErr := sweepEsaEdgeContainerAppVersionsForApp(client, appId); sweepErr != nil {
				return sweepErr
			}
		}
		if len(apps) < helper.PageSizeXLarge {
			break
		}
		pageNumber++
	}
	return nil
}

func sweepEsaEdgeContainerAppVersionsForApp(client *connectivity.AliyunClient, appId string) error {
	service := esa.NewService(client)
	listAction := "ListEdgeContainerAppVersions"
	pageNumber := 1
	for {
		query := map[string]interface{}{
			"AppId":      appId,
			"RegionId":   client.RegionId,
			"PageNumber": pageNumber,
			"PageSize":   helper.PageSizeXLarge,
		}
		var resp map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err := retry.RetryContext(context.Background(), 2*time.Minute, func() *retry.RetryError {
			var callErr error
			resp, callErr = client.RpcGet("ESA", "2024-09-10", listAction, query, map[string]interface{}{})
			if callErr != nil {
				if helper.NeedRetry(callErr) {
					wait()
					return retry.RetryableError(callErr)
				}
				return retry.NonRetryableError(callErr)
			}
			return nil
		})
		logs.AddDebug(listAction, resp, query)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, appId, listAction, errs.AlibabaCloudSdkGoERROR)
		}
		versions, ok := resp["Versions"].([]interface{})
		if !ok {
			break
		}
		for _, raw := range versions {
			v, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			versionId := fmt.Sprint(v["VersionId"])
			if versionId == "" {
				continue
			}
			id := appId + ":" + versionId
			if delErr := service.DeleteEsaEdgeContainerAppVersion(id); delErr != nil {
				return delErr
			}
		}
		if len(versions) < helper.PageSizeXLarge {
			break
		}
		pageNumber++
	}
	return nil
}

// ---- unit tests ----

func TestUnitEsaEdgeContainerAppVersionIdParse(t *testing.T) {
	cases := map[string]struct {
		id      string
		wantApp string
		wantVer string
		wantErr bool
	}{
		"valid":      {id: "app-123:ver-456", wantApp: "app-123", wantVer: "ver-456"},
		"singlePart": {id: "app-123", wantErr: true},
		"empty":      {id: "", wantErr: true},
		"threeParts": {id: "a:b:c", wantErr: true},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			app, ver, err := esa.EdgeContainerAppVersionId(c.id)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected an error for %q, got %s:%s", c.id, app, ver)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if app != c.wantApp || ver != c.wantVer {
				t.Fatalf("EdgeContainerAppVersionId(%q) = (%s,%s), want (%s,%s)", c.id, app, ver, c.wantApp, c.wantVer)
			}
		})
	}
}

// Internal unit tests covering the unexported response-mapping and request
// serialization helpers live in resource_alicloud_esa_edge_container_app_version_unit_test.go
// (package esa) so they can reach the lowercase conversion functions.
