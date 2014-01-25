//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package examples

import (
    "testing"
    //"github.com/prataprc/goparsec"
    "reflect"
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
    var refs = [][2]interface{}{
        {`-10000`, -10000},
        {`-10.11231`, -10.11231},
        {`"hello world"`, "hello world"},
        {`true`, true},
        {`false`, false},
        {`null`, nil},
        {`[79.1253026404128,null,"hello world", false, 10]`,
            []interface{}{79.1253026404128, nil, "hello world", false, 10},
        },
        {jsonText, jsonVal},
    }
    for _, x := range refs {
        v := Value(JsonParse([]byte(x[0].(string))))
        if !reflect.DeepEqual(v, x[1]) {
            t.Fatalf("Mismatch for %v %v\n", v, x[1])
        }
    }
}

func BenchmarkJsonInt(b *testing.B) {
    for i := 0; i < b.N; i++ {
        JsonParse([]byte(`10000`))
    }
}

func BenchmarkJsonFloat(b *testing.B) {
    for i := 0; i < b.N; i++ {
        JsonParse([]byte(`-10.11231`))
    }
}

func BenchmarkJsonString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        JsonParse([]byte(`"hello world"`))
    }
}

func BenchmarkJsonBool(b *testing.B) {
    for i := 0; i < b.N; i++ {
        JsonParse([]byte(`true`))
    }
}

func BenchmarkJsonNull(b *testing.B) {
    for i := 0; i < b.N; i++ {
        JsonParse([]byte(`null`))
    }
}

func BenchmarkJsonArray(b *testing.B) {
    for i := 0; i < b.N; i++ {
        JsonParse([]byte(`[79.1253026404128,null,"hello world", false, 10]`))
    }
}

func BenchmarkJsonMap(b *testing.B) {
    for i := 0; i < b.N; i++ {
        JsonParse([]byte(jsonText))
    }
}
