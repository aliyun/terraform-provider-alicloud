package arms

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// A dispatch rule read has to merge two endpoints: DescribeDispatchRule owns the
// rule itself (name, group rules, notify objects/channels) while the notify time
// window only exists on the ListNotificationPolicies detail payload.
func TestUnitBuildArmsDispatchRuleNotifyRulesReadsWindowFromPolicyDetail(t *testing.T) {
	describe := map[string]interface{}{
		"NotifyRules": []interface{}{
			map[string]interface{}{
				"NotifyChannels": []interface{}{"dingTalk", "wechat"},
				"NotifyObjects": []interface{}{
					map[string]interface{}{
						"Name":           "contact",
						"NotifyObjectId": 42,
						"NotifyType":     "CONTACT",
					},
				},
				// ARMS never returns the window here; a value would be stale.
				"NotifyStartTime": "01:00",
				"NotifyEndTime":   "02:00",
			},
		},
	}

	cases := map[string]map[string]interface{}{
		"parsed object": {
			"NotifyRule": map[string]interface{}{
				"NotifyStartTime": "08:30",
				"NotifyEndTime":   "20:15",
			},
		},
		"json string": {
			"NotifyRule": `{"NotifyStartTime":"08:30","NotifyEndTime":"20:15"}`,
		},
	}
	for name, detail := range cases {
		t.Run(name, func(t *testing.T) {
			rules := buildArmsDispatchRuleNotifyRules(describe, detail)
			require.Len(t, rules, 1)
			require.Equal(t, "08:30", rules[0]["notify_start_time"])
			require.Equal(t, "20:15", rules[0]["notify_end_time"])
			require.Equal(t, []string{"dingTalk", "wechat"}, rules[0]["notify_channels"])

			objects, ok := rules[0]["notify_objects"].([]map[string]interface{})
			require.True(t, ok)
			require.Len(t, objects, 1)
			require.Equal(t, "ARMS_CONTACT", objects[0]["notify_type"])
			require.Equal(t, 42, objects[0]["notify_object_id"])
			require.Equal(t, "contact", objects[0]["name"])
		})
	}
}

// A confirmed detail that omits both window fields is the verified all-day case:
// refresh normalizes it to the values a default configuration uses instead of
// leaving empty strings behind.
func TestUnitBuildArmsDispatchRuleNotifyRulesNormalizesConfirmedAllDayWindow(t *testing.T) {
	describe := map[string]interface{}{
		"NotifyRules": []interface{}{
			map[string]interface{}{
				"NotifyChannels": []interface{}{"dingTalk"},
			},
		},
	}
	cases := map[string]map[string]interface{}{
		"absent window": {"NotifyRule": map[string]interface{}{"NotifyChannels": []interface{}{"dingTalk"}}},
		"empty window":  {"NotifyRule": map[string]interface{}{"NotifyStartTime": "", "NotifyEndTime": ""}},
		"null window":   {"NotifyRule": map[string]interface{}{"NotifyStartTime": nil, "NotifyEndTime": nil}},
	}

	for name, detail := range cases {
		t.Run(name, func(t *testing.T) {
			rules := buildArmsDispatchRuleNotifyRules(describe, detail)
			require.Len(t, rules, 1)
			require.Equal(t, "00:00", rules[0]["notify_start_time"])
			require.Equal(t, "23:59", rules[0]["notify_end_time"])
		})
	}
}

// F01 regression: a missing or unparseable detail, or a partially published window,
// is NOT a confirmed all-day response. The builder must return nil so Read keeps the
// existing notify_rules (e.g. a non-default 08:30-20:15 window) instead of silently
// rewriting it to the 00:00-23:59 business default.
func TestUnitBuildArmsDispatchRuleNotifyRulesKeepsWindowWhenDetailUnconfirmed(t *testing.T) {
	describe := map[string]interface{}{
		"NotifyRules": []interface{}{
			map[string]interface{}{
				"NotifyChannels": []interface{}{"dingTalk"},
				"NotifyObjects": []interface{}{
					map[string]interface{}{"Name": "contact", "NotifyObjectId": 42, "NotifyType": "CONTACT"},
				},
			},
		},
	}
	cases := map[string]map[string]interface{}{
		"nil detail":            nil,
		"missing notify rule":   {},
		"null notify rule":      {"NotifyRule": nil},
		"notify rule scalar":    {"NotifyRule": "unexpected"},
		"notify rule number":    {"NotifyRule": 830},
		"only start published":  {"NotifyRule": map[string]interface{}{"NotifyStartTime": "08:30"}},
		"only end published":    {"NotifyRule": map[string]interface{}{"NotifyEndTime": "20:15"}},
		"empty start only":      {"NotifyRule": map[string]interface{}{"NotifyStartTime": "", "NotifyEndTime": "20:15"}},
		"empty end only":        {"NotifyRule": map[string]interface{}{"NotifyStartTime": "08:30", "NotifyEndTime": ""}},
		"non string start":      {"NotifyRule": map[string]interface{}{"NotifyStartTime": 830, "NotifyEndTime": "20:15"}},
		"unparseable json rule": {"NotifyRule": "{not json"},
	}
	for name, detail := range cases {
		t.Run(name, func(t *testing.T) {
			require.Nil(t, buildArmsDispatchRuleNotifyRules(describe, detail),
				"an unconfirmed window must leave notify_rules untouched, not fill business defaults")
		})
	}
}

// A partial DescribeDispatchRule payload must leave notify_rules untouched instead
// of panicking or clearing the attribute, even when the window itself is confirmed.
func TestUnitBuildArmsDispatchRuleNotifyRulesToleratesPartialPayload(t *testing.T) {
	detail := map[string]interface{}{
		"NotifyRule": map[string]interface{}{"NotifyStartTime": "08:30", "NotifyEndTime": "20:15"},
	}
	cases := map[string]map[string]interface{}{
		"missing notify rules":      {},
		"notify rules not a list":   {"NotifyRules": "unexpected"},
		"notify rule not an object": {"NotifyRules": []interface{}{"unexpected"}},
		"empty notify rules":        {"NotifyRules": []interface{}{}},
	}
	for name, describe := range cases {
		t.Run(name, func(t *testing.T) {
			require.Nil(t, buildArmsDispatchRuleNotifyRules(describe, detail))
		})
	}

	t.Run("missing notify objects and channels", func(t *testing.T) {
		rules := buildArmsDispatchRuleNotifyRules(map[string]interface{}{
			"NotifyRules": []interface{}{map[string]interface{}{}},
		}, detail)
		require.Len(t, rules, 1)
		require.Empty(t, rules[0]["notify_objects"])
		require.Empty(t, rules[0]["notify_channels"])
	})

	t.Run("wrongly typed notify objects", func(t *testing.T) {
		rules := buildArmsDispatchRuleNotifyRules(map[string]interface{}{
			"NotifyRules": []interface{}{
				map[string]interface{}{
					"NotifyChannels": "dingTalk",
					"NotifyObjects":  []interface{}{"unexpected", map[string]interface{}{"Name": "kept"}},
				},
			},
		}, detail)
		require.Len(t, rules, 1)
		require.Empty(t, rules[0]["notify_channels"])
		objects := rules[0]["notify_objects"].([]map[string]interface{})
		require.Len(t, objects, 1)
		require.Equal(t, "kept", objects[0]["name"])
	})
}

// The notify template only exists on the policy detail and may be served as a
// parsed object, as a JSON string, or wrapped in an array. Anything else - and an
// empty template - means the rule has no template to publish.
func TestUnitBuildArmsDispatchRuleNotifyTemplateAcceptsPublishedShapes(t *testing.T) {
	template := map[string]interface{}{
		"EmailTitle":          "title",
		"EmailContent":        "content",
		"EmailRecoverTitle":   "recover title",
		"EmailRecoverContent": "recover content",
		"SmsContent":          "sms",
		"SmsRecoverContent":   "sms recover",
		"TtsContent":          "tts",
		"TtsRecoverContent":   "tts recover",
		"RobotContent":        "robot",
	}
	expected := map[string]interface{}{
		"email_title":           "title",
		"email_content":         "content",
		"email_recover_title":   "recover title",
		"email_recover_content": "recover content",
		"sms_content":           "sms",
		"sms_recover_content":   "sms recover",
		"tts_content":           "tts",
		"tts_recover_content":   "tts recover",
		"robot_content":         "robot",
	}

	cases := map[string]struct {
		detail map[string]interface{}
	}{
		"parsed object":     {map[string]interface{}{"NotifyTemplate": template}},
		"json string":       {map[string]interface{}{"NotifyTemplate": `{"EmailTitle":"title","EmailContent":"content","EmailRecoverTitle":"recover title","EmailRecoverContent":"recover content","SmsContent":"sms","SmsRecoverContent":"sms recover","TtsContent":"tts","TtsRecoverContent":"tts recover","RobotContent":"robot"}`}},
		"single item array": {map[string]interface{}{"NotifyTemplate": []interface{}{template}}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, expected, buildArmsDispatchRuleNotifyTemplate(tc.detail))
		})
	}

	t.Run("robot content served as ding content", func(t *testing.T) {
		got := buildArmsDispatchRuleNotifyTemplate(map[string]interface{}{
			"NotifyTemplate": map[string]interface{}{"DingContent": "ding"},
		})
		require.Equal(t, "ding", got["robot_content"])
	})

	for name, detail := range map[string]map[string]interface{}{
		"absent template":  {},
		"null template":    {"NotifyTemplate": nil},
		"empty template":   {"NotifyTemplate": map[string]interface{}{}},
		"empty array":      {"NotifyTemplate": []interface{}{}},
		"unexpected shape": {"NotifyTemplate": "unexpected"},
	} {
		t.Run(name, func(t *testing.T) {
			require.Nil(t, buildArmsDispatchRuleNotifyTemplate(detail))
		})
	}
}

func TestUnitNormalizeArmsNotifyTimeRange(t *testing.T) {
	start, end := normalizeArmsNotifyTimeRange("", "")
	require.Equal(t, "00:00", start)
	require.Equal(t, "23:59", end)

	start, end = normalizeArmsNotifyTimeRange(nil, nil)
	require.Equal(t, "00:00", start)
	require.Equal(t, "23:59", end)

	start, end = normalizeArmsNotifyTimeRange("08:30", "20:15")
	require.Equal(t, "08:30", start)
	require.Equal(t, "20:15", end)
}
