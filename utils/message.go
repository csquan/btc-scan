package utils

import (
	"encoding/json"

	"github.com/btc-scan/types"
	"github.com/sirupsen/logrus"
)

type message struct {
	Status  string
	Content string
}

func GenTxMessage(state int, content string) string {
	status := types.TransactionStateName[types.TransactionState(state)]
	return genMessage(status, content)
}

func genMessage(status, content string) string {
	message := &message{status, content}
	data, err := json.Marshal(message)
	if err != nil {
		logrus.Warnf("genMessage marshal err:%v, data:%+v", err, message)
	}
	return string(data)
}
