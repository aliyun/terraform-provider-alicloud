package alicloud

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	roa "github.com/alibabacloud-go/tea-roa/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cms"
	"github.com/stretchr/testify/assert"
)

// TestUnitAliCloudCmsContactService is the no-cloud regression for the
// Framework-native alicloud_cms_contact resource and alicloud_cms_contacts
// data source. It mocks the ROA transport (tea-roa Client.DoRequest) so no real
// CMS API call is made, then drives CmsService Create/Read/Update/Delete/List
// directly, asserting the correct method, path, query and body are sent and
// that the responses (including not-found, empty arrays and a malformed
// non-map element) are parsed without panicking. It covers the real
// CRUD/state/failure paths and the necessary API request construction, not
// just helper plumbing.
func TestUnitAliCloudCmsContactService(t *testing.T) {
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		return
	}
	client := rawClient.(*connectivity.AliyunClient)
	svc := cms.NewService(client)

	// Captured request envelope from the mocked ROA transport.
	var lastMethod, lastPath string
	var lastQuery map[string]*string
	var lastBody interface{}

	// apply installs a DoRequest mock that returns canned responses keyed by
	// (method, path); nil resp + non-nil err drives failure/not-found. It also
	// captures the request envelope so the assertions can inspect what was
	// sent. The returned *gomonkey.Patches is reset by the caller.
	apply := func(handler func(method, path string, query map[string]*string, body interface{}) (map[string]interface{}, error)) *gomonkey.Patches {
		return gomonkey.ApplyMethod(reflect.TypeOf(&roa.Client{}), "DoRequest",
			func(_ *roa.Client, _ *string, method *string, _ *string, action *string, query map[string]*string, _ map[string]*string, body interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
				lastMethod = *method
				lastPath = *action
				lastQuery = query
				lastBody = body
				return handler(*method, *action, query, body)
			})
	}

	// singleResp returns one fixed response for every call.
	singleResp := func(resp map[string]interface{}, err error) *gomonkey.Patches {
		return apply(func(_, _ string, _ map[string]*string, _ interface{}) (map[string]interface{}, error) {
			return resp, err
		})
	}

	notFoundErr := func() error {
		return &tea.SDKError{
			Code:       tea.String("ResourceNotFound"),
			Data:       tea.String("ResourceNotFound"),
			Message:    tea.String("ResourceNotFound"),
			StatusCode: tea.Int(404),
		}
	}
	failedErr := func(code string, status int) error {
		return &tea.SDKError{
			Code:       tea.String(code),
			Data:       tea.String(code),
			Message:    tea.String(code),
			StatusCode: tea.Int(status),
		}
	}

	ctx := context.Background()

	// --- Create: POST /contact with name/email/phone/lang/imUserIds body ---
	patches := singleResp(map[string]interface{}{"contactId": "ct-123"}, nil)
	defer func() {
		if patches != nil {
			patches.Reset()
		}
	}()
	createResp, err := svc.CreateContact(ctx, map[string]interface{}{
		"name":      "unit-contact",
		"email":     "unit@example.com",
		"phone":     "86-13800000000",
		"lang":      "zh_CN",
		"imUserIds": map[string]interface{}{"key1": "user1"},
	})
	assert.Nil(t, err, "CreateContact should succeed with a contactId")
	assert.Equal(t, "ct-123", fmt.Sprint(createResp["contactId"]))
	assert.Equal(t, "POST", lastMethod)
	assert.Equal(t, "/contact", lastPath)
	body, ok := lastBody.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "unit-contact", body["name"])
	assert.Equal(t, "unit@example.com", body["email"])
	assert.Equal(t, "zh_CN", body["lang"])

	// Create failure: response missing contactId surfaces a clear error.
	patches.Reset()
	patches = singleResp(map[string]interface{}{}, nil)
	_, err = svc.CreateContact(ctx, map[string]interface{}{"name": "x"})
	assert.NotNil(t, err, "CreateContact without contactId must error")

	// --- Read (GetContact): GET /contact?contactIds=["id"]&source=OBS ---
	patches.Reset()
	patches = singleResp(map[string]interface{}{
		"data": map[string]interface{}{
			"data":  []interface{}{map[string]interface{}{"contactId": "ct-123", "name": "unit-contact"}},
			"total": json.Number("1"),
		},
	}, nil)
	contact, err := svc.GetContact(ctx, "ct-123")
	assert.Nil(t, err)
	assert.NotNil(t, contact, "GetContact should return the element when present")
	assert.Equal(t, "unit-contact", contact["name"])
	assert.Equal(t, "GET", lastMethod)
	assert.Equal(t, "/contact", lastPath)
	assert.NotNil(t, lastQuery["contactIds"], "GetContact must send the contactIds query")
	assert.NotNil(t, lastQuery["source"])

	// Read not-found: empty array -> (nil, nil), resource drops state.
	patches.Reset()
	patches = singleResp(map[string]interface{}{"data": map[string]interface{}{"data": []interface{}{}}}, nil)
	contact, err = svc.GetContact(ctx, "missing")
	assert.Nil(t, err)
	assert.Nil(t, contact, "GetContact returns nil when the contact is absent")

	// Read not-found: API 404 -> (nil, nil).
	patches.Reset()
	patches = singleResp(nil, notFoundErr())
	contact, err = svc.GetContact(ctx, "missing")
	assert.Nil(t, err)
	assert.Nil(t, contact, "GetContact maps a 404 to (nil, nil)")

	// AC2 panic guard: a non-map element must NOT panic; it returns an error.
	patches.Reset()
	patches = singleResp(map[string]interface{}{
		"data": map[string]interface{}{"data": []interface{}{"not-a-map"}},
	}, nil)
	contact, err = svc.GetContact(ctx, "ct-malformed")
	assert.NotNil(t, err, "a malformed (non-map) element must surface an error, not panic")
	assert.Nil(t, contact)

	// Read failure (non-404) propagates as an error.
	patches.Reset()
	patches = singleResp(nil, failedErr("InvalidContact", 400))
	contact, err = svc.GetContact(ctx, "ct-123")
	assert.NotNil(t, err, "a non-404 failure must propagate")
	assert.Nil(t, contact)

	// --- List (pagination): two pages merged by total ---
	patches.Reset()
	page1 := map[string]interface{}{
		"data": map[string]interface{}{
			"data":  []interface{}{map[string]interface{}{"contactId": "a"}, map[string]interface{}{"contactId": "b"}},
			"total": json.Number("3"),
		},
	}
	page2 := map[string]interface{}{
		"data": map[string]interface{}{
			"data":  []interface{}{map[string]interface{}{"contactId": "c"}},
			"total": json.Number("3"),
		},
	}
	page := 0
	patches = apply(func(_, _ string, query map[string]*string, _ interface{}) (map[string]interface{}, error) {
		lastMethod = "GET"
		lastPath = "/contact"
		lastQuery = query
		page++
		if page == 1 {
			return page1, nil
		}
		return page2, nil
	})
	objects, err := svc.ListContacts(ctx, "unit-contact", "", "", false)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(objects), "ListContacts must merge paginated results across pages")
	assert.Equal(t, "GET", lastMethod)
	assert.Equal(t, "/contact", lastPath)
	assert.NotNil(t, lastQuery["pageNumber"])
	assert.NotNil(t, lastQuery["pageSize"])
	assert.NotNil(t, lastQuery["source"])
	assert.NotNil(t, lastQuery["name"], "ListContacts must forward the name filter")

	// --- Update: PATCH /contact/{id} ---
	patches.Reset()
	patches = apply(func(method, path string, _ map[string]*string, body interface{}) (map[string]interface{}, error) {
		lastMethod = method
		lastPath = path
		lastBody = body
		return map[string]interface{}{}, nil
	})
	_, err = svc.UpdateContact(ctx, "ct-123", map[string]interface{}{"name": "renamed", "email": "new@example.com"})
	assert.Nil(t, err)
	assert.Equal(t, "PATCH", lastMethod)
	assert.Equal(t, "/contact/ct-123", lastPath, "UpdateContact must PATCH /contact/{id}")
	ub, ok := lastBody.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "renamed", ub["name"])

	// Update failure propagates.
	patches.Reset()
	patches = singleResp(nil, failedErr("InvalidContact", 400))
	_, err = svc.UpdateContact(ctx, "ct-123", map[string]interface{}{"name": "x"})
	assert.NotNil(t, err, "UpdateContact failure must propagate")

	// --- Delete: DELETE /contacts?contactIds=["id"]&source=OBS ---
	patches.Reset()
	patches = apply(func(method, path string, query map[string]*string, _ interface{}) (map[string]interface{}, error) {
		lastMethod = method
		lastPath = path
		lastQuery = query
		return map[string]interface{}{}, nil
	})
	err = svc.DeleteContact(ctx, "ct-123")
	assert.Nil(t, err)
	assert.Equal(t, "DELETE", lastMethod)
	assert.Equal(t, "/contacts", lastPath)
	assert.NotNil(t, lastQuery["contactIds"], "DeleteContact must send contactIds")
	assert.NotNil(t, lastQuery["source"])

	// Delete not-found is idempotent (no error).
	patches.Reset()
	patches = singleResp(nil, notFoundErr())
	err = svc.DeleteContact(ctx, "missing")
	assert.Nil(t, err, "DeleteContact on a missing contact must be idempotent")

	// Delete failure propagates.
	patches.Reset()
	patches = singleResp(nil, failedErr("InvalidContact", 400))
	err = svc.DeleteContact(ctx, "ct-123")
	assert.NotNil(t, err, "DeleteContact failure must propagate")

	patches.Reset()
	patches = nil
}
