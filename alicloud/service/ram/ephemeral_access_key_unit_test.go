package ram

import (
	"errors"
	"testing"
)

type fakeRamAccessKeyCall struct {
	Product   string
	Version   string
	Action    string
	Body      map[string]interface{}
	AutoRetry bool
}

type fakeRamAccessKeyClient struct {
	calls []fakeRamAccessKeyCall
	errs  []error
}

func (f *fakeRamAccessKeyClient) RpcPost(apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error) {
	f.calls = append(f.calls, fakeRamAccessKeyCall{
		Product:   apiProductCode,
		Version:   apiVersion,
		Action:    apiName,
		Body:      body,
		AutoRetry: autoRetry,
	})
	idx := len(f.calls) - 1
	if idx < len(f.errs) && f.errs[idx] != nil {
		return nil, f.errs[idx]
	}
	return map[string]interface{}{}, nil
}

func TestUnitRamAccessKeyEphemeralCreateIsNotRetried(t *testing.T) {
	client := &fakeRamAccessKeyClient{}
	request := map[string]interface{}{"UserName": "tf-testacc-ramak"}

	if _, err := createRamAccessKey(client, request); err != nil {
		t.Fatalf("createRamAccessKey returned error: %s", err)
	}

	if len(client.calls) != 1 {
		t.Fatalf("expected exactly 1 CreateAccessKey call, got %d", len(client.calls))
	}
	call := client.calls[0]
	if call.Product != "Ram" || call.Version != "2015-05-01" {
		t.Errorf("called %s %s, want Ram 2015-05-01", call.Product, call.Version)
	}
	if call.Action != "CreateAccessKey" {
		t.Errorf("called %s, want CreateAccessKey", call.Action)
	}
	if call.AutoRetry {
		t.Error("CreateAccessKey must not enable automatic retries: the operation has no idempotency token")
	}
	if got := call.Body["UserName"]; got != "tf-testacc-ramak" {
		t.Errorf("sent UserName %v, want tf-testacc-ramak", got)
	}
}

func TestUnitRamAccessKeyEphemeralCreateSurfacesError(t *testing.T) {
	client := &fakeRamAccessKeyClient{errs: []error{errors.New("Throttling.User")}}

	if _, err := createRamAccessKey(client, map[string]interface{}{}); err == nil {
		t.Fatal("expected CreateAccessKey error to surface, got nil")
	}
}

func TestUnitRamAccessKeyEphemeralCloseIdempotent(t *testing.T) {
	client := &fakeRamAccessKeyClient{errs: []error{nil, errors.New("EntityNotExist.User.AccessKey")}}
	cleanup := ramAccessKeyCleanup{AccessKeyId: "LTAI5tTestAccessKey", UserName: "tf-testacc-ramak"}

	for i := 0; i < 2; i++ {
		if err := deleteRamAccessKey(client, cleanup); err != nil {
			t.Fatalf("Close call %d returned error: %s", i+1, err)
		}
	}

	if len(client.calls) != 2 {
		t.Fatalf("expected 2 DeleteAccessKey calls, got %d", len(client.calls))
	}
	for i, call := range client.calls {
		if call.Action != "DeleteAccessKey" {
			t.Errorf("call %d called %s, want DeleteAccessKey", i+1, call.Action)
		}
		if !call.AutoRetry {
			t.Errorf("call %d disabled automatic retries; DeleteAccessKey is idempotent and may be retried", i+1)
		}
		if got := call.Body["UserAccessKeyId"]; got != cleanup.AccessKeyId {
			t.Errorf("call %d sent UserAccessKeyId %v, want %s", i+1, got, cleanup.AccessKeyId)
		}
		if got := call.Body["UserName"]; got != cleanup.UserName {
			t.Errorf("call %d sent UserName %v, want %s", i+1, got, cleanup.UserName)
		}
	}
}

func TestUnitRamAccessKeyEphemeralCloseWithoutCleanup(t *testing.T) {
	client := &fakeRamAccessKeyClient{}

	if err := deleteRamAccessKey(client, ramAccessKeyCleanup{}); err != nil {
		t.Fatalf("Close without cleanup data returned error: %s", err)
	}
	if len(client.calls) != 0 {
		t.Fatalf("expected no DeleteAccessKey call without cleanup data, got %d", len(client.calls))
	}
}

func TestUnitRamAccessKeyEphemeralCloseSurfacesUnexpectedError(t *testing.T) {
	client := &fakeRamAccessKeyClient{errs: []error{errors.New("Forbidden.RAM")}}

	if err := deleteRamAccessKey(client, ramAccessKeyCleanup{AccessKeyId: "LTAI5tTestAccessKey"}); err == nil {
		t.Fatal("expected a non-NotFound error to surface, got nil")
	}
}

func TestUnitRamAccessKeyEphemeralCloseSurfacesOtherEntityNotExist(t *testing.T) {
	client := &fakeRamAccessKeyClient{errs: []error{errors.New("EntityNotExist.Role")}}

	if err := deleteRamAccessKey(client, ramAccessKeyCleanup{AccessKeyId: "LTAI5tTestAccessKey"}); err == nil {
		t.Fatal("expected an unrelated EntityNotExist error to surface, got nil")
	}
}
