// Copyright 2018. bolaxy.org authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package middleware

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"fmt"
)

func TestSortInfo(t *testing.T) {

	//var testParam = "中文chineseChinese"
	Convey("TestLoadConfigFile", t, func(){
		Convey("Should return true when config is correct.", func(){
			//data := SortByValue(testParam)
			fmt.Println("&&&&&&&&&", string(111))
			//So(data, ShouldNotBeNil)
		})
	})
}
