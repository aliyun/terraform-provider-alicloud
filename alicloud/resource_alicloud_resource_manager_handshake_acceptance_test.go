package alicloud

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var ResourceManagerHandshakeAcceptanceMap = map[string]string{
	"status": CHECKSET,
}

var resourceManagerHandshakeInvitedAccountId string

const resourceManagerHandshakeConsistencyDelay = 90 * time.Second

func TestAccAliCloudResourceManagerHandshakeAcceptance_schema(t *testing.T) {
	resourceSchema := resourceAliCloudResourceManagerHandshakeAcceptance().Schema
	expected := map[string]struct {
		required bool
		computed bool
		forceNew bool
	}{
		"create_time":               {computed: true},
		"expire_time":               {computed: true},
		"handshake_id":              {required: true, forceNew: true},
		"invited_account_real_name": {computed: true},
		"master_account_id":         {computed: true},
		"master_account_name":       {computed: true},
		"master_account_real_name":  {computed: true},
		"modify_time":               {computed: true},
		"note":                      {computed: true},
		"resource_directory_id":     {computed: true},
		"status":                    {computed: true},
		"target_entity":             {computed: true},
		"target_type":               {computed: true},
	}

	if len(resourceSchema) != len(expected) {
		t.Fatalf("schema field count = %d, want %d", len(resourceSchema), len(expected))
	}

	for name, want := range expected {
		got, ok := resourceSchema[name]
		if !ok {
			t.Fatalf("schema missing field %q", name)
		}
		if got.Type != schema.TypeString {
			t.Fatalf("schema field %q type = %v, want TypeString", name, got.Type)
		}
		if got.Required != want.required || got.Computed != want.computed || got.ForceNew != want.forceNew {
			t.Fatalf("schema field %q flags = required:%t computed:%t force_new:%t, want required:%t computed:%t force_new:%t",
				name, got.Required, got.Computed, got.ForceNew, want.required, want.computed, want.forceNew)
		}
	}
}

// Accepting a handshake must be performed by the invited account, so this test needs two accounts.
// The management (master) account is supplied via ALICLOUD_ACCESS_KEY_1/SECRET_KEY_1 and sends the
// invitation; the invited account is supplied via ALICLOUD_ACCESS_KEY_2/SECRET_KEY_2 and accepts it.
// The invited account id is read from ALICLOUD_ACCOUNT_ID_2 / INVITED_ALICLOUD_ACCOUNT_ID, or derived
// from ALICLOUD_ACCESS_KEY_2 / ALICLOUD_SECRET_KEY_2 via STS GetCallerIdentity. Both providers are configured explicitly so generic
// credential environment variables cannot override the account used by the provider alias.
func ResourceManagerHandshakeAcceptanceBasicdependence(name string, invitedAccountId string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

# Management account (default provider) sends the invitation, explicit credentials to avoid env fallback.
provider "alicloud" {
  access_key = "%s"
  secret_key = "%s"
}

# Invited account (accepts the invitation), explicit credentials so it overrides the default keys.
provider "alicloud" {
  alias      = "invited"
  access_key = "%s"
  secret_key = "%s"
}

# Management account (default provider) sends the invitation.
resource "alicloud_resource_manager_handshake" "default" {
  target_entity = "%s"
  target_type   = "Account"
  note          = "tf-test-acceptance"
}
`, name, os.Getenv("ALICLOUD_ACCESS_KEY_1"), os.Getenv("ALICLOUD_SECRET_KEY_1"), os.Getenv("ALICLOUD_ACCESS_KEY_2"), os.Getenv("ALICLOUD_SECRET_KEY_2"), invitedAccountId)
}

func ResourceManagerHandshakeAcceptanceBasic(name string, invitedAccountId string) string {
	return ResourceManagerHandshakeAcceptanceBasicdependence(name, invitedAccountId) + `
resource "alicloud_resource_manager_handshake_acceptance" "default" {
  provider     = alicloud.invited
  handshake_id = "${alicloud_resource_manager_handshake.default.id}"
}
`
}

func testAccPreCheckHandshakeAcceptanceProfiles(t *testing.T) {
	testAccUseResourceManagerHandshakeManagementAccount(t)
	if os.Getenv("ALICLOUD_ACCESS_KEY_2") == "" || os.Getenv("ALICLOUD_SECRET_KEY_2") == "" {
		t.Skipf("Skipping: set ALICLOUD_ACCESS_KEY_1/SECRET_KEY_1 and ALICLOUD_ACCESS_KEY_2/SECRET_KEY_2 for the management and invited accounts")
	}
	testAccResolveResourceManagerHandshakeInvitedAccountId(t)
}

func testAccResourceManagerHandshakeManagementCredentialValues() (string, string) {
	ak := strings.TrimSpace(os.Getenv("ALICLOUD_ACCESS_KEY_1"))
	sk := strings.TrimSpace(os.Getenv("ALICLOUD_SECRET_KEY_1"))
	if ak != "" && sk != "" {
		return ak, sk
	}
	return strings.TrimSpace(os.Getenv("ALICLOUD_ACCESS_KEY")), strings.TrimSpace(os.Getenv("ALICLOUD_SECRET_KEY"))
}

func testAccResourceManagerHandshakeManagementCredentials(t *testing.T) (string, string) {
	ak, sk := testAccResourceManagerHandshakeManagementCredentialValues()
	if ak == "" || sk == "" {
		t.Skipf("Skipping: set ALICLOUD_ACCESS_KEY_1/SECRET_KEY_1 or ALICLOUD_ACCESS_KEY/SECRET_KEY for the management account")
	}
	return ak, sk
}

func testAccUseResourceManagerHandshakeManagementAccount(t *testing.T) {
	ak, sk := testAccResourceManagerHandshakeManagementCredentials(t)
	os.Setenv("ALICLOUD_ACCESS_KEY", ak)
	os.Setenv("ALICLOUD_SECRET_KEY", sk)
	os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
}

func testAccResolveResourceManagerHandshakeInvitedAccountId(t *testing.T) string {
	if resourceManagerHandshakeInvitedAccountId != "" {
		return resourceManagerHandshakeInvitedAccountId
	}
	for _, envName := range []string{"ALICLOUD_ACCOUNT_ID_2", "INVITED_ALICLOUD_ACCOUNT_ID"} {
		if v := strings.TrimSpace(os.Getenv(envName)); v != "" {
			resourceManagerHandshakeInvitedAccountId = v
			return resourceManagerHandshakeInvitedAccountId
		}
	}

	ak := strings.TrimSpace(os.Getenv("ALICLOUD_ACCESS_KEY_2"))
	sk := strings.TrimSpace(os.Getenv("ALICLOUD_SECRET_KEY_2"))
	if ak == "" || sk == "" {
		t.Skipf("Skipping: set ALICLOUD_ACCOUNT_ID_2 or INVITED_ALICLOUD_ACCOUNT_ID, or set ALICLOUD_ACCESS_KEY_2/SECRET_KEY_2 so the invited account id can be derived")
	}

	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	client, err := sdk.NewClientWithAccessKey(region, ak, sk)
	if err != nil {
		t.Fatalf("failed to build STS client for invited account: %s", err)
	}
	result, err := testAccResourceManagerHandshakeCommonRequest(client, "Sts", "2015-04-01", "sts.aliyuncs.com", "GetCallerIdentity", nil)
	if err != nil {
		t.Skipf("Skipping: failed to derive invited account id from ALICLOUD_ACCESS_KEY_2: %s", err)
	}
	accountId := strings.TrimSpace(fmt.Sprint(result["AccountId"]))
	if accountId == "" || accountId == "<nil>" {
		t.Skipf("Skipping: STS GetCallerIdentity did not return AccountId for ALICLOUD_ACCESS_KEY_2")
	}
	resourceManagerHandshakeInvitedAccountId = accountId
	return resourceManagerHandshakeInvitedAccountId
}

func testAccResourceManagerHandshakeInvitedAccountId() string {
	if resourceManagerHandshakeInvitedAccountId != "" {
		return resourceManagerHandshakeInvitedAccountId
	}
	for _, envName := range []string{"ALICLOUD_ACCOUNT_ID_2", "INVITED_ALICLOUD_ACCOUNT_ID"} {
		if v := strings.TrimSpace(os.Getenv(envName)); v != "" {
			return v
		}
	}
	return ""
}

func testAccPreCheckHandshakeAcceptanceAccountDetached(t *testing.T) {
	testAccPreCheckResourceManagerHandshakeAccountDetached(t, testAccResolveResourceManagerHandshakeInvitedAccountId(t))
}

func testAccPreCheckResourceManagerHandshakeAccountDetached(t *testing.T, accountId string) {
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	ak, sk := testAccResourceManagerHandshakeManagementCredentials(t)
	client, err := sdk.NewClientWithAccessKey(region, ak, sk)
	if err != nil {
		t.Fatalf("failed to build Resource Manager client: %s", err)
	}

	testAccPreCheckHandshakeAcceptancePendingHandshakesCanceled(t, client, accountId)
	if _, err := testAccPreCheckHandshakeAcceptanceResourceManagerRequest(client, "GetAccount", map[string]string{"AccountId": accountId}); err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.Account"}) || NotFoundError(err) {
			time.Sleep(resourceManagerHandshakeConsistencyDelay)
			return
		}
		t.Fatalf("failed to check Resource Manager member account %s: %s", accountId, err)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = testAccPreCheckHandshakeAcceptanceResourceManagerRequest(client, "RemoveCloudAccount", map[string]string{"AccountId": accountId})
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"EntityNotExists.Account"}) || NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to remove Resource Manager member account %s before handshake acceptance test: %s", accountId, err)
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = testAccPreCheckHandshakeAcceptanceResourceManagerRequest(client, "GetAccount", map[string]string{"AccountId": accountId})
		if err != nil {
			if IsExpectedErrors(err, []string{"EntityNotExists.Account"}) || NotFoundError(err) {
				return nil
			}
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		wait()
		return resource.RetryableError(fmt.Errorf("Resource Manager member account %s is still attached", accountId))
	})
	if err != nil {
		t.Fatalf("failed waiting for Resource Manager member account %s to detach: %s", accountId, err)
	}
	// Resource Manager may report the account as detached before new invitations are acceptable.
	time.Sleep(resourceManagerHandshakeConsistencyDelay)
}

func testAccPreCheckHandshakeAcceptancePendingHandshakesCanceled(t *testing.T, client *sdk.Client, accountId string) {
	t.Helper()
	wait := incrementalWait(3*time.Second, 5*time.Second)
	pageNumber := 1
	for {
		response, err := testAccPreCheckHandshakeAcceptanceResourceManagerRequest(client, "ListHandshakesForResourceDirectory", map[string]string{
			"PageSize":   "100",
			"PageNumber": fmt.Sprint(pageNumber),
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"EntityNotExists.ResourceDirectory"}) || NotFoundError(err) {
				return
			}
			t.Fatalf("failed to list Resource Manager handshakes before test: %s", err)
		}

		handshakes, _ := response["Handshakes"].(map[string]interface{})
		items, _ := handshakes["Handshake"].([]interface{})
		for _, raw := range items {
			item, _ := raw.(map[string]interface{})
			if fmt.Sprint(item["TargetEntity"]) != accountId || fmt.Sprint(item["Status"]) != "Pending" {
				continue
			}
			handshakeId := strings.TrimSpace(fmt.Sprint(item["HandshakeId"]))
			if handshakeId == "" || handshakeId == "<nil>" {
				continue
			}
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err := testAccPreCheckHandshakeAcceptanceResourceManagerRequest(client, "CancelHandshake", map[string]string{"HandshakeId": handshakeId})
				if err != nil {
					if IsExpectedErrors(err, []string{"EntityNotExists.Handshake"}) || NotFoundError(err) {
						return nil
					}
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				t.Fatalf("failed to cancel pending Resource Manager handshake %s before test: %s", handshakeId, err)
			}
		}
		if len(items) < 100 {
			return
		}
		pageNumber++
	}
}

func testAccPreCheckHandshakeAcceptanceResourceManagerRequest(client *sdk.Client, action string, params map[string]string) (map[string]interface{}, error) {
	return testAccResourceManagerHandshakeCommonRequest(client, "ResourceManager", "2020-03-31", "resourcemanager.aliyuncs.com", action, params)
}

func testAccResourceManagerHandshakeCommonRequest(client *sdk.Client, product string, version string, domain string, action string, params map[string]string) (map[string]interface{}, error) {
	request := requests.NewCommonRequest()
	request.Method = requests.POST
	request.Scheme = "https"
	request.Domain = domain
	request.Version = version
	request.Product = product
	request.ApiName = action
	for key, value := range params {
		request.QueryParams[key] = value
	}
	request.TransToAcsRequest()

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	_ = json.Unmarshal([]byte(response.GetHttpContentString()), &result)
	if response.IsSuccess() {
		return result, nil
	}

	code := fmt.Sprint(result["Code"])
	if code == "" || code == "<nil>" {
		code = fmt.Sprint(response.GetHttpStatus())
	}
	return nil, &ProviderError{
		errorCode: code,
		message:   fmt.Sprint(result["Message"]),
	}
}

func TestAccAliCloudResourceManagerHandshakeAcceptance_basic(t *testing.T) {
	resourceId := "alicloud_resource_manager_handshake_acceptance.default"
	ra := resourceAttrInit(resourceId, ResourceManagerHandshakeAcceptanceMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccResourceManagerHandshakeAcceptance%d", rand)
	testAccUseResourceManagerHandshakeManagementAccount(t)
	invitedAccountId := testAccResolveResourceManagerHandshakeInvitedAccountId(t)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckEnterpriseAccountEnabled(t)
			testAccPreCheckHandshakeAcceptanceProfiles(t)
			testAccPreCheckHandshakeAcceptanceAccountDetached(t)
		},

		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckResourceManagerHandshakeAcceptanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: ResourceManagerHandshakeAcceptanceBasic(name, invitedAccountId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"handshake_id": CHECKSET,
						"status":       "Accepted",
					}),
				),
			},
		},
	})
}

func testAccCheckResourceManagerHandshakeAcceptanceDestroy(s *terraform.State) error {
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	ak, sk := testAccResourceManagerHandshakeManagementCredentialValues()
	if ak == "" || sk == "" {
		return WrapError(fmt.Errorf("management account credentials are not set"))
	}
	client, err := sdk.NewClientWithAccessKey(region, ak, sk)
	if err != nil {
		return WrapError(err)
	}

	accountId := testAccResourceManagerHandshakeInvitedAccountId()
	wait := incrementalWait(3*time.Second, 5*time.Second)
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := testAccPreCheckHandshakeAcceptanceResourceManagerRequest(client, "GetAccount", map[string]string{"AccountId": accountId})
		if err != nil {
			if IsExpectedErrors(err, []string{"EntityNotExists.Account"}) || NotFoundError(err) {
				return nil
			}
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		wait()
		return resource.RetryableError(fmt.Errorf("Resource Manager member account %s is still attached", accountId))
	})
}
