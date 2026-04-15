package mapping

import (
	"testing"

	"github.com/ompgo-dev/ompgo/tools/codegen/model"
)

func TestStringOutputParamInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		fn        model.APIFunction
		wantCount int
		wantType  string
	}{
		{
			name: "string view output only",
			fn: model.APIFunction{
				Params: []model.APIParam{{Name: "name", Type: "CAPIStringView*"}},
			},
			wantCount: 1,
			wantType:  "CAPIStringView*",
		},
		{
			name: "string buffer output with string input",
			fn: model.APIFunction{
				Params: []model.APIParam{
					{Name: "value", Type: "const char*"},
					{Name: "buf", Type: "CAPIStringBuffer*"},
				},
			},
			wantCount: 1,
			wantType:  "CAPIStringBuffer*",
		},
		{
			name: "multiple string outputs disqualify candidate",
			fn: model.APIFunction{
				Params: []model.APIParam{
					{Name: "first", Type: "CAPIStringView*"},
					{Name: "second", Type: "CAPIStringBuffer*"},
				},
			},
			wantCount: 2,
			wantType:  "CAPIStringBuffer*",
		},
		{
			name: "non string pointer out disqualifies",
			fn: model.APIFunction{
				Params: []model.APIParam{
					{Name: "id", Type: "int32_t*"},
					{Name: "name", Type: "CAPIStringView*"},
				},
			},
			wantCount: 0,
			wantType:  "",
		},
		{
			name: "void pointer allowed alongside string output",
			fn: model.APIFunction{
				Params: []model.APIParam{
					{Name: "entity", Type: "void*"},
					{Name: "name", Type: "CAPIStringView*"},
				},
			},
			wantCount: 1,
			wantType:  "CAPIStringView*",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotCount, gotType := StringOutputParamInfo(tt.fn)
			if gotCount != tt.wantCount || gotType != tt.wantType {
				t.Fatalf(
					"StringOutputParamInfo() = (%d, %q), want (%d, %q)",
					gotCount,
					gotType,
					tt.wantCount,
					tt.wantType,
				)
			}
		})
	}
}

func TestIsStringOutputCandidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   model.APIFunction
		want bool
	}{
		{
			name: "single string output",
			fn: model.APIFunction{
				Params: []model.APIParam{{Name: "name", Type: "CAPIStringView*"}},
			},
			want: true,
		},
		{
			name: "multiple string outputs",
			fn: model.APIFunction{
				Params: []model.APIParam{
					{Name: "first", Type: "CAPIStringView*"},
					{Name: "second", Type: "CAPIStringBuffer*"},
				},
			},
			want: false,
		},
		{
			name: "extra pointer output",
			fn: model.APIFunction{
				Params: []model.APIParam{
					{Name: "player", Type: "void*"},
					{Name: "count", Type: "uint32_t*"},
					{Name: "name", Type: "CAPIStringView*"},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := IsStringOutputCandidate(tt.fn); got != tt.want {
				t.Fatalf("IsStringOutputCandidate() = %v, want %v", got, tt.want)
			}
		})
	}
}
