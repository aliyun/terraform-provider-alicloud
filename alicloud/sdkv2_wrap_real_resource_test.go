package alicloud

import (
	"context"
	"reflect"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/features"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/sdkv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var crudFields = map[string]bool{
	"Create": true, "Read": true, "Update": true, "Delete": true,
	"CreateContext": true, "ReadContext": true, "UpdateContext": true, "DeleteContext": true,
	"CreateWithoutTimeout": true, "ReadWithoutTimeout": true,
	"UpdateWithoutTimeout": true, "DeleteWithoutTimeout": true,
}

type recordingInterceptor struct {
	before []intercept.Call
	after  []intercept.Call
}

func (r *recordingInterceptor) Before(_ context.Context, call intercept.Call) error {
	r.before = append(r.before, call)
	return nil
}

func (r *recordingInterceptor) After(_ context.Context, call intercept.Call, err error) error {
	r.after = append(r.after, call)
	return err
}

func TestUnitAliCloudSdkV2WrapPreservesCustomizeDiff(t *testing.T) {
	rec := &recordingInterceptor{}
	wrapped := sdkv2.WrapResource("alicloud_instance", resourceAliCloudInstance(), []intercept.Interceptor{rec})

	state := &terraform.InstanceState{
		ID:         "i-unit-test",
		Attributes: map[string]string{"image_id": "image-before"},
	}
	config := terraform.NewResourceConfigRaw(map[string]interface{}{"image_id": "image-after"})
	client := &connectivity.AliyunClient{
		Features: features.Features{
			EcsInstance: features.EcsInstance{ReplaceOnImageUpdate: true},
		},
	}

	diff, err := wrapped.Diff(context.Background(), state, config, client)
	if err != nil {
		t.Fatalf("diff through the wrapper: %s", err)
	}
	attribute, ok := diff.Attributes["image_id"]
	if !ok {
		t.Fatalf("expected image_id in the diff, got %#v", diff.Attributes)
	}
	if !attribute.RequiresNew {
		t.Error("image_id RequiresNew = false through the wrapper, want true: the CustomizeDiff did not run")
	}

	unwrapped := resourceAliCloudInstance()
	controlState := &terraform.InstanceState{
		ID:         "i-unit-test",
		Attributes: map[string]string{"image_id": "image-before"},
	}
	controlDiff, err := unwrapped.Diff(context.Background(), controlState, config, client)
	if err != nil {
		t.Fatalf("diff without the wrapper: %s", err)
	}
	controlAttribute, ok := controlDiff.Attributes["image_id"]
	if !ok {
		t.Fatalf("expected image_id in the unwrapped diff, got %#v", controlDiff.Attributes)
	}
	if controlAttribute.RequiresNew != attribute.RequiresNew {
		t.Errorf("wrapped RequiresNew = %t, unwrapped = %t: the wrapper changed the plan",
			attribute.RequiresNew, controlAttribute.RequiresNew)
	}

	if len(rec.before) != 0 || len(rec.after) != 0 {
		t.Errorf("the chain ran during Diff: before=%v after=%v; only CRUD is intercepted", rec.before, rec.after)
	}
}

func TestUnitAliCloudSdkV2WrapPreservesStateUpgraders(t *testing.T) {
	rec := &recordingInterceptor{}
	original := resourceAliCloudApiGatewayInstance()
	wrapped := sdkv2.WrapResource("alicloud_api_gateway_instance", original, []intercept.Interceptor{rec})

	if wrapped.SchemaVersion != original.SchemaVersion {
		t.Fatalf("SchemaVersion = %d, want %d", wrapped.SchemaVersion, original.SchemaVersion)
	}
	if len(wrapped.StateUpgraders) != len(original.StateUpgraders) {
		t.Fatalf("StateUpgraders = %d, want %d", len(wrapped.StateUpgraders), len(original.StateUpgraders))
	}
	if len(wrapped.StateUpgraders) == 0 {
		t.Fatal("alicloud_api_gateway_instance no longer declares a StateUpgrader; pick another resource for this test")
	}

	upgrader := wrapped.StateUpgraders[0]
	if upgrader.Version != original.StateUpgraders[0].Version {
		t.Errorf("upgrader Version = %d, want %d", upgrader.Version, original.StateUpgraders[0].Version)
	}
	if !upgrader.Type.Equals(original.StateUpgraders[0].Type) {
		t.Error("upgrader Type changed: the v0 state would no longer decode")
	}

	rawState := map[string]interface{}{
		"to_connect_vpc_ip_block": map[string]interface{}{"cidr_block": "10.0.0.0/24"},
	}
	upgraded, err := upgrader.Upgrade(context.Background(), rawState, nil)
	if err != nil {
		t.Fatalf("upgrade through the wrapper: %s", err)
	}
	block, ok := upgraded["to_connect_vpc_ip_block"].([]interface{})
	if !ok {
		t.Fatalf("to_connect_vpc_ip_block = %#v, want []interface{}: the upgrade did not run",
			upgraded["to_connect_vpc_ip_block"])
	}
	if len(block) != 1 {
		t.Errorf("to_connect_vpc_ip_block has %d element(s), want 1", len(block))
	}

	if len(rec.before) != 0 || len(rec.after) != 0 {
		t.Errorf("the chain ran during the state upgrade: before=%v after=%v; only CRUD is intercepted",
			rec.before, rec.after)
	}
}

func TestUnitAliCloudSdkV2WrapCopiesEveryNonCRUDField(t *testing.T) {
	cases := []struct {
		name    string
		factory func() *schema.Resource
	}{
		{"alicloud_instance", resourceAliCloudInstance},
		{"alicloud_api_gateway_instance", resourceAliCloudApiGatewayInstance},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			original := tc.factory()
			originalCreate := reflect.ValueOf(original.Create).Pointer()

			rec := &recordingInterceptor{}
			wrapped := sdkv2.WrapResource(tc.name, original, []intercept.Interceptor{rec})

			if wrapped == original {
				t.Fatal("wrapper returned the input pointer: the shared map literal value would be mutated")
			}
			if reflect.ValueOf(original.Create).Pointer() != originalCreate {
				t.Error("the original resource was mutated by the wrapper")
			}

			originalValue := reflect.ValueOf(*original)
			wrappedValue := reflect.ValueOf(*wrapped)
			resourceType := originalValue.Type()

			var checked int
			for i := 0; i < resourceType.NumField(); i++ {
				field := resourceType.Field(i)
				originalField, wrappedField := originalValue.Field(i), wrappedValue.Field(i)
				checked++

				if crudFields[field.Name] {
					switch {
					case originalField.IsNil() && !wrappedField.IsNil():
						t.Errorf("%s: the wrapper populated a CRUD form the resource does not use", field.Name)
					case !originalField.IsNil() && wrappedField.IsNil():
						t.Errorf("%s: the wrapper dropped a CRUD form the resource uses", field.Name)
					case !originalField.IsNil() && originalField.Pointer() == wrappedField.Pointer():
						t.Errorf("%s: set CRUD field was not wrapped", field.Name)
					}
					continue
				}

				if field.Type.Kind() == reflect.Func {
					if originalField.Pointer() != wrappedField.Pointer() {
						t.Errorf("%s: non-CRUD func field changed", field.Name)
					}
					continue
				}
				if !reflect.DeepEqual(originalField.Interface(), wrappedField.Interface()) {
					t.Errorf("%s: non-CRUD field changed", field.Name)
				}
			}
			if checked != resourceType.NumField() {
				t.Fatalf("checked %d of %d fields", checked, resourceType.NumField())
			}
		})
	}
}

func TestUnitAliCloudSdkV2WrapRewritesThePlainCRUDForm(t *testing.T) {
	resource := resourceAliCloudInstance()
	if resource.Create == nil || resource.CreateContext != nil || resource.CreateWithoutTimeout != nil {
		t.Fatalf("alicloud_instance no longer uses the plain CRUD form (Create=%t CreateContext=%t CreateWithoutTimeout=%t)",
			resource.Create != nil, resource.CreateContext != nil, resource.CreateWithoutTimeout != nil)
	}

	var inner []string
	stub := func(op string) func(*schema.ResourceData, interface{}) error {
		return func(*schema.ResourceData, interface{}) error {
			inner = append(inner, op)
			return nil
		}
	}
	resource.Create = stub("Create")
	resource.Read = stub("Read")
	resource.Update = stub("Update")
	resource.Delete = stub("Delete")

	rec := &recordingInterceptor{}
	wrapped := sdkv2.WrapResource("alicloud_instance", resource, []intercept.Interceptor{rec})

	client := &connectivity.AliyunClient{}
	data := wrapped.Data(nil)
	for _, call := range []struct {
		op intercept.Op
		fn func(*schema.ResourceData, interface{}) error
	}{
		{intercept.OpCreate, wrapped.Create},
		{intercept.OpRead, wrapped.Read},
		{intercept.OpUpdate, wrapped.Update},
		{intercept.OpDelete, wrapped.Delete},
	} {
		if err := call.fn(data, client); err != nil {
			t.Fatalf("%s through the wrapper: %s", call.op, err)
		}
	}

	want := []intercept.Call{
		{Name: "alicloud_instance", Op: intercept.OpCreate, Meta: client},
		{Name: "alicloud_instance", Op: intercept.OpRead, Meta: client},
		{Name: "alicloud_instance", Op: intercept.OpUpdate, Meta: client},
		{Name: "alicloud_instance", Op: intercept.OpDelete, Meta: client},
	}
	if !reflect.DeepEqual(rec.before, want) {
		t.Errorf("Before calls = %v, want %v", rec.before, want)
	}
	if !reflect.DeepEqual(rec.after, want) {
		t.Errorf("After calls = %v, want %v", rec.after, want)
	}
	if got := []string{"Create", "Read", "Update", "Delete"}; !reflect.DeepEqual(inner, got) {
		t.Errorf("inner functions ran %v, want %v", inner, got)
	}
}
