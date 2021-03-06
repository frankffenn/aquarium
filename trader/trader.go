package trader

import (
	"context"
	"fmt"
	"time"

	"github.com/frankffenn/aquarium/api"
	"github.com/frankffenn/aquarium/comm"
	"github.com/frankffenn/aquarium/sdk"
	"github.com/frankffenn/aquarium/sdk/mod"
	"github.com/frankffenn/aquarium/utils/log"
	"github.com/robertkrimen/otto"
	"golang.org/x/xerrors"
)

var (
	errHalt  = xerrors.New("HALT")
	Executor = make(map[int64]*Global)
)

func Switch(id int64) error {
	if GetTraderStatus(id) {
		return stop(id)
	}
	return run(id)
}

func GetTraderStatus(id int64) bool {
	if t, ok := Executor[id]; ok && t != nil {
		return t.Status == mod.JSRunning
	}
	return false
}

func run(id int64) error {
	job, err := initialize(id)
	if err != nil {
		log.Err("initialize trader failed, %v", err)
		return err
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errw("recover", "err", err)
			}
		}()
		job.LastRunAt = time.Now()
		job.Status = mod.JSRunning
		if job.Algorithm.Script == "" {
			log.Errw("empty script", "aligorithm id", job.Algorithm.ID)
			return
		}
		if _, err := job.Ctx.Run(job.Algorithm.Script); err != nil {
			log.Err("run script failed, %v", err)
			job.log <- &mod.JobLog{
				Type:    mod.LogTypeError,
				UserID:  job.UserID,
				JobID:   job.ID,
				Content: err.Error(),
			}
			return
		}
		main, err := job.Ctx.Get("main")
		if err != nil {
			log.Err("Can not get the main function")
			return
		}
		if _, err := main.Call(main); err != nil {
			log.Err("call main function failed,%v", err)
			return
		}
	}()
	Executor[job.ID] = job
	Executor[id].log <- &mod.JobLog{
		Type:    mod.LogTypeStart,
		UserID:  job.UserID,
		JobID:   job.ID,
		Content: "Server Start",
	}
	return nil
}

func stop(id int64) error {
	t, ok := Executor[id]
	if !ok || t == nil {
		return xerrors.New("Can not found the Trader")
	}
	Executor[id].Ctx.Interrupt <- func() { panic(errHalt) }
	Executor[id].Job.Status = mod.JSStop

	Executor[id].log <- &mod.JobLog{
		Type:    mod.LogTypeStop,
		UserID:  t.UserID,
		JobID:   t.ID,
		Content: "Server Stop",
	}

	return nil
}

func initialize(id int64) (*Global, error) {
	if t := Executor[id]; t != nil && t.Status == mod.JSRunning {
		return nil, nil
	}

	ctx := context.Background()
	job, err := sdk.GetJobByID(ctx, id)
	if err != nil {
		log.Err("get job by id failed,%v", err)
		return nil, err
	}

	if job.AlgorithmID <= 0 {
		return nil, xerrors.New("Please select a algorithm")
	}

	job.Algorithm, err = sdk.GetAlgorithmByID(ctx, job.AlgorithmID)
	if err != nil {
		log.Err("get algorithm by id failed,%v", err)
		return nil, err
	}

	e, err := sdk.GetExchangeByID(ctx, job.ExchangeID)
	if err != nil {
		log.Err("get exchange by id failed,%v", err)
		return nil, err
	}

	ex := createExchange(
		comm.ExchangeType(e.Type),
		api.JobID(job.ID),
		api.Name(e.Name),
		api.Type(e.Type),
		api.AccessKey(e.AccessKey),
		api.SecretKey(e.SecretKey),
	)

	global := &Global{
		Job:   job,
		tasks: make([]task, 0),
		Ctx:   otto.New(),
		ex:    ex,
		log:   make(chan *mod.JobLog),
	}
	for _, c := range comm.Consts {
		global.Ctx.Set(c, c)
	}
	go global.RecordLog()

	global.Ctx.Interrupt = make(chan func(), 1)
	global.Ctx.Set("Global", global)
	global.Ctx.Set("G", global)
	global.Ctx.Set("Exchange", global.ex)
	global.Ctx.Set("E", global.ex)
	global.Ctx.Set("__log__", func(call otto.FunctionCall) otto.Value {
		m := ""
		for _, v := range call.ArgumentList {
			m = fmt.Sprintf("%s %s", m, v.String())
		}
		global.log <- &mod.JobLog{
			Type:    mod.LogTypeInfo,
			UserID:  job.UserID,
			JobID:   job.ID,
			Content: m,
		}
		return otto.Value{}
	})
	global.Ctx.Run("console.log = __log__;")

	return global, nil
}

func createExchange(t comm.ExchangeType, opts ...api.Option) api.Exchange {
	switch t {
	case comm.Huobi:
		return api.NewHuobi(opts...)
	default:
	}
	return nil
}

func clean(userID int64) {
	for _, t := range Executor {
		if t != nil && t.UserID == userID {
			stop(t.ID)
		}
	}
}

func (g *Global) RecordLog() {
	for {
		select {
		case l := <-g.log:
			sdk.AddJobLog(context.Background(), l)
			if l.Type == mod.LogTypeStop {
				return
			}
		}
	}
}
