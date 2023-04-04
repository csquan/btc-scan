package db

import (
	"time"

	"github.com/btc-scan/config"
	"github.com/btc-scan/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	"xorm.io/core"
)

type Mysql struct {
	conf   *config.DataBaseConf
	engine *xorm.Engine
}

func NewMysql(conf *config.DataBaseConf) (m *Mysql, err error) {
	//"test:123@/test?charset=utf8"
	engine, err := xorm.NewEngine("mysql", conf.DB)
	if err != nil {
		logrus.Errorf("create engine error: %v", err)
		return
	}
	engine.ShowSQL(false)
	engine.Logger().SetLevel(core.LOG_DEBUG)
	location, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}
	engine.SetTZLocation(location)
	engine.SetTZDatabase(location)
	m = &Mysql{
		conf:   conf,
		engine: engine,
	}
	return
}

func (m *Mysql) GetEngine() *xorm.Engine {
	return m.engine
}

func (m *Mysql) GetSession() *xorm.Session {
	return m.engine.NewSession()
}

func (m *Mysql) SaveTxTask(itf xorm.Interface, tasks *types.TransactionTask) (err error) {
	_, err = itf.Insert(tasks)
	if err != nil {
		logrus.Errorf("insert transaction task error:%v, tasks:%v", err, tasks)
	}
	return
}

func (m *Mysql) UpdateTransactionTask(itf xorm.Interface, task *types.TransactionTask) error {
	_, err := itf.Table("t_transaction_task").Where("f_id = ?", task.ID).Update(task)
	return err
}
func (m *Mysql) UpdateTransactionTaskMessage(taskID uint64, message string) error {
	_, err := m.engine.Exec("update t_transaction_task set f_message = ? where f_id = ?", message, taskID)
	return err
}

func (m *Mysql) UpdateTransactionTaskState(taskID uint64, state int) error {
	_, err := m.engine.Exec("update t_transaction_task set f_state = ? where f_id = ?", state, taskID)
	return err
}

func (m *Mysql) UpdateAccount(amount string, receiver string, contractAddr string) error {
	_, err := m.engine.Exec("update t_account set f_balance = ? where f_addr = ? and f_contractAddr = ?", amount, receiver, contractAddr)
	return err
}

func (m *Mysql) InsertCollectTx(itf xorm.Interface, task *types.CollectTxDB) (err error) {
	_, err = itf.Insert(task)
	if err != nil {
		logrus.Errorf("insert collect task error:%v, tasks:%v", err, task)
	}
	return
}

func (m *Mysql) InsertCollectSubTx(itf xorm.Interface, task *types.TransactionTask) (err error) {
	_, err = itf.Insert(task)
	if err != nil {
		logrus.Errorf("insert collect sub task error:%v, tasks:%v", err, task)
	}
	return
}

func (m *Mysql) UpdateCollectTx(itf xorm.Interface, task *types.CollectTxDB) error {
	_, err := itf.Table("t_src_tx").Where("f_id = ?", task.Base.ID).Update(task)
	return err
}

func (m *Mysql) UpdateCollectTxState(ID uint64, state int) error {
	_, err := m.engine.Exec("update t_src_tx set f_collect_state = ? where f_id = ?", state, ID)
	return err
}

func (m *Mysql) UpdateCollectSubTask(itf xorm.Interface, task *types.CollectTxDB) error {
	_, err := itf.Table("t_src_tx").Where("f_id = ?", task.Base.ID).Update(task)
	return err
}
