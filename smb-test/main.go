package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type TagMap struct {
	m   map[string]interface{}
	has map[string]bool
}

type Metadata struct {
	Tags       *TagMap
	Lens       map[string]uint64
	Offsets    map[string]uint64
	Parent     interface{}
	ParentBuf  []byte
	CurrOffset uint64
	CurrField  string
}

func (t *TagMap) Has(key string) bool {
	return t.has[key]
}

func (t *TagMap) Set(key string, val interface{}) {
	t.m[key] = val
	t.has[key] = true
}

func (t *TagMap) Get(key string) interface{} {
	return t.m[key]
}

func (t *TagMap) GetInt(key string) (int, error) {
	if !t.Has(key) {
		return 0, errors.New("Key does not exist in tag")
	}
	return t.Get(key).(int), nil
}

func (t *TagMap) GetString(key string) (string, error) {
	if !t.Has(key) {
		return "", errors.New("Key does not exist in tag")
	}
	return t.Get(key).(string), nil
}

func parseTags(sf reflect.StructField) (*TagMap, error) {
	ret := &TagMap{
		m:   make(map[string]interface{}),
		has: make(map[string]bool),
	}
	tag := sf.Tag.Get("smb")
	smbTags := strings.Split(tag, ",")
	for _, smbTag := range smbTags {
		tokens := strings.Split(smbTag, ":")
		switch tokens[0] {
		case "len", "offset", "count":
			if len(tokens) != 2 {
				return nil, errors.New("Missing required tag data. Expecting key:val")
			}
			ret.Set(tokens[0], tokens[1])
		case "fixed":
			if len(tokens) != 2 {
				return nil, errors.New("Missing required tag data. Expecting key:val")
			}
			i, err := strconv.Atoi(tokens[1])
			if err != nil {
				return nil, err
			}
			ret.Set(tokens[0], i)
		}
	}
	return ret, nil
}

func unmarshal(buf []byte, v interface{}, meta *Metadata) (interface{}, error) {
	typev := reflect.TypeOf(v)
	valuev := reflect.ValueOf(v)
	r := bytes.NewBuffer(buf)
	switch typev.Kind() {
	case reflect.Struct:
		m := &Metadata{
			Tags:       &TagMap{},
			Lens:       make(map[string]uint64),
			Parent:     v,
			ParentBuf:  buf,
			Offsets:    make(map[string]uint64),
			CurrOffset: 0,
		}
		for i := 0; i < typev.NumField(); i++ {
			m.CurrField = typev.Field(i).Name
			tags, err := parseTags(typev.Field(i))
			if err != nil {
				return nil, err
			}
			m.Tags = tags
			var data interface{}
			switch typev.Field(i).Type.Kind() {
			case reflect.Struct:
				data, err = unmarshal(buf[m.CurrOffset:], valuev.Field(i).Addr().Interface(), m)
			default:
				data, err = unmarshal(buf[m.CurrOffset:], valuev.Field(i).Interface(), m)
			}
			if err != nil {
				return nil, err
			}
			valuev.Field(i).Set(reflect.ValueOf(data))
		}
		v = reflect.Indirect(reflect.ValueOf(v)).Interface()
		meta.CurrOffset += m.CurrOffset
		return v, nil
	case reflect.Uint16:
		var ret uint16
		if err := binary.Read(r, binary.LittleEndian, &ret); err != nil {
			return nil, err
		}
		if meta.Tags.Has("len") {
			ref, err := meta.Tags.GetString("len")
			if err != nil {
				return nil, err
			}
			meta.Lens[ref] = uint64(ret)
		}
		meta.CurrOffset += uint64(binary.Size(ret))
		return ret, nil
	case reflect.Slice, reflect.Array:
		switch typev.Elem().Kind() {
		case reflect.Uint8:
			var length, offset int
			var err error
			if meta.Tags.Has("fixed") {
				if length, err = meta.Tags.GetInt("fixed"); err != nil {
					return nil, err
				}
				// Fixed length fields advance current offset
				meta.CurrOffset += uint64(length)
			} else {
				if val, ok := meta.Lens[meta.CurrField]; ok {
					length = int(val)
				} else {
					return nil, errors.New("Variable length field missing length reference in struct")
				}
				if val, ok := meta.Offsets[meta.CurrField]; ok {
					offset = int(val)
				} else {
					// No offset found in map. Use current offset
					offset = int(meta.CurrOffset)
				}
				// Variable length data is relative to parent/outer struct.
				// Reset reader to point to beginning of data
				r = bytes.NewBuffer(meta.ParentBuf[offset : offset+length])
				// Variable length data fields do NOT advance current offset.
			}
			data := make([]byte, length)
			if err := binary.Read(r, binary.LittleEndian, &data); err != nil {
				return nil, err
			}
			return data, nil
		}
	default:
		return errors.New("Unmarshal not implemented for kind:" + typev.Kind().String()), nil
	}
	return nil, nil
}

type user struct {
	age int32
	gender int32
	some []int32
}

func main() {
	str := int32(123)
	f, err := os.OpenFile("lalala.lyz", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic("open file failed")
	}
	defer f.Close()
	var binBuf bytes.Buffer
	err = binary.Write(&binBuf, binary.LittleEndian, str)
	f.Write(binBuf.Bytes())
	if err != nil {
		log.Fatalln(err)
	}
}
