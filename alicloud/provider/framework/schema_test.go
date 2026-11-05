package framework

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var schemaCmpOptions = []cmp.Option{
	cmpopts.SortSlices(func(i, j *tfprotov5.SchemaAttribute) bool { return i.Name < j.Name }),
	cmpopts.SortSlices(func(i, j *tfprotov5.SchemaNestedBlock) bool { return i.TypeName < j.TypeName }),
	cmpopts.IgnoreFields(tfprotov5.SchemaNestedBlock{}, "MinItems", "MaxItems"),
}

func TestUnitProviderSchema(t *testing.T) {
	testCases := map[string]map[string]*sdkschema.Schema{
		"primitives": {
			"a_string": {Type: sdkschema.TypeString, Optional: true},
			"an_int":   {Type: sdkschema.TypeInt, Optional: true},
			"a_bool":   {Type: sdkschema.TypeBool, Optional: true},
		},
		"required": {
			"a_string": {Type: sdkschema.TypeString, Required: true},
		},
		"sensitive": {
			"a_string": {Type: sdkschema.TypeString, Optional: true, Sensitive: true},
		},
		"described": {
			"a_string": {Type: sdkschema.TypeString, Optional: true, Description: "what it does"},
		},
		"deprecated": {
			"a_string": {Type: sdkschema.TypeString, Optional: true, Deprecated: "use something else"},
		},
		// A Required attribute with a DefaultFunc that yields a value is reported as
		// Optional by the SDK v2, because Terraform has no conditional required-ness.
		"required with a default func that yields a value": {
			"a_string": {
				Type:        sdkschema.TypeString,
				Required:    true,
				DefaultFunc: func() (interface{}, error) { return "value", nil },
			},
		},
		// ... while one that yields nothing stays Required.
		"required with a default func that yields nothing": {
			"a_string": {
				Type:        sdkschema.TypeString,
				Required:    true,
				DefaultFunc: func() (interface{}, error) { return nil, nil },
			},
		},
		"set block": {
			"a_block": {
				Type:     sdkschema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &sdkschema.Resource{
					Schema: map[string]*sdkschema.Schema{
						"a_string": {Type: sdkschema.TypeString, Required: true},
						"an_int":   {Type: sdkschema.TypeInt, Optional: true},
					},
				},
			},
		},
		"list block": {
			"a_block": {
				Type:        sdkschema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "a described block",
				Deprecated:  "a deprecated block",
				Elem: &sdkschema.Resource{
					Schema: map[string]*sdkschema.Schema{
						"a_string": {Type: sdkschema.TypeString, Optional: true},
					},
				},
			},
		},
		"nested blocks": {
			"a_block": {
				Type:     sdkschema.TypeSet,
				Optional: true,
				Elem: &sdkschema.Resource{
					Schema: map[string]*sdkschema.Schema{
						"a_nested_block": {
							Type:     sdkschema.TypeList,
							Optional: true,
							Elem: &sdkschema.Resource{
								Schema: map[string]*sdkschema.Schema{
									"a_string": {Type: sdkschema.TypeString, Optional: true},
								},
							},
						},
					},
				},
			},
		},
	}

	for name, sdk := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			expected := sdkProtoSchema(ctx, t, sdk)
			actual := frameworkProtoSchema(ctx, t, sdk)

			if diff := cmp.Diff(expected, actual, schemaCmpOptions...); diff != "" {
				t.Errorf("unexpected difference between the SDK v2 and framework provider schemas (-sdk +framework):\n%s", diff)
			}
		})
	}
}

func TestUnitProviderSchemaUnsupported(t *testing.T) {
	testCases := map[string]map[string]*sdkschema.Schema{
		"map": {
			"a_map": {Type: sdkschema.TypeMap, Optional: true},
		},
		"float": {
			"a_float": {Type: sdkschema.TypeFloat, Optional: true},
		},
		"list of primitives": {
			"a_list": {
				Type:     sdkschema.TypeList,
				Optional: true,
				Elem:     &sdkschema.Schema{Type: sdkschema.TypeString},
			},
		},
		"computed": {
			"a_string": {Type: sdkschema.TypeString, Optional: true, Computed: true},
		},
		"config mode": {
			"a_block": {
				Type:       sdkschema.TypeList,
				Optional:   true,
				ConfigMode: sdkschema.SchemaConfigModeAttr,
				Elem: &sdkschema.Resource{
					Schema: map[string]*sdkschema.Schema{
						"a_string": {Type: sdkschema.TypeString, Optional: true},
					},
				},
			},
		},
		"unsupported field nested in a block": {
			"a_block": {
				Type:     sdkschema.TypeSet,
				Optional: true,
				Elem: &sdkschema.Resource{
					Schema: map[string]*sdkschema.Schema{
						"a_float": {Type: sdkschema.TypeFloat, Optional: true},
					},
				},
			},
		},
	}

	for name, sdk := range testCases {
		t.Run(name, func(t *testing.T) {
			if _, diags := ProviderSchema(sdk); !diags.HasError() {
				t.Error("expected an error diagnostic, got none")
			}
		})
	}
}

func sdkProtoSchema(ctx context.Context, t *testing.T, sdk map[string]*sdkschema.Schema) *tfprotov5.Schema {
	t.Helper()

	provider := &sdkschema.Provider{Schema: sdk}

	resp, err := provider.GRPCProvider().GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema on the SDK v2 provider: %s", err)
	}

	return resp.Provider
}

func frameworkProtoSchema(ctx context.Context, t *testing.T, sdk map[string]*sdkschema.Schema) *tfprotov5.Schema {
	t.Helper()

	server := providerserver.NewProtocol5(NewProvider(&sdkschema.Provider{Schema: sdk}))()

	resp, err := server.GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("GetProviderSchema on the framework provider: %s", err)
	}
	for _, d := range resp.Diagnostics {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			t.Fatalf("GetProviderSchema on the framework provider: %s: %s", d.Summary, d.Detail)
		}
	}

	return resp.Provider
}
