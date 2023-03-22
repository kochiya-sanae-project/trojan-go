package hydra

import (
	"context"
	"github.com/p4gefau1t/trojan-go/config"
	"github.com/p4gefau1t/trojan-go/log"
	"github.com/p4gefau1t/trojan-go/statistic"
	"github.com/p4gefau1t/trojan-go/statistic/memory"
	"strconv"
	"time"
)

const Name = "HYDRA"

type Authenticator struct {
	*memory.Authenticator
	client         *RequestClient
	updateDuration time.Duration
	ctx            context.Context
}

func (a *Authenticator) updater() {
	for {
		//for _, user := range a.ListUsers() {
		//	hash := user.Hash()
		//	sent, recv := user.ResetTraffic()
		//	s := a.client.UpdateTraffic(hash, sent, recv)
		//
		//	if _, ok := s["id"]; ok {
		//		log.Info(fmt.Sprintf("transfer user#%s updated: %d/%d", hash, sent, recv))
		//	} else {
		//		a.DelUser(hash)
		//	}
		//}
		log.Info("buffered data has been written into the database")

		var subscriptions = a.client.PullSubscriptions()
		// todo make it type safe and do better error handle
		var list = subscriptions["list"].([]interface{})
		for _, subscription := range list {
			var sub = subscription.(map[string]interface{})
			var passwordHash = sub["passwordHash"].(string)
			var uploaded, _ = strconv.ParseInt(sub["uploaded"].(string), 10, 64)
			var downloaded, _ = strconv.ParseInt(sub["downloaded"].(string), 10, 64)
			var transferEnable, _ = strconv.ParseInt(sub["transferEnable"].(string), 10, 64)
			if downloaded+uploaded < transferEnable || transferEnable < 0 {
				a.AddUser(passwordHash)
			} else {
				a.DelUser(passwordHash)
			}
		}
		log.Debug("user list updated")

		select {
		case <-time.After(a.updateDuration):
		case <-a.ctx.Done():
			log.Debug("hydra exiting...")
			return
		}
	}
}

func NewAuthenticator(ctx context.Context) (statistic.Authenticator, error) {
	cfg := config.FromContext(ctx, Name).(*Config)

	client, err := NewRequestClient(ctx, cfg.Hydra.BaseUrl, cfg.Hydra.Username, cfg.Hydra.Password)
	if err != nil {
		return nil, err
	}
	client.Login()

	memoryAuth, err := memory.NewAuthenticator(ctx)
	if err != nil {
		return nil, err
	}
	a := &Authenticator{
		ctx:            ctx,
		client:         client,
		updateDuration: time.Duration(cfg.Hydra.CheckRate) * time.Second,
		Authenticator:  memoryAuth.(*memory.Authenticator),
	}
	go a.updater()
	log.Info("hydra authenticator created")
	return a, nil
}

func init() {
	statistic.RegisterAuthenticatorCreator(Name, NewAuthenticator)
}
