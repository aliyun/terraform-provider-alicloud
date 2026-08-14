package alicloud

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type testAccCmsContact struct {
	Identifier string
	Name       string
	Workspace  string
}

func testAccCmsUniqueResourceName() string {
	return acctest.RandomWithPrefix("tfacccms")
}

func testAccCmsCreateContact(t *testing.T, rand int, supportedRegions []connectivity.Region) *testAccCmsContact {
	t.Helper()
	contact := &testAccCmsContact{
		Name: fmt.Sprintf("tfacccms%x%x", rand, time.Now().UnixNano()),
	}
	if os.Getenv("TF_ACC") == "" {
		return contact
	}

	testAccPreCheckWithRegions(t, true, supportedRegions)
	rawClient, err := sharedClientForRegion(defaultRegionToTest)
	if err != nil {
		t.Fatalf("failed to get AliCloud client for CMS contact: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	contact.Workspace, err = testAccCmsFindWorkspace(client)
	if err != nil {
		t.Fatalf("failed to discover a CMS workspace for the contact fixture: %s", err)
	}
	query := map[string]*string{
		"RegionId": StringPointer(client.RegionId),
	}
	body := map[string]interface{}{
		"name":      contact.Name,
		"email":     "hello.uuuu@aaa.com",
		"source":    "OBS",
		"workspace": contact.Workspace,
	}
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", "/contact", query, nil, body, true)
		if err != nil {
			if testAccCmsContactRetryable(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to create CMS contact: %s", err)
	}
	contact.Identifier = testAccCmsContactIdentifier(response)

	// Register cleanup immediately after Create succeeds. CreateContact normally
	// returns the identifier in contactId, but if the response is malformed after the
	// service has committed the contact, recover it by the unique exact name.
	t.Cleanup(func() {
		if contact.Identifier == "" {
			recovered, recoverErr := testAccCmsWaitContact(client, "", contact.Name, contact.Workspace, 2*time.Minute, false)
			if recoverErr != nil {
				t.Errorf("failed to recover CMS contact %q for cleanup: %s", contact.Name, recoverErr)
				return
			}
			contact.Identifier = recovered.Identifier
		}
		if err := testAccCmsDeleteContact(client, contact.Identifier); err != nil {
			t.Errorf("failed to delete CMS contact %q: %s", contact.Identifier, err)
		}
	})

	found, err := testAccCmsWaitContact(client, contact.Identifier, contact.Name, contact.Workspace, 2*time.Minute, true)
	if err != nil {
		t.Fatalf("failed to read CMS contact after creation: %s", err)
	}
	contact.Identifier = found.Identifier
	contact.Workspace = found.Workspace
	return contact
}

func testAccCmsContactIdentifier(response map[string]interface{}) string {
	if identifier := testAccCmsContactString(response["contactId"]); identifier != "" {
		return identifier
	}
	data := response["data"]
	if object, ok := data.(map[string]interface{}); ok {
		for _, key := range []string{"identifier", "contactId"} {
			if identifier := strings.TrimSpace(fmt.Sprint(object[key])); identifier != "" && identifier != "<nil>" {
				return identifier
			}
		}
	}
	identifier := strings.TrimSpace(fmt.Sprint(data))
	if identifier == "" || identifier == "<nil>" || strings.HasPrefix(identifier, "map[") {
		return ""
	}
	return identifier
}

func testAccCmsFindWorkspace(client *connectivity.AliyunClient) (string, error) {
	regionId := strings.TrimSpace(client.RegionId)
	if regionId == "" {
		return "", fmt.Errorf("effective test region is empty")
	}

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		var getErr error
		response, getErr = client.RoaGet("Cms", "2024-03-30", "/workspace", testAccCmsWorkspaceQuery(regionId), nil, nil)
		if getErr != nil {
			if testAccCmsContactRetryable(getErr) {
				wait()
				return resource.RetryableError(getErr)
			}
			return resource.NonRetryableError(getErr)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	workspace := testAccCmsWorkspaceFromResponse(response, regionId)
	if workspace == "" {
		return "", fmt.Errorf("no CMS workspace is available in the effective test region")
	}
	return workspace, nil
}

func testAccCmsWorkspaceQuery(regionId string) map[string]*string {
	return map[string]*string{
		"region":     StringPointer(regionId),
		"maxResults": StringPointer("200"),
	}
}

func testAccCmsWorkspaceFromResponse(response map[string]interface{}, regionId string) string {
	workspaces, ok := response["workspaces"].([]interface{})
	if !ok {
		if data, dataOK := response["data"].(map[string]interface{}); dataOK {
			workspaces, ok = data["workspaces"].([]interface{})
		}
	}
	if !ok {
		return ""
	}

	var candidates []string
	for _, rawWorkspace := range workspaces {
		workspace, ok := rawWorkspace.(map[string]interface{})
		if !ok {
			continue
		}
		name := testAccCmsContactString(workspace["workspaceName"])
		actualRegion := testAccCmsContactString(workspace["regionId"])
		if name == "" || actualRegion != "" && actualRegion != regionId {
			continue
		}
		candidates = append(candidates, name)
	}
	if len(candidates) == 0 {
		return ""
	}
	sort.Strings(candidates)
	return candidates[0]
}

func testAccCmsWaitContact(client *connectivity.AliyunClient, identifier, name, workspace string, timeout time.Duration, requireWorkspace bool) (*testAccCmsContact, error) {
	query := testAccCmsListContactQuery(client.RegionId, identifier, workspace)
	if identifier == "" && name != "" {
		query["name"] = StringPointer(name)
	}

	var found *testAccCmsContact
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(timeout, func() *resource.RetryError {
		response, getErr := client.RoaGet("Cms", "2024-03-30", "/contact", query, nil, nil)
		if getErr != nil {
			if testAccCmsContactRetryable(getErr) {
				wait()
				return resource.RetryableError(getErr)
			}
			return resource.NonRetryableError(getErr)
		}
		found = testAccCmsContactFromResponse(response, identifier, name)
		readiness := testAccCmsContactReadiness(found, workspace, requireWorkspace)
		if readiness != "ready" {
			wait()
			return resource.RetryableError(fmt.Errorf("CMS contact %q is not ready: %s", name, readiness))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return found, nil
}

func testAccCmsListContactQuery(regionId, identifier, workspace string) map[string]*string {
	query := map[string]*string{
		"RegionId": StringPointer(regionId),
		"source":   StringPointer("OBS"),
	}
	if identifier != "" {
		query["contactIds"] = StringPointer(convertObjectToJsonString([]string{identifier}))
	}
	if workspace != "" {
		query["workspace"] = StringPointer(workspace)
	}
	return query
}

func testAccCmsContactFromResponse(response map[string]interface{}, identifier, name string) *testAccCmsContact {
	contacts, ok := response["contacts"].([]interface{})
	if !ok {
		data, dataOK := response["data"].(map[string]interface{})
		if !dataOK {
			return nil
		}
		contacts, ok = data["data"].([]interface{})
		if !ok {
			return nil
		}
	}
	for _, rawContact := range contacts {
		object, ok := rawContact.(map[string]interface{})
		if !ok {
			continue
		}
		actualIdentifier := testAccCmsContactString(object["contactId"])
		if actualIdentifier == "" {
			actualIdentifier = testAccCmsContactString(object["identifier"])
		}
		actualName := testAccCmsContactString(object["name"])
		if identifier != "" && actualIdentifier != identifier {
			continue
		}
		if identifier == "" && actualName != name {
			continue
		}
		return &testAccCmsContact{
			Identifier: actualIdentifier,
			Name:       actualName,
			Workspace:  testAccCmsContactString(object["workspace"]),
		}
	}
	return nil
}

func testAccCmsContactReadiness(contact *testAccCmsContact, expectedWorkspace string, requireWorkspace bool) string {
	if contact == nil {
		return "contact not found"
	}
	if contact.Identifier == "" {
		return "identifier is empty"
	}
	if !requireWorkspace {
		return "ready"
	}
	if contact.Workspace == "" {
		return "workspace is empty"
	}
	if expectedWorkspace != "" && contact.Workspace != expectedWorkspace {
		return "workspace does not match the selected workspace"
	}
	return "ready"
}

func testAccCmsContactString(value interface{}) string {
	result := strings.TrimSpace(fmt.Sprint(value))
	if result == "<nil>" {
		return ""
	}
	return result
}

func testAccCmsContactRetryable(err error) bool {
	return NeedRetry(err) || IsExpectedErrors(err, []string{"FrequencyLimit"})
}

func testAccCmsDeleteContact(client *connectivity.AliyunClient, identifier string) error {
	if identifier == "" {
		return fmt.Errorf("empty identifier")
	}
	deleteQuery := testAccCmsDeleteContactQuery(client.RegionId, identifier)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.RoaDelete("Cms", "2024-03-30", "/contacts", deleteQuery, nil, nil, true)
		if err == nil || NotFoundError(err) || IsExpectedErrors(err, []string{"NotFound", "ResourceNotFound"}) {
			return nil
		}
		if testAccCmsContactRetryable(err) {
			wait()
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})
}

func testAccCmsDeleteContactQuery(regionId, identifier string) map[string]*string {
	return map[string]*string{
		"RegionId":   StringPointer(regionId),
		"contactIds": StringPointer(convertObjectToJsonString([]string{identifier})),
		"source":     StringPointer("OBS"),
	}
}

func TestUnitCmsContactListQuery(t *testing.T) {
	query := testAccCmsListContactQuery("cn-hangzhou", "contact-123", "workspace-123")
	if _, ok := query["identifiers"]; ok {
		t.Fatal("ListContacts must not expose the backend query name identifiers")
	}
	contactIds, ok := query["contactIds"]
	if !ok || contactIds == nil || *contactIds != `["contact-123"]` {
		t.Fatalf("ListContacts contactIds = %v, want JSON array with contact-123", contactIds)
	}
	workspace, ok := query["workspace"]
	if !ok || workspace == nil || *workspace != "workspace-123" {
		t.Fatalf("ListContacts workspace = %v, want workspace-123", workspace)
	}
}

func TestUnitCmsContactDeleteQuery(t *testing.T) {
	query := testAccCmsDeleteContactQuery("cn-hangzhou", "contact-123")
	if _, ok := query["identifiers"]; ok {
		t.Fatal("DeleteContacts must not expose the backend query name identifiers")
	}
	contactIds, ok := query["contactIds"]
	if !ok || contactIds == nil || *contactIds != `["contact-123"]` {
		t.Fatalf("DeleteContacts contactIds = %v, want JSON array with contact-123", contactIds)
	}
}

func TestUnitCmsContactResponseParsing(t *testing.T) {
	if got := testAccCmsContactIdentifier(map[string]interface{}{"contactId": "contact-123"}); got != "contact-123" {
		t.Fatalf("CreateContact identifier = %q, want contact-123", got)
	}

	contact := testAccCmsContactFromResponse(map[string]interface{}{
		"contacts": []interface{}{
			map[string]interface{}{
				"contactId": "contact-123",
				"name":      "tfacccmscontact",
				"workspace": "workspace-from-service",
			},
		},
	}, "contact-123", "tfacccmscontact")
	if contact == nil || contact.Identifier != "contact-123" || contact.Workspace != "workspace-from-service" {
		t.Fatalf("ListContacts contact = %#v, want service identifier and workspace", contact)
	}
}

func TestUnitCmsWorkspaceSelection(t *testing.T) {
	query := testAccCmsWorkspaceQuery("cn-test-1")
	region, ok := query["region"]
	if !ok || region == nil || *region != "cn-test-1" {
		t.Fatalf("ListWorkspaces region = %v, want cn-test-1", region)
	}
	maxResults, ok := query["maxResults"]
	if !ok || maxResults == nil || *maxResults != "200" {
		t.Fatalf("ListWorkspaces maxResults = %v, want 200", maxResults)
	}

	workspace := testAccCmsWorkspaceFromResponse(map[string]interface{}{
		"workspaces": []interface{}{
			map[string]interface{}{"workspaceName": "workspace-z", "regionId": "cn-test-1"},
			map[string]interface{}{"workspaceName": "workspace-other", "regionId": "cn-test-2"},
			map[string]interface{}{"workspaceName": "workspace-a", "regionId": "cn-test-1"},
		},
	}, "cn-test-1")
	if workspace != "workspace-a" {
		t.Fatalf("selected workspace = %q, want deterministic workspace-a", workspace)
	}
	if workspace := testAccCmsWorkspaceFromResponse(map[string]interface{}{"workspaces": []interface{}{}}, "cn-test-1"); workspace != "" {
		t.Fatalf("selected workspace from empty response = %q, want empty", workspace)
	}
}

func TestUnitCmsContactReadiness(t *testing.T) {
	testCases := []struct {
		name              string
		contact           *testAccCmsContact
		expectedWorkspace string
		want              string
	}{
		{name: "contact missing", want: "contact not found"},
		{name: "identifier missing", contact: &testAccCmsContact{Workspace: "workspace-123"}, expectedWorkspace: "workspace-123", want: "identifier is empty"},
		{name: "workspace empty", contact: &testAccCmsContact{Identifier: "contact-123"}, expectedWorkspace: "workspace-123", want: "workspace is empty"},
		{name: "workspace mismatch", contact: &testAccCmsContact{Identifier: "contact-123", Workspace: "workspace-other"}, expectedWorkspace: "workspace-123", want: "workspace does not match the selected workspace"},
		{name: "ready", contact: &testAccCmsContact{Identifier: "contact-123", Workspace: "workspace-123"}, expectedWorkspace: "workspace-123", want: "ready"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := testAccCmsContactReadiness(testCase.contact, testCase.expectedWorkspace, true); got != testCase.want {
				t.Fatalf("readiness = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestUnitCmsContactFrequencyLimitIsRetryable(t *testing.T) {
	if !testAccCmsContactRetryable(fmt.Errorf("SDKError: FrequencyLimit")) {
		t.Fatal("CreateContact FrequencyLimit must be retryable")
	}
	if testAccCmsContactRetryable(fmt.Errorf("InvalidParameter")) {
		t.Fatal("CreateContact InvalidParameter must not be retryable")
	}
}

func TestUnitCmsUniqueResourceName(t *testing.T) {
	seen := make(map[string]struct{})
	for i := 0; i < 32; i++ {
		name := testAccCmsUniqueResourceName()
		if !strings.HasPrefix(name, "tfacccms-") {
			t.Fatalf("CMS resource name %q does not use the expected prefix", name)
		}
		if _, exists := seen[name]; exists {
			t.Fatalf("CMS resource name %q was generated more than once", name)
		}
		seen[name] = struct{}{}
	}
}
