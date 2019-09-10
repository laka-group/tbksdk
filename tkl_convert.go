package tbksdk

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// https://open.taobao.com/api.htm?docId=43873&docType=2&scopeId=16401
// taobao.tbk.sc.tpwd.convert
type respDataTpwdConvert struct {
	TbkScTpwdConvertResponse struct {
		Data struct {
			ClickURL string `json:"click_url"`
			NumIid   string `json:"num_iid"`
		} `json:"data"`
		RequestID string `json:"request_id"`
	} `json:"tbk_sc_tpwd_convert_response"`
}

func GetItemIdByTkl(session string, tkl string, adzoneId string, siteId string) (int64, error) {
	var paramsMap = map[string]string{
		"session":          session,
		"password_content": tkl,
		"adzone_id":        adzoneId,
		"site_id":          siteId,
	}
	var bodyByte, err = apply(Constants.AlimamaTpwdConvertUrl, paramsMap)
	if err != nil {
		return 0, err
	}

	var respData = respDataTpwdConvert{}
	err = json.Unmarshal(*bodyByte, &respData)
	if err != nil {
		return 0, err
	}

	if respData.TbkScTpwdConvertResponse.Data.NumIid == "" {
		return 0, errors.New("cannot find this tkl fanli info")
	}
	itemIdInt64, err := strconv.ParseInt(respData.TbkScTpwdConvertResponse.Data.NumIid, 10, 64)
	if err != nil {
		return 0, err
	}
	return itemIdInt64, nil
}

// https://open.taobao.com/api.htm?docId=24518&docType=2&scopeId=16189
// taobao.tbk.item.info.get( 淘宝客-公用-淘宝客商品详情查询(简版) )
type respDataItemInfo struct {
	TbkItemInfoGetResponse struct {
		Results struct {
			NTbkItem []itemInfo `json:"n_tbk_item"`
		} `json:"results"`
	} `json:"tbk_item_info_get_response"`
}
type itemInfo struct {
	CatName     string `json:"cat_name"`
	NumIid      int64  `json:"num_iid"`
	Title       string `json:"title"`
	PictURL     string `json:"pict_url"`
	SmallImages struct {
		String []string `json:"string"`
	} `json:"small_images"`
	ReservePrice    string `json:"reserve_price"`
	ZkFinalPrice    string `json:"zk_final_price"`
	UserType        int64  `json:"user_type"`
	Provcity        string `json:"provcity"`
	ItemURL         string `json:"item_url"`
	SellerID        int64  `json:"seller_id"`
	Volume          int64  `json:"volume"`
	Nick            string `json:"nick"`
	CatLeafName     string `json:"cat_leaf_name"`
	IsPrepay        bool   `json:"is_prepay"`
	ShopDsr         int64  `json:"shop_dsr"`
	Ratesum         int64  `json:"ratesum"`
	IRfdRate        bool   `json:"i_rfd_rate"`
	HGoodRate       bool   `json:"h_good_rate"`
	HPayRate30      bool   `json:"h_pay_rate30"`
	FreeShipment    bool   `json:"free_shipment"`
	MaterialLibType string `json:"material_lib_type"`
}

func ItemInfoGet(itemId ...int64) ([]itemInfo, error) {
	var itemInfo = respDataItemInfo{}
	var itemIds = ""
	for _, v := range itemId {
		itemIds = itemIds + fmt.Sprintf("%v", v) + ","
	}
	var paramsMap = map[string]string{
		"num_iids": strings.TrimRight(itemIds, ","),
	}
	var bodyByte, err = apply(Constants.AlimamaItemInfoGetUrl, paramsMap)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(*bodyByte, &itemInfo)
	if err != nil {
		return nil, err
	}
	if len(itemInfo.TbkItemInfoGetResponse.Results.NTbkItem) == 0 {
		return nil, errors.New("itemInfo.TbkItemInfoGetResponse.Results.NTbkItem len is 0")
	}
	return itemInfo.TbkItemInfoGetResponse.Results.NTbkItem, nil
}

// https://open.taobao.com/api.htm?docId=28625&docType=2&scopeId=12403
// taobao.tbk.privilege.get( 淘宝客-服务商-单品券高效转链 )
type respDataPrivilegeGet struct {
	TbkPrivilegeGetResponse struct {
		Result struct {
			Data privilegeData `json:"data"`
		} `json:"result"`
	} `json:"tbk_privilege_get_response"`
}
type privilegeData struct {
	CategoryID          int64  `json:"category_id"`
	CouponClickURL      string `json:"coupon_click_url"`
	CouponEndTime       string `json:"coupon_end_time"`
	CouponInfo          string `json:"coupon_info"`
	CouponStartTime     string `json:"coupon_start_time"`
	ItemID              int64  `json:"item_id"`
	MaxCommissionRate   string `json:"max_commission_rate"`
	CouponTotalCount    int64  `json:"coupon_total_count"`
	CouponRemainCount   int64  `json:"coupon_remain_count"`
	MmCouponRemainCount int64  `json:"mm_coupon_remain_count"`
	MmCouponTotalCount  int64  `json:"mm_coupon_total_count"`
	MmCouponClickURL    string `json:"mm_coupon_click_url"`
	MmCouponEndTime     string `json:"mm_coupon_end_time"`
	MmCouponStartTime   string `json:"mm_coupon_start_time"`
	MmCouponInfo        string `json:"mm_coupon_info"`
	CouponType          int64  `json:"coupon_type"`
	ItemURL             string `json:"item_url"`
}

func PrivilegeGet(session string, itemId int64, adzoneId string, siteId string) (*privilegeData, error) {
	var paramsMap = map[string]string{
		"session":   session,
		"item_id":   fmt.Sprintf("%v", itemId),
		"adzone_id": adzoneId,
		"site_id":   siteId,
	}
	var bodyByte, err = apply(Constants.AlimamaPrivilegeGet, paramsMap)
	if err != nil {
		return nil, err
	}
	var respData = respDataPrivilegeGet{}
	err = json.Unmarshal(*bodyByte, &respData)
	if err != nil {
		return nil, err
	}
	if respData.TbkPrivilegeGetResponse.Result.Data.ItemID == 0 {
		return nil, errors.New("respData.TbkPrivilegeGetResponse.Result.Data.ItemID is 0")
	}
	return &respData.TbkPrivilegeGetResponse.Result.Data, nil
}

// https://open.taobao.com/api.htm?docId=31127&docType=2&scopeId=11655
// taobao.tbk.tpwd.create( 淘宝客-公用-淘口令生成 )
type respDataTklCreate struct {
	TbkTpwdCreateResponse struct {
		Data struct {
			Model string `json:"model"`
		} `json:"data"`
	} `json:"tbk_tpwd_create_response"`
}

func TklCreate(title, url, logo string) (string, error) {
	var paramsMap = map[string]string{
		"text": title,
		"url":  url,
		"logo": logo,
	}
	var bodyByte, err = apply(Constants.AlimamaTklCreate, paramsMap)
	if err != nil {
		return "", err
	}
	var respData = respDataTklCreate{}
	err = json.Unmarshal(*bodyByte, &respData)
	if err != nil {
		return "", err
	}
	return respData.TbkTpwdCreateResponse.Data.Model, nil
}

// https://open.taobao.com/api.htm?docId=35263&docType=2&scopeId=13991
// taobao.tbk.sc.material.optional( 淘宝客-服务商-物料搜索 )
type respDataMaterialOptional struct {
	TbkScMaterialOptionalResponse struct {
		TotalResults int `json:"total_results"`
		ResultList   struct {
			MapData []SearchResultList `json:"map_data"`
		} `json:"result_list"`
	} `json:"tbk_sc_material_optional_response"`
}
type SearchResultList struct {
	CouponStartTime string `json:"coupon_start_time"`
	CouponEndTime   string `json:"coupon_end_time"`
	InfoDxjh        string `json:"info_dxjh"`
	TkTotalSales    string `json:"tk_total_sales"`
	TkTotalCommi    string `json:"tk_total_commi"`
	CouponID        string `json:"coupon_id"`
	NumIid          int64  `json:"num_iid"`
	Title           string `json:"title"`
	PictURL         string `json:"pict_url"`
	SmallImages     struct {
		String []string `json:"string"`
	} `json:"small_images"`
	ReservePrice          string `json:"reserve_price"`
	ZkFinalPrice          string `json:"zk_final_price"`
	UserType              int    `json:"user_type"`
	Provcity              string `json:"provcity"`
	ItemURL               string `json:"item_url"`
	IncludeMkt            string `json:"include_mkt"`
	IncludeDxjh           string `json:"include_dxjh"`
	CommissionRate        string `json:"commission_rate"`
	Volume                int    `json:"volume"`
	SellerID              int    `json:"seller_id"`
	CouponTotalCount      int    `json:"coupon_total_count"`
	CouponRemainCount     int    `json:"coupon_remain_count"`
	CouponInfo            string `json:"coupon_info"`
	CommissionType        string `json:"commission_type"`
	ShopTitle             string `json:"shop_title"`
	URL                   string `json:"url"`
	CouponShareURL        string `json:"coupon_share_url"`
	ShopDsr               int    `json:"shop_dsr"`
	WhiteImage            string `json:"white_image"`
	ShortTitle            string `json:"short_title"`
	CategoryID            int    `json:"category_id"`
	CategoryName          string `json:"category_name"`
	LevelOneCategoryID    int    `json:"level_one_category_id"`
	LevelOneCategoryName  string `json:"level_one_category_name"`
	Oetime                string `json:"oetime"`
	Ostime                string `json:"ostime"`
	JddNum                int    `json:"jdd_num"`
	JddPrice              string `json:"jdd_price"`
	UvSumPreSale          int    `json:"uv_sum_pre_sale"`
	CouponAmount          string `json:"coupon_amount"`
	CouponStartFee        string `json:"coupon_start_fee"`
	ItemDescription       string `json:"item_description"`
	Nick                  string `json:"nick"`
	XID                   string `json:"x_id"`
	OrigPrice             string `json:"orig_price"`
	TotalStock            int    `json:"total_stock"`
	SellNum               int    `json:"sell_num"`
	Stock                 int    `json:"stock"`
	TmallPlayActivityInfo string `json:"tmall_play_activity_info"`
	ItemID                int64  `json:"item_id"`
	RealPostFee           string `json:"real_post_fee"`
}

func SearchTitle(session, title, adzoneId, siteId string) ([]SearchResultList, error) {
	var paramsMap = map[string]string{
		"q":         title,
		"adzone_id": adzoneId,
		"site_id":   siteId,
		"session":   session,
		"page_size": "1",
	}
	var bodyByte, err = apply(Constants.AlimamaScMaterialOptional, paramsMap)
	if err != nil {
		return nil, err
	}
	var respData = respDataMaterialOptional{}
	err = json.Unmarshal(*bodyByte, &respData)
	if err != nil {
		return nil, err
	}
	return respData.TbkScMaterialOptionalResponse.ResultList.MapData, nil
}

func apply(api string, params map[string]string) (*[]byte, error) {
	paramMap := map[string]string{
		"app_key":     Constants.AlimamaKey,
		"method":      api,
		"format":      "json",
		"v":           "2.0",
		"sign_method": "md5",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
	}
	// 合并参数
	for k, v := range params {
		paramMap[k] = v
	}
	// 生成签名
	sign := createSign(paramMap)
	// 准备发起http请求
	req, err := http.NewRequest("GET", Constants.AlimamaApiUrl, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	for k, v := range paramMap {
		query.Add(k, v)
	}
	query.Add("sign", sign)
	req.URL.RawQuery = query.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &bodyByte, nil
}

// 制造签名
func createSign(paramMap map[string]string) string {
	sign := Constants.AlimamaSecret
	var paramKeySlice []string

	for k, v := range paramMap {
		//&& k != "session"
		if k != "" && v != "" {
			paramKeySlice = append(paramKeySlice, k)
		}
	}
	// 排序key
	sort.Strings(paramKeySlice)
	for i := 0; i < len(paramKeySlice); i++ {
		sign += paramKeySlice[i] + paramMap[paramKeySlice[i]]
	}
	sign += Constants.AlimamaSecret
	md5New := md5.New()
	_, _ = io.WriteString(md5New, sign)
	return strings.ToUpper(fmt.Sprintf("%x", md5New.Sum(nil)))
}
