package binstruct

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
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

func Test_parseReadDataFromTags(t *testing.T) {
	type args struct {
		structValue reflect.Value
		tags        []tag
	}
	tests := []struct {
		name string
		args args
		want *fieldReadData
	}{
		{
			name: "calc len 5+2",
			args: args{
				structValue: reflect.ValueOf(struct{}{}),
				tags: []tag{
					{
						Type:  "len",
						Value: "5+2",
					},
				},
			},
			want: &fieldReadData{
				Length: func() *int64 {
					i := int64(7)
					return &i
				}(),
			},
		},
		{
			name: "calc len 1-2",
			args: args{
				structValue: reflect.ValueOf(struct{}{}),
				tags: []tag{
					{
						Type:  "len",
						Value: "5-2",
					},
				},
			},
			want: &fieldReadData{
				Length: func() *int64 {
					i := int64(3)
					return &i
				}(),
			},
		},
		{
			name: "calc len 5+FieldValue",
			args: args{
				structValue: reflect.ValueOf(struct {
					FieldValue int
				}{
					FieldValue: 2,
				}),
				tags: []tag{
					{
						Type:  "len",
						Value: "5+FieldValue",
					},
				},
			},
			want: &fieldReadData{
				Length: func() *int64 {
					i := int64(7)
					return &i
				}(),
			},
		},
		{
			name: "calc len 5+FieldValue",
			args: args{
				structValue: reflect.ValueOf(struct {
					FieldValue int
				}{
					FieldValue: 2,
				}),
				tags: []tag{
					{
						Type:  "len",
						Value: "5-FieldValue",
					},
				},
			},
			want: &fieldReadData{
				Length: func() *int64 {
					i := int64(3)
					return &i
				}(),
			},
		},
		{
			name: "calc len many 5+FieldValue+10-5-FieldSub",
			args: args{
				structValue: reflect.ValueOf(struct {
					FieldAdd int
					FieldSub uint
				}{
					FieldAdd: 5,
					FieldSub: 10,
				}),
				tags: []tag{
					{
						Type:  "len",
						Value: "10 + FieldAdd + 10 - 5 - FieldSub",
					},
				},
			},
			want: &fieldReadData{
				Length: func() *int64 {
					i := int64(10)
					return &i
				}(),
			},
		},
		{
			name: "calc offset -10",
			args: args{
				structValue: reflect.ValueOf(struct{}{}),
				tags: []tag{
					{
						Type:  "len",
						Value: "-10",
					},
				},
			},
			want: &fieldReadData{
				Length: func() *int64 {
					i := int64(-10)
					return &i
				}(),
			},
		},
		{
			name: "calc offset -10 + -5",
			args: args{
				structValue: reflect.ValueOf(struct{}{}),
				tags: []tag{
					{
						Type:  "len",
						Value: "-10 + -5",
					},
				},
			},
			want: &fieldReadData{
				Length: func() *int64 {
					i := int64(-15)
					return &i
				}(),
			},
		},
		{
			name: "calc offset -10 + -5 + 5",
			args: args{
				structValue: reflect.ValueOf(struct{}{}),
				tags: []tag{
					{
						Type:  "len",
						Value: "-10 + -5 + 5",
					},
				},
			},
			want: &fieldReadData{
				Length: func() *int64 {
					i := int64(-10)
					return &i
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseReadDataFromTags(tt.args.structValue, tt.args.tags)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}