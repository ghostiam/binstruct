package binstruct

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_parseTag(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		want []tag
	}{
		{
			name: "empty",
			tag:  "",
			want: nil,
		},
		{
			name: "ignore",
			tag:  "-",
			want: []tag{{Type: "-"}},
		},
		{
			name: "func",
			tag:  "TestFunc",
			want: []tag{{Type: "func", Value: "TestFunc"}},
		},
		{
			name: "element ignore",
			tag:  "[-]",
			want: []tag{{Type: "elem", ElemTags: []tag{{Type: "-"}}}},
		},
		{
			name: "len",
			tag:  "len:42",
			want: []tag{{Type: "len", Value: "42"}},
		},
		{
			name: "offset",
			tag:  "offset:42",
			want: []tag{{Type: "offset", Value: "42"}},
		},
		{
			name: "offsetStart",
			tag:  "offsetStart:42",
			want: []tag{{Type: "offsetStart", Value: "42"}},
		},
		{
			name: "offsetEnd",
			tag:  "offsetEnd:42",
			want: []tag{{Type: "offsetEnd", Value: "42"}},
		},
		{
			name: "multi tag",
			tag:  "len:1, offset:2, offsetStart:3, offsetEnd:4",
			want: []tag{
				{Type: "len", Value: "1"},
				{Type: "offset", Value: "2"},
				{Type: "offsetStart", Value: "3"},
				{Type: "offsetEnd", Value: "4"},
			},
		},
		{
			name: "multi element",
			tag:  "[len:1, [len:2, [len3]]], offset:42",
			want: []tag{
				{
					Type: "elem", Value: "", ElemTags: []tag{
					{
						Type: "len", Value: "1", ElemTags: []tag(nil),
					},
					{
						Type: "elem", Value: "", ElemTags: []tag{
						{
							Type: "len", Value: "2", ElemTags: []tag(nil),
						},
						{
							Type: "elem", Value: "", ElemTags: []tag{
							{
								Type: "func", Value: "len3", ElemTags: []tag(nil),
							},
						},
						},
					},
					},
				},
				},
				{
					Type: "offset", Value: "42", ElemTags: []tag(nil),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTag(tt.tag)
			require.Equal(t, tt.want, got)
		})
	}
}
