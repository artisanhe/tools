package client_cashdesk

import (
	"fmt"

	golib_tools_courier "github.com/artisanhe/tools/courier"
	golib_tools_courier_client "github.com/artisanhe/tools/courier/client"
	golib_tools_courier_enumeration "github.com/artisanhe/tools/courier/enumeration"
	golib_tools_courier_status_error "github.com/artisanhe/tools/courier/status_error"
	golib_tools_timelib "github.com/artisanhe/tools/timelib"
)

type ClientCashdeskInterface interface {
	CancelTrans(req CancelTransRequest) (resp *CancelTransResponse, err error)
	CheckHealth() (resp *CheckHealthResponse, err error)
	CreateTrans(req CreateTransRequest) (resp *CreateTransResponse, err error)
	EBankPay(req EBankPayRequest) (resp *EBankPayResponse, err error)
	GetAreaBankByBankCode(req GetAreaBankByBankCodeRequest) (resp *GetAreaBankByBankCodeResponse, err error)
	GetAreaBankList(req GetAreaBankListRequest) (resp *GetAreaBankListResponse, err error)
	GetAreaList(req GetAreaListRequest) (resp *GetAreaListResponse, err error)
	GetBankByID(req GetBankByIDRequest) (resp *GetBankByIDResponse, err error)
	GetBankList(req GetBankListRequest) (resp *GetBankListResponse, err error)
	GetPabWhiteListTrans(req GetPabWhiteListTransRequest) (resp *GetPabWhiteListTransResponse, err error)
	GetSign(req GetSignRequest) (resp *GetSignResponse, err error)
	GetTransByOrderNo(req GetTransByOrderNoRequest) (resp *GetTransByOrderNoResponse, err error)
	OldTransQuery(req OldTransQueryRequest) (resp *OldTransQueryResponse, err error)
	PabPay(req PabPayRequest) (resp *PabPayResponse, err error)
	PabQueryByCertificate(req PabQueryByCertificateRequest) (resp *PabQueryByCertificateResponse, err error)
	PabQueryByUser(req PabQueryByUserRequest) (resp *PabQueryByUserResponse, err error)
	TransListQuery(req TransListQueryRequest) (resp *TransListQueryResponse, err error)
	UpdateTrans(req UpdateTransRequest) (resp *UpdateTransResponse, err error)
	WithSwagger() (resp *WithSwaggerResponse, err error)
}

type ClientCashdesk struct {
	golib_tools_courier_client.Client
}

func (ClientCashdesk) MarshalDefaults(v interface{}) {
	if cl, ok := v.(*ClientCashdesk); ok {
		cl.Name = "cashdesk"
		cl.Client.MarshalDefaults(&cl.Client)
	}
}

func (c ClientCashdesk) Init() {
	c.CheckService()
}

func (c ClientCashdesk) CheckService() {
	err := c.Request(c.Name+".Check", "HEAD", "/", nil).
		Do().
		Into(nil)
	statusErr := golib_tools_courier_status_error.FromError(err)
	if statusErr.Code == int64(golib_tools_courier_status_error.RequestTimeout) {
		panic(fmt.Errorf("service %s have some error %s", c.Name, statusErr))
	}
}

type CancelTransRequest struct {
	// ?????????id
	TransID uint64 `in:"path" name:"transID" validate:"@uint64[1,]"`
	// ????????????id
	UserID string `in:"path" name:"userID" validate:"@string[1,64]"`
	//
	Body CancelTransBody `in:"body" name:"body,omitempty"`
}

func (c ClientCashdesk) CancelTrans(req CancelTransRequest) (resp *CancelTransResponse, err error) {
	resp = &CancelTransResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".CancelTrans", "PUT", "/cashdesk/v0/trans/:transID/user/:userID/cancel", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type CancelTransResponse struct {
	Meta golib_tools_courier.Metadata
	Body string
}

func (c ClientCashdesk) CheckHealth() (resp *CheckHealthResponse, err error) {
	resp = &CheckHealthResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".CheckHealth", "HEAD", "/cashdesk", nil).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type CheckHealthResponse struct {
	Meta golib_tools_courier.Metadata
	Body string
}

type CreateTransRequest struct {
	// ?????????access key
	AccessKey string `in:"header" name:"AccessKey" validate:"@string[1,128]"`
	// ????????????
	SignAlgorithm CashdeskSignAlgorithm `in:"query" name:"signAlgorithm"`
	// ??????
	Sign string `in:"query" name:"sign" validate:"@string[1,32]"`
	//
	Body CreateTransReqBody `in:"body" name:"body,omitempty"`
}

func (c ClientCashdesk) CreateTrans(req CreateTransRequest) (resp *CreateTransResponse, err error) {
	resp = &CreateTransResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".CreateTrans", "POST", "/cashdesk/v0/trans", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type CreateTransResponse struct {
	Meta golib_tools_courier.Metadata
	Body CreateTransRespBody
}

type EBankPayRequest struct {
	// ????????????
	TransID uint64 `in:"path" name:"transID" validate:"@uint64[0,]"`
	//
	Body EBankPayReqBody `in:"body" name:"body,omitempty"`
}

func (c ClientCashdesk) EBankPay(req EBankPayRequest) (resp *EBankPayResponse, err error) {
	resp = &EBankPayResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".EBankPay", "POST", "/cashdesk/v0/trans/:transID/ebankpay", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type EBankPayResponse struct {
	Meta golib_tools_courier.Metadata
	Body string
}

type GetAreaBankByBankCodeRequest struct {
	// ???????????????
	BankCode string `in:"path" name:"bankCode" validate:"@string[1,]"`
}

func (c ClientCashdesk) GetAreaBankByBankCode(req GetAreaBankByBankCodeRequest) (resp *GetAreaBankByBankCodeResponse, err error) {
	resp = &GetAreaBankByBankCodeResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".GetAreaBankByBankCode", "GET", "/cashdesk/v0/bank/bankCode/:bankCode", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type GetAreaBankByBankCodeResponse struct {
	Meta golib_tools_courier.Metadata
	Body GetAreaBankByBankCodeRespBody
}

type GetAreaBankListRequest struct {
	// ??????ID
	BankID uint32 `in:"query" name:"bankID" validate:"@uint32[1,9999]"`
	// G7????????????
	AreaCode string `in:"query" name:"areaCode" validate:"@string[1,]"`
}

func (c ClientCashdesk) GetAreaBankList(req GetAreaBankListRequest) (resp *GetAreaBankListResponse, err error) {
	resp = &GetAreaBankListResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".GetAreaBankList", "GET", "/cashdesk/v0/bank/area-bank", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type GetAreaBankListResponse struct {
	Meta golib_tools_courier.Metadata
	Body GetAreaBankListRespBody
}

type GetAreaListRequest struct {
	// ??????????????????
	ParentCode string `in:"query" name:"parentCode,omitempty" validate:"@string[0,10]"`
	// ??????(??????????????????)
	Depth int32 `default:"1" in:"query" name:"depth,omitempty" validate:"@int32[1,4]"`
	// ????????????(?????????????????????????????????)
	DepthNow int32 `default:"0" in:"query" name:"depthNow,omitempty" validate:"@int32[0,4]"`
}

func (c ClientCashdesk) GetAreaList(req GetAreaListRequest) (resp *GetAreaListResponse, err error) {
	resp = &GetAreaListResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".GetAreaList", "GET", "/cashdesk/v0/bank/area", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type GetAreaListResponse struct {
	Meta golib_tools_courier.Metadata
	Body GetAreaListByParentCodeDepthRespBody
}

type GetBankByIDRequest struct {
	// ??????id
	BankID uint32 `in:"path" name:"bankID" validate:"@uint32[1,]"`
}

func (c ClientCashdesk) GetBankByID(req GetBankByIDRequest) (resp *GetBankByIDResponse, err error) {
	resp = &GetBankByIDResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".GetBankByID", "GET", "/cashdesk/v0/bank/bankID/:bankID", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type GetBankByIDResponse struct {
	Meta golib_tools_courier.Metadata
	Body Bank
}

type GetBankListRequest struct {
	// ????????????????????????
	SupportUnionpayAuth golib_tools_courier_enumeration.Bool `in:"query" name:"supportUnionpayAuth,omitempty" validate:"@string{,TRUE,FALSE}"`
	// ??????????????????????????????
	HaveSuperBankCode golib_tools_courier_enumeration.Bool `in:"query" name:"haveSuperBankCode,omitempty" validate:"@string{,TRUE,FALSE}"`
}

func (c ClientCashdesk) GetBankList(req GetBankListRequest) (resp *GetBankListResponse, err error) {
	resp = &GetBankListResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".GetBankList", "GET", "/cashdesk/v0/bank", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type GetBankListResponse struct {
	Meta golib_tools_courier.Metadata
	Body []Bank
}

type GetPabWhiteListTransRequest struct {
	// ???????????????????????? eg:2016-01-12T00:00:00Z
	CreateStartDate golib_tools_timelib.MySQLTimestamp `in:"query" name:"createStartDate,omitempty"`
	// ????????????????????????
	CreateEndDate golib_tools_timelib.MySQLTimestamp `in:"query" name:"createEndDate,omitempty"`
	// ????????????????????????10???-1???????????????
	Size int32 `default:"10" in:"query" name:"size,omitempty" validate:"@int32[-1,100]"`
	// ????????????????????????0
	Offset int32 `default:"0" in:"query" name:"offset,omitempty" validate:"@int32[0,]"`
	// ????????????ID
	BuyerAccountID uint64 `default:"0" in:"query" name:"buyerAccountID,omitempty" validate:"@uint64[0,]"`
	// ????????????ID
	SellerAccountID uint64 `default:"0" in:"query" name:"sellerAccountID,omitempty" validate:"@uint64[0,]"`
	// ??????????????????
	TransIDExt string `in:"query" name:"transIDExt,omitempty" validate:"@string[0,64]"`
	// ??????????????????
	TransID uint64 `default:"0" in:"path" name:"transID,omitempty" validate:"@uint64[0,]"`
}

func (c ClientCashdesk) GetPabWhiteListTrans(req GetPabWhiteListTransRequest) (resp *GetPabWhiteListTransResponse, err error) {
	resp = &GetPabWhiteListTransResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".GetPabWhiteListTrans", "GET", "/cashdesk/v0/trans/:transID/pabWhiteListTrans", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type GetPabWhiteListTransResponse struct {
	Meta golib_tools_courier.Metadata
	Body PabWhiteListTransResp
}

type GetSignRequest struct {
	// ????????????
	Secret string `in:"query" name:"secret" validate:"@string[1,128]"`
	// ????????????
	SignAlgorithm CashdeskSignAlgorithm `in:"query" name:"signAlgorithm"`
	//
	Body []KVPair `in:"body" name:"body,omitempty"`
}

func (c ClientCashdesk) GetSign(req GetSignRequest) (resp *GetSignResponse, err error) {
	resp = &GetSignResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".GetSign", "POST", "/cashdesk/v0/sign", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type GetSignResponse struct {
	Meta golib_tools_courier.Metadata
	Body GetSignRespBody
}

type GetTransByOrderNoRequest struct {
	// ?????????
	OrderNo string `in:"path" name:"orderNo" validate:"@string[6,32]"`
	// ??????AccessKey
	AccessKey string `in:"header" name:"accessKey" validate:"@string[1,128]"`
	// ????????????Mock(?????????????????????)
	UseMock bool `default:"false" in:"query" name:"useMock,omitempty"`
	// ??????????????????????????????(???UserMock???true??????)
	MockState CashdeskPayState `in:"query" name:"mockState,omitempty"`
}

func (c ClientCashdesk) GetTransByOrderNo(req GetTransByOrderNoRequest) (resp *GetTransByOrderNoResponse, err error) {
	resp = &GetTransByOrderNoResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".GetTransByOrderNo", "GET", "/cashdesk/v0/trans/0/orderNo/:orderNo", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type GetTransByOrderNoResponse struct {
	Meta golib_tools_courier.Metadata
	Body ReturnTransModel
}

type OldTransQueryRequest struct {
	// ?????????????????????AccessKey
	AccessKey string `in:"header" name:"accessKey" validate:"@string[1,128]"`
	// ?????????id
	TransID uint64 `in:"path" name:"transID" validate:"@uint64[1,]"`
	// ????????????Mock(?????????????????????)
	UseMock bool `default:"false" in:"query" name:"useMock,omitempty"`
	// ???????????????(UseMock???true?????????)
	BuyerID string `in:"query" name:"buyerID,omitempty"`
	// ???????????????????????????(UseMock???true?????????)
	PayStatus CashdeskPayState `in:"query" name:"payStatus,omitempty"`
}

func (c ClientCashdesk) OldTransQuery(req OldTransQueryRequest) (resp *OldTransQueryResponse, err error) {
	resp = &OldTransQueryResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".OldTransQuery", "GET", "/cashdesk/v0/user/0/oldtrans/:transID", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type OldTransQueryResponse struct {
	Meta golib_tools_courier.Metadata
	Body OldTransQueryRespBody
}

type PabPayRequest struct {
	// ?????????
	TransID uint64 `default:"0" in:"path" name:"transID,omitempty" validate:"@uint64[1,]"`
	//
	Body PabPayReqBody `in:"body" name:"body,omitempty"`
}

func (c ClientCashdesk) PabPay(req PabPayRequest) (resp *PabPayResponse, err error) {
	resp = &PabPayResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".PabPay", "POST", "/cashdesk/v0/trans/:transID/pabpay", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type PabPayResponse struct {
	Meta golib_tools_courier.Metadata
	Body PabPayRespBody
}

type PabQueryByCertificateRequest struct {
	// ????????????
	CertificateType CashdeskCertificateType `in:"query" name:"certificateType" validate:"@string{ID_CARD,ORG_INSITITUTE_CODE,UNITY_SOCIAL_CREDIT_CODE}"`
	// ????????????
	CertificateCode string `in:"query" name:"certificateCode" validate:"@string[1,32]"`
}

func (c ClientCashdesk) PabQueryByCertificate(req PabQueryByCertificateRequest) (resp *PabQueryByCertificateResponse, err error) {
	resp = &PabQueryByCertificateResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".PabQueryByCertificate", "GET", "/cashdesk/v0/user/0/pab/certificate", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type PabQueryByCertificateResponse struct {
	Meta golib_tools_courier.Metadata
	Body PabLinkModel
}

type PabQueryByUserRequest struct {
	// ????????????ID
	UserID string `in:"path" name:"userID" validate:"@string[1,32]"`
}

func (c ClientCashdesk) PabQueryByUser(req PabQueryByUserRequest) (resp *PabQueryByUserResponse, err error) {
	resp = &PabQueryByUserResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".PabQueryByUser", "GET", "/cashdesk/v0/user/:userID/pab", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type PabQueryByUserResponse struct {
	Meta golib_tools_courier.Metadata
	Body PabLinkModel
}

type TransListQueryRequest struct {
	// ????????????id
	UserID string `in:"path" name:"userID" validate:"@string[1,64]"`
	// ????????????
	StartTime golib_tools_timelib.MySQLTimestamp `in:"query" name:"startTime,omitempty"`
	// ????????????
	EndTime golib_tools_timelib.MySQLTimestamp `in:"query" name:"endTime,omitempty"`
	// ????????????
	TransState CashdeskTransState `in:"query" name:"transState,omitempty"`
	// ????????????
	TransType CashdeskTransType `in:"query" name:"transType,omitempty"`
	// ??????,?????????0
	Offset int32 `default:"0" in:"query" name:"offset,omitempty" validate:"@int32[0,]"`
	// ????????????
	Size int32 `default:"10" in:"query" name:"size,omitempty" validate:"@int32[1,50]"`
}

func (c ClientCashdesk) TransListQuery(req TransListQueryRequest) (resp *TransListQueryResponse, err error) {
	resp = &TransListQueryResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".TransListQuery", "GET", "/cashdesk/v0/user/:userID", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type TransListQueryResponse struct {
	Meta golib_tools_courier.Metadata
	Body QueryTransResponseBody
}

type UpdateTransRequest struct {
	// ????????????
	TransID uint64 `in:"path" name:"transID" validate:"@uint64[1,]"`
	// ?????????access key
	AccessKey string `in:"header" name:"AccessKey" validate:"@string[1,128]"`
	// ????????????
	SignAlgorithm CashdeskSignAlgorithm `in:"query" name:"signAlgorithm"`
	// ????????????Mock(?????????????????????)
	UseMock bool `default:"false" in:"query" name:"useMock,omitempty"`
	// ??????????????????????????????(???UserMock???true??????)
	MockState CashdeskPayState `in:"query" name:"mockState,omitempty"`
	// ??????
	Sign string `in:"query" name:"sign" validate:"@string[1,32]"`
	//
	Body UpdateTransReqBody `in:"body" name:"body,omitempty"`
}

func (c ClientCashdesk) UpdateTrans(req UpdateTransRequest) (resp *UpdateTransResponse, err error) {
	resp = &UpdateTransResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".UpdateTrans", "PUT", "/cashdesk/v0/trans/:transID", req).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type UpdateTransResponse struct {
	Meta golib_tools_courier.Metadata
	Body string
}

func (c ClientCashdesk) WithSwagger() (resp *WithSwaggerResponse, err error) {
	resp = &WithSwaggerResponse{}
	resp.Meta = golib_tools_courier.Metadata{}

	err = c.Request(c.Name+".WithSwagger", "GET", "/cashdesk", nil).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type WithSwaggerResponse struct {
	Meta golib_tools_courier.Metadata
	Body string
}
