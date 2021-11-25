package bridge

type Task struct {
	TaskNo         uint64
	FromAccountId  uint64
	ToAccountId    uint64
	FromCurrencyId uint64
	ToCurrencyId   uint64
	Amount         uint64
}

type AccountAdd struct {
	AccounType uint8  `json:"type"` //(1中转账户,2业务账户,3合约账户,4出入口钱包)
	ChainId    uint8  `json:"chainId"`
	Address    string `json:"address"`
	Account    uint64 `json:"account"` //交易所账号
	APIKey     string `json:"apiKey"`  //交易所apikey
}

type AccountRet struct {
	Account     string `json:"account"`
	Address     string `json:"address"`
	ChainId     uint8  `json:"chainId"`
	ChainName   string `json:"chainName"`
	ChainType   uint8  `json:"chainType"`
	AccountId   uint64 `json:"accountId"`
	AccountType uint8  `json:"type"`
}

type Chain struct {
	ChainId   int    `json:"chainId"`
	IsNode    int    `json:"isNode"`
	Name      string `json:"name"`
	Status    int    `json:"status"`
	ChainType int    `json:"type"`
}

type ChainListRet struct {
	Code int                 `json:"code"`
	Data map[string][]*Chain `json:"data"`
}

type Currency struct {
	CurrencyId uint64           `json:"currencyId"`
	Currency   string           `json:"currency"`
	Tokens     []*CurrencyToken `json:"chainList"`
}

type CurrencyToken struct {
	ChainId         uint64 `json:"chainId"`
	ContractAddress string `json:"contractAddress"`
	Decimals        uint64 `json:"decimals"`
	Symbol          string `json:"symbol"`
}

type CurrencyList struct {
	Code int                    `json:"code"`
	Data map[string][]*Currency `json:"data"`
}

type AccountAddResult struct {
	AccountId uint64 `json:"accountId"`
}

type AccountAddRet struct {
	Code int               `json:"code"`
	Data *AccountAddResult `json:"data"`
}

type TaskAddResult struct {
	TaskId uint64 `json:"taskId"`
}

type TaskAddRet struct {
	Code int            `json:"code"`
	Data *TaskAddResult `json:"data"`
}

type EstimateTaskResult struct {
	TotalQuota  uint64   `json:"totalQuota"`
	SingleQuota uint64   `json:"singleQuota"`
	Routes      []string `json:"routes"`
}

type EstimateTaskRet struct {
	Code int                 `json:"code"`
	Data *EstimateTaskResult `json:"data"`
}

type TaskDetailResult struct {
	Amount        string `json:"amount"`
	CurrencyId    uint64 `json:"currencyId"`
	DstAmount     string `json:"dstAmount"`
	FromAccountId uint64 `json:"fromAccountId"`
	Status        int    `json:"status"`
	TaskId        uint64 `json:"taskId"`
	TaskNo        uint64 `json:"taskNo"`
	ToAccountId   uint64 `json:"toAccountId"`
}

type TaskDetailRet struct {
	Code int               `json:"code"`
	Data *TaskDetailResult `json:"data"`
}
