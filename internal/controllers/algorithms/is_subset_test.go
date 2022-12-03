package algorithms

import "testing"

func TestIssubset(t *testing.T) {
	type args struct {
		first  []string
		second []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_success",
			args: args{
				first:  []string{"1", "2", "3"},
				second: []string{"2", "1"},
			},
			want: true,
		},
		{
			name: "test_fail",
			args: args{
				first:  []string{"1", "2", "3"},
				second: []string{"5", "1"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Issubset(tt.args.first, tt.args.second); got != tt.want {
				t.Errorf("Issubset() = %v, want %v", got, tt.want)
			}
		})
	}
}