package tongo

import (
	"errors"
	"fmt"
	"strconv"

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
	Code  int32  `json:"code"`
}

// TonCenterClient is the client for toncenter.com HTTP API
type TonCenterClient struct {
	url   string
	token string
}

// Balance returns a wallet balance in nanotons
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

// AddressState represents the state of an address
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

// PackAddress converts an address from raw to human-readable format
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

// UnpackAddress converts an address from human-readable to raw format
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

// AddressInformation gets basic information about the address: balance, code, data, last_transaction_id.
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

// ExtendedAddressInformation is similar to previous one (AddressInformation) but tries to parse additional information for known contract types.
// This method is based on tonlib's function getAccountState. For detecting wallets we recommend to use getWalletInformation.
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

// MasterchainInfo returns up-to-date masterchain state
func (t *TonCenterClient) MasterchainInfo() (MasterchainInfo, error) {
	var resp struct {
		baseTonCenterResponse
		Result MasterchainInfo `json:"result"`
	}
	_, _, errs := t.withAuth(gorequest.New().Get(t.url + "getMasterChainInfo")).
		EndStruct(&resp)

	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return MasterchainInfo{}, err
	}

	return resp.Result, nil
}

// ConsensusBlock consensus block and its update timestamp
func (t *TonCenterClient) ConsensusBlock() (consensusBlock int64, timestamp float64, err error) {
	var resp struct {
		baseTonCenterResponse
		Result struct {
			ConsensusBlock int64   `json:"consensus_block_id"`
			Timestamp      float64 `json:"timestamp"`
		} `json:"result"`
	}
	_, _, errs := t.withAuth(gorequest.New().Get(t.url + "getConsensusBlock")).
		EndStruct(&resp)

	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return consensusBlock, timestamp, err
	}

	return resp.Result.ConsensusBlock, resp.Result.Timestamp, nil
}

// LookupBlockRequestParameters
type LookupBlockRequestParameters struct {
	SeqNo    *int64
	Lt       int64
	UnixTime int64
}

// LookupBlock by either seqno, lt or unix time
func (t *TonCenterClient) LookupBlock(
	workchain int32,
	shard int64,
	options LookupBlockRequestParameters,
) (BlockID, error) {
	var resp struct {
		baseTonCenterResponse
		Result BlockID `json:"result"`
	}
	req := t.withAuth(gorequest.New().Get(t.url + "lookupBlock"))
	if options.SeqNo != nil {
		req.Param("seqno", strconv.FormatInt(*options.SeqNo, 10))
	}
	if options.Lt != 0 {
		req.Param("lt", strconv.FormatInt(options.Lt, 10))
	}
	if options.UnixTime != 0 {
		req.Param("unixtime", strconv.FormatInt(options.UnixTime, 10))
	}

	_, _, errs := req.EndStruct(&resp)
	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return BlockID{}, err
	}

	return resp.Result, nil
}

// Shards returns shards information
func (t *TonCenterClient) Shards(seqno int64) ([]BlockID, error) {
	var resp struct {
		baseTonCenterResponse
		Result []BlockID `json:"result"`
	}
	_, _, errs := t.withAuth(gorequest.New().Get(t.url+"shards")).
		Param("seqno", strconv.FormatInt(seqno, 10)).
		EndStruct(&resp)

	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// BlockTransactionsRequestParameters
type BlockTransactionsRequestParameters struct {
	RootHash  string
	FileHash  string
	AfterLt   int64
	AfterHash string
	Count     int32
}

func (b BlockTransactionsRequestParameters) fillRequest(req *gorequest.SuperAgent) {
	if b.RootHash != "" {
		req.Param("root_hash", b.RootHash)
	}
	if b.FileHash != "" {
		req.Param("file_hash", b.FileHash)
	}
	if b.AfterLt != 0 {
		req.Param("after_lt", strconv.FormatInt(b.AfterLt, 10))
	}
	if b.AfterHash != "" {
		req.Param("after_hash", b.AfterHash)
	}
	if b.Count != 0 {
		req.Param("count", strconv.Itoa(int(b.Count)))
	}
}

// BlockTransactions returns transactions of the given block
func (t *TonCenterClient) BlockTransactions(
	workchain int64,
	shard int64,
	seqno int64,
	parameters BlockTransactionsRequestParameters,
) ([]BlockID, error) {
	var resp struct {
		baseTonCenterResponse
		Result []BlockID `json:"result"`
	}
	req := t.withAuth(gorequest.New().Get(t.url+"blockTransactions")).
		Param("workchain", strconv.FormatInt(workchain, 10)).
		Param("shard", strconv.FormatInt(shard, 10)).
		Param("seqno", strconv.FormatInt(seqno, 10))

	parameters.fillRequest(req)

	_, _, errs := req.EndStruct(&resp)
	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// WalletInformation retrieves wallet information. This method parses contract state and currently supports more wallet
// types than getExtendedAddressInformation: simple wallet, standart wallet, v3 wallet, v4 wallet.
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

// TransactionsRequestOptions contains options for transactions request
type TransactionsRequestOptions struct {
	// Maximum number of transactions in response
	Limit int32
	// Logical time of transaction to start with, must be sent with hash
	Lt int64
	// Logical time of transaction to finish with (to get tx from lt to to_lt).
	ToLt int64
	// Hash of transaction to start with, in base64 or hex encoding , must be sent with lt
	Hash string
	// By default a request is processed by any available liteserver.
	// If archival=true only liteservers with full history are used
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

func (t TransactionsRequestOptions) fillRequest(req *gorequest.SuperAgent) error {
	if err := t.validate(); err != nil {
		return err
	}
	if t.Limit > 0 {
		req.Param("limit", fmt.Sprintf("%d", t.Limit))
	}
	if t.Lt != 0 {
		req.Param("lt", fmt.Sprintf("%d", t.Lt)).Param("hash", t.Hash)
	}
	if t.ToLt != 0 {
		req.Param("to_lt", fmt.Sprintf("%d", t.ToLt))
	}
	if t.Archival {
		req.Param("archival", "true")
	}
	return nil
}

// Transactions gets transaction history of a given address
func (t *TonCenterClient) Transactions(address string, parameters TransactionsRequestOptions) ([]Transaction, error) {
	var resp struct {
		baseTonCenterResponse
		Result []Transaction `json:"result"`
	}

	req := t.withAuth(gorequest.New().Get(t.url+"getTransactions")).Param("address", address)
	if err := parameters.fillRequest(req); err != nil {
		return nil, err
	}

	_, _, errs := req.EndStruct(&resp)
	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return nil, err
	}
	return resp.Result, nil
}

// BlockHeaderRequestParameters
type BlockHeaderRequestParameters struct {
	RootHash string
	FileHash string
}

func (b BlockHeaderRequestParameters) fillRequest(req *gorequest.SuperAgent) {
	if b.RootHash != "" {
		req.Param("root_hash", b.RootHash)
	}
	if b.FileHash != "" {
		req.Param("file_hash", b.FileHash)
	}
}

// BlockHeader returns metadata of a given block
func (t *TonCenterClient) BlockHeader(
	workchain int64,
	shard int64,
	seqno int64,
	parameters BlockHeaderRequestParameters,
) (BlockHeader, error) {
	type responseData struct {
		baseTonCenterResponse
		Result BlockHeader `json:"result"`
	}

	var resp responseData
	req := t.withAuth(gorequest.New().Get(t.url+"getBlockHeader")).
		Param("workchain", strconv.FormatInt(workchain, 10)).
		Param("shard", strconv.FormatInt(shard, 10)).
		Param("seqno", strconv.FormatInt(seqno, 10))

	parameters.fillRequest(req)

	_, _, errs := req.EndStruct(&resp)
	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return BlockHeader{}, err
	}

	return resp.Result, nil
}

// TryLocateTx needs to locate outcoming transaction of destination address by incoming message
func (t *TonCenterClient) TryLocateTx(source string, destination string, createdLt int64) (Transaction, error) {
	var resp struct {
		baseTonCenterResponse
		Result Transaction `json:"result"`
	}
	_, _, errs := t.withAuth(gorequest.New().Get(t.url+"tryLocateTx")).
		Param("source", source).
		Param("destination", destination).
		Param("created_lt", strconv.FormatInt(createdLt, 10)).
		EndStruct(&resp)

	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return Transaction{}, err
	}

	return resp.Result, nil
}

// TryLocateSourceTx needs to locate incoming transaction of source address by outcoming message.
func (t *TonCenterClient) TryLocateSourceTx(source string, destination string, createdLt int64) (Transaction, error) {
	var resp struct {
		baseTonCenterResponse
		Result Transaction `json:"result"`
	}
	_, _, errs := t.withAuth(gorequest.New().Get(t.url+"tryLocateSourceTx")).
		Param("source", source).
		Param("destination", destination).
		Param("created_lt", strconv.FormatInt(createdLt, 10)).
		EndStruct(&resp)

	if err := checkErrors(errs, resp.OK, resp.Error); err != nil {
		return Transaction{}, err
	}

	return resp.Result, nil
}

func (t *TonCenterClient) withAuth(req *gorequest.SuperAgent) *gorequest.SuperAgent {
	if t.token != "" {
		return req.Set("X-API-Key", t.token)
	}
	return req
}

// NewTonCenterAnonimousClient creates new anonymous toncenter HTTP client
func NewTonCenterAnonimousClient(url string) *TonCenterClient {
	return NewTonCenterClient(url, "")
}

// NewTonCenterClient creates new toncenter HTTP client
func NewTonCenterClient(url string, token string) *TonCenterClient {
	return &TonCenterClient{
		url:   url,
		token: token,
	}
}
