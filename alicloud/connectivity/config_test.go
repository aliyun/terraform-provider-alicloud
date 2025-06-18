package connectivity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnitCommonRefreshAuthCredential_ECS_Role(t *testing.T) {
	client := NewTestClient(t)
	err := client.config.RefreshAuthCredential()

	assert.NoError(t, err)
	assert.NotEmpty(t, client.config.AccessKey)
	assert.NotEmpty(t, client.config.SecretKey)
}

func TestUnitCommonRefreshAuthCredential_OIDC(t *testing.T) {
	client := NewTestClient(t)
	config := client.config
	err := config.RefreshAuthCredential()

	assert.NoError(t, err)
	assert.NotEmpty(t, config.AccessKey)
	assert.NotEmpty(t, config.SecretKey)
}

func TestUnitCommonRefreshAuthCredential_AssumeRole(t *testing.T) {
	client := NewTestClient(t)
	config := client.config
	err := config.RefreshAuthCredential()

	assert.NoError(t, err)
}

func TestUnitCommonRefreshAuthCredential_Error(t *testing.T) {
	config := &Config{
		AssumeRoleWithOidc: &AssumeRoleWithOidc{
			RoleARN: "invalid-arn",
		},
	}

	err := config.RefreshAuthCredential()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "refreshing credential failed")
}
