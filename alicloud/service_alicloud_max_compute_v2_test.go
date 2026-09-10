package alicloud

import (
	"strings"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
)

// TestDescribeMaxComputeRoleUserAttachment_ExternallyDeletedMember verifies
// that when MaxCompute returns a non-404 error whose message indicates the
// linked RAM/ALIYUN user referenced by a role has been deleted externally
// ("baseId and tenantId do not match"), the service-layer Read path in
// DescribeMaxComputeRoleUserAttachment surfaces a NotFoundError so the
// resource Read drops the attachment from state instead of failing
// refresh/apply.
//
// The test exercises the production describeMaxComputeRoleUserAttachmentWithRoaGet
// method (the testable core of DescribeMaxComputeRoleUserAttachment) by
// injecting a fake roaGet func, so it covers the real error-conversion
// branch including the resource.Retry loop and the NeedRetry guard, instead
// of copy-pasting the matching logic into the test.
//
// The fake roaGet signature matches connectivity.AliyunClient.RoaGet so the
// injected func stands in for the real SDK call without a live AliyunClient.
func TestDescribeMaxComputeRoleUserAttachment_ExternallyDeletedMember(t *testing.T) {
	const (
		ramID    = "default_project&role_project_admin&RAM$openapiautomation@test.aliyunid.com:tf-example"
		aliyunID = "default_project&role_project_admin&ALIYUN$openapiautomation@test.aliyunid.com"
	)

	roaGetFunc := func(err error) func(apiProductCode string, apiVersion string, pathName string, query map[string]*string, headers map[string]*string, body interface{}) (map[string]interface{}, error) {
		return func(string, string, string, map[string]*string, map[string]*string, interface{}) (map[string]interface{}, error) {
			return nil, err
		}
	}

	cases := []struct {
		name         string
		id           string
		roaGet       func(string, string, string, map[string]*string, map[string]*string, interface{}) (map[string]interface{}, error)
		wantNotFound bool
		wantErr      bool
	}{
		{
			name: "externally_deleted_ram_member_message_in_message_field",
			id:   ramID,
			roaGet: roaGetFunc(&tea.SDKError{
				StatusCode: tea.Int(400),
				Code:       tea.String("ODPS-0000001"),
				Message:    tea.String("baseId and tenantId do not match"),
			}),
			wantNotFound: true,
			wantErr:      true,
		},
		{
			name: "externally_deleted_aliyun_member_message_in_data_field",
			id:   aliyunID,
			roaGet: roaGetFunc(&tea.SDKError{
				StatusCode: tea.Int(400),
				Code:       tea.String("InternalError.Project.User.BaseIdNotMatch"),
				Message:    tea.String("ODPS request failed"),
				Data:       tea.String(`{"message":"baseId and tenantId do not match"}`),
			}),
			wantNotFound: true,
			wantErr:      true,
		},
		{
			name: "unrelated_auth_error_does_not_match",
			id:   ramID,
			roaGet: roaGetFunc(&tea.SDKError{
				StatusCode: tea.Int(403),
				Code:       tea.String("Forbidden.Access"),
				Message:    tea.String("Access denied"),
			}),
			wantNotFound: false,
			wantErr:      true,
		},
		{
			name: "empty_users_response_is_notfound",
			id:   ramID,
			roaGet: func(string, string, string, map[string]*string, map[string]*string, interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{
					"data": map[string]interface{}{
						"users": []interface{}{},
					},
				}, nil
			},
			wantNotFound: true,
			wantErr:      true,
		},
		{
			name: "matching_aliyun_user_response_succeeds",
			id:   aliyunID,
			roaGet: func(string, string, string, map[string]*string, map[string]*string, interface{}) (map[string]interface{}, error) {
				return map[string]interface{}{
					"data": map[string]interface{}{
						"users": []interface{}{
							map[string]interface{}{
								"name": "ALIYUN$openapiautomation@test.aliyunid.com",
							},
						},
					},
				}, nil
			},
			wantNotFound: false,
			wantErr:      false,
		},
	}

	s := &MaxComputeServiceV2{}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := s.describeMaxComputeRoleUserAttachmentWithRoaGet(c.id, c.roaGet)
			if (err != nil) != c.wantErr {
				t.Fatalf("expected err != nil = %v, got err = %v", c.wantErr, err)
			}
			if NotFoundError(err) != c.wantNotFound {
				t.Fatalf("expected NotFoundError = %v, got %v. err: %v", c.wantNotFound, NotFoundError(err), err)
			}
			if c.wantNotFound && !strings.Contains(err.Error(), c.id) {
				t.Fatalf("expected NotFoundError to contain resource ID %q, got %v", c.id, err)
			}
		})
	}
}
