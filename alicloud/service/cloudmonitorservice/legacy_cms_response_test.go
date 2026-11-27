package cloudmonitorservice

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

// The CMS envelope always carries a top-level Code, so a payload without it
// is a structure mismatch: it must surface as an error instead of decoding
// into an empty list that callers would misread as "resource absent".
func TestUnitLegacyCmsResponse(t *testing.T) {
	resp := &cms.DescribeContactListResponse{}
	raw := map[string]interface{}{
		"Code": "200",
		"Contacts": map[string]interface{}{
			"Contact": []interface{}{
				map[string]interface{}{"Name": "alarm-contact-1"},
			},
		},
	}
	if err := legacyCmsResponse(raw, resp); err != nil {
		t.Fatalf("decoding: %s", err)
	}
	if resp.Code != "200" || len(resp.Contacts.Contact) != 1 || resp.Contacts.Contact[0].Name != "alarm-contact-1" {
		t.Fatalf("typed fields lost after rehydration: %+v", resp)
	}

	if err := legacyCmsResponse(map[string]interface{}{"Contacts": map[string]interface{}{}}, resp); err == nil {
		t.Fatal("a payload without the top-level Code must be a structure error, not a silent empty list")
	}
}
