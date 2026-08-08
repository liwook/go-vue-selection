package idconv

import (
	"testing"
)

func TestToStr(t *testing.T) {
	tests := []struct {
		name string
		in   int64
		want string
	}{
		{"positive", 123456, "123456"},
		{"zero", 0, "0"},
		{"negative", -789, "-789"},
		{"max int64", 9223372036854775807, "9223372036854775807"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToStr(tt.in); got != tt.want {
				t.Errorf("ToStr(%d) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestToStrPtr(t *testing.T) {
	tests := []struct {
		name string
		in   *int64
		want string
	}{
		{"non-nil", int64Ptr(42), "42"},
		{"zero value ptr", int64Ptr(0), "0"},
		{"nil -> empty", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToStrPtr(tt.in); got != tt.want {
				t.Errorf("ToStrPtr(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    int64
		wantErr bool
	}{
		{"normal", "123456", 123456, false},
		{"negative", "-789", -789, false},
		{"zero", "0", 0, false},
		{"empty -> zero no err", "", 0, false},
		{"invalid", "abc", 0, true},
		// strconv.ParseInt 溢出时返回 ErrRange，且值被截断为 int64 上界（非 0）
		{"overflow", "999999999999999999999", 9223372036854775807, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt64(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt64(%q) err = %v, wantErr %v", tt.in, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToInt64(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestToInt64Safe(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int64
	}{
		{"normal", "123456", 123456},
		{"empty -> 0", "", 0},
		{"invalid -> 0", "abc", 0},
		{"negative", "-5", -5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt64Safe(tt.in); got != tt.want {
				t.Errorf("ToInt64Safe(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func int64Ptr(v int64) *int64 { return &v }
