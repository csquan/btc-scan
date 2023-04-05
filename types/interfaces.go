package types

import (
	"github.com/go-xorm/xorm"
)

//go:generate mockgen -source=$GOFILE -destination=./mock/mock_db.go -package=mock
type IReader interface {
	//查询可以进行交易组装的任务--状态为Init
	GetOpenedAssemblyTasks() ([]*TransactionTask, error)
	//查询可以进行签名的任务--状态为Assembly
	GetOpenedSignTasks() ([]*TransactionTask, error)
	//查询可以进行广播的任务--状态为sign
	GetOpenedBroadcastTasks() ([]*TransactionTask, error)
	//查询可以进行广播的任务--状态为boradcast
	GetOpenedCheckTasks() ([]*TransactionTask, error)

	GetOpenedUpdateAccountTasks() ([]*TransactionTask, error)

	GetOpenedCollectTask() ([]*CollectTxDB, error)

	GetTaskHeight(taskName string) (taskHeight uint64, err error)

	GetCollectTask(id uint64) (*CollectTxDB, error)

	UpdateTransactionTaskState(taskID uint64, state int) error

	//查询非完成状态的task
	GetTaskNonce(from string) (*TransactionTask, error)

	GetTokenInfo(contratAddr string, chain string) (*Token, error)
}

type IWriter interface {
	GetSession() *xorm.Session
	GetEngine() *xorm.Engine
	SaveTxTask(itf xorm.Interface, task *TransactionTask) (err error)
	UpdateTransactionTask(itf xorm.Interface, task *TransactionTask) error
	UpdateTransactionTaskMessage(taskID uint64, message string) error
	InsertCollectTx(itf xorm.Interface, task *CollectTxDB) (err error)
	InsertCollectSubTx(itf xorm.Interface, task *TransactionTask) (err error)

	UpdateCollectTx(itf xorm.Interface, task *CollectTxDB) error
	UpdateCollectSubTask(itf xorm.Interface, tasks *CollectTxDB) error
	UpdateCollectTxState(taskID uint64, state int) error
	UpdateAccount(amount string, receiver string, contractAddr string) error
}

type IDB interface {
	IReader
	IWriter
}

type IAsyncService interface {
	Name() string
	Run() error
}
