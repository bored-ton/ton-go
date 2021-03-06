package tongo

import "github.com/shopspring/decimal"

// TonID
type TonID struct {
	Type string `json:"@type"`
	Hash string `json:"hash"`
	Lt   string `json:"lt"`
}

// BlockID
type BlockID struct {
	Type       string `json:"@type"`
	Workchaing int32  `json:"workchain"`
	Shard      string `json:"shard"`
	SeqNo      int64  `json:"seqno"`
	RootHash   string `json:"root_hash"`
	FileHash   string `json:"file_hash"`
}

// Message describes transaction message
type Message struct {
	Source      string          `json:"source"`
	Destination string          `json:"destination"`
	Value       decimal.Decimal `json:"value"`
	ForwardFee  decimal.Decimal `json:"fwd_fee"`
	IhrFee      decimal.Decimal `json:"ihr_fee"`
	CreatedLt   string          `json:"created_lt"`
	BodyHash    string          `json:"body_hash"`
	Message     string          `json:"message"`
	MessageData MessageData     `json:"msg_data"`
}

// MessageData
type MessageData struct {
	Body      string `json:"body"`
	InitState string `json:"init_state"`
}

// Transaction
type Transaction struct {
	Address       string            `json:"address"`
	Utime         int64             `json:"utime"`
	Data          string            `json:"data"`
	TransactionID TonID             `json:"transaction_id"`
	InMessage     Message           `json:"in_msg"`
	OurMessages   []Message         `json:"out_msgs"`
	Fee           decimal.Decimal   `json:"fee"`
	StorageFee    decimal.Decimal   `json:"storage_fee"`
	OtherFees     []decimal.Decimal `json:"other_fees"`
}

// AddressInformation
type AddressInformation struct {
	Type              string          `json:"@type"`
	Balance           decimal.Decimal `json:"balance"`
	Code              string          `json:"code"`
	Data              string          `json:"data"`
	LastTransactionID TonID           `json:"last_transaction_id"`
	BlockID           BlockID         `json:"block_id"`
	FrozenHash        string          `json:"frozen_hash"`
	SyncUtime         int64           `json:"sync_utime"`
	Extra             string          `json:"@extra"`
	State             string          `json:"state"`
}

// ExtendedAddress
type ExtendedAddress struct {
	Type    string `json:"@type"`
	Address string `json:"address"`
}

// ExtendedAccountState
type ExtendedAccountState struct {
	Type     string `json:"@type"`
	WalledID string `json:"walled_id"`
	SeqNo    int64  `json:"seqno"`
}

// ExtendedAddressInformation
type ExtendedAddressInformation struct {
	Type              string               `json:"@type"`
	Address           ExtendedAddress      `json:"address"`
	Balance           decimal.Decimal      `json:"balance"`
	LastTransactionID TonID                `json:"last_transaction_id"`
	BlockID           BlockID              `json:"block_id"`
	SyncUtime         int64                `json:"sync_utime"`
	AccountState      ExtendedAccountState `json:"account_state"`
	Revision          int64                `json:"revision"`
	Extra             string               `json:"@extra"`
}

// WalletInformation
type WalletInformation struct {
	Wallet            bool            `json:"wallet"`
	Balance           decimal.Decimal `json:"balance"`
	AccountState      string          `json:"account_state"`
	WalletType        string          `json:"wallet_type"`
	SeqNo             int64           `json:"seqno"`
	LastTransactionID TonID           `json:"last_transaction_id"`
	WalletID          uint64          `json:"wallet_id"`
}

// MasterchainInfo
type MasterchainInfo struct {
	Type          string  `json:"@type"`
	LastBlockID   BlockID `json:"last"`
	StateRootHash string  `json:"state_root_hash"`
	InitBlockID   BlockID `json:"init"`
	Extra         string  `json:"@extra"`
}

// BlockHeader information
type BlockHeader struct {
	Type                   string    `json:"@type"`
	ID                     BlockID   `json:"id"`
	GlobalID               int64     `json:"global_id"`
	Version                int64     `json:"version"`
	AfterMerge             bool      `json:"after_merge"`
	AfterSplit             bool      `json:"after_split"`
	BeforeSplit            bool      `json:"before_split"`
	WantMerge              bool      `json:"want_merge"`
	WantSplit              bool      `json:"want_split"`
	ValidatorListHashShort int64     `json:"validator_list_hash_short"`
	CatchainSeqNo          int64     `json:"catchain_seqno"`
	MinRefMcSeqNo          int64     `json:"min_ref_mc_seqno"`
	IsKeyBlock             bool      `json:"is_key_block"`
	PrevKeyBlockSeqNo      int64     `json:"prev_key_block_seqno"`
	StartLt                string    `json:"start_lt"`
	EndLt                  string    `json:"end_lt"`
	VertSeqNo              int64     `json:"vert_seqno"`
	PrevBlocks             []BlockID `json:"prev_blocks"`
	Extra                  string    `json:"@extra"`
}
