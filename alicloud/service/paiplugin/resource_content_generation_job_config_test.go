package paiplugin_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// jobConfigObject is the canonical GetContentGenerationJobConfig Data object the
// mock returns for the configured id; it echoes every field the resource maps.
func jobConfigObject(id string) map[string]interface{} {
	return map[string]interface{}{
		"ContentGenerationJobConfigId": id,
		"Name":                         "tf-test-cfg",
		"Purpose":                      "UserTouch",
		"Type":                         "Template",
		"Model":                        "QWen",
		"IndustryCode":                 "100",
		"IndustryName":                 "E-commerce",
		"SignatureId":                  "sig-1",
		"SignatureName":                "tf-sig",
		"RecommendedKeywords":          []interface{}{"a", "b"},
		"CustomKeywords":               []interface{}{"c"},
		"JobCount":                     2,
		"GmtCreateTime":                "2026-01-01 00:00:00",
		"GmtModifiedTime":              "2026-01-02 00:00:00",
	}
}

func TestUnitPaiPluginContentGenerationJobConfig(t *testing.T) {
	const resourceId = "alicloud_pai_plugin_content_generation_job_config.default"

	paiPluginTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/contentgeneration/jobconfigs":
			writeJSON(t, w, http.StatusOK, paiPluginEnvelope(map[string]interface{}{
				"ContentGenerationJobConfigId": "cfg-mock-1",
			}))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/contentgeneration/jobconfigs/cfg-mock-1":
			writeJSON(t, w, http.StatusOK, paiPluginEnvelope(jobConfigObject("cfg-mock-1")))
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v2/contentgeneration/jobconfigs/cfg-mock-1":
			writeJSON(t, w, http.StatusOK, paiPluginEnvelope(map[string]interface{}{}))
		default:
			paiPluginNotFound(w)
		}
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: factories(),
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				// Create + Read: POST returns the new id, then Read re-fetches the
				// full object so the apply result carries the server-side fields.
				Config: paiPluginProviderConfig() + `
resource "alicloud_pai_plugin_content_generation_job_config" "default" {
  name                 = "tf-test-cfg"
  type                 = "Template"
  model                = "QWen"
  industry_code        = "100"
  signature_id         = "sig-1"
  recommended_keywords = ["a", "b"]
  custom_keywords      = ["c"]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "id", "cfg-mock-1"),
					resource.TestCheckResourceAttr(resourceId, "content_generation_job_config_id", "cfg-mock-1"),
					resource.TestCheckResourceAttr(resourceId, "name", "tf-test-cfg"),
					resource.TestCheckResourceAttr(resourceId, "industry_name", "E-commerce"),
					resource.TestCheckResourceAttr(resourceId, "job_count", "2"),
					resource.TestCheckResourceAttr(resourceId, "recommended_keywords.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "recommended_keywords.0", "a"),
					resource.TestCheckResourceAttr(resourceId, "recommended_keywords.1", "b"),
				),
			},
			{
				// No-op update path: every input is ForceNew, so an identical plan
				// is a no-op and Update never reaches the API. The handler returns
				// not-found for any unrecognized call, so a stray update call would
				// fail the test rather than silently pass.
				Config: paiPluginProviderConfig() + `
resource "alicloud_pai_plugin_content_generation_job_config" "default" {
  name                 = "tf-test-cfg"
  type                 = "Template"
  model                = "QWen"
  industry_code        = "100"
  signature_id         = "sig-1"
  recommended_keywords = ["a", "b"]
  custom_keywords      = ["c"]
}
`,
				PlanOnly: true,
			},
			{
				// Import path: id passthrough, then Read re-fetches by id.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: func(*terraform.State) error { return nil },
	})
}

// TestUnitPaiPluginContentGenerationJobConfigCreateFails covers the failure
// path: when CreateContentGenerationJobConfig returns HTTP 400, RoaPost surfaces
// the error and Create aborts the apply, so the step reports a creation error
// rather than writing a partial instance to state.
func TestUnitPaiPluginContentGenerationJobConfigCreateFails(t *testing.T) {
	paiPluginTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Every call returns the not-found/forbidden envelope, including POST.
		paiPluginNotFound(w)
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: factories(),
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: paiPluginProviderConfig() + `
resource "alicloud_pai_plugin_content_generation_job_config" "default" {
  name                 = "tf-test-cfg"
  type                 = "Template"
  model                = "QWen"
  industry_code        = "100"
  signature_id         = "sig-1"
  recommended_keywords = ["a", "b"]
}
`,
				ExpectError: regexCreateFails,
			},
		},
	})
}

// jobObject is the canonical GetContentGenerationJob Data object.
func jobObject(id string) map[string]interface{} {
	return map[string]interface{}{
		"ContentGenerationJobId":         id,
		"ContentGenerationJobConfigId":   "jcfg-1",
		"ContentGenerationJobConfigName": "tf-test-cfg",
		"Status":                         "Success",
		"Results":                        []interface{}{"result-a", "result-b"},
		"Contents":                       []interface{}{"content-1"},
		"SignatureId":                    "sig-1",
		"SignatureName":                  "tf-sig",
		"GmtCreateTime":                  "2026-01-01 00:00:00",
		"GmtModifiedTime":                "2026-01-02 00:00:00",
	}
}

func TestUnitPaiPluginContentGenerationJob(t *testing.T) {
	const resourceId = "alicloud_pai_plugin_content_generation_job.default"

	// putSeen tracks whether UpdateContentGenerationJob (PUT) has run, so the
	// GET Read returns the post-update contents after an update and the initial
	// contents before. The mock follows the repo convention of canned responses
	// keyed on method+path (the request body is not inspected): the create step
	// plans a single content and the update step plans two, so the fixture echoes
	// exactly those, keeping the apply result consistent with the plan at each
	// step and exercising the real create→read and update→read flows.
	var putSeen bool

	paiPluginTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/contentgeneration/jobs":
			writeJSON(t, w, http.StatusOK, paiPluginEnvelope(map[string]interface{}{
				"ContentGenerationJobId": "job-mock-1",
			}))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/contentgeneration/jobs/job-mock-1":
			obj := jobObject("job-mock-1")
			if putSeen {
				obj["Contents"] = []interface{}{"content-1", "content-2"}
			}
			writeJSON(t, w, http.StatusOK, paiPluginEnvelope(obj))
		case r.Method == http.MethodPut && r.URL.Path == "/api/v2/contentgeneration/jobs/job-mock-1":
			putSeen = true
			writeJSON(t, w, http.StatusOK, paiPluginEnvelope(map[string]interface{}{}))
		default:
			paiPluginNotFound(w)
		}
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: factories(),
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: paiPluginProviderConfig() + `
resource "alicloud_pai_plugin_content_generation_job" "default" {
  content_generation_job_config_id = "jcfg-1"
  contents                         = ["content-1"]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "id", "job-mock-1"),
					resource.TestCheckResourceAttr(resourceId, "content_generation_job_id", "job-mock-1"),
					resource.TestCheckResourceAttr(resourceId, "status", "Success"),
					resource.TestCheckResourceAttr(resourceId, "results.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "results.0", "result-a"),
					resource.TestCheckResourceAttr(resourceId, "results.1", "result-b"),
				),
			},
			{
				// In-place update via UpdateContentGenerationJob: Contents changes,
				// the PUT returns success, then Read refreshes.
				Config: paiPluginProviderConfig() + `
resource "alicloud_pai_plugin_content_generation_job" "default" {
  content_generation_job_config_id = "jcfg-1"
  contents                         = ["content-1", "content-2"]
}
`,
				Check: resource.TestCheckResourceAttr(resourceId, "id", "job-mock-1"),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: func(*terraform.State) error {
			// PAIPlugin exposes no DeleteContentGenerationJob: destroy is a no-op
			// that only removes state, so there is nothing to assert server-side.
			return nil
		},
	})
}

// TestUnitPaiPluginContentGenerationJobCreateFails covers the failure path:
// when CreateContentGenerationJob returns HTTP 400 the apply aborts with a
// creation error.
func TestUnitPaiPluginContentGenerationJobCreateFails(t *testing.T) {
	paiPluginTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		paiPluginNotFound(w)
	})
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: factories(),
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: paiPluginProviderConfig() + `
resource "alicloud_pai_plugin_content_generation_job" "default" {
  content_generation_job_config_id = "jcfg-1"
}
`,
				ExpectError: regexCreateFails,
			},
		},
	})
}

// messageListEnvelope builds the ListMessages Data object for a single message,
// keyed by the MessageId query param so findMessage matches the created one.
func messageListEnvelope(messageId string) []byte {
	return paiPluginEnvelope(map[string]interface{}{
		"Messages": []interface{}{
			map[string]interface{}{
				"Id":             messageId,
				"PhoneNumber":    "13800000000",
				"GroupId":        "grp-1",
				"OutId":          "out-1",
				"ScheduleId":     "sch-1",
				"ErrorCode":      "0",
				"Signature":      "tf-sig",
				"SignatureId":    "sig-1",
				"TemplateCode":   "SMS_1",
				"TemplateParams": []interface{}{`{"code":"1234"}`},
				"Status":         1,
				"TemplateType":   0,
			},
		},
	})
}

func TestUnitPaiPluginMessage(t *testing.T) {
	const resourceId = "alicloud_pai_plugin_message.default"

	paiPluginTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/messages":
			// SendMessage returns the first message id in Data.Messages[0].Id.
			writeJSON(t, w, http.StatusOK, paiPluginEnvelope(map[string]interface{}{
				"Messages": []interface{}{
					map[string]interface{}{"Id": "msg-mock-1"},
				},
			}))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/messages":
			// ListMessages is queried with Datetime + MessageId; echo the matching
			// message object keyed by the MessageId query value, defaulting to the
			// created id when the query omits it.
			mid := r.URL.Query().Get("MessageId")
			if mid == "" {
				mid = "msg-mock-1"
			}
			writeJSON(t, w, http.StatusOK, messageListEnvelope(mid))
		default:
			paiPluginNotFound(w)
		}
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: factories(),
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: paiPluginProviderConfig() + `
resource "alicloud_pai_plugin_message" "default" {
  phone_number  = "13800000000"
  signature     = "tf-sig"
  template_code = "SMS_1"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "id", "msg-mock-1"),
					resource.TestCheckResourceAttr(resourceId, "message_id", "msg-mock-1"),
					resource.TestCheckResourceAttrSet(resourceId, "send_date"),
					resource.TestCheckResourceAttr(resourceId, "status", "1"),
					resource.TestCheckResourceAttr(resourceId, "group_id", "grp-1"),
				),
			},
			{
				// No-op update path: every input is ForceNew, so an unchanged
				// config is a no-op plan and Update/SendMessage is not recalled.
				Config: paiPluginProviderConfig() + `
resource "alicloud_pai_plugin_message" "default" {
  phone_number  = "13800000000"
  signature     = "tf-sig"
  template_code = "SMS_1"
}
`,
				PlanOnly: true,
			},
		},
		CheckDestroy: func(*terraform.State) error {
			// A sent message cannot be unsent: Delete is a no-op removing state.
			return nil
		},
	})
}

// TestUnitPaiPluginMessageCreateFails covers the failure path: when SendMessage
// returns HTTP 400 the apply aborts with a creation error.
func TestUnitPaiPluginMessageCreateFails(t *testing.T) {
	paiPluginTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		paiPluginNotFound(w)
	})
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV5ProviderFactories: factories(),
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: paiPluginProviderConfig() + `
resource "alicloud_pai_plugin_message" "default" {
  phone_number = "13800000000"
}
`,
				ExpectError: regexCreateFails,
			},
		},
	})
}
