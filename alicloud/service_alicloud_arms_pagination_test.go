package alicloud

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListArmsNotificationPoliciesPaginates(t *testing.T) {
	pageItems := func(start, count int) []interface{} {
		items := make([]interface{}, count)
		for i := range items {
			items[i] = map[string]interface{}{"Id": start + i}
		}
		return items
	}
	page := 0
	items, err := listArmsNotificationPolicies("cn-hangzhou", func(_ map[string]interface{}) (map[string]interface{}, error) {
		page++
		if page == 1 {
			return map[string]interface{}{"PageBean": map[string]interface{}{"NotificationPolicies": pageItems(0, 100), "Total": 250}}, nil
		}
		return map[string]interface{}{"PageBean": map[string]interface{}{"NotificationPolicies": pageItems(100, 50), "Total": 250}}, nil
	})
	require.NoError(t, err)
	require.Len(t, items, 150)
	require.Equal(t, 2, page)
}

func TestListArmsNotificationPoliciesRejectsMalformedItems(t *testing.T) {
	_, err := listArmsNotificationPolicies("cn-hangzhou", func(_ map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{"PageBean": map[string]interface{}{"NotificationPolicies": []interface{}{"bad"}}}, nil
	})
	require.Error(t, err)
}
