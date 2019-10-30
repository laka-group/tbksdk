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

type ErrorResponse struct {
	ErrorResponse struct {
		SubMsg  string `json:"sub_msg"`
		Code    int64  `json:"code"`
		SubCode string `json:"sub_code"`
		Msg     string `json:"msg"`
	} `json:"error_response"`
}

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

type ItemInfoMap map[int64]itemInfo

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
	// 得到了返回,但不是item信息,尝试解析为ErrorResponse
	if respData.TbkPrivilegeGetResponse.Result.Data.ItemID == 0 {
		errorResponse := ErrorResponse{}
		_ = json.Unmarshal(*bodyByte, &errorResponse)
		if errorResponse.ErrorResponse.Code != 0 {
			return nil, errors.New(string(*bodyByte))
		}
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

type TbkCouponGetResponse struct {
	TbkCouponGetResponse struct {
		Data TbkCouponGetData `json:"data"`
	} `json:"tbk_coupon_get_response"`
}

type TbkCouponGetData struct {
	CouponStartFee    string `json:"coupon_start_fee"`
	CouponRemainCount int64  `json:"coupon_remain_count"`
	CouponTotalCount  int64  `json:"coupon_total_count"`
	CouponEndTime     string `json:"coupon_end_time"`
	CouponStartTime   string `json:"coupon_start_time"`
	CouponAmount      string `json:"coupon_amount"`
	CouponSrcScene    int64  `json:"coupon_src_scene"`
	CouponType        int64  `json:"coupon_type"`
	CouponActivityId  string `json:"coupon_activity_id"`
}

/**
淘宝客-公用-阿里妈妈推广券详情查询
https://open.taobao.com/api.htm?spm=a219a.7386797.0.0.40b2669alwkQgI&source=search&docId=31106&docType=2
*/
func TbkCouponGet(itemId int64, couponId string) (*TbkCouponGetData, error) {
	var paramsMap = map[string]string{
		"item_id":     strconv.FormatInt(itemId, 10),
		"activity_id": couponId,
	}
	var bodyByte, err = apply(Constants.AlimamaTbkCouponGet, paramsMap)
	if err != nil {
		return nil, err
	}
	var respData = TbkCouponGetResponse{}
	err = json.Unmarshal(*bodyByte, &respData)
	if err != nil {
		return nil, err
	}
	if respData.TbkCouponGetResponse.Data.CouponActivityId == "" {
		errorResponse := ErrorResponse{}
		_ = json.Unmarshal(*bodyByte, &errorResponse)
		if errorResponse.ErrorResponse.Code != 0 {
			return nil, errors.New(string(*bodyByte))
		}
		return nil, errors.New("data is empty")
	}
	return &respData.TbkCouponGetResponse.Data, nil
}

type tbkScOrderDetailsGetResponseData struct {
	TbkScOrderDetailsGetResponse TbkScOrderDetailsGetResponse `json:"tbk_sc_order_details_get_response"`
}
type TbkScOrderDetailsGetResponse struct {
	Data struct {
		Results struct {
			PublisherOrderDto []PublisherOrderDto `json:"publisher_order_dto"`
		} `json:"results"`
		HasPre        bool   `json:"has_pre"`
		PositionIndex string `json:"position_index"`
		HasNext       bool   `json:"has_next"`
		PageNo        int64  `json:"page_no"`
		PageSize      int64  `json:"page_size"`
	} `json:"data"`
}
type PublisherOrderDto struct {
	TbPaidTime                         string `json:"tb_paid_time"`
	TkPaidTime                         string `json:"tk_paid_time"`
	PayPrice                           string `json:"pay_price"`
	PubShareFee                        string `json:"pub_share_fee"`
	TradeId                            string `json:"trade_id"`
	TkOrderRole                        int64  `json:"tk_order_role"`
	TkEarningTime                      string `json:"tk_earning_time"`
	AdzoneId                           int64  `json:"adzone_id"`
	PubShareRate                       string `json:"pub_share_rate"`
	RefundTag                          int64  `json:"refund_tag"`
	SubsidyRate                        string `json:"subsidy_rate"`
	TkTotalRate                        string `json:"tk_total_rate"`
	ItemCategoryName                   string `json:"item_category_name"`
	SellerNick                         string `json:"seller_nick"`
	PubId                              int64  `json:"pub_id"`
	AlimamaRate                        string `json:"alimama_rate"`
	SubsidyType                        string `json:"subsidy_type"`
	ItemImg                            string `json:"item_img"`
	PubSharePreFee                     string `json:"pub_share_pre_fee"`
	AlipayTotalPrice                   string `json:"alipay_total_price"`
	ItemTitle                          string `json:"item_title"`
	SiteName                           string `json:"site_name"`
	ItemNum                            int64  `json:"item_num"`
	SubsidyFee                         string `json:"subsidy_fee"`
	AlimamaShareFee                    string `json:"alimama_share_fee"`
	TradeParentId                      string `json:"trade_parent_id"`
	OrderType                          string `json:"order_type"`
	TkCreateTime                       string `json:"tk_create_time"`
	FlowSource                         string `json:"flow_source"`
	TerminalType                       string `json:"terminal_type"`
	ClickTime                          string `json:"click_time"`
	TkStatus                           int8   `json:"tk_status"`
	ItemPrice                          string `json:"item_price"`
	ItemId                             int64  `json:"item_id"`
	AdzoneName                         string `json:"adzone_name"`
	TotalCommissionRate                string `json:"total_commission_rate"`
	ItemLink                           string `json:"item_link"`
	SiteId                             int64  `json:"site_id"`
	SellerShopTitle                    string `json:"seller_shop_title"`
	IncomeRate                         string `json:"income_rate"`
	TotalCommissionFee                 string `json:"total_commission_fee"`
	TkCommissionPreFeeForMediaPlatform string `json:"tk_commission_pre_fee_for_media_platform"`
	TkCommissionFeeForMediaPlatform    string `json:"tk_commission_fee_for_media_platform"`
	TkCommissionRateForMediaPlatform   string `json:"tk_commission_rate_for_media_platform"`
	SpecialId                          int64  `json:"special_id"`
	RelationId                         int64  `json:"relation_id"`
}

type TbkScActivityResponse struct {
	TbkScActivitylinkToolgetResponse struct {
		ResultMsg    string `json:"result_msg"`
		Data         string `json:"data"`
		ResultCode   int64  `json:"result_code"`
		BizErrorDesc string `json:"biz_error_desc"`
		BizErrorCode int64  `json:"biz_error_code"`
	} `json:"tbk_sc_activitylink_toolget_response"`
}

/**
淘宝客-服务商-所有订单查询
https://open.taobao.com/api.htm?spm=a219a.7386653.0.0.2b43669aJY6jMz&source=search&docId=43755&docType=2
*/
func TbkScOrderDetailsGet(startTime string, endTime string, session string, pageNo int64, pageSize int64, positionIndex string, queryType int64) (*TbkScOrderDetailsGetResponse, error) {
	var paramsMap = map[string]string{
		"start_time":     startTime,
		"end_time":       endTime,
		"session":        session,
		"page_no":        strconv.FormatInt(pageNo, 10),
		"page_size":      strconv.FormatInt(pageSize, 10),
		"position_index": positionIndex,
		"query_type":     strconv.FormatInt(queryType, 10),
		"jump_type":      "1",
		"order_scene":    "1",
	}
	var bodyByte, err = apply(Constants.AlimamaTbkScOrderDetailsGet, paramsMap)
	if err != nil {
		return nil, err
	}
	var respData = tbkScOrderDetailsGetResponseData{}
	err = json.Unmarshal(*bodyByte, &respData)
	if err != nil {
		return nil, err
	}
	if respData.TbkScOrderDetailsGetResponse.Data.PageNo == 0 {
		errorResponse := ErrorResponse{}
		_ = json.Unmarshal(*bodyByte, &errorResponse)
		if errorResponse.ErrorResponse.Code != 0 {
			return nil, errors.New(string(*bodyByte))
		}
		return nil, errors.New("data is empty")
	}
	return &respData.TbkScOrderDetailsGetResponse, nil
}

// 官方活动转链 https://open.taobao.com/api.htm?docId=41921&docType=2&source=search
func TbkScActivityLinkToolGet(activityId int64, adzoneId, siteId, session string) (string, error) {
	paramMap := map[string]string{
		"adzone_id":          adzoneId,
		"site_id":            siteId,
		"promotion_scene_id": strconv.FormatInt(activityId, 10),
		"platform":           "2",
		"session":            session,
	}
	var bodyByte, err = apply(Constants.AlimamaTbkScActivityLinkToolGet, paramMap)
	if err != nil {
		return "", err
	}
	var respData = TbkScActivityResponse{}
	err = json.Unmarshal(*bodyByte, &respData)
	if err != nil {
		return "", err
	}
	return respData.TbkScActivitylinkToolgetResponse.Data, nil
}

type TbkSpreadGetResponse struct {
	TbkSpreadGetResponse struct {
		Results struct {
			TbkSpread []struct {
				Content string `json:"content"`
				ErrMsg  string `json:"err_msg"`
			} `json:"tbk_spread"`
		} `json:"results"`
		TotalResults int    `json:"total_results"`
		RequestID    string `json:"request_id"`
	} `json:"tbk_spread_get_response"`
}

// 长链转短链 https://open.taobao.com/api.htm?spm=a219a.7386797.0.0.723f669aQgtoLl&source=search&docId=27832&docType=2
func TbkSpreadGet(url, session string) (string, error) {
	paramMap := map[string]string{
		"app_key":     Constants.AlimamaKey,
		"method":      Constants.AlimamaTbkSpreadGet,
		"session":     session,
		"format":      "json",
		"v":           "2.0",
		"sign_method": "md5",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"requests":    `[{"url":"` + url + `"}]`,
	}
	var bodyByte, err = apply(Constants.AlimamaTbkSpreadGet, paramMap)
	if err != nil {
		return "", err
	}
	var respData = TbkSpreadGetResponse{}
	err = json.Unmarshal(*bodyByte, &respData)
	if err != nil {
		return "", err
	}
	if len(respData.TbkSpreadGetResponse.Results.TbkSpread) > 0 {
		return respData.TbkSpreadGetResponse.Results.TbkSpread[0].Content, nil
	}
	return "", nil
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

	for k, _ := range paramMap {
		// && k != "session"
		if k != "" {
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
