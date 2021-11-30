package services

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	"github.com/starslabhq/hermes-rebalance/config"
	signer "github.com/starslabhq/hermes-rebalance/sign"
	"github.com/starslabhq/hermes-rebalance/types"
	"github.com/starslabhq/hermes-rebalance/utils"
	"strings"
	"time"
)

type Transaction struct {
	db        types.IDB
	config    *config.Config
	clientMap map[string]*ethclient.Client
}

func NewTransactionService(db types.IDB, conf *config.Config) (p *Transaction, err error) {
	p = &Transaction{
		db:        db,
		config:    conf,
		clientMap: utils.ClientMap,
	}
	return
}

func (t *Transaction) Name() string {
	return "transaction"
}
func (t *Transaction) Run() (err error) {
	tasks, err := t.db.GetOpenedTransactionTask()
	if err != nil {
		return
	}
	if len(tasks) == 0 {
		logrus.Infof("no available part Transaction task.")
		return
	}

	for _, task := range tasks {
		switch types.TransactionState(task.State) {
		case types.TxUnInitState:
			return t.handleSign(task)
		case types.TxAuditState:
			return t.handleAudit(task)
		case types.TxValidatorState:
			return t.handleValidator(task)
		case types.TxSignedState:
			return t.handleTransactionSigned(task)
		case types.TxCheckReceiptState:
			return t.handleTransactionCheck(task)
		default:
			logrus.Errorf("unkonwn task state [%v] for task [%v]", tasks[0].State, tasks[0].ID)
		}
	}
	return
}

func (t *Transaction) handleSign(task *types.TransactionTask) (err error) {
	nonce := task.Nonce
	input := task.InputData
	decimal := 18
	from := task.From
	to := task.To
	GasLimit := "2000000"
	GasPrice := task.GasPrice
	Amount := "0"
	quantity := "0"
	receiver := task.To //和to一致

	signRet, err := signer.SignTx(input, decimal, int(nonce), from, to, GasLimit, GasPrice, Amount, quantity, receiver)

	if err == nil && signRet.Result == true {
		err = utils.CommitWithSession(t.db, func(session *xorm.Session) (execErr error) {
			task.State = int(types.TxAuditState)
			task.Cipher = signRet.Data.Extra.Cipher
			task.EncryptData = signRet.Data.EncryptData
			task.Hash = signRet.Data.Extra.TxHash
			execErr = t.db.UpdateTransactionTask(session, task)
			if execErr != nil {
				logrus.Errorf("update sign task error:%v task:[%v]", err, task)
				return
			}
			return
		})
	}

	return err
}

func (t *Transaction) handleAudit(task *types.TransactionTask) (err error) {
	input := task.InputData
	quantity := "0"
	receiver := task.To
	orderID := time.Now().UnixNano() / 1e6 //毫秒

	auditRet, err := signer.AuditTx(input, receiver, quantity, orderID)

	if err == nil && auditRet.Success == true {
		err = utils.CommitWithSession(t.db, func(session *xorm.Session) (execErr error) {
			task.State = int(types.TxValidatorState)
			task.OrderId = orderID
			execErr = t.db.UpdateTransactionTask(session, task)
			if execErr != nil {
				logrus.Errorf("update  audit task error:%v task:[%v]", err, task)
				return
			}
			return
		})
	}
	return err
}

func (t *Transaction) handleValidator(task *types.TransactionTask) (err error) {
	vRet, err := signer.ValidatorTx(task)

	if err ==  nil && vRet.OK == true {
		err = utils.CommitWithSession(t.db, func(session *xorm.Session) (execErr error) {
			task.State = int(types.TxSignedState)
			task.SignData = vRet.RawTx
			execErr = t.db.UpdateTransactionTask(session, task)
			if execErr != nil {
				logrus.Errorf("update  validator task error:%v task:[%v]", err, task)
				return
			}
			return
		})
	}else{
		_ = utils.CommitWithSession(t.db, func(session *xorm.Session) (execErr error) {
			task.State = int(types.TxAuditState) //失败了则退回安审状态，下次重新安审

			execErr = t.db.UpdateTransactionTask(session, task)
			if execErr != nil {
				logrus.Errorf("update  validator task error:%v task:[%v]", err, task)
				return
			}
			return
		})
	}

	return err
}

func (t *Transaction) handleTransactionSigned(task *types.TransactionTask) error {
	client, ok := t.clientMap[task.ChainName]
	if !ok {
		logrus.Fatalf("not find chain client, task:%v", task)
	}
	transaction := &etypes.Transaction{}
	var input string

	input = task.SignData

	if strings.Contains(task.SignData, "0x") {
		input = task.SignData[2:]
	}

	rawTxBytes, err := hex.DecodeString(input)

	rlp.DecodeBytes(rawTxBytes, &transaction)

	if err := client.SendTransaction(context.Background(), transaction); err != nil {
		return err
	}

	err = utils.CommitWithSession(t.db, func(session *xorm.Session) (execErr error) {
		task.State = int(types.TxCheckReceiptState)
		execErr = t.db.UpdateTransactionTask(session, task)
		if execErr != nil {
			logrus.Errorf("update part audit task error:%v task:[%v]", execErr, task)
			return execErr
		}
		return nil
	})
	return err
}

func (t *Transaction) handleTransactionCheck(task *types.TransactionTask) error {
	client, ok := t.clientMap[task.ChainName]
	if !ok {
		logrus.Fatalf("not find chain client, task:%v", task)
	}
	receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(task.Hash))
	if err != nil {
		return err
	}
	if receipt == nil {
		transaction := &etypes.Transaction{}
		var input string

		input = task.SignData
		if strings.Contains(task.SignData, "0x") {
			input = task.SignData[2:]
		}
		rawTxBytes, _ := hex.DecodeString(input)

		rlp.DecodeBytes(rawTxBytes, &transaction)

		if err := client.SendTransaction(context.Background(), transaction); err != nil {
			return err
		}
		return nil
	}
	if receipt.Status == 1 {
		task.State = int(types.TxSuccessState)
	} else if receipt.Status == 0 {
		task.State = int(types.TxFailedState)
	}
	err = utils.CommitWithSession(t.db, func(session *xorm.Session) (execErr error) {
		if task.TransactionType == int(types.Approve) && task.State == int(types.TxSuccessState) {
			execErr = t.db.SaveApprove(&types.ApproveRecord{Spender: task.ContractAddress, Token: task.To, From: task.To})
			if execErr != nil {
				logrus.Fatalf("SaveApprove err:%v", err)
				return
			}
		}
		execErr = t.db.UpdateTransactionTask(session, task)
		if execErr != nil {
			logrus.Errorf("update part audit task error:%v task:[%v]", err, task)
			return execErr
		}
		return nil
	})
	return err
}
