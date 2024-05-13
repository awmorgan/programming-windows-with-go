package win32

import (
	"testing"
)

func TestPtInRect(t *testing.T) {
	rect := RECT{Left: 10, Top: 10, Right: 100, Bottom: 100}
	testCases := []struct {
		name  string
		point POINT
		want  bool
	}{
		{"Point Inside", POINT{X: 50, Y: 50}, true},
		{"Point Outside", POINT{X: 150, Y: 150}, false},
		{"Point On Left Edge", POINT{X: 10, Y: 50}, true},
		{"Point On Top Edge", POINT{X: 50, Y: 10}, true},
		{"Point On Right Edge", POINT{X: 100, Y: 50}, false},
		{"Point On Bottom Edge", POINT{X: 50, Y: 100}, false},
		{"Point On Top Left Corner", POINT{X: 10, Y: 10}, true},
		{"Point On Bottom Right Corner", POINT{X: 100, Y: 100}, false},
		{"Point Just Inside Right Edge", POINT{X: 99, Y: 50}, true},
		{"Point Just Inside Bottom Edge", POINT{X: 50, Y: 99}, true},
		{"Point Just Outside Right Edge", POINT{X: 101, Y: 50}, false},
		{"Point Just Outside Bottom Edge", POINT{X: 50, Y: 101}, false},
		{"Point At Origin", POINT{X: 0, Y: 0}, false},
		{"Point Negative", POINT{X: -50, Y: -50}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := PtInRect(&rect, tc.point); got != tc.want {
				t.Errorf("PtInRect was incorrect for %s, got: %v, want: %v", tc.name, got, tc.want)
			}
		})
	}
}

// // Benchmark for the syscall version of PtInRect
// func BenchmarkPtInRectSyscall(b *testing.B) {
// 	rect := RECT{Left: 10, Top: 10, Right: 100, Bottom: 100}
// 	point := POINT{X: 50, Y: 50} // A point inside the rectangle
// 	for i := 0; i < b.N; i++ {
// 		PtInRect(&rect, point)
// 	}
// }

// Benchmark for the pure Go version of PtInRect
func BenchmarkPtInRectPureGo(b *testing.B) {
	rect := RECT{Left: 10, Top: 10, Right: 100, Bottom: 100}
	point := POINT{X: 50, Y: 50} // A point inside the rectangle
	for i := 0; i < b.N; i++ {
		PtInRect(&rect, point)
	}
}
