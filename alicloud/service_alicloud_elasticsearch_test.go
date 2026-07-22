package alicloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnitElasticsearchJsonToMap verifies jsonToMap returns an error (instead of
// terminating the whole provider process via log.Fatalf/os.Exit) when the content
// is empty or not valid JSON.
func TestUnitElasticsearchJsonToMap(t *testing.T) {
	t.Run("valid json", func(t *testing.T) {
		m, err := jsonToMap(`{"Result":[{"pvlId":"pvl-123"}]}`)
		assert.NoError(t, err)
		assert.NotNil(t, m)
		assert.Contains(t, m, "Result")
	})

	t.Run("empty content", func(t *testing.T) {
		m, err := jsonToMap("")
		assert.Error(t, err)
		assert.Nil(t, m)
	})

	t.Run("illegal json", func(t *testing.T) {
		m, err := jsonToMap("not-a-json")
		assert.Error(t, err)
		assert.Nil(t, m)
	})
}

// TestUnitElasticsearchParseKibanaPvlId verifies parseKibanaPvlId returns an error
// (instead of panicking on an unchecked type assertion) for every malformed shape
// of the ListKibanaPvlNetwork response Result.
func TestUnitElasticsearchParseKibanaPvlId(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		result := []interface{}{
			map[string]interface{}{"pvlId": "pvl-123"},
		}
		pvlId, err := parseKibanaPvlId(result)
		assert.NoError(t, err)
		assert.Equal(t, "pvl-123", pvlId)
	})

	t.Run("result is nil", func(t *testing.T) {
		pvlId, err := parseKibanaPvlId(nil)
		assert.Error(t, err)
		assert.Empty(t, pvlId)
	})

	t.Run("result is not a list", func(t *testing.T) {
		pvlId, err := parseKibanaPvlId(map[string]interface{}{"pvlId": "pvl-123"})
		assert.Error(t, err)
		assert.Empty(t, pvlId)
	})

	t.Run("empty list", func(t *testing.T) {
		pvlId, err := parseKibanaPvlId([]interface{}{})
		assert.Error(t, err)
		assert.Empty(t, pvlId)
	})

	t.Run("item is not an object", func(t *testing.T) {
		pvlId, err := parseKibanaPvlId([]interface{}{"pvl-123"})
		assert.Error(t, err)
		assert.Empty(t, pvlId)
	})

	t.Run("pvlId missing", func(t *testing.T) {
		result := []interface{}{
			map[string]interface{}{"otherField": "x"},
		}
		pvlId, err := parseKibanaPvlId(result)
		assert.Error(t, err)
		assert.Empty(t, pvlId)
	})

	t.Run("pvlId not a string", func(t *testing.T) {
		result := []interface{}{
			map[string]interface{}{"pvlId": 123},
		}
		pvlId, err := parseKibanaPvlId(result)
		assert.Error(t, err)
		assert.Empty(t, pvlId)
	})
}
