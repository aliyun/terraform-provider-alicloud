package registry

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/action"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/list"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
)

func sdkFactory() *sdkschema.Resource               { return &sdkschema.Resource{} }
func resourceFactory() fwresource.Resource          { return nil }
func dataSourceFactory() fwdatasource.DataSource    { return nil }
func functionFactory() function.Function            { return nil }
func ephemeralFactory() ephemeral.EphemeralResource { return nil }
func listFactory() list.ListResource                { return nil }
func actionFactory() action.Action                  { return nil }

func TestValidateRealDeclarations(t *testing.T) {
	if err := validate(servicePackages); err != nil {
		t.Fatalf("the shipped service package declarations are invalid: %v", err)
	}
}

func TestValidateAcceptsOneOfEachCategory(t *testing.T) {
	sp := conns.ServicePackage{
		Name:               "Vpc",
		SDKResources:       []conns.SDKResource{{TypeName: "alicloud_vpc", Factory: sdkFactory}},
		Resources:          []conns.Resource{{TypeName: "alicloud_vpc_v2", Factory: resourceFactory}},
		SDKDataSources:     []conns.SDKDataSource{{TypeName: "alicloud_vpcs", Factory: sdkFactory}},
		DataSources:        []conns.DataSource{{TypeName: "alicloud_vpcs_v2", Factory: dataSourceFactory}},
		Functions:          []conns.Function{{TypeName: "arn_build", Factory: functionFactory}},
		EphemeralResources: []conns.EphemeralResource{{TypeName: "alicloud_vpc_token", Factory: ephemeralFactory}},
		ListResources:      []conns.ListResource{{TypeName: "alicloud_vpc", Factory: listFactory}},
		Actions:            []conns.Action{{TypeName: "alicloud_vpc_flush", Factory: actionFactory}},
	}
	if err := validate([]conns.ServicePackage{sp}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateNamespaces(t *testing.T) {
	cases := []struct {
		name string
		sps  []conns.ServicePackage
		want string
	}{
		{
			name: "two SDK resources",
			sps: []conns.ServicePackage{{Name: "Vpc",
				SDKResources: []conns.SDKResource{
					{TypeName: "alicloud_vpc", Factory: sdkFactory},
					{TypeName: "alicloud_vpc", Factory: sdkFactory},
				}}},
			want: `resource "alicloud_vpc" is declared twice: Vpc.SDKResources and Vpc.SDKResources`,
		},
		{
			name: "SDK and framework resource share the managed resource namespace",
			sps: []conns.ServicePackage{{Name: "Vpc",
				SDKResources: []conns.SDKResource{{TypeName: "alicloud_vpc", Factory: sdkFactory}},
				Resources:    []conns.Resource{{TypeName: "alicloud_vpc", Factory: resourceFactory}},
			}},
			want: `resource "alicloud_vpc" is declared twice: Vpc.SDKResources and Vpc.Resources`,
		},
		{
			name: "SDK and framework data source share the data source namespace",
			sps: []conns.ServicePackage{{Name: "Vpc",
				SDKDataSources: []conns.SDKDataSource{{TypeName: "alicloud_vpcs", Factory: sdkFactory}},
				DataSources:    []conns.DataSource{{TypeName: "alicloud_vpcs", Factory: dataSourceFactory}},
			}},
			want: `data source "alicloud_vpcs" is declared twice`,
		},
		{
			name: "two service packages claim the same resource",
			sps: []conns.ServicePackage{
				{Name: "Vpc", SDKResources: []conns.SDKResource{{TypeName: "alicloud_vpc", Factory: sdkFactory}}},
				{Name: "Vpc2", Resources: []conns.Resource{{TypeName: "alicloud_vpc", Factory: resourceFactory}}},
			},
			want: `resource "alicloud_vpc" is declared twice: Vpc.SDKResources and Vpc2.Resources`,
		},
		{
			name: "a resource and a data source may share a name",
			sps: []conns.ServicePackage{{Name: "Vpc",
				SDKResources:   []conns.SDKResource{{TypeName: "alicloud_vpc", Factory: sdkFactory}},
				SDKDataSources: []conns.SDKDataSource{{TypeName: "alicloud_vpc", Factory: sdkFactory}},
			}},
		},
		{
			name: "a list resource shares its managed resource's name",
			sps: []conns.ServicePackage{{Name: "Vpc",
				SDKResources:  []conns.SDKResource{{TypeName: "alicloud_vpc", Factory: sdkFactory}},
				ListResources: []conns.ListResource{{TypeName: "alicloud_vpc", Factory: listFactory}},
			}},
		},
		{
			name: "two list resources may not",
			sps: []conns.ServicePackage{{Name: "Vpc",
				ListResources: []conns.ListResource{
					{TypeName: "alicloud_vpc", Factory: listFactory},
					{TypeName: "alicloud_vpc", Factory: listFactory},
				}}},
			want: `list resource "alicloud_vpc" is declared twice`,
		},
		{
			name: "two functions",
			sps: []conns.ServicePackage{{Name: "Arn",
				Functions: []conns.Function{
					{TypeName: "arn_build", Factory: functionFactory},
					{TypeName: "arn_build", Factory: functionFactory},
				}}},
			want: `function "arn_build" is declared twice`,
		},
		{
			name: "two ephemeral resources",
			sps: []conns.ServicePackage{{Name: "Sts",
				EphemeralResources: []conns.EphemeralResource{
					{TypeName: "alicloud_sts_token", Factory: ephemeralFactory},
					{TypeName: "alicloud_sts_token", Factory: ephemeralFactory},
				}}},
			want: `ephemeral resource "alicloud_sts_token" is declared twice`,
		},
		{
			name: "two actions",
			sps: []conns.ServicePackage{{Name: "Ecs",
				Actions: []conns.Action{
					{TypeName: "alicloud_instance_reboot", Factory: actionFactory},
					{TypeName: "alicloud_instance_reboot", Factory: actionFactory},
				}}},
			want: `action "alicloud_instance_reboot" is declared twice`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate(tc.sps)
			if tc.want == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected an error containing %q", tc.want)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error %q does not contain %q", err, tc.want)
			}
		})
	}
}

func TestValidateEmptyTypeName(t *testing.T) {
	err := validate([]conns.ServicePackage{{Name: "Vpc",
		SDKResources: []conns.SDKResource{{Factory: sdkFactory}}}})
	if err == nil || !strings.Contains(err.Error(), "resource declared in Vpc.SDKResources has an empty TypeName") {
		t.Fatalf("got %v", err)
	}
}

func TestValidateEmptyTypeNameDoesNotCollide(t *testing.T) {
	err := validate([]conns.ServicePackage{{Name: "Vpc",
		SDKResources: []conns.SDKResource{{Factory: sdkFactory}, {Factory: sdkFactory}}}})
	if err == nil {
		t.Fatal("expected an error")
	}
	if got := strings.Count(err.Error(), "empty TypeName"); got != 2 {
		t.Fatalf("expected 2 empty-TypeName problems, got %d: %v", got, err)
	}
	if strings.Contains(err.Error(), "declared twice") {
		t.Fatalf("an empty name must not be reported as a duplicate: %v", err)
	}
}

func TestValidateNilFactory(t *testing.T) {
	err := validate([]conns.ServicePackage{{Name: "Vpc",
		Resources: []conns.Resource{{TypeName: "alicloud_vpc"}}}})
	if err == nil || !strings.Contains(err.Error(), `resource "alicloud_vpc" declared in Vpc.Resources has a nil Factory`) {
		t.Fatalf("got %v", err)
	}
}

func TestValidateReportsEveryProblem(t *testing.T) {
	err := validate([]conns.ServicePackage{{Name: "Vpc",
		SDKResources: []conns.SDKResource{
			{TypeName: "alicloud_vpc", Factory: sdkFactory},
			{TypeName: "alicloud_vpc", Factory: sdkFactory},
			{TypeName: "alicloud_route"},
			{Factory: sdkFactory},
		}}})
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.HasPrefix(err.Error(), "3 invalid service package declaration(s):") {
		t.Fatalf("expected a 3-problem summary line, got %v", err)
	}
	for _, want := range []string{"declared twice", "nil Factory", "empty TypeName"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error is missing %q: %v", want, err)
		}
	}
}
