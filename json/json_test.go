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

package json

import "encoding/json"
import "io/ioutil"
import "reflect"
import "testing"

var jsonText = []byte(`{ "inelegant":27.53096820876087,
"horridness":true,
"iridodesis":[79.1253026404128,null,"hello world", false, 10],
"arrogantness":null,
"unagrarian":false
}`)

var jsonVal = map[string]interface{}{
	"inelegant":  27.53096820876087,
	"horridness": True("true"),
	"iridodesis": []interface{}{
		79.1253026404128, Null("null"), "hello world", False("false"), 10},
	"arrogantness": Null("null"),
	"unagrarian":   False("false"),
}

func TestJson(t *testing.T) {
	var largeVal []interface{}
	var smallMap map[string]interface{}

	largeText, err := ioutil.ReadFile("./../testdata/large.json")
	if err != nil {
		t.Fatal(err)
	}
	json.Unmarshal(largeText, &largeVal)

	smallText := jsonText
	json.Unmarshal(smallText, &smallMap)

	var refs = [][2]interface{}{
		{[]byte(`-10000`), Num("-10000")},
		{[]byte(`-10.11231`), Num("-10.11231")},
		{[]byte(`"hello world"`), String("hello world")},
		{[]byte(`true`), True("true")},
		{[]byte(`false`), False("false")},
		{[]byte(`null`), Null("null")},
		{[]byte(`[79.1253026404128,null,"hello world", false, 10]`),
			[]interface{}{
				Num("79.1253026404128"), Null("null"), String("hello world"),
				False("false"), Num("10")},
		},
	}
	for _, x := range refs {
		s := NewJSONScanner(x[0].([]byte))
		if v, _ := Y(s); !reflect.DeepEqual(v, x[1]) {
			t.Fatalf("Mismatch `%v`: %v vs %v", string(x[0].([]byte)), x[1], v)
		}
	}
	s := NewJSONScanner(smallText)
	if v, _ := Y(s); !reflect.DeepEqual(nativeValue(v), smallMap) {
		t.Fatalf("Mismatch `%v`: %v vs %v", string(smallText), smallMap, v)
	}
	s = NewJSONScanner(largeText)
	if v, _ := Y(s); !reflect.DeepEqual(nativeValue(v), largeVal) {
		t.Fatalf("Mismatch `%v`: %v vs %v", string(largeText), largeVal, v)
	}
}

func BenchmarkJSONInt(b *testing.B) {
	text := []byte(`10000`)
	for i := 0; i < b.N; i++ {
		Y(NewJSONScanner(text))
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkJSONFloat(b *testing.B) {
	text := []byte(`-10.11231`)
	for i := 0; i < b.N; i++ {
		Y(NewJSONScanner(text))
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkJSONString(b *testing.B) {
	text := []byte(`"hello world"`)
	for i := 0; i < b.N; i++ {
		Y(NewJSONScanner(text))
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkJSONBool(b *testing.B) {
	text := []byte(`true`)
	for i := 0; i < b.N; i++ {
		Y(NewJSONScanner(text))
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkJSONNull(b *testing.B) {
	text := []byte(`null`)
	for i := 0; i < b.N; i++ {
		Y(NewJSONScanner(text))
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkJSONArray(b *testing.B) {
	text := []byte(`[79.1253026404128,null,"hello world", false, 10]`)
	for i := 0; i < b.N; i++ {
		Y(NewJSONScanner(text))
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkJSONMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Y(NewJSONScanner(jsonText))
	}
	b.SetBytes(int64(len(jsonText)))
}

func BenchmarkJSONMedium(b *testing.B) {
	text, _ := ioutil.ReadFile("./../testdata/medium.json")
	for i := 0; i < b.N; i++ {
		Y(NewJSONScanner(text))
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkJSONLarge(b *testing.B) {
	text, _ := ioutil.ReadFile("./../testdata/large.json")
	for i := 0; i < b.N; i++ {
		Y(NewJSONScanner(text))
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
		json.Unmarshal(jsonText, &val)
	}
	b.SetBytes(int64(len(jsonText)))
}

func BenchmarkEncJSONMedium(b *testing.B) {
	var val []interface{}

	text, _ := ioutil.ReadFile("./../testdata/medium.json")
	for i := 0; i < b.N; i++ {
		json.Unmarshal(text, &val)
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkEncJSONLarge(b *testing.B) {
	var val []interface{}

	text, _ := ioutil.ReadFile("./../testdata/large.json")
	for i := 0; i < b.N; i++ {
		json.Unmarshal(text, &val)
	}
	b.SetBytes(int64(len(text)))
}
