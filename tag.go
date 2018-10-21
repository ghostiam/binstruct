package binstruct

import (
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

const (
	tagName = "bin"
	/*
		len:10,offset:20,skip:5
		len:Len,offset:Offset,skip:Skip
		func // func from struct method or FuncMap
	*/
)

const (
	tagTypeIgnore  = "-"
	tagTypeEmpty   = ""
	tagTypeFunc    = "func"
	tagTypeElement = "elem"

	tagTypeLength            = "len"
	tagTypeOffsetFromCurrent = "offset"
	tagTypeOffsetFromStart   = "offsetStart"
	tagTypeOffsetFromEnd     = "offsetEnd"
)

type tag struct {
	Type  string
	Value string

	ElemTags []tag
}

func parseTag(t string) []tag {
	var tags []tag

	split := strings.Split(t, ",")
	for _, v := range split {
		v = strings.TrimSpace(v)

		switch {
		case v == tagTypeIgnore:
			tags = append(tags, tag{Type: tagTypeIgnore})
		case v == tagTypeEmpty:
			tags = append(tags, tag{Type: tagTypeEmpty})
		case strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]"):
			elem := v[1 : len(v)-1]
			tags = append(tags, tag{Type: tagTypeElement, ElemTags: parseTag(elem)})
		default:
			t := strings.Split(v, ":")

			if len(t) == 2 {
				tags = append(tags, tag{
					Type:  t[0],
					Value: t[1],
				})
			} else {
				tags = append(tags, tag{
					Type:  tagTypeFunc,
					Value: v,
				})
			}
		}
	}

	return tags
}

type fieldReadData struct {
	Ignore            bool
	Length            int64
	OffsetFromCurrent *int64
	OffsetFromStart   *int64
	OffsetFromEnd     *int64
	FuncName          string

	ElemFieldData *fieldReadData // if type Element
}

func parseReadDataFromTags(structValue reflect.Value, tags []tag) (*fieldReadData, error) {
	parseValue := func(v string) (int64, error) {
		l, err := strconv.ParseInt(v, 10, 0)
		if err != nil {
			lenVal := structValue.FieldByName(v)
			switch lenVal.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				l = lenVal.Int()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				l = int64(lenVal.Uint())
			default:
				return 0, errors.New("can't get field len from " + v + " field")
			}
		}
		return l, nil
	}

	var data fieldReadData
	var err error
	for _, t := range tags {
		switch t.Type {
		case tagTypeIgnore:
			return &fieldReadData{Ignore: true}, nil

		case tagTypeLength:
			data.Length, err = parseValue(t.Value)

		case tagTypeOffsetFromCurrent:
			var offset int64
			offset, err = parseValue(t.Value)
			data.OffsetFromCurrent = &offset

		case tagTypeOffsetFromStart:
			var offset int64
			offset, err = parseValue(t.Value)
			data.OffsetFromStart = &offset

		case tagTypeOffsetFromEnd:
			var offset int64
			offset, err = parseValue(t.Value)
			data.OffsetFromEnd = &offset

		case tagTypeFunc:
			data.FuncName = t.Value

		case tagTypeElement:
			data.ElemFieldData, err = parseReadDataFromTags(structValue, t.ElemTags)
		}

		if err != nil {
			return nil, err
		}
	}

	return &data, nil
}
