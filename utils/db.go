package utils

import (
	"github.com/btc-scan/types"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
)

func CommitWithSession(db types.IDB, executeFunc func(*xorm.Session) error) (err error) {
	session := db.GetSession()
	err = session.Begin()
	if err != nil {
		logrus.Errorf("begin session error:%v", err)
		return
	}

	defer session.Close()

	err = executeFunc(session)
	if err != nil {
		logrus.Errorf("execute func error:%v", err)
		err1 := session.Rollback()
		if err1 != nil {
			logrus.Errorf("session rollback error:%v", err1)
		}
		return
	}

	err = session.Commit()
	if err != nil {
		logrus.Errorf("commit session error:%v", err)
	}

	return
}
