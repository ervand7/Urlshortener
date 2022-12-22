package algorithms

import (
	"strconv"
	"testing"
)

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

func BenchmarkIssubset(b *testing.B) {
	N := 1000

	b.Run("Size 100", func(b *testing.B) {
		b.ReportAllocs()
		length := 100
		first := make([]string, length)
		second := make([]string, length)
		for i := 0; i < length; i++ {
			first = append(first, strconv.Itoa(i))
			second = append(second, strconv.Itoa(i))
		}

		b.ResetTimer()
		for i := 0; i < N; i++ {
			Issubset(first, second)
		}
	})

	b.Run("Size 1_000", func(b *testing.B) {
		b.ReportAllocs()
		length := 1_000
		first := make([]string, length)
		second := make([]string, length)
		for i := 0; i < length; i++ {
			first = append(first, strconv.Itoa(i))
			second = append(second, strconv.Itoa(i))
		}

		b.ResetTimer()
		for i := 0; i < N; i++ {
			Issubset(first, second)
		}
	})

	b.Run("Size 5_000", func(b *testing.B) {
		b.ReportAllocs()
		length := 5_000
		first := make([]string, length)
		second := make([]string, length)
		for i := 0; i < length; i++ {
			first = append(first, strconv.Itoa(i))
			second = append(second, strconv.Itoa(i))
		}

		b.ResetTimer()
		for i := 0; i < N; i++ {
			Issubset(first, second)
		}
	})
}
