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
package template

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/boxproject/apiServer/middleware"
	"github.com/boxproject/apiServer/utils"
	"time"
	"encoding/hex"
	"fmt"
	"github.com/boxproject/apiServer/common"
)

var tc = `{
  "name": "单人审批",
  "period": 0,
  "approvalInfo": [
    {
      "require" : 1,
      "approvers" : [
        {
          "account": "employee_2-1",
          "pubkey": "MIIBCgKCAQEAyRjJ2hJWv4uZpqcmpRKrz+KWhsF0FGc9oWFbSOr6sp+hrOqy0Ezb+osTGsvqQlvKTPZvLgB0tcsmu3z3IktwrirMn/q29iIBFLSdE5cKME2Hm/S6Z/oZYpEaz5ss56ADAXZT8FlB72NYOnPDcOXtiHRIDS2JGeyQtxAGSOJDIcGw+8qlIeKP0XpK4m8eQwlab672pQ3K55G5NtBg2j8OJawo0uswNMmKY/FoB6elrLUZpO9fObzCCHSZjTPOXi279uLI5u8ZnD8LmQ9T7q7GmoDn8Ds5rRpXT2JzNV4w8a7cz8oG9WXRhJGzj+6ujYXiJJLUFGffwacMh4bp8wcfVwIDAQAB"
        }
      ]
    }
  ],
  "limitInfo" : [
    {
      "tokenAddress": "0xe1A178B681BD05964d3e3Ed33AE731577d9d96dD",
      "symbol": "BOX",
      "name": "BOX Token",
      "precise": 18,
      "limit": "10000"
    },
	{
      "tokenAddress": "",
      "symbol": "ETH",
      "name": "以太坊",
      "precise": 18,
      "limit": "100"
    },
	{
      "tokenAddress": "",
      "symbol": "BTC",
      "name": "比特币",
      "precise": 8,
      "limit": "10"
    },
	{
      "tokenAddress": "",
      "symbol": "LTC",
      "name": "莱特币",
      "precise": 18,
      "limit": "200"
    },
	{
      "tokenAddress": "",
      "symbol": "USDT",
      "name": "泰达币",
      "precise": 0,
      "limit": "10000"
    }
  ]
}`

func TestCreateNewTemplate(t *testing.T) {

	Convey("TestTemplateOperation", t, func() {
		Convey("Should return false when param is not full.", func() {
			private_key := "MIIEowIBAAKCAQEAovUhmHbdWNs6eVI/4YS3cgwNdF0sVs2UxLGv90JHlPzlVR7vHliEsE7e+YHwqHlHcDFHoXBdp3n2OYfYFCrVAfD6zwyEuCA5GYYHKX5jZyO6wTsMdaIXVFAZgESENfHx2FZgVZnqe13TItg1S1jNoHHVNfaR44n0KrJSPTzBypIdDBeHYorPjXrvrZkQnE8g3NMyKne+gdn16YcKl7RWNaUxoQ3If6mgeidVj1MEYbH7Vfy3dIruf/l5GRvWjveMP8AywiUpP8CEG+9M5sjFB/pXOV67s9YiExalDPEzpQWCLz9op07iILwZaopAsQwOAfoO7jLtboLo6tcAj7X/9QIDAQABAoIBAAX4+gw+fwpcrp33t8OqO5cAfhW/vHpJ+qPi51Imqbz6L+WYxqbUE7jqix2V954VI9sm0ztFhQk4XR/qrK7AiyIRFQA1sz+UXsNiuCx7J2WGV7fxVBmToUtDzEt5N2dKwNRvBgTdKUzwOWbTidCDJrY0td3OdbZjPgG6m20HJwsnlM3SXFKBssOLEr1aL4ieHcelQ8sXLeTFU+hOlCP/CSPDEMa2GwvWC0Aa6ihCUu2J00zXNUPlVvLmz96vBxuP+MQgLO0nwjHQW+TObdQxnFp2UjwyLhE6l/qCXeZ6dAfUg4TFFXo9m5lqmCJlIUEe24qytb0Q6SWkMKpD5Ami+JkCgYEA2E3yr7RjR9MY7DCGXf6zIwGRg2FBe/xqq9tY+2AHdebVtpWhqToS1FEG6YGoZzO3EOPqXfiUdqbOwpuqxdTDIhscy2iWEmmSlpm+XVpcUCdAQQUiJT0HVnD7qKyvpy0RTxV6ry6OWWXKxnR7+Kqip7CVNLUO5Ua297aYdEn/09MCgYEAwNztaBuAR9y34p4DeCAZElhEg3o9nYtF3WKJqy+pxV1/EEcDD+5KzFuggC2WqfvO6pZ3llRAV9lnzmrNK7w4pImS4Ot/Yo1vLz/dhi7odHGZ025iScqB7kKTKXmVelH0RN5ya4dYGtiM+cGyqWkXQ570r52uWOHzxQhGZwU7KBcCgYBX6XZVSy9+paFffWlfEoGkHKMTjpea2MBSENhdcm4V0TfS+YW2zc+RU01H4labez1QNBGxF7LfdnRaTAJgXVThq7tMQLmdyiK16s6gCnWybgGDa56xG9i7nTfrGiRaAcsjJUuzn2xjkSeU+QrZyiBJn09FGMFxFgiPLTtRKDI92QKBgBhXyl9pmTd41Hz2FNoEsyVtnrg5pa1M9vSKi3Xf/j27H16el/Raz3Yb2pZTKsEp12QkudNvru9nsYKMWHk8uKmL884P63Q1BeOg3AUjxpNsA058kEtgFCZOoOSLRCK8VWib6zVHUAeTbbfYiwR3D0ipt5iy/l6ZpnsoIDrAnkbhAoGBAJKebsuGx/YfWuvCRnHPAvPL3yHtrVxRuYhlZIKayIDUPP+ZV74gzwz4GvWDUcS9aJku65b5TYRBBu2wA8Y0vrcvvbGGkEZnJ7XHl7vEddrLCI8hER7ZV/QErMP9STut8WMF7HS3M6uOaa1O/p2+sD/IGAdL7+hITHj+xCnr200y"
			template_sign, _ := utils.RsaSign([]byte(private_key), []byte(tc))
			param_add_template := &AddTmp{Content: tc, TemplateSign: string(template_sign)}
			err := CreateNewTemplate(param_add_template, &middleware.CustomClaims{PubKey: "MIIBCgKCAQEAovUhmHbdWNs6eVI/4YS3cgwNdF0sVs2UxLGv90JHlPzlVR7vHliEsE7e+YHwqHlHcDFHoXBdp3n2OYfYFCrVAfD6zwyEuCA5GYYHKX5jZyO6wTsMdaIXVFAZgESENfHx2FZgVZnqe13TItg1S1jNoHHVNfaR44n0KrJSPTzBypIdDBeHYorPjXrvrZkQnE8g3NMyKne+gdn16YcKl7RWNaUxoQ3If6mgeidVj1MEYbH7Vfy3dIruf/l5GRvWjveMP8AywiUpP8CEG+9M5sjFB/pXOV67s9YiExalDPEzpQWCLz9op07iILwZaopAsQwOAfoO7jLtboLo6tcAj7X/9QIDAQAB", ID: 99})
			So(err.Code, ShouldEqual, 0)
		})
	})
}

func TestVerifyTemplate(t *testing.T) {
	// 股东私钥
	prvKey1, _ := hex.DecodeString("4D4949456F77494241414B434151454177673962494D58752F456857422B6C7563645463686D57305437357558346C764E65394B4C5942546271383553776F6A35745A4279442F7978626A526D726157465871746E4B626C48467144566567774569634741476B4971625157494D436F43504E65625246553449515A524B6F6979442F6E764B5A56724F787A6A485645736D517836614B3337376A6A503059784B65693256717150615A376E6953554550426379445737335878716C42734B724D2F652B6C44736935724E5A70632B736C385741787A464762652B592B304A474C4264727350386A37466564594B6B747850654D516D6E3134343143752F4A2B756B6336723458625934705357632F572B4130794D364D522B49303646685475535977425061303331736B6276697A41664E636470332F3848337542776633634F7049634E33526B38577A4268776B784D76526F49624F5242526B503364525564636C704477494441514142416F494241435A512B65464373366C456278676F3078484F77617758734C352B447A543134657434542B456B45565574613961524C66324B4748723968575863314359454879596D6A746A3761776C714A6278306C2B6B396F305735643268716B6D50744A6B422B4C5172442F5570485375587634655A41462F5637736A6275695049577858576A71756443706162446A7273546C4F6C3964574D7655686D38736D4B31324677314345443278416B573154786D7A336873436472622F47313251514F615831417762646650594372425A67646A79715A37786E7672664C6B716347334A716C675870457A3376546E4738567238645778696F3230686954504334706E515454357078484966446A4368444E4335347034455665476736445642396C51676E47437954694C30612F673543383450444A6374334F7842686A336F344743567A4379474C50386E302B4F6E3356377373464B7A4B4945436759454139747A534737653375673459766E6C7A685376356B72424D397764425A4E6B6969564835756678656F44537537445666554D6341485662426D59477067486B314F593354774334715659617053766A5953754967654C677632536D674A59517157416D50436A2F7A586E5A576638776E61416D53454F7932796442784671796C2F413657364D395279564276747349783179306A4E7865426736427A31784F707242616451706E6742314543675945417954347965325A2F634439755339334E2F464B74536776677438306733616C59556F4F6F6B3552397663306D56384254434F2B5047762F763641374845414F66763753496E6A356D7A33733474616E587A77666A78686C737778546332704F556A51583742645277523339743872647247333053704D5A63424A6E2B73657456395074485A79594B435274766E75773466775153314D6E3545786133746F7A776342466E436B597A456C3843675941704B4E4A4457744E675378396E43727276466C44725730664A334554542B7277655A7A4F624266414643434F754D497675334F387739526B6362686B443262473949642B5061474D596C58592B6B4E554854304B59766955454D554F4A586372746D4E4A6E704267417850785248496E706438634A32563330736A4E59303370745630566B4663554F4B73496C6C36565675696E6E38707865684F387269685339493577653061636673514B4267514358645551784D35632B4E373866516A524262306755727050575159413230376B786750316939534D64736B546C546B34774C3377593666373550415839753379446E3741697950772F733547657839457763394F7479487A6535646A43654D675979794E393438454E546F3746576434327467394F44437739476C6C694E716865314B2B4D706B6876536B4C623752556F497A793541626C2F384630622B7A384B575536666F4F33584B514B42674172376C536C302F766B75776E3965565134346B6B5638416C6F3357764C6853366F4749666345486270335232437478637067306C5144683649517675547466473137504B6E444C7477676F5A464A5736336D4571554B45386F34642B3335645A327237515A4D33394F524961576765716A74494D646B5148506F68594F3354747A45454D4E6747695569786A354A616A77566B3130574571624C6E33784D6D635073544C793055396839")
	pubkey1, _ := hex.DecodeString("4D49494243674B434151454177673962494D58752F456857422B6C7563645463686D57305437357558346C764E65394B4C5942546271383553776F6A35745A4279442F7978626A526D726157465871746E4B626C48467144566567774569634741476B4971625157494D436F43504E65625246553449515A524B6F6979442F6E764B5A56724F787A6A485645736D517836614B3337376A6A503059784B65693256717150615A376E6953554550426379445737335878716C42734B724D2F652B6C44736935724E5A70632B736C385741787A464762652B592B304A474C4264727350386A37466564594B6B747850654D516D6E3134343143752F4A2B756B6336723458625934705357632F572B4130794D364D522B49303646685475535977425061303331736B6276697A41664E636470332F3848337542776633634F7049634E33526B38577A4268776B784D76526F49624F5242526B503364525564636C704477494441514142")
	Convey("审批模板", t, func() {
		Convey("股东1审批签名", func() {
			param_time := time.Now().UnixNano()
			param_sign, _ := utils.RsaSign(prvKey1, []byte(utils.ToStr(param_time)))
			template_sign, _ := utils.RsaSign(prvKey1, []byte(tc))
			admin1 := &TemplateVo{
				Timestamp:    param_time,
				Sign:         string(param_sign),
				AdminStatus:  common.TemplateApprovaled,
				Id:           "a0b49b81-cae9-43d3-b17e-73989003dd46",
				TemplateSign: string(template_sign),
			}
			claims := &middleware.CustomClaims{ID: 98, AppID: "001", Account: "app-002", PubKey: string(pubkey1)}
			//claims := &middleware.CustomClaims{ID:82, PubKey:admin2_pub}
			err := VerifyTemplate(admin1, claims)
			So(err.Code, ShouldEqual, 0)
		})
	})
}

func TestCancelTemp(t *testing.T) {
	priv_key := "MIIEowIBAAKCAQEAovUhmHbdWNs6eVI/4YS3cgwNdF0sVs2UxLGv90JHlPzlVR7vHliEsE7e+YHwqHlHcDFHoXBdp3n2OYfYFCrVAfD6zwyEuCA5GYYHKX5jZyO6wTsMdaIXVFAZgESENfHx2FZgVZnqe13TItg1S1jNoHHVNfaR44n0KrJSPTzBypIdDBeHYorPjXrvrZkQnE8g3NMyKne+gdn16YcKl7RWNaUxoQ3If6mgeidVj1MEYbH7Vfy3dIruf/l5GRvWjveMP8AywiUpP8CEG+9M5sjFB/pXOV67s9YiExalDPEzpQWCLz9op07iILwZaopAsQwOAfoO7jLtboLo6tcAj7X/9QIDAQABAoIBAAX4+gw+fwpcrp33t8OqO5cAfhW/vHpJ+qPi51Imqbz6L+WYxqbUE7jqix2V954VI9sm0ztFhQk4XR/qrK7AiyIRFQA1sz+UXsNiuCx7J2WGV7fxVBmToUtDzEt5N2dKwNRvBgTdKUzwOWbTidCDJrY0td3OdbZjPgG6m20HJwsnlM3SXFKBssOLEr1aL4ieHcelQ8sXLeTFU+hOlCP/CSPDEMa2GwvWC0Aa6ihCUu2J00zXNUPlVvLmz96vBxuP+MQgLO0nwjHQW+TObdQxnFp2UjwyLhE6l/qCXeZ6dAfUg4TFFXo9m5lqmCJlIUEe24qytb0Q6SWkMKpD5Ami+JkCgYEA2E3yr7RjR9MY7DCGXf6zIwGRg2FBe/xqq9tY+2AHdebVtpWhqToS1FEG6YGoZzO3EOPqXfiUdqbOwpuqxdTDIhscy2iWEmmSlpm+XVpcUCdAQQUiJT0HVnD7qKyvpy0RTxV6ry6OWWXKxnR7+Kqip7CVNLUO5Ua297aYdEn/09MCgYEAwNztaBuAR9y34p4DeCAZElhEg3o9nYtF3WKJqy+pxV1/EEcDD+5KzFuggC2WqfvO6pZ3llRAV9lnzmrNK7w4pImS4Ot/Yo1vLz/dhi7odHGZ025iScqB7kKTKXmVelH0RN5ya4dYGtiM+cGyqWkXQ570r52uWOHzxQhGZwU7KBcCgYBX6XZVSy9+paFffWlfEoGkHKMTjpea2MBSENhdcm4V0TfS+YW2zc+RU01H4labez1QNBGxF7LfdnRaTAJgXVThq7tMQLmdyiK16s6gCnWybgGDa56xG9i7nTfrGiRaAcsjJUuzn2xjkSeU+QrZyiBJn09FGMFxFgiPLTtRKDI92QKBgBhXyl9pmTd41Hz2FNoEsyVtnrg5pa1M9vSKi3Xf/j27H16el/Raz3Yb2pZTKsEp12QkudNvru9nsYKMWHk8uKmL884P63Q1BeOg3AUjxpNsA058kEtgFCZOoOSLRCK8VWib6zVHUAeTbbfYiwR3D0ipt5iy/l6ZpnsoIDrAnkbhAoGBAJKebsuGx/YfWuvCRnHPAvPL3yHtrVxRuYhlZIKayIDUPP+ZV74gzwz4GvWDUcS9aJku65b5TYRBBu2wA8Y0vrcvvbGGkEZnJ7XHl7vEddrLCI8hER7ZV/QErMP9STut8WMF7HS3M6uOaa1O/p2+sD/IGAdL7+hITHj+xCnr200y"
	pub_key := "MIIBCgKCAQEAovUhmHbdWNs6eVI/4YS3cgwNdF0sVs2UxLGv90JHlPzlVR7vHliEsE7e+YHwqHlHcDFHoXBdp3n2OYfYFCrVAfD6zwyEuCA5GYYHKX5jZyO6wTsMdaIXVFAZgESENfHx2FZgVZnqe13TItg1S1jNoHHVNfaR44n0KrJSPTzBypIdDBeHYorPjXrvrZkQnE8g3NMyKne+gdn16YcKl7RWNaUxoQ3If6mgeidVj1MEYbH7Vfy3dIruf/l5GRvWjveMP8AywiUpP8CEG+9M5sjFB/pXOV67s9YiExalDPEzpQWCLz9op07iILwZaopAsQwOAfoO7jLtboLo6tcAj7X/9QIDAQAB"
	Convey("作废审批模板", t, func() {
		param_timestamp := time.Now().UnixNano()
		param_sign, _ := utils.RsaSign([]byte(priv_key), []byte(utils.ToStr(param_timestamp)))
		template_sign, _ := utils.RsaSign([]byte(priv_key), []byte(tc))
		template_info := &ParamCancelTemplate{
			TemplateId:   "027d74ef-c1f0-43b3-aca8-6c62000be19c",
			TemplateSign: string(template_sign),
			Sign:         string(param_sign),
			Timestamp:    param_timestamp,
		}

		claims := &middleware.CustomClaims{ID: 99, PubKey: pub_key}
		fmt.Println("templateInfo", template_info)
		err := CancelTemp(template_info, claims)
		So(err.Code, ShouldEqual, 0)
	})
}

func TestGetATemplateListByUserName(t *testing.T) {
	Convey("获取可用的审批流模板", t, func() {
		err := DisableTemplateByUserName("gudong_2")
		So(err, ShouldBeNil)
	})
}

func TestTxNumByTempID(t *testing.T) {
	Convey("统计转账数目", t, func() {
		err := TxNumByTempID("1c4d2458-478c-4fdf-8570-882ba6a71048")
		So(err, ShouldBeNil)
	})
}
