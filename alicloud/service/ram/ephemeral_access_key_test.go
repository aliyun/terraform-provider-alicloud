package ram_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	sdk "github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	sdkerrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	sdkram "github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	sdksts "github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	ramservice "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ram"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	ephemeralschema "github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestUnitRamAccessKeyEphemeral_schema(t *testing.T) {
	ctx := context.Background()
	ephemeralResource := ramservice.NewAccessKeyEphemeralResource()

	var metadata ephemeral.MetadataResponse
	ephemeralResource.Metadata(ctx, ephemeral.MetadataRequest{}, &metadata)
	if metadata.TypeName != "alicloud_ram_access_key" {
		t.Fatalf("expected ephemeral type name %q, got %q", "alicloud_ram_access_key", metadata.TypeName)
	}

	var schemaResponse ephemeral.SchemaResponse
	ephemeralResource.Schema(ctx, ephemeral.SchemaRequest{}, &schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("schema returned diagnostics: %v", schemaResponse.Diagnostics)
	}

	attributes := schemaResponse.Schema.Attributes
	expected := map[string]struct {
		required  bool
		optional  bool
		computed  bool
		sensitive bool
	}{
		"user_name":         {optional: true},
		"access_key_id":     {computed: true},
		"access_key_secret": {computed: true, sensitive: true},
		"status":            {computed: true},
		"create_time":       {computed: true},
	}
	if len(attributes) != len(expected) {
		t.Errorf("expected %d attributes, got %d: %v", len(expected), len(attributes), attributes)
	}
	for name, flags := range expected {
		raw, ok := attributes[name]
		if !ok {
			t.Errorf("attribute %q is missing from the ephemeral schema", name)
			continue
		}
		attribute, ok := raw.(ephemeralschema.StringAttribute)
		if !ok {
			t.Errorf("attribute %q is %T, want schema.StringAttribute", name, raw)
			continue
		}
		if attribute.Required != flags.required {
			t.Errorf("attribute %q Required = %t, want %t", name, attribute.Required, flags.required)
		}
		if attribute.Optional != flags.optional {
			t.Errorf("attribute %q Optional = %t, want %t", name, attribute.Optional, flags.optional)
		}
		if attribute.Computed != flags.computed {
			t.Errorf("attribute %q Computed = %t, want %t", name, attribute.Computed, flags.computed)
		}
		if attribute.Sensitive != flags.sensitive {
			t.Errorf("attribute %q Sensitive = %t, want %t", name, attribute.Sensitive, flags.sensitive)
		}
	}
	for _, name := range []string{"pgp_key", "secret_file", "encrypted_secret", "key_fingerprint"} {
		if _, ok := attributes[name]; ok {
			t.Errorf("attribute %q must not be part of the ephemeral schema", name)
		}
	}
}

func TestAccAliCloudRamAccessKeyEphemeral_basic(t *testing.T) {
	rand := sdkacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-ram-access-key-ephemeral-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		CheckDestroy: testAccCheckRamUserDestroyed(name),
		Steps: []resource.TestStep{
			{
				Config: testAccRamAccessKeyEphemeralUserConfig(name),
			},
			{
				Config: testAccRamAccessKeyEphemeralConfig(name),
				Check:  testAccCheckNoAccessKeyEphemeralInState,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.open", tfjsonpath.New("data").AtMapKey("user_name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("echo.open", tfjsonpath.New("data").AtMapKey("access_key_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", tfjsonpath.New("data").AtMapKey("access_key_secret"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", tfjsonpath.New("data").AtMapKey("status"), knownvalue.StringExact("Active")),
					statecheck.ExpectKnownValue("echo.open", tfjsonpath.New("data").AtMapKey("create_time"), knownvalue.NotNull()),
					accessKeyUsableCheck{userName: name},
					// Runs after the usability check and before the post-apply plans:
					// every plan and apply that opens the ephemeral resource creates a
					// new permanent AccessKey, and a RAM user only has room for two.
					deleteRamAccessKeysCheck{userName: name},
				},
			},
			{
				// Remove the ephemeral resource before the final destroy so the
				// destroy plan does not open it again.
				Config: testAccRamAccessKeyEphemeralUserConfig(name),
				Check:  testAccDeleteRamAccessKeys(name),
			},
		},
	})
}

func TestAccAliCloudRamAccessKeyEphemeral_userNotExist(t *testing.T) {
	rand := sdkacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-ram-access-key-ephemeral-missing-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		CheckDestroy: acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				Config:      testAccRamAccessKeyEphemeralConfigForMissingUser(name),
				ExpectError: regexp.MustCompile(`Opening alicloud_ram_access_key: calling CreateAccessKey`),
			},
		},
	})
}

// accessKeyUsableCheck proves the credential handed out by the ephemeral
// resource actually authenticates, using the values echoed into test state.
type accessKeyUsableCheck struct {
	userName string
}

func (c accessKeyUsableCheck) CheckState(_ context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	state := req.State
	if state == nil || state.Values == nil || state.Values.RootModule == nil {
		resp.Error = errors.New("no state available to verify the opened access key")
		return
	}
	var echoData interface{}
	found := false
	for _, resourceState := range state.Values.RootModule.Resources {
		if resourceState.Address == "echo.open" {
			echoData = resourceState.AttributeValues["data"]
			found = true
			break
		}
	}
	if !found {
		resp.Error = errors.New(`resource "echo.open" is missing from state`)
		return
	}
	data, ok := echoData.(map[string]interface{})
	if !ok {
		resp.Error = fmt.Errorf(`attribute "data" of "echo.open" has type %T, want object`, echoData)
		return
	}
	accessKeyId, _ := data["access_key_id"].(string)
	accessKeySecret, _ := data["access_key_secret"].(string)
	if accessKeyId == "" || accessKeySecret == "" {
		resp.Error = fmt.Errorf("opened credential is incomplete: access_key_id set=%t, access_key_secret set=%t", accessKeyId != "", accessKeySecret != "")
		return
	}

	client, err := sdksts.NewClientWithOptions(testAccEphemeralRegion(), testAccSdkConfig(), credentials.NewAccessKeyCredential(accessKeyId, accessKeySecret))
	if err != nil {
		resp.Error = fmt.Errorf("creating an STS client with the opened credential: %w", err)
		return
	}
	response, err := client.GetCallerIdentity(sdksts.CreateGetCallerIdentityRequest())
	if err != nil {
		resp.Error = fmt.Errorf("calling STS GetCallerIdentity with the opened credential: %w", err)
		return
	}
	if response.IdentityType != "RAMUser" {
		resp.Error = fmt.Errorf("opened credential acts as %q, want RAMUser", response.IdentityType)
		return
	}
	if !strings.HasSuffix(response.Arn, ":user/"+c.userName) {
		resp.Error = fmt.Errorf("opened credential belongs to %q, want the RAM user %q", response.Arn, c.userName)
		return
	}
}

type deleteRamAccessKeysCheck struct {
	userName string
}

func (c deleteRamAccessKeysCheck) CheckState(_ context.Context, _ statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	resp.Error = deleteRamAccessKeys(c.userName)
}

func testAccRamAccessKeyEphemeralUserConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ram_user" "default" {
  name  = var.name
  force = true
}
`, name)
}

func testAccRamAccessKeyEphemeralConfig(name string) string {
	return fmt.Sprintf(`
%s

ephemeral "alicloud_ram_access_key" "default" {
  user_name = alicloud_ram_user.default.name
}

provider "echo" {
  data = ephemeral.alicloud_ram_access_key.default
}

resource "echo" "open" {}
`, testAccRamAccessKeyEphemeralUserConfig(name))
}

func testAccRamAccessKeyEphemeralConfigForMissingUser(name string) string {
	return fmt.Sprintf(`
ephemeral "alicloud_ram_access_key" "default" {
  user_name = "%s"
}

provider "echo" {
  data = ephemeral.alicloud_ram_access_key.default
}

resource "echo" "open" {}
`, name)
}

func testAccCheckNoAccessKeyEphemeralInState(state *terraform.State) error {
	for address, resourceState := range state.RootModule().Resources {
		if resourceState.Type == "alicloud_ram_access_key" {
			return fmt.Errorf("ephemeral alicloud_ram_access_key must not be persisted to state, found %s", address)
		}
	}
	return nil
}

func testAccDeleteRamAccessKeys(userName string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		return deleteRamAccessKeys(userName)
	}
}

func deleteRamAccessKeys(userName string) error {
	client, err := testAccRamClient()
	if err != nil {
		return err
	}
	listRequest := sdkram.CreateListAccessKeysRequest()
	listRequest.UserName = userName
	listResponse, err := client.ListAccessKeys(listRequest)
	if err != nil {
		if isRamEntityNotExist(err) {
			return nil
		}
		return fmt.Errorf("listing access keys of RAM user %s: %w", userName, err)
	}
	for _, accessKey := range listResponse.AccessKeys.AccessKey {
		deleteRequest := sdkram.CreateDeleteAccessKeyRequest()
		deleteRequest.UserAccessKeyId = accessKey.AccessKeyId
		deleteRequest.UserName = userName
		if _, err := client.DeleteAccessKey(deleteRequest); err != nil && !isRamEntityNotExist(err) {
			return fmt.Errorf("deleting access key %s of RAM user %s: %w", accessKey.AccessKeyId, userName, err)
		}
	}
	return nil
}

func testAccCheckRamUserDestroyed(userName string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		client, err := testAccRamClient()
		if err != nil {
			return err
		}
		getUserRequest := sdkram.CreateGetUserRequest()
		getUserRequest.UserName = userName
		if _, err := client.GetUser(getUserRequest); err != nil {
			if isRamEntityNotExist(err) {
				return nil
			}
			return fmt.Errorf("reading RAM user %s after destroy: %w", userName, err)
		}
		return fmt.Errorf("RAM user %s still exists after destroy", userName)
	}
}

// testAccRamClient builds a RAM client from the same environment credentials the
// provider under test uses. It must not reuse acctest.Provider.Meta(): provider
// blocks in the test configuration can re-configure that shared instance.
func testAccRamClient() (*sdkram.Client, error) {
	accessKey := os.Getenv("ALICLOUD_ACCESS_KEY")
	secretKey := os.Getenv("ALICLOUD_SECRET_KEY")
	if accessKey == "" || secretKey == "" {
		return nil, errors.New("ALICLOUD_ACCESS_KEY and ALICLOUD_SECRET_KEY must be set for the RAM test client")
	}
	if securityToken := os.Getenv("ALICLOUD_SECURITY_TOKEN"); securityToken != "" {
		return sdkram.NewClientWithOptions(testAccEphemeralRegion(), testAccSdkConfig(), credentials.NewStsTokenCredential(accessKey, secretKey, securityToken))
	}
	return sdkram.NewClientWithOptions(testAccEphemeralRegion(), testAccSdkConfig(), credentials.NewAccessKeyCredential(accessKey, secretKey))
}

// RAM and STS reject plain HTTP with InvalidProtocol.NeedSsl, and the SDK
// defaults to HTTP.
func testAccSdkConfig() *sdk.Config {
	return sdk.NewConfig().WithScheme("HTTPS")
}

func testAccEphemeralRegion() string {
	if region := os.Getenv("ALICLOUD_REGION"); region != "" {
		return region
	}
	return "cn-beijing"
}

func isRamEntityNotExist(err error) bool {
	var serverError *sdkerrors.ServerError
	return errors.As(err, &serverError) && strings.HasPrefix(serverError.ErrorCode(), "EntityNotExist")
}
