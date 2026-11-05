package ecs

// White-box unit tests for alicloud_instance, split out of the external
// ecs_test file. These tests call the resourceAliCloudInstance factory and
// unexported CRUD funcs (resourceAliCloudInstance Create/Read), so they must
// live in the internal package ecs test binary (mirroring disk_unit_test.go).
// Provider().ResourcesMap is replaced with unitTestProviderResources to avoid
// importing acctest (import cycle: alicloud -> service/ecs).

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/features"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestUnitAliCloudInstanceImageUpdateForcesNew exercises the diff decision without the API.
func TestUnitAliCloudInstanceImageUpdateForcesNew(t *testing.T) {
	testCases := []struct {
		name            string
		replaceEnabled  bool
		stateImageId    string
		configImageId   string
		wantRequiresNew bool
	}{
		{
			name:            "enabled forces a replacement",
			replaceEnabled:  true,
			stateImageId:    "image-before",
			configImageId:   "image-after",
			wantRequiresNew: true,
		},
		{
			name:            "disabled keeps the update in place",
			replaceEnabled:  false,
			stateImageId:    "image-before",
			configImageId:   "image-after",
			wantRequiresNew: false,
		},
		{
			// An instance created from a launch template has no image_id of its own, so there is no
			// user chosen image being replaced.
			name:            "an image_id that was empty in state is left alone",
			replaceEnabled:  true,
			stateImageId:    "",
			configImageId:   "image-after",
			wantRequiresNew: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// A fresh resource per case: ForceNew is recorded on the schema, so cases would
			// otherwise leak into each other through a shared schema map.
			instanceResource := resourceAliCloudInstance()
			state := &sdktf.InstanceState{
				ID:         "i-unit-test",
				Attributes: map[string]string{"image_id": testCase.stateImageId},
			}
			config := sdktf.NewResourceConfigRaw(map[string]interface{}{
				"image_id": testCase.configImageId,
			})
			client := &connectivity.AliyunClient{
				Features: features.Features{
					EcsInstance: features.EcsInstance{ReplaceOnImageUpdate: testCase.replaceEnabled},
				},
			}

			diff, err := instanceResource.Diff(context.Background(), state, config, client)
			if err != nil {
				t.Fatalf("diff: %s", err)
			}
			attribute, ok := diff.Attributes["image_id"]
			if !ok {
				t.Fatalf("expected image_id to be part of the diff, got %#v", diff.Attributes)
			}
			if attribute.RequiresNew != testCase.wantRequiresNew {
				t.Fatalf("image_id RequiresNew = %t, want %t", attribute.RequiresNew, testCase.wantRequiresNew)
			}
		})
	}
}

const unitTestMockInstanceID = "i-mock-managed-host-1"
const unitTestMockManagedHostID = "mh-mock-1234567890"

// The MANAGED_HOST_ID query attribute is a plain passthrough value: the provider
// requests it only when the user configures it in additional_attributes. The
// trigger value is MANAGED_HOST_ID per the ManagedHost OpenAPI change doc; the
// PascalCase ManagedHostId is only the response field inside AdditionalInfo.
const unitTestManagedHostIdAdditionalAttribute = "MANAGED_HOST_ID"

func TestUnitAliCloudInstanceManagedHostId(t *testing.T) {
	var mu sync.Mutex
	// instanceOnManagedHost simulates the cloud-side state: DescribeInstances returns
	// AdditionalInfo.ManagedHostId only for instances created on a managed host and only
	// when the query requests the ManagedHostId attribute.
	var instanceOnManagedHost bool
	type capturedRequest struct {
		action          string
		managedHostID   string
		additionalAttrs []string
	}
	var requests []capturedRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request form failed: %v", err)
		}
		action := r.FormValue("Action")
		managedHostID := r.FormValue("ManagedHostId")
		var additionalAttrs []string
		if vs := r.Form["AdditionalAttributes"]; len(vs) > 0 {
			additionalAttrs = append(additionalAttrs, vs...)
		}
		for i := 1; ; i++ {
			v := r.FormValue(fmt.Sprintf("AdditionalAttributes.%d", i))
			if v == "" {
				break
			}
			additionalAttrs = append(additionalAttrs, v)
		}
		mu.Lock()
		requests = append(requests, capturedRequest{action: action, managedHostID: managedHostID, additionalAttrs: additionalAttrs})
		if action == "RunInstances" && managedHostID != "" {
			instanceOnManagedHost = true
		}
		onManagedHost := instanceOnManagedHost
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		switch action {
		case "RunInstances":
			fmt.Fprintf(w, `{"RequestId":"mockrid","InstanceIdSets":{"InstanceIdSet":["%s"]}}`, unitTestMockInstanceID)
		case "DescribeInstances":
			var additionalInfo string
			hasTrigger := false
			for _, attr := range additionalAttrs {
				if attr == unitTestManagedHostIdAdditionalAttribute {
					hasTrigger = true
					break
				}
			}
			if onManagedHost && hasTrigger {
				additionalInfo = fmt.Sprintf(`,"AdditionalInfo":{"ManagedHostId":"%s"}`, unitTestMockManagedHostID)
			}
			fmt.Fprintf(w, `{"RequestId":"mockrid","TotalCount":1,"PageNumber":1,"PageSize":10,"Instances":{"Instance":[{"InstanceId":"%s","Status":"Running","InstanceName":"%s","ImageId":"aliyun-2193317896-mock","InstanceType":"ecs.n4.small","ZoneId":"cn-hangzhou-i","RegionId":"cn-hangzhou","InstanceChargeType":"PostPaid","SpotInterruptionBehavior":""%s}]}}`, unitTestMockInstanceID, unitTestMockInstanceID, additionalInfo)
		case "DescribeDisks":
			fmt.Fprintf(w, `{"RequestId":"mockrid","TotalCount":1,"Disks":{"Disk":[{"DiskId":"d-mock-system-disk","InstanceId":"%s","Category":"cloud_efficiency","DiskName":"mock-system-disk","Size":40,"Type":"system","Status":"In_use"}]}}`, unitTestMockInstanceID)
		case "DescribeInstanceMaintenanceAttributes":
			fmt.Fprintf(w, `{"RequestId":"mockrid","MaintenanceAttributes":{"MaintenanceAttribute":[{"InstanceId":"%s","NotifyOnMaintenance":false}]}}`, unitTestMockInstanceID)
		case "DescribeInstanceAttribute":
			fmt.Fprintf(w, `{"RequestId":"mockrid","InstanceId":"%s"}`, unitTestMockInstanceID)
		case "DescribeInstanceAttachmentAttributes":
			fmt.Fprintf(w, `{"RequestId":"mockrid","Instances":{"Instance":[{"InstanceId":"%s","PrivatePoolOptionsMatchCriteria":"","PrivatePoolOptionsId":""}]}}`, unitTestMockInstanceID)
		default:
			fmt.Fprint(w, `{"RequestId":"mockrid"}`)
		}
	}))
	defer server.Close()

	t.Setenv("NO_PROXY", "127.0.0.1")
	t.Setenv("no_proxy", "127.0.0.1")

	credential, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").
		SetAccessKeyId("test-access-key-id").
		SetAccessKeySecret("test-access-key-secret"))
	if err != nil {
		t.Fatalf("build credential failed: %v", err)
	}
	endpoints := new(sync.Map)
	config := connectivity.Config{
		AccessKey:            "test-access-key-id",
		SecretKey:            "test-access-key-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "http",
		SkipRegionValidation: true,
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
	}
	rawClient, err := config.Client()
	if err != nil {
		t.Fatalf("build client failed: %v", err)
	}
	// Both the tea RPC layer and the legacy SDK prepend the protocol scheme themselves,
	// so the endpoint is stored without a scheme.
	endpoints.Store("ecs", strings.TrimPrefix(server.URL, "http://"))

	resources := unitTestProviderResources()

	createInstance := func(t *testing.T, managedHostID string, additionalAttributes []interface{}) *schema.ResourceData {
		t.Helper()
		raw := map[string]interface{}{"instance_name": "tf-mock-managed-host-instance"}
		if managedHostID != "" {
			raw["managed_host_id"] = managedHostID
		}
		if len(additionalAttributes) > 0 {
			raw["additional_attributes"] = additionalAttributes
		}
		d := schema.TestResourceDataRaw(t, resources["alicloud_instance"].Schema, raw)
		d.MarkNewResource()
		if err := resourceAliCloudInstanceCreate(d, rawClient); err != nil {
			t.Fatalf("create instance failed: %v", err)
		}
		return d
	}

	t.Run("SchemaManagedHostId", func(t *testing.T) {
		s, ok := resources["alicloud_instance"].Schema["managed_host_id"]
		if !ok {
			t.Fatal("alicloud_instance schema does not contain managed_host_id")
		}
		if s.Type != schema.TypeString {
			t.Errorf("managed_host_id type = %v, want TypeString", s.Type)
		}
		if !s.Optional {
			t.Error("managed_host_id should be Optional")
		}
		if !s.ForceNew {
			t.Error("managed_host_id should be ForceNew (ECS does not support changing the managed host after creation)")
		}
		if s.Computed {
			t.Error("managed_host_id should not be Computed")
		}
		if len(s.ConflictsWith) != 1 || s.ConflictsWith[0] != "dedicated_host_id" {
			t.Errorf("managed_host_id ConflictsWith = %v, want [dedicated_host_id]", s.ConflictsWith)
		}
	})

	t.Run("CreateWithManagedHostId", func(t *testing.T) {
		mu.Lock()
		requests = nil
		instanceOnManagedHost = false
		mu.Unlock()

		// additional_attributes carries the ManagedHostId query attribute, so
		// Read refreshes managed_host_id from AdditionalInfo.
		d := createInstance(t, unitTestMockManagedHostID, []interface{}{unitTestManagedHostIdAdditionalAttribute})
		if d.Id() != unitTestMockInstanceID {
			t.Fatalf("instance id = %s, want %s", d.Id(), unitTestMockInstanceID)
		}

		mu.Lock()
		runInstancesCount := 0
		describeWithTrigger := 0
		describeWithoutDefaults := 0
		for _, req := range requests {
			switch req.action {
			case "RunInstances":
				runInstancesCount++
				if req.managedHostID != unitTestMockManagedHostID {
					t.Errorf("RunInstances ManagedHostId = %q, want %q", req.managedHostID, unitTestMockManagedHostID)
				}
			case "DescribeInstances":
				hasTrigger := false
				hasDefault := false
				for _, attr := range req.additionalAttrs {
					if attr == unitTestManagedHostIdAdditionalAttribute {
						hasTrigger = true
					}
					if attr == "META_OPTIONS" {
						hasDefault = true
					}
				}
				if hasTrigger {
					describeWithTrigger++
				}
				// Configured values merge into the defaults instead of
				// replacing them, so every query still carries them.
				if !hasDefault {
					describeWithoutDefaults++
				}
			}
		}
		mu.Unlock()
		if runInstancesCount < 1 {
			t.Error("RunInstances was not called")
		}
		if describeWithTrigger < 1 {
			t.Error("no DescribeInstances query requested the ManagedHostId attribute")
		}
		if describeWithoutDefaults > 0 {
			t.Errorf("%d DescribeInstances queries dropped the provider default attributes", describeWithoutDefaults)
		}
		if v := d.Get("managed_host_id").(string); v != unitTestMockManagedHostID {
			t.Errorf("state managed_host_id = %q, want %q", v, unitTestMockManagedHostID)
		}

		// Read must refresh the attribute from the API response instead of
		// merely retaining the configured value.
		if err := d.Set("managed_host_id", "mh-stale"); err != nil {
			t.Fatalf("set managed_host_id failed: %v", err)
		}
		if err := resourceAliCloudInstanceRead(d, rawClient); err != nil {
			t.Fatalf("read instance failed: %v", err)
		}
		if v := d.Get("managed_host_id").(string); v != unitTestMockManagedHostID {
			t.Errorf("read managed_host_id = %q, want %q refreshed from AdditionalInfo", v, unitTestMockManagedHostID)
		}
	})

	t.Run("CreateWithManagedHostIdWithoutAdditionalAttributes", func(t *testing.T) {
		mu.Lock()
		requests = nil
		instanceOnManagedHost = false
		mu.Unlock()

		// The ManagedHostId query attribute is a pure passthrough: without it
		// in additional_attributes the provider must not request it.
		d := createInstance(t, unitTestMockManagedHostID, nil)
		if d.Id() != unitTestMockInstanceID {
			t.Fatalf("instance id = %s, want %s", d.Id(), unitTestMockInstanceID)
		}

		mu.Lock()
		describeWithTrigger := 0
		for _, req := range requests {
			if req.action != "DescribeInstances" {
				continue
			}
			for _, attr := range req.additionalAttrs {
				if attr == unitTestManagedHostIdAdditionalAttribute {
					describeWithTrigger++
					break
				}
			}
		}
		mu.Unlock()
		if describeWithTrigger > 0 {
			t.Errorf("%d DescribeInstances queries requested the ManagedHostId attribute although it was not configured", describeWithTrigger)
		}
		if v := d.Get("managed_host_id").(string); v != unitTestMockManagedHostID {
			t.Errorf("state managed_host_id = %q, want %q retained from config", v, unitTestMockManagedHostID)
		}

		// Read leaves the state value untouched: the response has no
		// AdditionalInfo because the attribute was not requested.
		if err := d.Set("managed_host_id", "mh-stale"); err != nil {
			t.Fatalf("set managed_host_id failed: %v", err)
		}
		if err := resourceAliCloudInstanceRead(d, rawClient); err != nil {
			t.Fatalf("read instance failed: %v", err)
		}
		if v := d.Get("managed_host_id").(string); v != "mh-stale" {
			t.Errorf("read managed_host_id = %q, want %q left untouched", v, "mh-stale")
		}
	})

	t.Run("CreateWithoutManagedHostId", func(t *testing.T) {
		mu.Lock()
		requests = nil
		instanceOnManagedHost = false
		mu.Unlock()

		d := createInstance(t, "", nil)
		if d.Id() != unitTestMockInstanceID {
			t.Fatalf("instance id = %s, want %s", d.Id(), unitTestMockInstanceID)
		}

		mu.Lock()
		runInstancesCount := 0
		describeWithTrigger := 0
		for _, req := range requests {
			switch req.action {
			case "RunInstances":
				runInstancesCount++
				if req.managedHostID != "" {
					t.Errorf("RunInstances should not send ManagedHostId, got %q", req.managedHostID)
				}
			case "DescribeInstances":
				for _, attr := range req.additionalAttrs {
					if attr == unitTestManagedHostIdAdditionalAttribute {
						describeWithTrigger++
						break
					}
				}
			}
		}
		mu.Unlock()
		if runInstancesCount < 1 {
			t.Error("RunInstances was not called")
		}
		if describeWithTrigger > 0 {
			t.Errorf("%d DescribeInstances queries requested the ManagedHostId attribute although it was not configured", describeWithTrigger)
		}
		if v := d.Get("managed_host_id").(string); v != "" {
			t.Errorf("state managed_host_id = %q, want empty", v)
		}
	})
}
