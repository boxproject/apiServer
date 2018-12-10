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
package accmanage

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	//"fmt"
	"fmt"
	"github.com/boxproject/apiServer/errors"
)

func TestTransferRecord(t *testing.T) {
	Convey("查询账户的转账记录列表", t, func() {
		msg := TransferRecord(2)
		fmt.Println("******", msg.Data)
		So(msg.Code, ShouldEqual, errors.Success)
	})
}
