//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not
//  use this file except in compliance with the License. You may obtain a copy
//  of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package lib

import (
	"encoding/json"
	"fmt"
	gojson "github.com/dustin/gojson"
	"io/ioutil"
	"reflect"
	"testing"
)

var jsonText = `{ "inelegant":27.53096820876087,
"horridness":true,
"iridodesis":[79.1253026404128,null,"hello world", false, 10],
"arrogantness":null,
"unagrarian":false
}`

var jsonVal = map[string]interface{}{
	"inelegant":    27.53096820876087,
	"horridness":   true,
	"iridodesis":   []interface{}{79.1253026404128, nil, "hello world", false, 10},
	"arrogantness": nil,
	"unagrarian":   false,
}

func TestJson(t *testing.T) {
	var largeVal []interface{}
	var smallMap map[string]interface{}

	largeText, err := ioutil.ReadFile("./large.json")
	if err != nil {
		t.Fatal(err)
	}
	json.Unmarshal(largeText, &largeVal)

	smallText := []byte(jsonText)
	json.Unmarshal(smallText, &smallMap)

	var refs = [][2]interface{}{
		{[]byte(`-10000`), -10000.0},
		{[]byte(`-10.11231`), -10.11231},
		{[]byte(`"hello world"`), "hello world"},
		{[]byte(`true`), true},
		{[]byte(`false`), false},
		{[]byte(`null`), nil},
		{[]byte(`[79.1253026404128,null,"hello world", false, 10]`),
			[]interface{}{79.1253026404128, nil, "hello world", false, 10.0},
		},
		{smallText, smallMap},
		{largeText, largeVal},
	}
	for _, x := range refs {
		v := Value(JSONParse(x[0].([]byte)))
		if !reflect.DeepEqual(v, x[1]) {
			fmt.Println(v)
			fmt.Println(x[1])
			t.Fatal("Mismatch")
		}
	}
}

func BenchmarkJSONInt(b *testing.B) {
	text := []byte(`10000`)
	for i := 0; i < b.N; i++ {
		JSONParse(text)
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkJSONFloat(b *testing.B) {
	text := []byte(`-10.11231`)
	for i := 0; i < b.N; i++ {
		JSONParse(text)
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkJSONString(b *testing.B) {
	text := []byte(`"hello world"`)
	for i := 0; i < b.N; i++ {
		JSONParse(text)
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkJSONBool(b *testing.B) {
	text := []byte(`true`)
	for i := 0; i < b.N; i++ {
		JSONParse(text)
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkJSONNull(b *testing.B) {
	text := []byte(`null`)
	for i := 0; i < b.N; i++ {
		JSONParse(text)
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkJSONArray(b *testing.B) {
	text := []byte(`[79.1253026404128,null,"hello world", false, 10]`)
	for i := 0; i < b.N; i++ {
		Value(JSONParse(text))
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkJSONMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Value(JSONParse([]byte(jsonText)))
	}
	b.SetBytes(int64(len(jsonText)))
}

func BenchmarkJSONMedium(b *testing.B) {
	text, _ := ioutil.ReadFile("./medium.json")
	for i := 0; i < b.N; i++ {
		Value(JSONParse(text))
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkJSONLarge(b *testing.B) {
	text, _ := ioutil.ReadFile("./large.json")
	for i := 0; i < b.N; i++ {
		Value(JSONParse(text))
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkEncJSONArray(b *testing.B) {
	var val []interface{}

	text := []byte(`[79.1253026404128,null,"hello world", false, 10]`)
	for i := 0; i < b.N; i++ {
		json.Unmarshal(text, &val)
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkEncJSONMap(b *testing.B) {
	var val map[string]interface{}

	for i := 0; i < b.N; i++ {
		json.Unmarshal([]byte(jsonText), &val)
	}
	b.SetBytes(int64(len(jsonText)))
}

func BenchmarkEncJSONMedium(b *testing.B) {
	var val []interface{}

	text, _ := ioutil.ReadFile("./medium.json")
	for i := 0; i < b.N; i++ {
		json.Unmarshal(text, &val)
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkEncJSONLarge(b *testing.B) {
	var val []interface{}

	text, _ := ioutil.ReadFile("./large.json")
	for i := 0; i < b.N; i++ {
		json.Unmarshal(text, &val)
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkGoJSONArray(b *testing.B) {
	var val []interface{}

	text := []byte(`[79.1253026404128,null,"hello world", false, 10]`)
	for i := 0; i < b.N; i++ {
		gojson.Unmarshal(text, &val)
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkGoJSONMap(b *testing.B) {
	var val map[string]interface{}

	for i := 0; i < b.N; i++ {
		gojson.Unmarshal([]byte(jsonText), &val)
	}
	b.SetBytes(int64(len(jsonText)))
}

func BenchmarkGoJSONMedium(b *testing.B) {
	var val []interface{}

	text, _ := ioutil.ReadFile("./medium.json")
	for i := 0; i < b.N; i++ {
		gojson.Unmarshal(text, &val)
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkGoJSONLarge(b *testing.B) {
	var val []interface{}

	text, _ := ioutil.ReadFile("./large.json")
	for i := 0; i < b.N; i++ {
		gojson.Unmarshal(text, &val)
	}
	b.SetBytes(int64(len(text)))
}
