package cms_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The context store that an API key hangs off is not a managed Terraform
// resource, so the test stands it up with a raw Cms ROA call once the managed
// workspace exists, and tears it down in t.Cleanup. The SLS project and the
// Cms workspace are managed resources, so terraform owns their lifecycle.

func TestAccAliCloudCmsContextStoreAPIKey_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	workspaceName := fmt.Sprintf("tfacccmscsk%d", rand)
	contextStoreName := fmt.Sprintf("tfacccmscsk-cs%d", rand)
	keyName := fmt.Sprintf("tfacccmscsk-key%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// Bring up the managed parent chain (SLS project + Cms workspace),
				// then stand up the unmanaged context store the API key targets.
				Config: testAccCmsContextStoreAPIKeyParentConfig(workspaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_workspace.default", "workspace_name", workspaceName),
					func(_ *terraform.State) error {
						return createContextStoreForTest(t, workspaceName, contextStoreName)
					},
				),
			},
			{
				Config: testAccCmsContextStoreAPIKeyConfig(workspaceName, contextStoreName, keyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_context_store_api_key.default", "id"),
					resource.TestCheckResourceAttr("alicloud_cms_context_store_api_key.default", "workspace", workspaceName),
					resource.TestCheckResourceAttr("alicloud_cms_context_store_api_key.default", "context_store_name", contextStoreName),
					resource.TestCheckResourceAttr("alicloud_cms_context_store_api_key.default", "name", keyName),
					resource.TestCheckResourceAttrSet("alicloud_cms_context_store_api_key.default", "api_key"),
					resource.TestCheckResourceAttrSet("alicloud_cms_context_store_api_key.default", "create_time"),
					resource.TestCheckResourceAttrSet("alicloud_cms_context_store_api_key.default", "region_id"),
				),
			},
			{
				ResourceName:            "alicloud_cms_context_store_api_key.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
		},
	})
}

func TestAccAliCloudCmsContextStoreAPIKey_basic0_twin(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	workspaceName := fmt.Sprintf("tfacccmscsktwin%d", rand)
	contextStoreName := fmt.Sprintf("tfacccmscsktwin-cs%d", rand)
	keyName := fmt.Sprintf("tfacccmscsktwin-key%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCmsContextStoreAPIKeyParentConfig(workspaceName),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						return createContextStoreForTest(t, workspaceName, contextStoreName)
					},
				),
			},
			{
				Config: testAccCmsContextStoreAPIKeyConfig(workspaceName, contextStoreName, keyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_context_store_api_key.default", "id"),
					resource.TestCheckResourceAttr("alicloud_cms_context_store_api_key.default", "name", keyName),
				),
			},
			{
				ResourceName:            "alicloud_cms_context_store_api_key.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
		},
	})
}

func testAccCmsContextStoreAPIKeyParentConfig(workspaceName string) string {
	return fmt.Sprintf(`
resource "alicloud_log_project" "default" {
  project_name = "%[1]s"
}

resource "alicloud_cms_workspace" "default" {
  sls_project    = alicloud_log_project.default.project_name
  workspace_name = "%[1]s"
}
`, workspaceName)
}

func testAccCmsContextStoreAPIKeyConfig(workspaceName, contextStoreName, keyName string) string {
	return fmt.Sprintf(`
resource "alicloud_log_project" "default" {
  project_name = "%[1]s"
}

resource "alicloud_cms_workspace" "default" {
  sls_project    = alicloud_log_project.default.project_name
  workspace_name = "%[1]s"
}

resource "alicloud_cms_context_store_api_key" "default" {
  workspace          = "%[1]s"
  context_store_name = "%[2]s"
  name               = "%[3]s"
}
`, workspaceName, contextStoreName, keyName)
}

// cmsAccClient builds a live AliyunClient from the same env vars the rest of
// the suite uses. It mirrors alicloud.sharedClientForRegion, which lives in
// package alicloud and so cannot be imported from an external service test.
func cmsAccClient(t *testing.T) *connectivity.AliyunClient {
	t.Helper()
	accessKey := os.Getenv("ALICLOUD_ACCESS_KEY")
	secretKey := os.Getenv("ALICLOUD_SECRET_KEY")
	if accessKey == "" || secretKey == "" {
		t.Fatal("ALICLOUD_ACCESS_KEY and ALICLOUD_SECRET_KEY must be set for acceptance tests")
	}
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	securityToken := os.Getenv("ALICLOUD_SECURITY_TOKEN")

	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if v := securityToken; v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	cred, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		t.Fatalf("building credential: %s", err)
	}

	conf := &connectivity.Config{
		Region:     connectivity.Region(region),
		RegionId:   region,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		Protocol:   "HTTPS",
		Credential: cred,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	client, err := conf.Client()
	if err != nil {
		t.Fatalf("building AliyunClient: %s", err)
	}
	return client
}

// createContextStoreForTest stands up the unmanaged context store the API key
// targets, and registers a cleanup that deletes it (and tolerates it already
// being gone). The Cms backend has intermittently returned 503 here; the test
// surfaces that as a real failure rather than masking it.
func createContextStoreForTest(t *testing.T, workspaceName, contextStoreName string) error {
	t.Helper()
	client := cmsAccClient(t)

	action := fmt.Sprintf("/workspace/%s/contextstore", workspaceName)
	body := map[string]interface{}{
		"contextStoreName": contextStoreName,
		"contextType":      "memory",
		"description":      "tf-acc context store for context_store_api_key test",
	}
	if _, err := client.RoaPost("Cms", "2024-03-30", action, map[string]*string{}, nil, body, true); err != nil {
		return fmt.Errorf("creating test context store %s: %w", contextStoreName, err)
	}

	t.Cleanup(func() {
		deleteAction := fmt.Sprintf("/workspace/%s/contextstore/%s", workspaceName, contextStoreName)
		_, _ = client.RoaDelete("Cms", "2024-03-30", deleteAction, map[string]*string{}, nil, nil, true)
	})
	return nil
}
