package types

//LP接口参数
type LPResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Ts   int64  `json:"ts"`
	Data *Data  `json:"data"`
}
type Data struct {
	VaultInfoList []*VaultInfo `json:"vaultInfoList"`
	LiquidityProviderList []*LiquidityProvider `json:"liquidityProviderList"`
}
type LiquidityProvider struct {
	Chain          string    `json:"chain"`
	ChainId        int       `json:"chainId"`
	LpSymbol       string    `json:"lpSymbol"`
	LpAmount       string    `json:"lpAmount"`
	LpTokenAddress string    `json:"lpTokenAddress"`
	LpPlatform     string    `json:"lpPlatform"`
	LpInfoList     []*LpInfo `json:"lpInfoList"`
}
type LpInfo struct {
	LpIndex           int    `json:"lpIndex"`
	LpAmount          string `json:"lpAmount"`
	BaseTokenAddress  string `json:"baseTokenAddress"`
	QuoteTokenAddress string `json:"quoteTokenAddress"`
	BaseTokenSymbol   string `json:"baseTokenSymbol"`
	QuoteTokenSymbol  string `json:"quoteTokenSymbol"`
	BaseTokenAmount   string `json:"baseTokenAmount"`
	QuoteTokenAmount  string `json:"quoteTokenAmount"`
}
type VaultInfo struct {
	TokenSymbol  string `json:"tokenSymbol"`
	ActiveAmount struct {
		BSC  *ControllerInfo `json:"BSC"`
		Heco *ControllerInfo `json:"Heco"`
		Polygon *ControllerInfo `json:"Polygon"`
	} `json:"activeAmount"`
}

type ControllerInfo struct {
	Amount            string `json:"amount"`
	ControllerAddress string `json:"controllerAddress"`
}




//平无常接口参数
type ImpermanectLostReq struct {
	BizNo  string `json:"bizNo"`
	LpList []*LpReq  `json:"lpList"`
}

type LpReq struct {
	Chain              string `json:"chain"`
	LpTokenAddress     string `json:"lpTokenAddress"`
	LpAmount           string    `json:"lpAmount"`
	Token0OriginAmount string    `json:"token0OriginAmount"`
	Token1OriginAmount string    `json:"token1OriginAmount"`
}

type NormalResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data map[string]interface{} `json:"data"`
}



