package geom

import (
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"
)

var (
	_ = []T{
		&LineString{},
		&LinearRing{},
		&MultiLineString{},
		&MultiPoint{},
		&MultiPolygon{},
		&Point{},
		&Polygon{},
	}
	_ = []interface {
		Area() float64
		Empty() bool
		Length() float64
	}{
		&LineString{},
		&LinearRing{},
		&MultiLineString{},
		&MultiPoint{},
		&MultiPolygon{},
		&Point{},
		&Polygon{},
	}
)

func aliases(x, y []float64) bool {
	// http://golang.org/src/pkg/math/big/nat.go#L340
	return cap(x) > 0 && cap(y) > 0 && &x[0:cap(x)][cap(x)-1] == &y[0:cap(y)][cap(y)-1]
}

func TestArea(t *testing.T) {
	for _, tc := range []struct {
		g interface {
			Area() float64
		}
		want float64
	}{
		{
			g:    NewPoint(XY),
			want: 0,
		},
		{
			g:    NewLineString(XY),
			want: 0,
		},
		{
			g:    NewLinearRing(XY),
			want: 0,
		},
		{
			g: NewLinearRing(XY).MustSetCoords([]Coord{
				{0, 0},
				{1, 0},
				{1, 1},
				{0, 1},
				{0, 0},
			}),
			want: 1,
		},
		{
			g: NewLinearRing(XY).MustSetCoords([]Coord{
				{0, 0},
				{1, 1},
				{1, 0},
				{0, 0},
			}),
			want: -0.5,
		},
		{
			g: NewLinearRing(XY).MustSetCoords([]Coord{
				{-3, -2},
				{-1, 4},
				{6, 1},
				{3, 10},
				{-4, 9},
				{-3, -2},
			}),
			want: 60,
		},
		{
			g:    NewMultiLineString(XY),
			want: 0,
		},
		{
			g:    NewPolygon(XY),
			want: 0,
		},
		{
			g: NewPolygon(XY).MustSetCoords([][]Coord{
				{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
			}),
			want: 1,
		},
		{
			g: NewPolygon(XY).MustSetCoords([][]Coord{
				{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
			}),
			want: -0.5,
		},
		{
			g: NewPolygon(XY).MustSetCoords([][]Coord{
				{{-3, -2}, {-1, 4}, {6, 1}, {3, 10}, {-4, 9}, {-3, -2}},
			}),
			want: 60,
		},
		{
			g: NewPolygon(XY).MustSetCoords([][]Coord{
				{{-3, -2}, {-1, 4}, {6, 1}, {3, 10}, {-4, 9}, {-3, -2}},
				{{0, 6}, {2, 6}, {2, 8}, {0, 8}, {0, 6}},
			}),
			want: 56,
		},
		{
			g:    NewMultiPolygon(XY),
			want: 0,
		},
		{
			g: NewMultiPolygon(XY).MustSetCoords([][][]Coord{
				{
					{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
				},
			}),
			want: 1,
		},
		{
			g: NewMultiPolygon(XY).MustSetCoords([][][]Coord{
				{
					{{0, 0}, {1, 1}, {1, 0}, {0, 0}},
				},
			}),
			want: -0.5,
		},
		{
			g: NewMultiPolygon(XY).MustSetCoords([][][]Coord{
				{
					{{-3, -2}, {-1, 4}, {6, 1}, {3, 10}, {-4, 9}, {-3, -2}},
				},
			}),
			want: 60,
		},
		{
			g: NewMultiPolygon(XY).MustSetCoords([][][]Coord{
				{
					{{-3, -2}, {-1, 4}, {6, 1}, {3, 10}, {-4, 9}, {-3, -2}},
					{{0, 6}, {2, 6}, {2, 8}, {0, 8}, {0, 6}},
				},
			}),
			want: 56,
		},
		{
			g: NewMultiPolygon(XY).MustSetCoords([][][]Coord{
				{
					{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
				},
				{
					{{-3, -2}, {-1, 4}, {6, 1}, {3, 10}, {-4, 9}, {-3, -2}},
					{{0, 6}, {2, 6}, {2, 8}, {0, 8}, {0, 6}},
				},
			}),
			want: 57,
		},
	} {
		assert.Equal(t, tc.want, tc.g.Area())
	}
}

func TestSet(t *testing.T) {
	for _, tc := range []struct {
		c, other, want Coord
	}{
		{Coord{1.0, 2.0}, Coord{2.0, 3.0}, Coord{2.0, 3.0}},
		{Coord{1.0, 2.0, 3.0}, Coord{2.0, 3.0, 4.0}, Coord{2.0, 3.0, 4.0}},
		{Coord{1.0, 2.0}, Coord{2.0, 3.0, 4.0}, Coord{2.0, 3.0}},
		{Coord{1.0, 2.0, 3.0}, Coord{2.0, 3.0}, Coord{2.0, 3.0, 3.0}},
	} {
		tc.c.Set(tc.other)
		assert.Equal(t, tc.want, tc.c)
	}
}

func TestLayoutString(t *testing.T) {
	for _, tc := range []struct {
		l    Layout
		want string
	}{
		{NoLayout, "NoLayout"},
		{XY, "XY"},
		{XYZ, "XYZ"},
		{XYM, "XYM"},
		{XYZM, "XYZM"},
		{Layout(5), "Layout(5)"},
	} {
		assert.Equal(t, tc.want, tc.l.String())
	}
}

func TestVerify(t *testing.T) {
	for _, tc := range []struct {
		v interface {
			verify() error
		}
		want error
	}{
		{
			&Geom0{},
			nil,
		},
		{
			&Geom0{NoLayout, 0, Coord{0, 0}, 0},
			errNonEmptyFlatCoords,
		},
		{
			&Geom0{XY, 1, Coord{0, 0}, 0},
			errStrideLayoutMismatch,
		},
		{
			&Geom0{XY, 2, Coord{0}, 0},
			errLengthStrideMismatch,
		},
		{
			&geom1{},
			nil,
		},
		{
			&geom1{Geom0{NoLayout, 0, Coord{0}, 0}},
			errNonEmptyFlatCoords,
		},
		{
			&geom1{Geom0{XY, 1, Coord{0, 0}, 0}},
			errStrideLayoutMismatch,
		},
		{
			&geom1{Geom0{XY, 2, Coord{0}, 0}},
			errLengthStrideMismatch,
		},
		{
			&geom2{},
			nil,
		},
		{
			&geom2{geom1{Geom0{NoLayout, 0, Coord{0}, 0}}, []int{}},
			errNonEmptyFlatCoords,
		},
		{
			&geom2{geom1{Geom0{NoLayout, 0, Coord{}, 0}}, []int{4}},
			errNonEmptyEnds,
		},
		{
			&geom2{geom1{Geom0{XY, 2, Coord{0}, 0}}, []int{4}},
			errLengthStrideMismatch,
		},
		{
			&geom2{geom1{Geom0{XY, 1, Coord{0, 0, 0, 0}, 0}}, []int{-1}},
			errStrideLayoutMismatch,
		},
		{
			&geom2{geom1{Geom0{XY, 2, Coord{0, 0, 0, 0}, 0}}, []int{-1}},
			errMisalignedEnd,
		},
		{
			&geom2{geom1{Geom0{XY, 2, Coord{0, 0, 0, 0}, 0}}, []int{3}},
			errMisalignedEnd,
		},
		{
			&geom2{geom1{Geom0{XY, 2, Coord{0, 0, 0, 0, 0, 0, 0, 0}, 0}}, []int{8, 4}},
			errOutOfOrderEnd,
		},
		{
			&geom2{geom1{Geom0{XY, 2, Coord{0, 0, 0, 0, 0, 0, 0, 0}, 0}}, []int{4, 4}},
			errIncorrectEnd,
		},
		{
			&geom2{geom1{Geom0{XY, 2, Coord{0, 0, 0, 0, 0, 0, 0, 0}, 0}}, []int{4, 12}},
			errIncorrectEnd,
		},
		{
			&geom3{},
			nil,
		},
		{
			&geom3{geom1{Geom0{XY, 3, Coord{}, 0}}, [][]int{}},
			errStrideLayoutMismatch,
		},
		{
			&geom3{geom1{Geom0{NoLayout, 0, Coord{0}, 0}}, [][]int{}},
			errNonEmptyFlatCoords,
		},
		{
			&geom3{geom1{Geom0{NoLayout, 0, Coord{}, 0}}, [][]int{{0}}},
			errNonEmptyEndss,
		},
		{
			&geom3{geom1{Geom0{XY, 2, Coord{0}, 0}}, [][]int{}},
			errLengthStrideMismatch,
		},
		{
			&geom3{geom1{Geom0{XY, 2, Coord{0, 0}, 0}}, [][]int{{1}}},
			errMisalignedEnd,
		},
		{
			&geom3{geom1{Geom0{XY, 2, Coord{0, 0, 0, 0}, 0}}, [][]int{{4, 2}}},
			errOutOfOrderEnd,
		},
		{
			&geom3{geom1{Geom0{XY, 2, Coord{0, 0, 0, 0}, 0}}, [][]int{{2}}},
			errIncorrectEnd,
		},
	} {
		assert.Equal(t, tc.want, tc.v.verify())
	}
}

func TestEqualCoords(t *testing.T) {
	for _, tc := range []struct {
		c1, c2 Coord
		layout Layout
		equal  bool
	}{
		{
			c1:     Coord{},
			c2:     Coord{0, 0},
			layout: XY,
			equal:  false,
		},
		{
			c1:     Coord{},
			c2:     Coord{},
			layout: XY,
			equal:  true,
		},
		{
			c1:     Coord{1, 0},
			c2:     Coord{},
			layout: XY,
			equal:  false,
		},
		{
			c1:     Coord{1, 0},
			c2:     Coord{1},
			layout: XY,
			equal:  false,
		},
		{
			c1:     Coord{1},
			c2:     Coord{},
			layout: XY,
			equal:  false,
		},
		{
			c1:     Coord{1},
			c2:     Coord{1},
			layout: XY,
			equal:  true,
		},
		{
			c1:     Coord{1},
			c2:     Coord{0},
			layout: XY,
			equal:  false,
		},
		{
			c1:     Coord{0, 0},
			c2:     Coord{0, 0},
			layout: XY,
			equal:  true,
		},
		{
			c1:     Coord{0, 0},
			c2:     Coord{1, 0},
			layout: XY,
			equal:  false,
		},
		{
			c1:     Coord{0, 1},
			c2:     Coord{0, 0},
			layout: XY,
			equal:  false,
		},
		{
			c1:     Coord{0, 0, 3},
			c2:     Coord{0, 0},
			layout: XY,
			equal:  true,
		},
		{
			c1:     Coord{0, 0, 3},
			c2:     Coord{0, 0, 3},
			layout: XYZ,
			equal:  true,
		},
		{
			c1:     Coord{0, 0, 3},
			c2:     Coord{0, 0, 4},
			layout: XYZ,
			equal:  false,
		},
		{
			c1:     Coord{0, 0, 3, 4, 5, 6, 7, 8, 9, 10},
			c2:     Coord{0, 0, 3, 4, 5, 6, 7, 8, 9, 10},
			layout: Layout(10),
			equal:  true,
		},
		{
			c1:     Coord{0, 0, 3, 4, 5, 6, 7, 8, 9, 10},
			c2:     Coord{0, 0, 3, 4, 5, 6, 8, 8, 9, 10},
			layout: Layout(10),
			equal:  false,
		},
	} {
		assert.Equal(t, tc.equal, tc.c1.Equal(tc.layout, tc.c2))
	}
}

func TestLength(t *testing.T) {
	for _, tc := range []struct {
		g interface {
			Length() float64
		}
		want float64
	}{
		{
			g:    NewPoint(XY),
			want: 0,
		},
		{
			g:    NewMultiPoint(XY),
			want: 0,
		},
		{
			g:    NewLineString(XY),
			want: 0,
		},
		{
			g: NewLineString(XY).MustSetCoords([]Coord{
				{0, 0},
				{1, 0},
			}),
			want: 1,
		},
		{
			g:    NewLinearRing(XY),
			want: 0,
		},
		{
			g: NewLinearRing(XY).MustSetCoords([]Coord{
				{0, 0},
				{1, 0},
				{1, 1},
				{0, 1},
				{0, 0},
			}),
			want: 4,
		},
		{
			g:    NewPolygon(XY),
			want: 0,
		},
		{
			g: NewPolygon(XY).MustSetCoords([][]Coord{
				{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
			}),
			want: 4,
		},
		{
			g: NewPolygon(XY).MustSetCoords([][]Coord{
				{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
				{{0.25, 0.25}, {0.75, 0.25}, {0.75, 0.75}, {0.25, 0.75}, {0.25, 0.25}},
			}),
			want: 6,
		},
		{
			g:    NewMultiPolygon(XY),
			want: 0,
		},
		{
			g: NewMultiPolygon(XY).MustSetCoords([][][]Coord{
				{
					{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
					{{0.25, 0.25}, {0.75, 0.25}, {0.75, 0.75}, {0.25, 0.75}, {0.25, 0.25}},
				},
			}),
			want: 6,
		},
		{
			g: NewMultiPolygon(XY).MustSetCoords([][][]Coord{
				{
					{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}},
					{{0.25, 0.25}, {0.75, 0.25}, {0.75, 0.75}, {0.25, 0.75}, {0.25, 0.25}},
				},
				{
					{{2, 2}, {4, 2}, {4, 4}, {2, 4}, {2, 2}},
				},
			}),
			want: 14,
		},
	} {
		assert.Equal(t, tc.want, tc.g.Length())
	}
}

func TestSetCoord(t *testing.T) {
	for _, tc := range []struct {
		src, dest Coord
		expected  Coord
		layout    Layout
	}{
		{
			src:      Coord{0, 0},
			dest:     Coord{1, 1},
			expected: Coord{0, 0},
			layout:   XY,
		},
		{
			src:      Coord{1, 0},
			dest:     Coord{},
			expected: Coord{},
			layout:   Layout(0),
		},
		{
			src:      Coord{},
			dest:     Coord{1, 2},
			expected: Coord{1, 2},
			layout:   XY,
		},
		{
			src:      Coord{3},
			dest:     Coord{1, 2},
			expected: Coord{3, 2},
			layout:   XY,
		},
	} {
		tc.dest.Set(tc.src)
		assert.True(t, tc.dest.Equal(tc.layout, tc.expected))
	}
}

func TestTransformInPlace(t *testing.T) {
	f := func(coord Coord) {
		for i := range coord {
			coord[i] += float64(i + 1)
		}
	}
	for i, tc := range []struct {
		g        T
		expected T
	}{
		{
			g:        NewPoint(XY).MustSetCoords(Coord{0, 0}),
			expected: NewPoint(XY).MustSetCoords(Coord{1, 2}),
		},
		{
			g:        NewPoint(XYZ).MustSetCoords(Coord{0, 0, 0}),
			expected: NewPoint(XYZ).MustSetCoords(Coord{1, 2, 3}),
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tc.expected, TransformInPlace(tc.g, f))
		})
	}
}

func TestSetSRID(t *testing.T) {
	_, err := SetSRID(nil, 4326)
	assert.Error(t, err)
}

func TestReverse(t *testing.T) {
	for _, tc := range []struct {
		g interface {
			Reverse()
		}
		want interface {
			Reverse()
		}
	}{
		{
			g:    NewLinearRing(XYZM).MustSetCoords([]Coord{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}),
			want: NewLinearRing(XYZM).MustSetCoords([]Coord{{9, 10, 11, 12}, {5, 6, 7, 8}, {1, 2, 3, 4}}),
		},
		{
			g:    NewLineString(XYZM).MustSetCoords([]Coord{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}),
			want: NewLineString(XYZM).MustSetCoords([]Coord{{9, 10, 11, 12}, {5, 6, 7, 8}, {1, 2, 3, 4}}),
		},
		{
			g:    NewMultiLineString(XY).MustSetCoords([][]Coord{{}, {}, {{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}, {}}),
			want: NewMultiLineString(XY).MustSetCoords([][]Coord{{}, {}, {{5, 6}, {3, 4}, {1, 2}}, {{11, 12}, {9, 10}, {7, 8}}, {}}),
		},
		{
			g:    NewMultiPoint(XY).MustSetCoords([]Coord{nil, {1, 2}, nil, {3, 4}, nil, {5, 6}, nil}),
			want: NewMultiPoint(XY).MustSetCoords([]Coord{nil, {1, 2}, nil, {3, 4}, nil, {5, 6}, nil}),
		},
		{
			g:    NewMultiPolygon(XY).MustSetCoords([][][]Coord{{{{1, 2}, {4, 5}, {6, 7}, {1, 2}}}, {}, {}, {{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}}, {}}),
			want: NewMultiPolygon(XY).MustSetCoords([][][]Coord{{{{1, 2}, {6, 7}, {4, 5}, {1, 2}}}, {}, {}, {{{5, 6}, {3, 4}, {1, 2}}, {{11, 12}, {9, 10}, {7, 8}}}, {}}),
		},
		{
			g:    NewPolygon(XY).MustSetCoords([][]Coord{{{1, 2}, {3, 4}, {5, 6}}, {{7, 8}, {9, 10}, {11, 12}}}),
			want: NewPolygon(XY).MustSetCoords([][]Coord{{{5, 6}, {3, 4}, {1, 2}}, {{11, 12}, {9, 10}, {7, 8}}}),
		},
	} {
		tc.g.Reverse()
		assert.Equal(t, tc.want, tc.g)
	}
}
