package hydra

import (
	"context"
	"github.com/p4gefau1t/trojan-go/config"
	"github.com/p4gefau1t/trojan-go/log"
	"github.com/p4gefau1t/trojan-go/statistic"
	"github.com/p4gefau1t/trojan-go/statistic/memory"
	"time"
)

const Name = "HYDRA"

type Authenticator struct {
	*memory.Authenticator
	updateDuration time.Duration
	ctx            context.Context
}

func (a *Authenticator) updater() {
	for {
		log.Info("should update traffic stats")
		log.Info("should user list")

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

	memoryAuth, err := memory.NewAuthenticator(ctx)
	if err != nil {
		return nil, err
	}
	a := &Authenticator{
		ctx:            ctx,
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
