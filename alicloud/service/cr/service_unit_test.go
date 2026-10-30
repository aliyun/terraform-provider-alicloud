package cr

import (
	"strings"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

func TestUnitRepositoryNotExistClassification(t *testing.T) {
	for _, code := range []string{"REPO_NOT_EXIST", "INSTANCE_NOT_EXIST", "ResourceNotFound.Repository", "ResourceNotFound"} {
		err := &tea.SDKError{Code: tea.String(code), Data: tea.String("")}
		if !repositoryNotExist(err) {
			t.Fatalf("code %s must classify as repository not found", code)
		}
	}
	for _, code := range []string{"NoPrivilege", "SystemError", "InvalidNamespaceName"} {
		err := &tea.SDKError{Code: tea.String(code), Data: tea.String("")}
		if repositoryNotExist(err) {
			t.Fatalf("code %s must not classify as repository not found", code)
		}
	}
	if repositoryNotExist(nil) {
		t.Fatal("nil error must not classify as not found")
	}
}

func TestUnitRepositoryResponseNotFoundClassification(t *testing.T) {
	for _, code := range []string{"REPO_NOT_EXIST", "INSTANCE_NOT_EXIST", "ResourceNotFound.Repository", "ResourceNotFound"} {
		response := map[string]interface{}{"IsSuccess": false, "Code": code}
		if !repositoryResponseNotFound(response) {
			t.Fatalf("response code %s must classify as not found", code)
		}
	}
	cases := []map[string]interface{}{
		nil,
		{"IsSuccess": true},
		{"IsSuccess": false, "Code": "NoPrivilege"},
		{"IsSuccess": false},
	}
	for _, response := range cases {
		if repositoryResponseNotFound(response) {
			t.Fatalf("response %v must not classify as not found", response)
		}
	}
}

func TestUnitRepositoryCheckResponse(t *testing.T) {
	if err := repositoryCheckResponse("GetRepository", nil); err == nil {
		t.Fatal("nil response must be rejected")
	}
	if err := repositoryCheckResponse("GetRepository", map[string]interface{}{"RepoId": "crr-1"}); err == nil {
		t.Fatal("response without IsSuccess must be rejected")
	}
	if err := repositoryCheckResponse("GetRepository", map[string]interface{}{"IsSuccess": "true"}); err == nil {
		t.Fatal("non-boolean IsSuccess must be rejected")
	}
	if err := repositoryCheckResponse("GetRepository", map[string]interface{}{"IsSuccess": true}); err != nil {
		t.Fatalf("successful response must pass, got %v", err)
	}
	err := repositoryCheckResponse("UpdateRepository", map[string]interface{}{"IsSuccess": false, "Code": "NoPrivilege"})
	if err == nil || !strings.Contains(err.Error(), "NoPrivilege") {
		t.Fatalf("failed response must surface its code, got %v", err)
	}
}

func TestUnitRepositorySelectReadResponse(t *testing.T) {
	identities := []fwReadIdentity{
		{Path: "$.InstanceId", Value: "cri-1", AllowMissing: true},
		{Path: "$.RepoId", Value: "crr-1", AllowMissing: false},
	}
	object, err := repositorySelectReadResponse(map[string]interface{}{"IsSuccess": true, "InstanceId": "cri-1", "RepoId": "crr-1"}, identities, "cri-1:crr-1")
	if err != nil || object == nil {
		t.Fatalf("matching identity must be selected, got object=%v err=%v", object, err)
	}
	_, err = repositorySelectReadResponse(map[string]interface{}{"IsSuccess": true, "InstanceId": "cri-1"}, identities, "cri-1:crr-1")
	if err == nil {
		t.Fatal("successful response without RepoId must be an error, not NotFound")
	}
}

func TestUnitRepositoryNotExistKeepsOtherErrors(t *testing.T) {
	err := &tea.SDKError{Code: tea.String("NoPrivilege"), Message: tea.String("denied"), Data: tea.String("")}
	if repositoryNotExist(err) || errs.IsExpectedErrors(err, []string{"REPO_NOT_EXIST"}) {
		t.Fatal("privilege errors must not be swallowed by the not-found classification")
	}
}
