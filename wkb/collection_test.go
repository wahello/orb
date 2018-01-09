package wkb

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/paulmach/orb"
)

func TestCollection(t *testing.T) {
	cases := []struct {
		name     string
		bytes    []byte
		bom      binary.ByteOrder
		expected orb.Collection
	}{
		{
			name: "collection",
			bytes: []byte{
				//01    02    03    04    05    06    07    08
				0x02, 0x00, 0x00, 0x00, // Number of Geometries in Collection
				0x01,                   // Byte order marker little
				0x01, 0x00, 0x00, 0x00, // Type (1) Point
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x40, // X1 4
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x40, // Y1 6
				0x01,                   // Byte order marker little
				0x02, 0x00, 0x00, 0x00, // Type (2) Line
				0x02, 0x00, 0x00, 0x00, // Number of Points (2)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x40, // X1 4
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x40, // Y1 6
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1c, 0x40, // X2 7
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y2 10
			},
			bom: binary.LittleEndian,
			expected: orb.Collection{
				orb.Point{4, 6},
				orb.LineString{{4, 6}, {7, 10}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := readCollection(bytes.NewReader(tc.bytes), tc.bom)
			if err != nil {
				t.Fatalf("read error: %v", err)
			}

			if len(c) != len(tc.expected) {
				t.Fatalf("incorrect length: %d != %d", len(c), len(tc.expected))
			}

			for i := range tc.expected {
				if c[i].GeoJSONType() != tc.expected[i].GeoJSONType() {
					t.Errorf("expected[%v]: %v != %v", i, c[i], tc.expected[i])
				}
			}
		})
	}
}
