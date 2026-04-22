package connectivity

import (
	"sync"
	"testing"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/stretchr/testify/assert"
)

// TestApplyOpenapiSignVersion verifies that the sign_version configured on
// the provider is actually wired into the openapi client for each supported
// product (notably sls, which previously only had the oss branch wired up and
// therefore returned SignatureVersionNotSupported against SLS private
// endpoints that only accept v4).
func TestApplyOpenapiSignVersion(t *testing.T) {
	t.Run("sls v4 is applied", func(t *testing.T) {
		signVersion := &sync.Map{}
		signVersion.Store("sls", "v4")

		c := &openapi.Client{}
		applyOpenapiSignVersion(c, signVersion, "sls")

		if assert.NotNil(t, c.SignatureVersion, "sls SignatureVersion must be set when sign_version.sls is configured") {
			assert.Equal(t, "v4", tea.StringValue(c.SignatureVersion))
		}
	})

	t.Run("oss v4 is applied", func(t *testing.T) {
		signVersion := &sync.Map{}
		signVersion.Store("oss", "v4")

		c := &openapi.Client{}
		applyOpenapiSignVersion(c, signVersion, "oss")

		if assert.NotNil(t, c.SignatureVersion) {
			assert.Equal(t, "v4", tea.StringValue(c.SignatureVersion))
		}
	})

	t.Run("key missing is a no-op", func(t *testing.T) {
		signVersion := &sync.Map{}
		signVersion.Store("oss", "v4")

		c := &openapi.Client{}
		applyOpenapiSignVersion(c, signVersion, "sls")

		assert.Nil(t, c.SignatureVersion, "SignatureVersion must remain unset when key is absent")
	})

	t.Run("nil sync map is a no-op", func(t *testing.T) {
		c := &openapi.Client{}
		assert.NotPanics(t, func() {
			applyOpenapiSignVersion(c, nil, "sls")
		})
		assert.Nil(t, c.SignatureVersion)
	})

	t.Run("nil openapi client is a no-op", func(t *testing.T) {
		signVersion := &sync.Map{}
		signVersion.Store("sls", "v4")

		assert.NotPanics(t, func() {
			applyOpenapiSignVersion(nil, signVersion, "sls")
		})
	})

	t.Run("non-string value is a no-op", func(t *testing.T) {
		signVersion := &sync.Map{}
		signVersion.Store("sls", 4)

		c := &openapi.Client{}
		assert.NotPanics(t, func() {
			applyOpenapiSignVersion(c, signVersion, "sls")
		})
		assert.Nil(t, c.SignatureVersion)
	})
}

// TestApplyLogClientSignVersion mirrors TestApplyOpenapiSignVersion for the
// legacy aliyun-log-go-sdk path used by WithLogClient. The helper only
// writes c.AuthVersion from the signVersion map; it does not inspect the
// endpoint or infer anything from other client state. These tests pin that
// contract.
func TestApplyLogClientSignVersion(t *testing.T) {
	t.Run("sls v4 is applied", func(t *testing.T) {
		signVersion := &sync.Map{}
		signVersion.Store("sls", "v4")

		c := &sls.Client{}
		applyLogClientSignVersion(c, signVersion, "sls")

		assert.Equal(t, sls.AuthV4, c.AuthVersion)
	})

	t.Run("sls v1 is applied", func(t *testing.T) {
		signVersion := &sync.Map{}
		signVersion.Store("sls", "v1")

		c := &sls.Client{}
		applyLogClientSignVersion(c, signVersion, "sls")

		assert.Equal(t, sls.AuthV1, c.AuthVersion)
	})

	t.Run("key missing is a no-op", func(t *testing.T) {
		signVersion := &sync.Map{}
		signVersion.Store("oss", "v4")

		c := &sls.Client{}
		applyLogClientSignVersion(c, signVersion, "sls")

		assert.Empty(t, string(c.AuthVersion), "AuthVersion must remain unset when key is absent")
	})

	t.Run("nil sync map is a no-op", func(t *testing.T) {
		c := &sls.Client{}
		assert.NotPanics(t, func() {
			applyLogClientSignVersion(c, nil, "sls")
		})
		assert.Empty(t, string(c.AuthVersion))
	})

	t.Run("nil sls client is a no-op", func(t *testing.T) {
		signVersion := &sync.Map{}
		signVersion.Store("sls", "v4")

		assert.NotPanics(t, func() {
			applyLogClientSignVersion(nil, signVersion, "sls")
		})
	})

	t.Run("non-string value is a no-op", func(t *testing.T) {
		signVersion := &sync.Map{}
		signVersion.Store("sls", 4)

		c := &sls.Client{}
		assert.NotPanics(t, func() {
			applyLogClientSignVersion(c, signVersion, "sls")
		})
		assert.Empty(t, string(c.AuthVersion))
	})
}
