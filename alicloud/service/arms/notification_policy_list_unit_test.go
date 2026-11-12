package arms

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/stretchr/testify/require"
)

// notificationPolicyPage builds one ListNotificationPolicies response page.
func notificationPolicyPage(firstID, count, total int) map[string]interface{} {
	policies := make([]interface{}, 0, count)
	for i := 0; i < count; i++ {
		policies = append(policies, map[string]interface{}{
			"Id":   firstID + i,
			"Name": fmt.Sprintf("rule-%d", firstID+i),
		})
	}
	return map[string]interface{}{
		"PageBean": map[string]interface{}{
			"NotificationPolicies": policies,
			"Total":                total,
		},
	}
}

// An account with more rules than a single page must still list every rule: this
// is the regression that made refresh report existing dispatch rules as missing.
func TestUnitListArmsNotificationPoliciesCollectsEveryPage(t *testing.T) {
	fetches := 0
	fetch := func(page int) (map[string]interface{}, error) {
		fetches++
		switch page {
		case 1:
			return notificationPolicyPage(1, helper.PageSizeXLarge, 205), nil
		case 2:
			return notificationPolicyPage(101, helper.PageSizeXLarge, 205), nil
		case 3:
			return notificationPolicyPage(201, 5, 205), nil
		}
		t.Fatalf("unexpected page %d", page)
		return nil, nil
	}

	policies, err := listArmsNotificationPolicies(fetch)
	require.NoError(t, err)
	require.Len(t, policies, 205)
	require.Equal(t, 3, fetches)
	require.Equal(t, 1, policies[0]["Id"])
	require.Equal(t, 205, policies[204]["Id"])
}

func TestUnitListArmsNotificationPoliciesStopsOnShortPage(t *testing.T) {
	fetches := 0
	fetch := func(int) (map[string]interface{}, error) {
		fetches++
		return notificationPolicyPage(1, 7, 7), nil
	}

	policies, err := listArmsNotificationPolicies(fetch)
	require.NoError(t, err)
	require.Len(t, policies, 7)
	require.Equal(t, 1, fetches)
}

// ARMS is observed to under-report PageBean.Total while a rule is still being
// indexed: a full page whose total already looks complete must not end the walk,
// or the last rule stays invisible and refresh reports it as missing.
func TestUnitListArmsNotificationPoliciesPagesPastStaleTotal(t *testing.T) {
	fetches := 0
	fetch := func(page int) (map[string]interface{}, error) {
		fetches++
		if page == 1 {
			return notificationPolicyPage(1, helper.PageSizeXLarge, helper.PageSizeXLarge), nil
		}
		return notificationPolicyPage(101, 1, helper.PageSizeXLarge), nil
	}

	policies, err := listArmsNotificationPolicies(fetch)
	require.NoError(t, err)
	require.Len(t, policies, helper.PageSizeXLarge+1)
	require.Equal(t, 2, fetches)
}

// Without a usable total the short-page rule still has to end the walk, and an
// empty page ends it too.
func TestUnitListArmsNotificationPoliciesStopsWithoutUsableTotal(t *testing.T) {
	t.Run("unusable total", func(t *testing.T) {
		fetches := 0
		fetch := func(page int) (map[string]interface{}, error) {
			fetches++
			response := notificationPolicyPage(1+(page-1)*helper.PageSizeXLarge, helper.PageSizeXLarge, 0)
			pageBean := response["PageBean"].(map[string]interface{})
			pageBean["Total"] = "not-a-number"
			if page == 2 {
				pageBean["NotificationPolicies"] = notificationPolicyPage(101, 3, 0)["PageBean"].(map[string]interface{})["NotificationPolicies"]
			}
			return response, nil
		}

		policies, err := listArmsNotificationPolicies(fetch)
		require.NoError(t, err)
		require.Len(t, policies, helper.PageSizeXLarge+3)
		require.Equal(t, 2, fetches)
	})

	t.Run("empty page", func(t *testing.T) {
		fetches := 0
		fetch := func(int) (map[string]interface{}, error) {
			fetches++
			return notificationPolicyPage(1, 0, 0), nil
		}

		policies, err := listArmsNotificationPolicies(fetch)
		require.NoError(t, err)
		require.Empty(t, policies)
		require.Equal(t, 1, fetches)
	})
}

func TestUnitListArmsNotificationPoliciesRejectsMalformedPage(t *testing.T) {
	cases := map[string]map[string]interface{}{
		"missing container":   {"PageBean": map[string]interface{}{}},
		"container is object": {"PageBean": map[string]interface{}{"NotificationPolicies": map[string]interface{}{"Id": 1}}},
		"element is scalar":   {"PageBean": map[string]interface{}{"NotificationPolicies": []interface{}{"unexpected"}}},
	}
	for name, response := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := listArmsNotificationPolicies(func(int) (map[string]interface{}, error) {
				return response, nil
			})
			require.Error(t, err)
		})
	}
}

func TestUnitListArmsNotificationPoliciesPropagatesFetchError(t *testing.T) {
	_, err := listArmsNotificationPolicies(func(int) (map[string]interface{}, error) {
		return nil, fmt.Errorf("transport failure")
	})
	require.EqualError(t, err, "transport failure")
}
