package tbksdk

type constants struct {
	AlimamaKey                string
	AlimamaSecret             string
	AlimamaApiUrl             string
	AlimamaTpwdConvertUrl     string
	AlimamaItemInfoGetUrl     string
	AlimamaPrivilegeGet       string
	AlimamaTklCreate          string
	AlimamaScMaterialOptional string
}

var Constants = constants{}

func init() {
	// alimama related
	Constants.AlimamaKey = ""
	Constants.AlimamaSecret = ""
	//api
	Constants.AlimamaApiUrl = "http://gw.api.taobao.com/router/rest?"
	Constants.AlimamaTpwdConvertUrl = "taobao.tbk.sc.tpwd.convert"
	Constants.AlimamaItemInfoGetUrl = "taobao.tbk.item.info.get"
	Constants.AlimamaPrivilegeGet = "taobao.tbk.privilege.get" // 淘宝客-服务商-单品券高效转链
	Constants.AlimamaTklCreate = "taobao.tbk.tpwd.create"
	Constants.AlimamaScMaterialOptional = "taobao.tbk.sc.material.optional"
}
