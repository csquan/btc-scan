package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/btc-scan/config"
	"github.com/btc-scan/pkg/util/ecies"
	"github.com/btc-scan/types"
	"github.com/btc-scan/utils"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"net/http"
	"strconv"
	"time"
)

type ScanService struct {
	db     types.IDB
	config *config.Config
}

func NewScanService(db types.IDB, c *config.Config) *ScanService {
	return &ScanService{
		db:     db,
		config: c,
	}
}

func (c *ScanService) getUidFromAddr(address string) (uid string, err error) {
	pubKey, err1 := ecies.PublicFromString(c.config.UserInfo.KycPubKey)
	if err1 != nil {
		logrus.Println(err)
	}

	cli := resty.New()
	cli.SetBaseURL(c.config.UserInfo.URL)

	nowStr := time.Now().UTC().Format(http.TimeFormat)
	ct, err1 := ecies.Encrypt(rand.Reader, pubKey, []byte(nowStr), nil, nil)
	if err1 != nil {
		logrus.Println(err1)
	}
	data := map[string]interface{}{
		"verified": hex.EncodeToString(ct),
		"addr":     address,
	}
	var result types.HttpData
	resp, er := cli.R().SetBody(data).SetResult(&result).Post("/api/v1/pub/i-q-user-by-addr")
	if er != nil {
		logrus.Println(err)
	}
	if resp.StatusCode() != http.StatusOK {
		logrus.Println(err)
	}
	if result.Code != 0 {
		logrus.Println(err)
	}

	return result.Data.UID, nil
}

func (c *ScanService) GetLastBlockHeight() (height uint64, err error) {
	str := c.config.Chains["btc"].RpcUrl + "/latestblock"
	res, err := utils.Get(str)
	if err != nil {
		logrus.Error(err)
	}
	height1 := gjson.Get(res, "block_index")
	return height1.Uint(), nil
}

func (c *ScanService) getBlock(height uint64) (blockStr string, err error) {
	heightStr := strconv.Itoa(int(height))
	str := c.config.Chains["btc"].RpcUrl + "/block-height/" + heightStr
	res, err := utils.Get(str)
	if err != nil {
		logrus.Error(err)
	}
	return res, nil
}

func (c *ScanService) Run() (err error) {
	taskHeight, err := c.db.GetTaskHeight("BTC")
	if err != nil {
		return
	}

	chainHeight, err := c.GetLastBlockHeight()
	if err != nil {
		return
	}

	startHeight := uint64(c.config.Chains["btc"].Delay) + taskHeight
	if uint64(startHeight) <= chainHeight {
		c.ParseBlock(startHeight)
		startHeight = startHeight + 1
	}
	return
}

func (c ScanService) ParseBlock(startHeight uint64) string {
	block, err := c.getBlock(startHeight)
	if err != nil {
		logrus.Error(err)
	}
	blocks := gjson.Get(block, "blocks")
	logrus.Info(blocks)
	ret := make([]*types.BtcBlocks, 0)
	bb := []byte(blocks.String())
	json.Unmarshal(bb, &ret)
	txs := gjson.Get(block, "tx")
	for _, tx := range txs.Array() {
		logrus.Info(tx)
	}
	return "btc-scan"
}

func (c ScanService) Name() string {
	return "btc-scan"
}
