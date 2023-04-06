package db

import (
	"github.com/btc-scan/types"
	"github.com/sirupsen/logrus"
)

func (m *Mysql) GetOpenedCollectTask() ([]*types.CollectTxDB, error) {
	tasks := make([]*types.CollectTxDB, 0)
	err := m.engine.Table("t_src_tx").Where("f_collect_state = ?", types.TxReadyCollectState).Find(&tasks)
	if err != nil {
		return nil, err
	}
	return tasks, err
}

func (m *Mysql) GetCollectTask(id uint64) (*types.CollectTxDB, error) {
	task := &types.CollectTxDB{}
	ok, err := m.engine.Table("t_src_tx").Where("f_id = ?", id).Limit(1).Get(task)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return task, nil
}

func (m *Mysql) GetOpenedAssemblyTasks() ([]*types.TransactionTask, error) {
	tasks := make([]*types.TransactionTask, 0)
	err := m.engine.Table("t_transaction_task").Where("f_state in (?)", types.TxInitState).Find(&tasks)
	if err != nil {
		return nil, err
	}
	return tasks, err
}

func (m *Mysql) GetOpenedSignTasks() ([]*types.TransactionTask, error) {
	tasks := make([]*types.TransactionTask, 0)
	err := m.engine.Table("t_transaction_task").Where("f_error = \"\"  and f_state in (?)", types.TxAssmblyState).Find(&tasks)
	if err != nil {
		return nil, err
	}
	return tasks, err
}

func (m *Mysql) GetOpenedBroadcastTasks() ([]*types.TransactionTask, error) {
	tasks := make([]*types.TransactionTask, 0)
	err := m.engine.Table("t_transaction_task").Where("f_error = \"\"  and f_state in (?)", types.TxSignState).Find(&tasks)
	if err != nil {
		return nil, err
	}
	return tasks, err
}

func (m *Mysql) GetOpenedCheckTasks() ([]*types.TransactionTask, error) {
	tasks := make([]*types.TransactionTask, 0)
	err := m.engine.Table("t_transaction_task").Where("f_error = \"\"  and f_state in (?)", types.TxBroadcastState).Find(&tasks)
	if err != nil {
		return nil, err
	}
	return tasks, err
}

func (m *Mysql) GetOpenedUpdateAccountTasks() ([]*types.TransactionTask, error) {
	tasks := make([]*types.TransactionTask, 0)
	err := m.engine.Table("t_transaction_task").Where("f_error = \"\" and f_state in (?)", types.TxCheckState).Find(&tasks)
	if err != nil {
		return nil, err
	}
	return tasks, err
}

func (m *Mysql) GetTaskNonce(from string) (*types.TransactionTask, error) {
	task := &types.TransactionTask{}
	ok, err := m.engine.Table("t_transaction_task").Where("f_from = ? and f_state >= ?", from, types.TxBroadcastState).Desc("f_nonce").Limit(1).Get(task)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return task, nil
}

func (m *Mysql) GetTokenInfo(contratAddr string, chain string) (*types.Token, error) {
	token := &types.Token{}
	ok, err := m.engine.Table("t_token").Where("f_address = ? and f_chain = ?", contratAddr, chain).Limit(1).Get(token)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return token, nil
}

// 得到当前存储的任务高度
func (m *Mysql) GetTaskHeight(taskName string) (taskHeight uint64, err error) {
	height := 0
	_, err = m.engine.SQL("select num from t_task where name = ?", taskName).Get(&height)
	if err != nil {
		logrus.Error(err)
	}
	return uint64(height), err
}

// 查询公钥hash，得到对应的UID
func (db *Mysql) GetMonitorInfo(pubhash string) (string, string, string, error) {
	monitor := &types.Monitor{}
	ok, err := db.engine.Table("t_monitor").Where("f_pubhash = ?", pubhash).Limit(1).Get(monitor)
	if err != nil {
		return "", "", "", err
	}
	if !ok {
		return "", "", "", nil
	}
	return monitor.Uid, monitor.Addr, monitor.AppId, nil
}
