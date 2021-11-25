package services

import (
	"encoding/json"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	"github.com/starslabhq/hermes-rebalance/config"
	"github.com/starslabhq/hermes-rebalance/types"
	"github.com/starslabhq/hermes-rebalance/utils"
)

type PartReBalance struct {
	db     types.IDB
	config *config.Config
}

func NewPartReBalanceService(db types.IDB, conf *config.Config) (p *PartReBalance, err error) {
	p = &PartReBalance{
		db:     db,
		config: conf,
	}

	return
}

func (p *PartReBalance) Name() string {
	return "part_rebalance"
}

func (p *PartReBalance) Run() (err error) {
	tasks, err := p.db.GetOpenedPartReBalanceTasks()
	if err != nil {
		return
	}

	if len(tasks) == 0 {
		logrus.Infof("no available part rebalance task.")
		return
	}

	if len(tasks) > 1 {
		logrus.Errorf("more than one rebalance tasks are being processed. tasks:%v", tasks)
	}

	switch tasks[0].State {
	case types.PartReBalanceInit:
		return p.handleInit(tasks[0])
	case types.PartReBalanceCross:
		return p.handleCross(tasks[0])
	case types.PartReBalanceTransferIn:
		return p.handleTransferIn(tasks[0])
	case types.PartReBalanceInvest:
		return p.handleInvest(tasks[0])
	default:
		logrus.Errorf("unkonwn task state [%v] for task [%v]", tasks[0].State, tasks[0].ID)
	}

	return
}

func (p *PartReBalance) handleInit(task *types.PartReBalanceTask) (err error) {
	crossBalances := make([]*types.CrossBalanceItem, 0)
	err = json.Unmarshal([]byte(task.Params), &crossBalances)
	if err != nil {
		logrus.Errorf("read task [%v] cross params error:%v", task, err)
		return
	}

	if len(crossBalances) == 0 {
		logrus.Errorf("no cross balance is found for rebalance task: [%v]", task)
		return
	}

	err = utils.CommitWithSession(p.db, func(session *xorm.Session) (execErr error) {
		baseTask := &types.BaseTask{State: int(AssetTransferInit)}
		assetTransfer := &types.AssetTransferTask{BaseTask: baseTask,
			RebalanceId: task.ID, TransferType: AssetTransferOut, Params: task.Params}
		execErr = p.db.InsertAssetTransfer(assetTransfer)
		if execErr != nil {
			logrus.Errorf("save assetTransfer task error:%v task:[%v]", err, task)
			return
		}
		task.State = types.PartReBalanceTransferIn
		execErr = p.db.UpdatePartReBalanceTask(session, task)
		if execErr != nil {
			logrus.Errorf("update part rebalance task error:%v task:[%v]", err, task)
			return
		}
		return
	})

	crossTasks := make([]*types.CrossTask, 0, len(crossBalances))
	for _, param := range crossBalances {
		crossTasks = append(crossTasks, &types.CrossTask{
			ReBalanceId:  task.ID,
			ChainFrom:    param.FromChain,
			ChainTo:      param.ToChain,
			CurrencyFrom: param.FromCurrency,
			CurrencyTo:   param.ToCurrency,
			Amount:       param.Amount,
		})
	}

	err = utils.CommitWithSession(p.db, func(session *xorm.Session) (execErr error) {
		execErr = p.db.SaveCrossTasks(session, crossTasks)
		if execErr != nil {
			logrus.Errorf("save cross task error:%v task:[%v]", err, task)
			return
		}
		task.State = types.PartReBalanceCross
		execErr = p.db.UpdatePartReBalanceTask(session, task)
		if execErr != nil {
			logrus.Errorf("update part rebalance task error:%v task:[%v]", err, task)
			return
		}

		return
	})

	return
}

func (p *PartReBalance) handleCross(task *types.PartReBalanceTask) (err error) {
	//TODO check cross task and create transferIn task
	return
}

func (p *PartReBalance) handleTransferIn(task *types.PartReBalanceTask) (err error) {

	atTasks, err := p.db.GetAssetTransferTasksWithReBalanceId(task.ID)
	if err != nil {
		logrus.Errorf("get asset transfer task error:%v", err)
		return err
	}

	if len(atTasks) == 0 {
		err = fmt.Errorf("part rebalance task [%v] has no transfer in task", task)
		return
	}

	for _, at := range atTasks {
		if at.State != AssetTransferFailed && at.State != AssetTransferSuccess {
			logrus.Debugf("asset transfer task [%v] is not finished", at)
			return
		}
	}

	//TODO check transferIn task and create farm task

	return
}

func (p *PartReBalance) handleInvest(task *types.PartReBalanceTask) (err error) {
	return
}
