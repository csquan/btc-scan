package full_rebalance

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"
	"github.com/starslabhq/hermes-rebalance/config"
	"github.com/starslabhq/hermes-rebalance/types"
	"github.com/starslabhq/hermes-rebalance/utils"
)

type impermanenceLostHandler struct {
	db   types.IDB
	conf *config.Config
}

func (i *impermanenceLostHandler) Do(task *types.FullReBalanceTask) (err error) {
	lpList, err := getLp(i.conf.ApiConf.LpUrl)
	if err != nil {
		return
	}
	lpReq := lp2Req(lpList)
	if err = callImpermanentLoss(i.conf.ApiConf.MarginUrl,
		&types.ImpermanectLostReq{BizNo: fmt.Sprintf("%d", task.ID), LpList: lpReq}); err != nil {
		return
	}
	task.State = types.FullReBalanceMarginIn
	err = i.db.UpdateFullReBalanceTask(i.db.GetSession(), task)
	return
}

func (i *impermanenceLostHandler) CheckFinished(task *types.FullReBalanceTask) (finished bool, nextState types.ReBalanceState, err error) {
	finished, err = checkMarginJobStatus(i.conf.ApiConf.MarginUrl, fmt.Sprintf("%d", task.ID))
	if err != nil {
		return
	}
	return true, types.FullReBalanceClaimLP, nil
}


func checkMarginJobStatus(url string, bizNo string) (finished bool, err error) {
	req := struct {
		BizNo string `json:"bizNo"`
	}{bizNo}
	data, err := utils.DoPost(url+"status/query", req)
	if err != nil {
		logrus.Errorf("request ImpermanentLoss api err:%v", err)
		return
	}
	resp := &types.NormalResponse{}
	if err = json.Unmarshal(data, resp); err != nil {
		logrus.Errorf("unmarshar lpResponse err:%v", err)
		return
	}
	if resp.Code != 200 {
		logrus.Errorf("callImpermanentLoss code not 200, msg:%s", resp.Msg)
		return
	}
	if v, ok := resp.Data["status"]; ok {
		return v.(string) == "SUCCESS", nil
	}
	return
}

func getLp(url string) (lpList []*types.LiquidityProvider, err error) {
	data, err := utils.DoPost(url, nil)
	if err != nil {
		logrus.Errorf("request lp err:%v", err)
		return
	}
	lpResponse := &types.LPResponse{}
	if err = json.Unmarshal(data, lpResponse); err != nil {
		logrus.Errorf("unmarshar lpResponse err:%v", err)
		return
	}
	if lpResponse.Code != 200 {
		logrus.Errorf("lpResponse code not 200, msg:%s", lpResponse.Msg)
		return
	}
	lpList = lpResponse.Data.LiquidityProviderList
	return
}
func callImpermanentLoss(url string, req *types.ImpermanectLostReq) (err error) {
	data, err := utils.DoPost(url+"submit", req)
	if err != nil {
		logrus.Errorf("request ImpermanentLoss api err:%v", err)
		return
	}
	resp := &types.NormalResponse{}
	if err = json.Unmarshal(data, resp); err != nil {
		logrus.Errorf("unmarshar lpResponse err:%v", err)
		return
	}
	if resp.Code != 200 {
		logrus.Errorf("callImpermanentLoss code not 200, msg:%s", resp.Msg)
		return
	}
	return
}

func lp2Req(lpList []*types.LiquidityProvider) (req []*types.LpReq) {
	for _, lp := range lpList {
		var totalBaseAmount, totalQuoteAmount *big.Int
		for _, lpinfo := range lp.LpInfoList {
			add(totalBaseAmount, lpinfo.BaseTokenAmount)
			add(totalQuoteAmount, lpinfo.QuoteTokenAmount)
		}
		r := &types.LpReq{
			Chain:              lp.Chain,
			LpTokenAddress:     lp.LpTokenAddress,
			LpAmount:           lp.LpAmount,
			Token0OriginAmount: totalBaseAmount.String(),
			Token1OriginAmount: totalQuoteAmount.String(),
		}
		req = append(req, r)
	}
	return
}

func add(x *big.Int, y string) {
	y1, ok := new(big.Int).SetString(y, 10)
	if !ok {
		logrus.Fatalf("lpinfo to request failed, amount:%s", y)
	}
	x.Add(x, y1)
}