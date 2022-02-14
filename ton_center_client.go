package tongo

import (
	"errors"
	"fmt"

	"github.com/parnurzeal/gorequest"
	"github.com/shopspring/decimal"
)

func joinErrors(errs []error) error {
	return fmt.Errorf("%v", errs)
}

func checkErrors(errs []error, ok bool, errorStr string) error {
	if len(errs) > 0 {
		return joinErrors(errs)
	}
	if !ok {
		return fmt.Errorf("reponse status is not ok, error: %s", errorStr)
	}
	return nil
}

type baseTonCenterResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

type TonCenterClient struct {
	url   string
	token string
}

func (t *TonCenterClient) Balance(address string) (balance decimal.Decimal, err error) {
	var resp struct {
		baseTonCenterResponse
		Result decimal.Decimal `json:"result"`
	}
	_, _, errs := t.withAuth(gorequest.New().Get(t.url+"getAddressBalance")).
		Param("address", address).
		EndStruct(&resp)

	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return balance, err
	}

	return resp.Result, nil
}

func (t *TonCenterClient) AddressState(address string) (string, error) {
	var resp struct {
		baseTonCenterResponse
		Result string `json:"result"`
	}
	_, _, errs := t.withAuth(gorequest.New().Get(t.url+"getAddressState")).
		Param("address", address).
		EndStruct(&resp)

	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return "", err
	}

	return resp.Result, nil
}

func (t *TonCenterClient) PackAddress(address string) (string, error) {
	var resp struct {
		baseTonCenterResponse
		Result string `json:"result"`
	}
	_, _, errs := t.withAuth(gorequest.New().Get(t.url+"packAddress")).
		Param("address", address).
		EndStruct(&resp)

	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return "", err
	}

	return resp.Result, nil
}

func (t *TonCenterClient) UnpackAddress(address string) (string, error) {
	var resp struct {
		baseTonCenterResponse
		Result string `json:"result"`
	}
	_, _, errs := t.withAuth(gorequest.New().Get(t.url+"unpackAddress")).
		Param("address", address).
		EndStruct(&resp)

	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return "", err
	}

	return resp.Result, nil
}

func (t *TonCenterClient) AddressInformation(address string) (AddressInformation, error) {
	var resp struct {
		baseTonCenterResponse
		Result AddressInformation `json:"result"`
	}
	_, _, errs := t.withAuth(gorequest.New().Get(t.url+"getAddressInformation")).
		Param("address", address).
		EndStruct(&resp)
	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return AddressInformation{}, err
	}
	return resp.Result, nil
}

func (t *TonCenterClient) ExtendedAddressInformation(address string) (ExtendedAddressInformation, error) {
	var resp struct {
		baseTonCenterResponse
		Result ExtendedAddressInformation `json:"result"`
	}
	_, _, errs := t.withAuth(gorequest.New().Get(t.url+"getExtendedAddressInformation")).
		Param("address", address).
		EndStruct(&resp)

	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return ExtendedAddressInformation{}, err
	}

	return resp.Result, nil
}

func (t *TonCenterClient) WalletInformation(address string) (WalletInformation, error) {
	var resp struct {
		baseTonCenterResponse
		Result WalletInformation `json:"result"`
	}
	_, _, errs := t.withAuth(gorequest.New().Get(t.url+"getWalletInformation")).
		Param("address", address).
		EndStruct(&resp)

	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return WalletInformation{}, err
	}

	return resp.Result, nil
}

type TransactionsRequestOptions struct {
	Limit    int32
	Lt       int64
	ToLt     int64
	Hash     string
	Archival bool
}

func (t TransactionsRequestOptions) validate() error {
	if t.Lt != 0 && t.Hash == "" {
		return errors.New("lt must be sent with hash")
	}
	if t.Lt > t.ToLt {
		return errors.New("lt must be > to_lt")
	}
	return nil
}

func (t TransactionsRequestOptions) fillRequest(req *gorequest.SuperAgent) (*gorequest.SuperAgent, error) {
	if err := t.validate(); err != nil {
		return req, err
	}
	if t.Limit > 0 {
		req = req.Param("limit", fmt.Sprintf("%d", t.Limit))
	}
	if t.Lt != 0 {
		req = req.Param("lt", fmt.Sprintf("%d", t.Lt)).Param("hash", t.Hash)
	}
	if t.ToLt != 0 {
		req = req.Param("to_lt", fmt.Sprintf("%d", t.ToLt))
	}
	if t.Archival {
		req = req.Param("archival", "true")
	}
	return req, nil
}

func (t *TonCenterClient) Transactions(address string, options TransactionsRequestOptions) ([]Transaction, error) {
	type responseData struct {
		baseTonCenterResponse
		Result []Transaction `json:"result"`
	}

	var resp responseData
	req := t.withAuth(gorequest.New().Get(t.url+"getTransactions")).Param("address", address)
	req, err := options.fillRequest(req)
	if err != nil {
		return nil, err
	}

	_, _, errs := req.EndStruct(&resp)
	if len(errs) > 0 {
		return nil, joinErrors(errs)
	}
	if !resp.OK {
		return nil, fmt.Errorf("reponse status is not ok: %v, error: %s", resp.OK, resp.Error)
	}
	return resp.Result, nil
}

func (t *TonCenterClient) withAuth(req *gorequest.SuperAgent) *gorequest.SuperAgent {
	if t.token != "" {
		return req.Set("X-API-Key", t.token)
	}
	return req
}

func NewTonCenterAnonimousClient(url string) *TonCenterClient {
	return NewTonCenterClient(url, "")
}

func NewTonCenterClient(url string, token string) *TonCenterClient {
	return &TonCenterClient{
		url:   url,
		token: token,
	}
}
