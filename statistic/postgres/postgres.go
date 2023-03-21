package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/p4gefau1t/trojan-go/common"
	"github.com/p4gefau1t/trojan-go/config"
	"github.com/p4gefau1t/trojan-go/log"
	"github.com/p4gefau1t/trojan-go/statistic"
	"github.com/p4gefau1t/trojan-go/statistic/memory"
	"time"

	_ "github.com/lib/pq"
)

const Name = "POSTGRES"

type Authenticator struct {
	*memory.Authenticator
	db             *sql.DB
	updateDuration time.Duration
	ctx            context.Context
}

func (a *Authenticator) updater() {
	for {
		for _, user := range a.ListUsers() {
			hash := user.Hash()
			sent, recv := user.ResetTraffic()

			s, err := a.db.Exec(
				fmt.Sprintf(
					"UPDATE %s SET uploaded=uploaded+$1, downloaded=downloaded+$2 WHERE %s = $3;",
					pq.QuoteIdentifier("subscription_entity"),
					pq.QuoteIdentifier("passwordHash")),
				recv,
				sent,
				hash)
			if err != nil {
				log.Error(common.NewError(fmt.Sprintf("failed to update transfer data for %s into subscription table", hash)).Base(err))
				continue
			}

			if affected, err := s.RowsAffected(); err != nil {
				if affected == 0 {
					log.Info(fmt.Sprintf("del user#%s", hash))
					a.DelUser(hash)
				} else {
					log.Info(fmt.Sprintf("transfer user#%s updated: %s/%s", hash, sent, recv))
				}
			}
		}
		log.Info("buffered data has been written into the database")

		rows, err := a.db.Query(
			//"SELECT password, uploaded, downloaded, \"transferEnable\" FROM subscription_entity;"
			fmt.Sprintf(
				"SELECT %s, %s, %s, %s FROM %s;",
				pq.QuoteIdentifier("password"),
				pq.QuoteIdentifier("uploaded"),
				pq.QuoteIdentifier("downloaded"),
				pq.QuoteIdentifier("transferEnable"),
				pq.QuoteIdentifier("subscription_entity")))
		if err != nil || rows.Err() != nil {
			log.Error(common.NewError("failed to pull data from database").Base(err))
			time.Sleep(a.updateDuration)
			continue
		}
		for rows.Next() {
			var hash string
			var quota, download, upload int64
			err := rows.Scan(&hash, &quota, &download, &upload)
			if err != nil {
				log.Error(common.NewError("failed to obtain data from the query result").Base(err))
				break
			}
			if download+upload < quota || quota < 0 {
				a.AddUser(hash)
			} else {
				a.DelUser(hash)
			}
		}

		select {
		case <-time.After(a.updateDuration):
		case <-a.ctx.Done():
			log.Debug("Postgres daemon exiting...")
			return
		}
	}
}

func connectDatabase(driverName, connectUrl string) (*sql.DB, error) {
	return sql.Open(driverName, connectUrl)
}

func NewAuthenticator(ctx context.Context) (statistic.Authenticator, error) {
	cfg := config.FromContext(ctx, Name).(*Config)
	db, err := connectDatabase("postgres", cfg.Postgres.Url)

	if err != nil {
		return nil, common.NewError("Failed to connect to database server").Base(err)
	}
	memoryAuth, err := memory.NewAuthenticator(ctx)
	if err != nil {
		return nil, err
	}
	a := &Authenticator{
		db:             db,
		ctx:            ctx,
		updateDuration: time.Duration(cfg.Postgres.CheckRate) * time.Second,
		Authenticator:  memoryAuth.(*memory.Authenticator),
	}
	go a.updater()
	log.Info("postgres authenticator created")
	return a, nil
}

func init() {
	statistic.RegisterAuthenticatorCreator(Name, NewAuthenticator)
}
