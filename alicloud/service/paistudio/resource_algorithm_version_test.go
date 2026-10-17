package paistudio_test

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/paistudio"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// AlgorithmVersion depends on a parent Algorithm that has no Terraform
// resource of its own, so the parent is provisioned out-of-band via a
// standalone client built from the same ALICLOUD_* env (the same pattern
// alicloud_sweeper_test.go's sharedClientForRegion uses).
//
// The parent ID is only known after the standalone RoaPost returns, so it
// cannot be created inside resource.TestCase.PreCheck: the Steps' Config
// strings are evaluated when the resource.TestCase literal is built, i.e.
// before PreCheck runs, and a PreCheck-assigned variable would be frozen as
// the empty string into both Configs. The parent is therefore created in the
// test body (after the ACC gate) and the real, non-empty algorithmId is
// substituted into every Step Config at literal-construction time.
//
// A plain `go test` (no TF_ACC) skips before any client is built or cloud call
// is made: neither resource.Test's own TF_ACC gate nor the explicit guard here
// touches PaiStudio. No placeholder ID, static mirror or account/region/cred
// override is used to paper over the timing.
func TestAccAliCloudPaiAlgorithmVersion(t *testing.T) {
	// Mirror terraform-plugin-testing: skip acceptance unless TF_ACC is set, so
	// an unconfigured `go test` run never builds the standalone client, never
	// calls PaiStudio, and only compiles the test.
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set; skipping PaiStudio acceptance test")
	}
	acctest.PreCheck(t)

	// PaiStudio (POP product "PaiStudio", 2022-01-12) is served from the unified
	// PAI endpoint pai.<region>.aliyuncs.com; it is not registered with the
	// location service and paistudio.<region>.aliyuncs.com does not resolve.
	// (client).loadEndpoint reads ALIBABA_CLOUD_ENDPOINT_<PRODUCT> first, so set
	// it here (before the first RoaPost) for both this standalone client and the
	// muxed provider's client.
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = os.Getenv("ALIBABA_CLOUD_REGION")
	}
	if region != "" {
		os.Setenv("ALIBABA_CLOUD_ENDPOINT_PAISTUDIO", fmt.Sprintf("pai.%s.aliyuncs.com", region))
	}
	client, err := newTestPaiStudioClient()
	if err != nil {
		t.Fatalf("building PaiStudio client: %s", err)
	}
	algorithmId, err := createTestPaiStudioAlgorithm(client)
	if err != nil {
		t.Fatalf("creating parent Algorithm: %s", err)
	}
	// t.Cleanup covers all exit paths, so a failed step still removes the parent
	// Algorithm. A cleanup failure is surfaced as a test failure (not logged and
	// swallowed): only an already-gone parent (typed NoSuchObject) is tolerated.
	t.Cleanup(func() {
		if err := deleteTestPaiStudioAlgorithm(client, algorithmId); err != nil {
			if isPaiStudioNotFound(err) {
				return
			}
			t.Errorf("cleanup: deleting parent Algorithm %s: %s", algorithmId, err)
		}
	})

	const algorithmVersion = "v1.0.0"
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// CheckDestroy proves the managed alicloud_pai_algorithm_version itself is
		// gone after a step teardown, not just that the parent cleanup ran. A
		// transient 5xx on the read-back is NOT treated as destroyed.
		CheckDestroy: testAccCheckPaiAlgorithmVersionDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiAlgorithmVersionConfig(algorithmId, algorithmVersion, "1.0"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_algorithm_version.default", "algorithm_id", algorithmId),
					resource.TestCheckResourceAttr("alicloud_pai_algorithm_version.default", "algorithm_version", algorithmVersion),
					resource.TestCheckResourceAttrSet("alicloud_pai_algorithm_version.default", "algorithm_spec"),
					resource.TestCheckResourceAttrSet("alicloud_pai_algorithm_version.default", "id"),
				),
			},
			{
				Config: testAccPaiAlgorithmVersionConfig(algorithmId, algorithmVersion, "2.0"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_algorithm_version.default", "algorithm_id", algorithmId),
					resource.TestCheckResourceAttr("alicloud_pai_algorithm_version.default", "algorithm_version", algorithmVersion),
					resource.TestCheckResourceAttrSet("alicloud_pai_algorithm_version.default", "algorithm_spec"),
				),
			},
			{
				// Config is provided so the resource exists in state to import; the
				// real algorithmId (not a frozen empty string) drives the import.
				Config:                  testAccPaiAlgorithmVersionConfig(algorithmId, algorithmVersion, "2.0"),
				ResourceName:            "alicloud_pai_algorithm_version.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"algorithm_spec"},
			},
		},
	})
}

// testAccCheckPaiAlgorithmVersionDestroy verifies every
// alicloud_pai_algorithm_version in state is gone after destroy: a RoaGet that
// returns the typed NoSuchObject proves the version is absent, an unexpected
// error (5xx/auth/malformed) is reported rather than mistaken for gone, and a
// 200 response means the resource leaked and fails the test.
func testAccCheckPaiAlgorithmVersionDestroy(client *connectivity.AliyunClient) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_pai_algorithm_version" {
				continue
			}
			id := rs.Primary.ID
			parts := strings.SplitN(id, ":", 2)
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
				return fmt.Errorf("destroy check: bad composite id %q", id)
			}
			_, err := client.RoaGet(paistudio.PaiStudioProduct, paistudio.PaiStudioVersion,
				fmt.Sprintf("/api/v1/algorithms/%s/versions/%s", parts[0], parts[1]), nil, nil, nil)
			if err == nil {
				return fmt.Errorf("algorithm version %s still exists after destroy", id)
			}
			if !isPaiStudioNotFound(err) {
				return fmt.Errorf("destroy check for algorithm version %s: %w", id, err)
			}
		}
		return nil
	}
}

func testAccPaiAlgorithmVersionConfig(algorithmId, algorithmVersion, specVersion string) string {
	return fmt.Sprintf(`
resource "alicloud_pai_algorithm_version" "default" {
  algorithm_id      = %[1]q
  algorithm_version = %[2]q
  algorithm_spec = jsonencode({
    JobType = "TfJob"
    Image   = "registry.cn-hangzhou.aliyuncs.com/pai-test/echo:%[3]s"
    Command = "python main.py"
  })
}
`, algorithmId, algorithmVersion, specVersion)
}

// newTestPaiStudioClient builds a standalone AliyunClient from ALICLOUD_*
// env vars, mirroring alicloud_sweeper_test.go's sharedClientForRegion. The
// muxed provider meta is nil before the first step, so acctest.Provider.Meta()
// cannot be used in PreCheck.
func newTestPaiStudioClient() (*connectivity.AliyunClient, error) {
	protocol := os.Getenv("ALICLOUD_PROTOCOL")
	if protocol == "" {
		protocol = "HTTPS"
	}
	conf := connectivity.Config{
		Region:        connectivity.Region(os.Getenv("ALICLOUD_REGION")),
		RegionId:      os.Getenv("ALICLOUD_REGION"),
		AccessKey:     os.Getenv("ALICLOUD_ACCESS_KEY"),
		SecretKey:     os.Getenv("ALICLOUD_SECRET_KEY"),
		SecurityToken: os.Getenv("ALICLOUD_SECURITY_TOKEN"),
		Protocol:      protocol,
	}
	if conf.AccessKey == "" {
		conf.AccessKey = os.Getenv("ALIBABA_CLOUD_ACCESS_KEY")
	}
	if conf.SecretKey == "" {
		conf.SecretKey = os.Getenv("ALIBABA_CLOUD_SECRET_KEY")
	}
	if conf.RegionId == "" {
		conf.RegionId = os.Getenv("ALIBABA_CLOUD_REGION")
		conf.Region = connectivity.Region(conf.RegionId)
	}
	if conf.SecurityToken == "" {
		conf.SecurityToken = os.Getenv("ALIBABA_CLOUD_SECURITY_TOKEN")
	}

	accountId := os.Getenv("ALICLOUD_ACCOUNT_ID")
	if accountId == "" {
		accountId = os.Getenv("ALIBABA_CLOUD_ACCOUNT_ID")
	}
	if accountId != "" {
		conf.AccountId = accountId
	}
	accountType := os.Getenv("ALICLOUD_ACCOUNT_TYPE")
	if accountType == "" {
		accountType = os.Getenv("ALIBABA_CLOUD_ACCOUNT_TYPE")
	}
	if accountType != "" {
		conf.AccountType = accountType
	}

	// Config.Credential is dereferenced by getAuthCredential during the first
	// RoaPost (describeEndpointForService). conf.Client() does not populate it,
	// so mirror alicloud_sweeper_test.go's sharedClientForRegion and set it
	// explicitly from the same ALICLOUD_* env to avoid a nil-pointer panic.
	credCfg := new(credentials.Config).SetType("access_key").
		SetAccessKeyId(conf.AccessKey).SetAccessKeySecret(conf.SecretKey)
	if strings.TrimSpace(conf.SecurityToken) != "" {
		credCfg.SetType("sts").SetSecurityToken(conf.SecurityToken)
	}
	cred, err := credentials.NewCredential(credCfg)
	if err != nil {
		return nil, fmt.Errorf("building PaiStudio credential: %w", err)
	}
	conf.Credential = cred

	conf.Endpoints = &sync.Map{}
	conf.SignVersion = &sync.Map{}
	return conf.Client()
}

func createTestPaiStudioAlgorithm(client *connectivity.AliyunClient) (string, error) {
	rand := tfacctest.RandIntRange(1000, 9999)
	body := map[string]interface{}{
		"AlgorithmName":        fmt.Sprintf("tf-acc-algo-%d", rand),
		"DisplayName":          fmt.Sprintf("tf-acc-algo-%d", rand),
		"AlgorithmDescription": "tf acc parent algorithm",
	}
	resp, err := client.RoaPost(paistudio.PaiStudioProduct, paistudio.PaiStudioVersion,
		"/api/v1/algorithms", nil, nil, body, true)
	if err != nil {
		return "", fmt.Errorf("RoaPost /api/v1/algorithms: %w", err)
	}
	id, _ := resp["AlgorithmId"].(string)
	if id == "" {
		return "", fmt.Errorf("CreateAlgorithm returned empty AlgorithmId: %v", resp)
	}
	return id, nil
}

func deleteTestPaiStudioAlgorithm(client *connectivity.AliyunClient, algorithmId string) error {
	_, err := client.RoaDelete(paistudio.PaiStudioProduct, paistudio.PaiStudioVersion,
		fmt.Sprintf("/api/v1/algorithms/%s", algorithmId), nil, nil, nil, true)
	return err
}

// isPaiStudioNotFound reports whether err is the PaiStudio not-found error,
// using the product error contract (typed *tea.SDKError Code == NoSuchObject)
// only. A loose substring match is intentionally avoided so a non-NotFound
// message that merely contains the word "NoSuchObject" cannot be misread as
// gone and mask a leaked resource during cleanup or destroy checks.
func isPaiStudioNotFound(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*tea.SDKError); ok && e.Code != nil && *e.Code == "NoSuchObject" {
		return true
	}
	return false
}
