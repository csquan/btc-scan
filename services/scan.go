package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/btc-scan/config"
	"github.com/btc-scan/kafka"
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
	db        types.IDB
	config    *config.Config
	monitorDb types.IDB
	kafka     *kafka.PushKafkaService
}

func NewScanService(db types.IDB, monitorDb types.IDB, c *config.Config) *ScanService {
	scanService := ScanService{
		db:        db,
		monitorDb: monitorDb,
		config:    c,
	}
	p, err := kafka.NewSyncProducer(c.Kafka)
	if err != nil {
		return nil
	}

	scanService.kafka, err = kafka.NewPushKafkaService(c, p)
	if err != nil {
		return nil
	}
	scanService.kafka.TopicTx = c.Kafka.TopicTx
	scanService.kafka.TopicMatch = c.Kafka.TopicMatch

	return &scanService
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
	height1 := gjson.Get(res, "height")
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

	for true {
		if uint64(startHeight) <= chainHeight {
			c.ParseBlock(startHeight)
			//更新数据库高度
			c.db.UpdateTaskHeight(startHeight, "BTC")
			startHeight = startHeight + 1
		} else {
			logrus.Info("开始休眠5分钟")
			time.Sleep(5 * time.Minute) //比特币平均出块10分钟
		}

		chainHeight, err = c.GetLastBlockHeight()
		if err != nil {
			return
		}
	}

	return
}

func (c ScanService) ParseBlock(startHeight uint64) error {
	block, err := c.getBlock(startHeight)
	if err != nil {
		logrus.Error(err)
	}
	blocks := gjson.Get(block, "blocks")
	ret := make([]*types.BtcBlocks, 0)
	bb := []byte(blocks.String())
	json.Unmarshal(bb, &ret)

	//优化：这里每个区块取一次db，将btc的都取出，下面直接匹配
	monitors, err := c.monitorDb.GetChainMonitor("btc")
	if err != nil {
		logrus.Error(err)
		return err
	}

	for _, btcBlock := range ret { //P2PKH 每个tx中的锁定脚本中格式 OP_DUP OP_HASH160 <Public Key Hash> OP_EQUALVERIFY OP_CHECKSIG
		logrus.Info("当前交易条数：", len(btcBlock.BtcTxs))
		for _, tx := range btcBlock.BtcTxs {
			outLen := len(tx.TxOut)
			count := 0
			for _, out := range tx.TxOut {
				count = count + 1
				//下面根据公钥hash找到UID
				if len(out.Script) < 44 { // 只关注长度
					logrus.Info("不是关注的P2PKH交易")
					continue
				}
				pubhash := out.Script[4:44]

				uid := ""
				addr := ""
				apiKey := ""

				for _, monitor := range monitors {
					if monitor.PubHash == pubhash {
						logrus.Info("匹配到地址公钥")
						uid = monitor.Uid
						addr = monitor.Addr
						apiKey = monitor.AppId
					}
				}

				if len(uid) > 0 {
					logrus.Info("get kafka data ++")
					logrus.Info("找到uid,当前out索引为:", count, " 每条交易的out条数:", outLen)
					//对于优化后的归集，这里仅仅是一个归集通知
					txKakfa := &types.TxKakfa{
						From:           addr,
						To:             addr,
						UID:            uid,
						ApiKey:         apiKey,
						Amount:         "0",
						TokenType:      1,
						TxHash:         "",
						Chain:          "btc",
						AssetSymbol:    "btc",
						Decimals:       18,
						TxHeight:       startHeight,
						CurChainHeight: startHeight + uint64(c.config.Chains["btc"].Delay),
					}
					bb, err := json.Marshal(txKakfa)
					if err != nil {
						logrus.Warnf("Marshal txErc20s err:%v", err)
					}

					//push tx to kafka
					err = c.PushKafka(bb, c.kafka.TopicTx)
					if err != nil {
						logrus.Error(err)
					}
					logrus.Info("push kafka success ++")
				} else {
					logrus.Info("找不到uid,当前out索引为:", count, " 每条交易的out条数:", outLen)
				}
			}
		}
	}

	return nil
}

func (c *ScanService) PushKafka(bb []byte, topic string) error {
	entool, err := utils.EnTool(c.config.Ery.PUB)
	if err != nil {
		return err
	}
	//加密
	out, err := entool.ECCEncrypt(bb)
	if err != nil {
		return err
	}

	err = c.kafka.Pushkafka(out, topic)
	return err
}

func (c ScanService) Name() string {
	return "btc-scan"
}
