package binstruct

import (
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
)

const (
	tagName = "bin"
)

const (
	tagTypeEmpty   = ""
	tagTypeIgnore  = "-"
	tagTypeFunc    = "func"
	tagTypeElement = "elem"

	tagTypeOrderLE = "le"
	tagTypeOrderBE = "be"

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

func parseTag(t string) ([]tag, error) {
	var tags []tag

	for {
		var v string

		index := strings.Index(t, ",")
		switch {
		case index == -1:
			v = t
		default:
			v = t[:index]
			t = t[index+1:]
		}

		v = strings.TrimSpace(v)

		switch {
		case v == tagTypeEmpty:
			// Just skip

		case v == tagTypeIgnore:
			tags = append(tags, tag{Type: tagTypeIgnore})

		case strings.HasPrefix(v, "["):
			v = v + "," + t
			var arrBalance int
			var closeIndex int
			for {
				in := v[closeIndex:]
				idx := strings.IndexAny(in, "[]")
				closeIndex += idx

				if idx == -1 {
					return nil, errors.New("unbalanced square bracket")
				}

				switch in[idx] {
				case '[':
					arrBalance--
				case ']':
					arrBalance++
				}

				closeIndex++

				if arrBalance == 0 {
					break
				}
			}

			t = v[closeIndex:]
			v = v[1 : closeIndex-1]

			pt, err := parseTag(v)
			if err != nil {
				return nil, err
			}

			tags = append(tags, tag{Type: tagTypeElement, ElemTags: pt})

		case v == tagTypeOrderLE:
			tags = append(tags, tag{Type: tagTypeOrderLE})

		case v == tagTypeOrderBE:
			tags = append(tags, tag{Type: tagTypeOrderBE})

		default:
			ts := strings.Split(v, ":")

			if len(ts) == 2 {
				tags = append(tags, tag{
					Type:  ts[0],
					Value: ts[1],
				})
			} else {
				tags = append(tags, tag{
					Type:  tagTypeFunc,
					Value: v,
				})
			}
		}

		if index == -1 {
			return tags, nil
		}
	}
}

type fieldOffset struct {
	Offset int64
	Whence int
}

type fieldReadData struct {
	Ignore   bool
	Length   *int64
	Offsets  []fieldOffset
	FuncName string
	Order    binary.ByteOrder

	ElemFieldData *fieldReadData // if type Element
}

func parseCalc(v string) (nums, ops []string) {
	cur := v
	for {
		idx := strings.IndexAny(cur, "+-/*")
		if idx == -1 {
			nums = append(nums, cur)
			break
		}

		nums = append(nums, cur[:idx])
		ops = append(ops, string(cur[idx]))

		cur = cur[idx+1:]
	}

	return nums, ops
}

func parseValue(structValue reflect.Value, v string) (int64, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, nil
	}

	// calculate
	mathIndex := strings.IndexAny(v, "+-/*")
	if mathIndex != -1 {
		nums, ops := parseCalc(v)

		var result int64
		for k := range nums {
			n, err := parseValue(structValue, nums[k])
			if err != nil {
				return 0, err
			}

			if k == 0 {
				result = n
				continue
			}

			switch ops[k-1] {
			case "+":
				result += n
			case "-":
				result -= n
			case "/":
				result /= n
			case "*":
				result *= n
			}
		}

		return result, nil
	}

	// parse value or get from field
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

func parseReadDataFromTags(structValue reflect.Value, tags []tag) (*fieldReadData, error) {
	var data fieldReadData
	var err error
	for _, t := range tags {
		switch t.Type {
		case tagTypeIgnore:
			return &fieldReadData{Ignore: true}, nil

		case tagTypeLength:
			var length int64
			length, err = parseValue(structValue, t.Value)
			data.Length = &length

		case tagTypeOffsetFromCurrent:
			var offset int64
			offset, err = parseValue(structValue, t.Value)
			data.Offsets = append(data.Offsets, fieldOffset{
				Offset: offset,
				Whence: io.SeekCurrent,
			})

		case tagTypeOffsetFromStart:
			var offset int64
			offset, err = parseValue(structValue, t.Value)
			data.Offsets = append(data.Offsets, fieldOffset{
				Offset: offset,
				Whence: io.SeekStart,
			})

		case tagTypeOffsetFromEnd:
			var offset int64
			offset, err = parseValue(structValue, t.Value)
			data.Offsets = append(data.Offsets, fieldOffset{
				Offset: offset,
				Whence: io.SeekEnd,
			})

		case tagTypeFunc:
			data.FuncName = t.Value

		case tagTypeElement:
			data.ElemFieldData, err = parseReadDataFromTags(structValue, t.ElemTags)

		case tagTypeOrderLE:
			data.Order = binary.LittleEndian

		case tagTypeOrderBE:
			data.Order = binary.BigEndian
		}

		if err != nil {
			return nil, err
		}
	}

	return &data, nil
}
