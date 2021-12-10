package full_rebalance

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/starslabhq/hermes-rebalance/config"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/starslabhq/hermes-rebalance/types"
	"github.com/starslabhq/hermes-rebalance/utils"
)

func getLpData(url string) (lpList *types.Data, err error) {
	data, err := utils.DoRequest(url, "GET", nil)
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
		err = fmt.Errorf("lpResponse code not 200, msg:%s", lpResponse.Msg)
		logrus.Error(err)
		return
	}
	lpList = lpResponse.Data
	return
}

func callMarginApi(url string, conf *config.Config, body interface{})(res *types.NormalResponse, err error){
	headers := make(map[string]string)
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	headers["timestamp"] = timestamp
	headers["appId"] = conf.Margin.AppID
	data, err := json.Marshal(body)
	if err != nil{
		return
	}
	sign := sign(timestamp, data, conf.Margin.SecretKey)
	headers["sign"] = sign
	resData, err := utils.DoRequestWithHeaders(url, "POST", data, headers)
	if err != nil {
		logrus.Errorf("DoRequestWithHeaders req%s: err:%v", string(data), err)
		return
	}
	lpResponse := &types.NormalResponse{}
	if err = json.Unmarshal(resData, lpResponse); err != nil {
		logrus.Errorf("unmarshar lpResponse err:%v", err)
		return
	}
	if lpResponse.Code != 200 {
		err =fmt.Errorf("do http request, code:%d, msg:%s, url:%s, requestBody:%+v", lpResponse.Code, lpResponse.Msg, url, body)
		logrus.Error(err)
		return
	}
	res = lpResponse
	return
}

func sign(timestamp string, body []byte, secretKey string) string {
	s := fmt.Sprintf("timestamp=%s&body=%s&secretKey=%s", timestamp, string(body), secretKey)
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}