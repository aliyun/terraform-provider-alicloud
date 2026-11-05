// Package ims_test.
package ims_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var imsTestEndpoints sync.Map
var imsTestSignVersion sync.Map

func init() {
	resource.AddTestSweepers("alicloud_ims_app_secret", &resource.Sweeper{
		Name: "alicloud_ims_app_secret",
		F:    testSweepImsAppSecret,
	})
}

// testSweepImsAppSecret cleans up leftover IMS applications (and the secrets
// attached to them) created by interrupted acceptance runs. Deleting the
// application cascades to its secrets, so sweeping at the application level is
// enough and avoids per-application ListAppSecretIds churn.
func testSweepImsAppSecret(region string) error {
	client, err := imsTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	prefixes := []string{"tf-testacc", "tf_testacc"}
	response, err := client.RpcPost("Ims", "2019-08-15", "ListApplications", nil, map[string]interface{}{}, true)
	if err != nil {
		log.Printf("[ERROR] Ims ListApplications error: %#v", err)
		return nil
	}
	apps, _ := response["Applications"].(map[string]interface{})
	list, _ := apps["Application"].([]interface{})
	for _, v := range list {
		app, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := app["AppName"].(string)
		skip := true
		for _, p := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(p)) {
				skip = false
				break
			}
		}
		if skip {
			continue
		}
		appId, _ := app["AppId"].(string)
		if _, err := client.RpcPost("Ims", "2019-08-15", "DeleteApplication", nil, map[string]interface{}{
			"AppId": appId,
		}, true); err != nil {
			log.Printf("[ERROR] Failed to delete IMS application %s (%s): %s", name, appId, err)
		} else {
			log.Printf("[INFO] Deleted IMS application %s (%s)", name, appId)
		}
	}
	return nil
}

// imsTestClientForRegion rebuilds the sweeper/setup client on the exported
// connectivity API; the root-suite sharedClientForRegion is not reachable from
// a service test package.
func imsTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	accessKey := os.Getenv("ALICLOUD_ACCESS_KEY")
	if accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}
	secretKey := os.Getenv("ALICLOUD_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}
	securityToken := os.Getenv("ALICLOUD_SECURITY_TOKEN")

	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &imsTestEndpoints,
		SignVersion: &imsTestSignVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	credentialConfig := new(credentials.Config).SetType("access_key").
		SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if v := strings.TrimSpace(securityToken); v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential
	return conf.Client()
}

func imsTestClient() (*connectivity.AliyunClient, error) {
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	return imsTestClientForRegion(region)
}

// createImsTestApplication provisions an IMS application to host the app
// secret under test; the caller registers its deletion. SecretRequired is
// forced true so the application can hold secrets.
func createImsTestApplication(t *testing.T, client *connectivity.AliyunClient, displayName string) string {
	t.Helper()
	response, err := client.RpcPost("Ims", "2019-08-15", "CreateApplication", nil, map[string]interface{}{
		"AppName":        displayName,
		"DisplayName":    displayName,
		"AppType":        "WebApp",
		"SecretRequired": true,
	}, true)
	if err != nil {
		t.Fatalf("creating IMS test application %q: %s", displayName, err)
	}
	app, _ := response["Application"].(map[string]interface{})
	appId, _ := app["AppId"].(string)
	if appId == "" {
		t.Fatalf("creating IMS test application %q: response carried no Application.AppId: %v", displayName, response)
	}
	return appId
}

func TestAccAliCloudImsAppSecret_basic(t *testing.T) {
	resourceId := "alicloud_ims_app_secret.default"
	r := tfacctest.RandIntRange(10000, 99999)

	client, err := imsTestClient()
	if err != nil {
		t.Fatalf("building IMS test client: %s", err)
	}
	displayName := fmt.Sprintf("tf-testacc%d", r)
	appId := createImsTestApplication(t, client, displayName)
	t.Cleanup(func() {
		if _, err := client.RpcPost("Ims", "2019-08-15", "DeleteApplication", nil, map[string]interface{}{
			"AppId": appId,
		}, true); err != nil {
			log.Printf("[ERROR] Failed to delete IMS test application %s (%s): %s", displayName, appId, err)
		}
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckImsAppSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudImsAppSecretBasic(appId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "app_id", appId),
					resource.TestCheckResourceAttrSet(resourceId, "app_secret_id"),
					resource.TestCheckResourceAttrSet(resourceId, "app_secret_value"),
					resource.TestCheckResourceAttrSet(resourceId, "create_date"),
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

func testAccAliCloudImsAppSecretBasic(appId string) string {
	return fmt.Sprintf(`
resource "alicloud_ims_app_secret" "default" {
  app_id = %q
}
`, appId)
}

func testAccCheckImsAppSecretDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ims_app_secret" {
			continue
		}
		appId, appSecretId, err := parseImsAppSecretIdForTest(rs.Primary.ID)
		if err != nil {
			return err
		}
		_, err = client.RpcPost("Ims", "2019-08-15", "GetAppSecret", nil, map[string]interface{}{
			"AppId":       appId,
			"AppSecretId": appSecretId,
		}, true)
		if err != nil {
			if isImsAppSecretNotExistForTest(err) {
				continue
			}
			return fmt.Errorf("reading app secret %s: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("IMS app secret %s still exists", rs.Primary.ID)
	}
	return nil
}

func parseImsAppSecretIdForTest(id string) (appId, appSecretId string, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid resource id %q: expected <app_id>:<app_secret_id>", id)
	}
	return parts[0], parts[1], nil
}

func isImsAppSecretNotExistForTest(err error) bool {
	if err == nil {
		return false
	}
	if msg := err.Error(); strings.Contains(msg, "EntityNotExist.AppSecret") ||
		strings.Contains(msg, "EntityNotExist.Application") {
		return true
	}
	return errs.IsExpectedErrors(err, []string{"EntityNotExist.AppSecret", "EntityNotExist.Application"}) ||
		errs.NotFoundError(err)
}
