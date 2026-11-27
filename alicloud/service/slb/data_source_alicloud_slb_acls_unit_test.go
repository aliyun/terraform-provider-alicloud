package slb

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	sdkauth "github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	sdkrequests "github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	sdkresponses "github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	slbsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// TestUnitAlicloudSlbAclsDataSourceTagQueryError guards the data-source
// error contract: when the auxiliary ListTagResources call of a single ACL
// fails, dataSourceAclsRead must return that error instead of the historical
// `return nil` inside slbAclsDescriptionAttributes, which silently ended the
// Read as a success before acls/ids/names were written. The Read is driven
// offline: the old-SDK list/attribute traffic is intercepted at
// (*sdk.Client).DoActionWithSigner, the deep entry every generated DoAction
// path converges on. The patch must be registered against the defining
// sdk.Client type, not *slb.Client: a method promoted through the embedded
// sdk.Client resolves reflect's method-table entry to a compiler-generated
// wrapper that the real (static) call chain never goes through, and the
// per-action Describe wrappers are thin enough to be inlined at the provider
// call sites — both leave the real RPC running. The tag query goes through
// the tea-rpc client whose DoRequest is patched per action, and the
// mock-only AliyunClient from unit_test_helpers_test.go supplies credentials
// and the pre-populated SLB endpoint, so no real RPC leaves the process.

// lintignore: R001
func TestUnitAlicloudSlbAclsDataSourceTagQueryError(t *testing.T) {
	p := DataSourceAcls()

	rawClient, err := unitTestMockClient()
	if err != nil {
		t.Fatalf("failed to build mock client: %s", err)
	}

	aclList := []slbsdk.Acl{
		{AclId: "acl-mock1", AclName: "tf-acl-1"},
		{AclId: "acl-mock2", AclName: "tf-acl-2"},
	}
	tagResponse := map[string]interface{}{
		"TagResources": map[string]interface{}{
			"TagResource": []interface{}{
				map[string]interface{}{
					"TagKey":   "Created",
					"TagValue": "TF",
				},
			},
		},
	}
	failedTagResponse := func() (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       helper.String("NonRetryableError"),
			Data:       helper.String("NonRetryableError"),
			Message:    helper.String("NonRetryableError"),
			StatusCode: tea.Int(400),
		}
	}

	// readWithTags drives one full dataSourceAclsRead run; tagCalls decides
	// which ListTagResources invocations fail: an entry true makes that ACL's
	// tag query return a non-retryable SDK error. The resource data is
	// returned so failure cases can assert no truncated output was written.
	readWithTags := func(tagCalls []bool) (error, *schema.ResourceData) {
		d, _ := schema.InternalMap(p.Schema).Data(nil, nil)

		patches := gomonkey.ApplyMethod(reflect.TypeOf(&sdk.Client{}), "DoActionWithSigner",
			func(_ *sdk.Client, request sdkrequests.AcsRequest, response sdkresponses.AcsResponse, _ sdkauth.Signer) error {
				switch req := request.(type) {
				case *slbsdk.DescribeAccessControlListsRequest:
					if resp, ok := response.(*slbsdk.DescribeAccessControlListsResponse); ok {
						resp.Acls = slbsdk.Acls{Acl: aclList}
					}
					return nil
				case *slbsdk.DescribeAccessControlListAttributeRequest:
					if resp, ok := response.(*slbsdk.DescribeAccessControlListAttributeResponse); ok {
						resp.AclId = req.AclId
						resp.AclName = "tf-acl-" + req.AclId[len(req.AclId)-1:]
						resp.AddressIPVersion = "ipv4"
						resp.ResourceGroupId = "rg-mock"
					}
					return nil
				default:
					return fmt.Errorf("unexpected slb SDK action under mock: %T", request)
				}
			})
		defer patches.Reset()

		callIndex := -1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest",
			func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				if *action == "ListTagResources" {
					callIndex++
					if callIndex < len(tagCalls) && tagCalls[callIndex] {
						return failedTagResponse()
					}
					return tagResponse, nil
				}
				return tagResponse, nil
			})
		defer patches.Reset()

		return dataSourceAclsRead(d, rawClient), d
	}

	t.Run("SingleAclTagQueryError", func(t *testing.T) {
		// The list returns two ACLs; the first ACL's tag query fails, so the
		// Read must abort with a wrapped error before writing any output.
		err, d := readWithTags([]bool{true, false})
		if err == nil {
			t.Fatal("expected the Read to fail when the first ACL's ListTagResources call fails, got nil")
		}
		assert.Contains(t, err.Error(), "alicloud_slb_acls")
		assert.Contains(t, err.Error(), "ListTagResources")
		assert.Empty(t, d.Id())
		assert.Equal(t, 0, len(d.Get("ids").([]interface{})))
		assert.Equal(t, 0, len(d.Get("names").([]interface{})))
		assert.Equal(t, 0, len(d.Get("acls").([]interface{})))
	})

	t.Run("MultipleAclsPartialTagQueryError", func(t *testing.T) {
		// The first ACL's tags resolve; the second ACL's tag query fails. The
		// partial result must not be exposed: Read returns an error and no
		// truncated acls/ids/names are written.
		err, d := readWithTags([]bool{false, true})
		if err == nil {
			t.Fatal("expected the Read to fail when a later ACL's ListTagResources call fails, got nil")
		}
		assert.Contains(t, err.Error(), "alicloud_slb_acls")
		assert.Contains(t, err.Error(), "ListTagResources")
		assert.Empty(t, d.Id())
		assert.Equal(t, 0, len(d.Get("ids").([]interface{})))
		assert.Equal(t, 0, len(d.Get("names").([]interface{})))
		assert.Equal(t, 0, len(d.Get("acls").([]interface{})))
	})

	t.Run("TagQuerySuccessKeepsOutputs", func(t *testing.T) {
		// Success path (multi ACL, no failing tag query): acls/ids/names are
		// complete and each ACL carries the queried tags, matching the
		// pre-fix behaviour.
		err, d := readWithTags(nil)
		assert.Nil(t, err)

		assert.Equal(t, 2, len(d.Get("ids").([]interface{})))
		assert.Equal(t, 2, len(d.Get("names").([]interface{})))
		acls := d.Get("acls").([]interface{})
		assert.Equal(t, 2, len(acls))
		for _, rawAcl := range acls {
			acl := rawAcl.(map[string]interface{})
			assert.Equal(t, "TF", acl["tags"].(map[string]interface{})["Created"])
		}
		assert.Equal(t, helper.DataResourceIdHash([]string{"acl-mock1", "acl-mock2"}), d.Id())
	})
}
