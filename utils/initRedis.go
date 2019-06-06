package utils

import (
	"git.ringcentral.com/archops/goFsync/models"
	"github.com/go-acme/lego/log"
	"github.com/gomodule/redigo/redis"
)

func InitRedis(cfg *models.Config) {
	conn, err := redis.DialURL("redis://redis")
	if err != nil {
		log.Fatal(err)
	}
	cfg.Web.Redis = conn
}
