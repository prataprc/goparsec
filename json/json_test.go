//  Copyright (c) 2013 Couchbase, Inc.

package json

import "encoding/json"
import "io/ioutil"
import "fmt"
import "reflect"
import "testing"

import "github.com/prataprc/goparsec"

var _ = fmt.Sprintf("dummy print")

var jsonText = []byte(`{ "inelegant":27.53096820876087,
"horridness":true,
"iridodesis":[79.1253026404128,null,"hello world", false, 10],
"arrogantness":null,
"unagrarian":false
}`)

var jsonVal = map[string]interface{}{
	"inelegant":  27.53096820876087,
	"horridness": True("true"),
	"iridodesis": []parsec.ParsecNode{
		79.1253026404128, Null("null"), "hello world", False("false"), 10},
	"arrogantness": Null("null"),
	"unagrarian":   False("false"),
}

func TestJson(t *testing.T) {
	var largeVal []interface{}
	var mediumVal []interface{}
	var smallMap map[string]interface{}

	largeText, err := ioutil.ReadFile("./../testdata/large.json")
	if err != nil {
		t.Fatal(err)
	}
	json.Unmarshal(largeText, &largeVal)

	mediumText, err := ioutil.ReadFile("./../testdata/medium.json")
	if err != nil {
		t.Fatal(err)
	}
	json.Unmarshal(mediumText, &mediumVal)

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
			[]parsec.ParsecNode{
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
	s = NewJSONScanner(mediumText)
	if v, _ := Y(s); !reflect.DeepEqual(nativeValue(v), mediumVal) {
		t.Fatalf("Mismatch `%v`: %v vs %v", string(mediumText), mediumVal, v)
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

func BenchmarkEncJSONInt(b *testing.B) {
	var val interface{}
	text := []byte(`10000`)
	for i := 0; i < b.N; i++ {
		json.Unmarshal(text, &val)
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkEncJSONFloat(b *testing.B) {
	var val interface{}
	text := []byte(`-10.11231`)
	for i := 0; i < b.N; i++ {
		json.Unmarshal(text, &val)
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkEncJSONString(b *testing.B) {
	var val interface{}
	text := []byte(`"hello world"`)
	for i := 0; i < b.N; i++ {
		json.Unmarshal(text, &val)
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkEncJSONBool(b *testing.B) {
	var val interface{}
	text := []byte(`true`)
	for i := 0; i < b.N; i++ {
		json.Unmarshal(text, &val)
	}
	b.SetBytes(int64(len(text)))
}

func BenchmarkEncJSONNull(b *testing.B) {
	var val interface{}
	text := []byte(`null`)
	for i := 0; i < b.N; i++ {
		json.Unmarshal(text, &val)
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
