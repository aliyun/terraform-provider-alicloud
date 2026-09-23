package alicloud

import (
	"testing"
)

func TestTagsMapEqual(t *testing.T) {
	lmap := map[string]interface{}{
		"Created": "TF6666",
		"For":     "acceptance test",
	}
	rmap := map[string]string{
		"Created": "TF6666",
		"For":     "acceptance test",
	}
	if !tagsMapEqual(lmap, rmap) {
		t.Fatal("Tag maps is not equal.")
	}

	rmap = map[string]string{
		"Created": "TF6666",
	}
	if tagsMapEqual(lmap, rmap) {
		t.Fatal("Tag maps is equal.")
	}

	rmap = map[string]string{
		"Created": "TF_Fake",
		"For":     "acceptance test",
	}
	if tagsMapEqual(lmap, rmap) {
		t.Fatal("Tag maps is equal.")
	}

	rmap = map[string]string{
		"Fake": "TF6666",
		"For":  "acceptance test",
	}
	if tagsMapEqual(lmap, rmap) {
		t.Fatal("Tag maps is equal.")
	}

	lmap = map[string]interface{}{
		"Created": 6666,
		"For":     "acceptance test",
	}
	if tagsMapEqual(lmap, rmap) {
		t.Fatal("Tag maps is equal.")
	}
}
