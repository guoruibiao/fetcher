package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	bucketCheckOK = iota
	bucketIsNil
	bucketDataIsNil
)

type Bucket struct {
	Data interface{}
}

func NewFetcher(data interface{}) *Bucket{
	return &Bucket{
		Data: data,
	}
}

func (b *Bucket) check() int {
	if b.Data == nil {
		return bucketIsNil
	}else if b.Data == nil {
		return bucketDataIsNil
	}

	return bucketCheckOK
}

func (b *Bucket) Array() ([]interface{}, error) {
	if errno := b.check(); errno != bucketCheckOK {
		return nil, errors.New("Bucket check failed, errno=" + fmt.Sprintf("%d", errno))
	}

	v := reflect.ValueOf(b.Data)
	if k := v.Kind(); k != reflect.Array && k != reflect.Slice {
		return nil, errors.New("Bucket.Data is not Array or Slice")
	}

	arr := make([]interface{}, v.Len())
	for idx:=0; idx<v.Len(); idx++ {
		arr[idx] = v.Index(idx).Interface()
	}

	return arr, nil
}

func (b *Bucket) Map() (map[string]interface{}, error) {
	if errno := b.check(); errno != bucketCheckOK {
		return nil, errors.New("Bucket check failed, errno=" + fmt.Sprintf("%d", errno))
	}

	if m, ok := b.Data.(map[string]interface{}); ok {
		return m, nil
	}
	return nil, errors.New("Bucket.Data is not map")
}

func (b *Bucket)String() (string, error) {
	if errno := b.check(); errno != bucketCheckOK {
		return "", errors.New("Bucket check failed, errno=" + fmt.Sprintf("%d", errno))
	}

	if s, ok := b.Data.(string); ok {
		return s, nil
	}
	return "", errors.New("reflect to type string failed")
}

func (b *Bucket) Bool() (bool, error) {
	if errno := b.check(); errno != bucketCheckOK {
		return false, errors.New("Bucket check failed, errno=" + fmt.Sprintf("%d", errno))
	}
	if b, ok := b.Data.(bool); ok {
		return b, nil
	}
	return false, errors.New("reflect to type bool failed")
}

func (b *Bucket)Int() (int, error) {
	if errno := b.check(); errno != bucketCheckOK {
		return 0, errors.New("Bucket check failed, errno=" + fmt.Sprintf("%d", errno))
	}

	switch b.Data.(type) {
	case json.Number:
		i, err := b.Data.(json.Number).Int64()
		return int(i), err
	case float32, float64:
		return int(reflect.ValueOf(b.Data).Float()), nil
	case int, int8, int16, int32, int64:
		return int(reflect.ValueOf(b.Data).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return int(reflect.ValueOf(b.Data).Uint()), nil
	}

    return 0, errors.New("reflect to type int failed")
}

func (b *Bucket)Int64() (int64, error) {
	if errno := b.check(); errno != bucketCheckOK {
		return 0, errors.New("Bucket check failed, errno=" + fmt.Sprintf("%d", errno))
	}

	switch b.Data.(type) {
	case json.Number:
		return b.Data.(json.Number).Int64()
	case float32, float64:
		return int64(reflect.ValueOf(b.Data).Float()), nil
	case int, int8, int16, int32, int64:
		return int64(reflect.ValueOf(b.Data).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(b.Data).Uint()), nil
	}

	return 0, errors.New("reflect to type int failed")
}

func (b *Bucket)Float() (float64, error) {
	if errno := b.check(); errno != bucketCheckOK {
		return 0, errors.New("Bucket check failed, errno=" + fmt.Sprintf("%d", errno))
	}

	switch b.Data.(type) {
	case json.Number:
		return b.Data.(json.Number).Float64()
	case float32, float64:
		return reflect.ValueOf(b.Data).Float(), nil
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(b.Data).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(b.Data).Uint()), nil
	}

	return 0, errors.New("reflect to type float failed")
}


// Fetch : works like xpath
func (b *Bucket) Fetch(path string) *Bucket {
	f := &Bucket{
		Data: b.Data,
	}
	paths := strings.Split(path, "/")
	for _, p := range paths {
		if p == "" {
			continue
		}

		if d, ok := f.Data.(map[string]interface{}); ok && len(d) > 0 {
			f.Data = d[p]
		}else if arr, err := f.Array(); err == nil {
			// get the item by index
			if idx, e := strconv.Atoi(p); e == nil && idx < len(arr) {
				f.Data = arr[idx]
			}else {
				return nil
			}
		}else{
			return nil
		}
	}
	return f
}
