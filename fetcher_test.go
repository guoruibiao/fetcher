package fetcher

import (
	"encoding/json"
	"testing"
)

var data string
var fetcher *Bucket

func init() {
	data = `{
    "whiteList": [
      "A390261C6FB8EC202B82FDCBE4D4FC8977692D3E2FHJGHIFBTP",
      "E14D19792551441F6804938A2C93F649CF4586E5AOSCCKOCMDM",
      "AB9727A254C4E755174D226E4FCBDD9E|823750740870568"
    ],
    "age": 26,
    "alive": true,
    "weather": 28.3,
    "switch": "on",
    "role":{
        "name": "worker",
        "privilege": null
    },
    "attachTypes": [
        {
            "name": "xm",
            "index": 1
        },
        {
            "name": "default",
            "index": 2
        }
    ],
    "address": {
        "countries":[
            {
                "name":"中国",
                "provinces":[
{
                        "name": "北京市",
                        "counties":[
                            {
                                "name": "北京市",
                                "detail": "江苏市海淀区上地街道66号"
                            }
                        ]
                    },
                    {
                        "name": "江苏省",
                        "counties":[
                            {
                                "name": "苏州市",
                                "detail": "江苏省苏州市姑苏区沧浪亭街3号"
                            }
                        ]
                    }
                ]
            }
        ]
    }
}`
	var tmp interface{}
	if err := json.Unmarshal([]byte(data), &tmp); err == nil {
		fetcher = NewFetcher(tmp)
	}else{
		panic("init fetcher failed" + err.Error())
	}

}


func TestBucket_Fetch(t *testing.T) {
	name := fetcher.Fetch("/address/countries/0/provinces/1/name")
	t.Log(name)

	detail, _ := fetcher.Fetch("/address/countries/0/provinces/1/counties/0/detail").String()
	t.Log(detail)
}

func TestBucket_Array(t *testing.T) {
	if arr, err := fetcher.Fetch("/whiteList").Array(); err != nil {
		t.Error(err)
	}else{
		t.Log(arr)
	}
}

func TestBucket_Bool(t *testing.T) {
	if alive, err := fetcher.Fetch("/alive").Bool(); err != nil {
		t.Error(err)
	}else{
		t.Log(alive)
	}
}

func TestBucket_Float(t *testing.T) {
	if weather, err := fetcher.Fetch("/weather").Float(); err != nil {
		t.Error(err)
	}else{
		t.Log(weather)
	}
}

func TestBucket_Int(t *testing.T) {
	if i, err := fetcher.Fetch("/age").Int(); err != nil {
		t.Error(err)
	}else{
		t.Log(i)
	}
}

func TestBucket_Int64(t *testing.T) {
	if i64, err := fetcher.Fetch("/age").Int64(); err != nil {
		t.Error(err)
	}else{
		t.Log(i64)
	}
}

func TestBucket_Map(t *testing.T) {
	if m, err := fetcher.Fetch("/role").Map(); err != nil {
		t.Error(err)
	}else{
		t.Log(m)
	}
}

func TestBucket_String(t *testing.T) {
	if s, err := fetcher.Fetch("/switch").String(); err != nil {
		t.Error(err)
	}else{
		t.Log(s)
	}
}