package alicloud

import (
	"testing"
)

func TestValidateRamPolicyDocumentStructure(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr int
	}{
		{
			name:    "valid single statement array with Action string and Effect Allow",
			input:   `{"Version":"1","Statement":[{"Effect":"Allow","Action":"ecs:Describe*","Resource":"*"}]}`,
			wantErr: 0,
		},
		{
			name:    "valid array with Action slice Deny plus NotAction-only statement",
			input:   `{"Version":"1","Statement":[{"Effect":"Deny","Action":["ecs:Describe*","oss:Get*"],"Resource":"*"},{"Effect":"Allow","NotAction":"ram:*","Resource":"*"}]}`,
			wantErr: 0,
		},
		{
			name:    "valid single statement object (not array)",
			input:   `{"Version":"1","Statement":{"Effect":"Allow","Action":"ecs:Describe*","Resource":"*"}}`,
			wantErr: 0,
		},
		{
			name:    "top-level JSON array is not an object",
			input:   `[]`,
			wantErr: 0,
		},
		{
			name:    "top-level JSON string is not an object",
			input:   `"foo"`,
			wantErr: 0,
		},
		{
			name:    "invalid JSON syntax left to StringIsJSON",
			input:   `{`,
			wantErr: 0,
		},
		{
			name:    "statement with no Effect key but with Action",
			input:   `{"Version":"1","Statement":[{"Action":"ecs:Describe*","Resource":"*"}]}`,
			wantErr: 0,
		},
		{
			name:    "statement missing both Action and NotAction",
			input:   `{"Version":"1","Statement":[{"Effect":"Allow","Resource":"*"}]}`,
			wantErr: 1,
		},
		{
			name:    "Effect typo Alow",
			input:   `{"Version":"1","Statement":[{"Effect":"Alow","Action":"ecs:Describe*","Resource":"*"}]}`,
			wantErr: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, errs := validateRamPolicyDocumentStructure(tc.input, "policy_document")
			if len(errs) != tc.wantErr {
				t.Fatalf("case %q: expected %d error(s), got %d: %v", tc.name, tc.wantErr, len(errs), errs)
			}
		})
	}
}
