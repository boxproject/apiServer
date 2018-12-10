// Copyright 2018. bolaxy.org authors.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
// 		 http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package utils

import (
	"testing"
	//. "github.com/smartystreets/goconvey/convey"
	"fmt"
)

func TestRsaVerify(t *testing.T) {
	msg := "1540522499021"
	sign := "P2NZsSr2tT+9WenOLvmoGjsHcn3XThyGC0GkfnCij/aG/MIDo8fiuxacZVBIsshilPa+iccOqcjPnti4O4SkKbeeYtmgJhrHdM4o8eVDVx7L2zTcConSN6MwZaBK+7pO2Sjb8kNdK3mS33tTH133BMVp4dRsu53fAj/PwpGslkZfjRFOsKfm+ZYsUFq1TL4nbohP7c/oqmHWupZs3noaPDT03WB/ZRw/ejHp9GUjAPbFeOTtVuJpCNgxFoh/awTsjdEJ0ZR+5/BIzlkeKLemCD9lfkVsqiONwHKeXfD94c8GqK0M5K8xEyDSQDYeN21Pg3/2/5M51ZoreVXxopPw4Q=="
	pubkey := "MIIBCgKCAQEAlU5agCX/It+2I7LqA6wFO9RBLBM6kKH+l1kqFI7PSN7JgnsgAGRp0NdT20FGYfturRpbd5W4Wu0B4cUg0JzgjQ5okbIohMZ0T+JfVWrECPPuj2i5Mtv6U3LnW/TuXnZwJ9tcgADdYPb2XI9lhOrpQU60NPB3idNrlT53SwaqLX4GK7ZmU2ePT7lfdQPQWL/N3b+pgi89VEz8lvgoIBkk4WI3tkxnHgRYydaj0Est5rfnp4eN1y8FWOQZPj58qjHWggfcwLSpuWLnDCB8cGA3q3o5mGm+cw+O54Ie66CLEsDI4nNllc5/vRYIraqpF02TemJw+8O2CfzkB6uFoHKMQQIDAQAB"
	err := RsaVerify([]byte(pubkey), []byte(msg), []byte(sign))
	fmt.Println("222", err)
}

func TestRsaSign(t *testing.T) {
	msg := "1540522499021"
	privte_key := "MIIEpAIBAAKCAQEAlU5agCX/It+2I7LqA6wFO9RBLBM6kKH+l1kqFI7PSN7JgnsgAGRp0NdT20FGYfturRpbd5W4Wu0B4cUg0JzgjQ5okbIohMZ0T+JfVWrECPPuj2i5Mtv6U3LnW/TuXnZwJ9tcgADdYPb2XI9lhOrpQU60NPB3idNrlT53SwaqLX4GK7ZmU2ePT7lfdQPQWL/N3b+pgi89VEz8lvgoIBkk4WI3tkxnHgRYydaj0Est5rfnp4eN1y8FWOQZPj58qjHWggfcwLSpuWLnDCB8cGA3q3o5mGm+cw+O54Ie66CLEsDI4nNllc5/vRYIraqpF02TemJw+8O2CfzkB6uFoHKMQQIDAQABAoIBAAJ7H9PbTccFkqA7O9o9xIR+/Wo/E40NtA5Nw/49dUJPuWd6pkh7Yqq+uTz+c63zIJ6tvnFZQov4PjzDKs1sneqaH7C0FipGpe4h49WmhYVkkEU/xzwKHbm/QgrsIf3d1/VBluLloMgYsnVNSLGRubiFF0d9626V3cbIF/FeNfu1guNQ4Y75sYem21xTtxxt0tToZllaKHAmMDfNNM7F5TJvoydQryK6Z5FdBBwLY01Ck433Nr1vtIrcifi+DeH5mtHzEc8rWoyz6+Lv9r4gfnlP0zAU9my0KlrangDmBvJsFaOzHC88LpXc/zFSwn0a+amWULLFHUEvIxcjNGjH9AECgYEAxae4AdQVvF6XkIGZpBj4eJizm2dzc+VwworPQIp9rJ2JQQN6icbnYKgNDXq573QafLMwq2LqY1Hdy+6oMYjEUcTU0SQjWl8qa89mS8HduAzr3qrNxZh9Dqof8ToEh4iZuV7f1MMyzqUi+QHZnnFI7BVfOK7XOLJ0NnA8A4l90mECgYEAwWEFeXP6BWcM4i8bUZe2lM3O39oz9gaSy8xamRK9F84cwzpDTSP5ZyS0ygmFLG3xJi8p0Zxcnamt864VFrOIIWfkWdYfmSdVQelPTguLhnt1jfZFlFGePbVHpe649vmSeg4Ypanaf/OcTeQ/GjIEANO9VapLOGTXXV2JeSj4xeECgYAdm38Dvxo2alD304IJQ7hMkEsvNzLjJHZcneBnwZcLuVlrBLNhWgskvmeeIkkh5lllXo4mzh1gHU5FEw3cxajqurpKTciB7Al1ts6TAIpO3JikdR93vtzUyoUYZGFzT/H/A4gx3b+JltEDTdSkWEHdI2JtIjuZAZpI4U0MFpdjYQKBgQCYOo87r9/HNvs+ZwTjc0Ho3CcZs0UqTLxssH81dzniLoAX85qddE9WMeAcF+h9NEMc4w3Rk3yZJMTqSkURrNziJ03sppITV5JXI5opfw8kG7ZLve3CN4oRW85+QnHbAlabvNMMPqziAt0tuBswvOTH3edzM26pg0DCn+qjtWw5IQKBgQCGmMQ7fYxOltO1Y1D+HSUzSDpHHzSqO5xvfmJ+/w9SqFm02KT01264hGIFuQtXw18Ql6yML0Po0lE1z3fsEmcWX04+Pdsi+y6k+cdx0UFzs6ncCuWN0X0Xc4ii9s+9qPANUp/iby/DO752yS3NHqWdvDLsPsE1fQvDhJ/4yAyGXQ=="
	sign, _ := RsaSign([]byte(privte_key), []byte(msg))
	fmt.Println("sign", string(sign))
}
