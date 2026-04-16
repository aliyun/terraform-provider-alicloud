package alicloud

import "testing"

func TestMergeSecurityGroupsWithAttach(t *testing.T) {
	attach := []string{"sg-4", "sg-3", "sg-4"}
	existing := []string{"sg-1", "sg-2", "sg-3"}

	got := mergeSecurityGroupsWithAttach(attach, existing)
	want := []string{"sg-4", "sg-3", "sg-1", "sg-2"}

	if len(got) != len(want) {
		t.Fatalf("unexpected length, got=%v want=%v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("unexpected merge result, got=%v want=%v", got, want)
		}
	}
}

func TestEqualStringSet(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		want bool
	}{
		{
			name: "same members different order",
			a:    []string{"sg-1", "sg-2", "sg-3"},
			b:    []string{"sg-3", "sg-1", "sg-2"},
			want: true,
		},
		{
			name: "duplicates are ignored",
			a:    []string{"sg-1", "sg-2", "sg-2"},
			b:    []string{"sg-2", "sg-1"},
			want: true,
		},
		{
			name: "different members",
			a:    []string{"sg-1", "sg-2"},
			b:    []string{"sg-1", "sg-3"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := equalStringSet(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("equalStringSet()=%v, want=%v, a=%v, b=%v", got, tt.want, tt.a, tt.b)
			}
		})
	}
}
