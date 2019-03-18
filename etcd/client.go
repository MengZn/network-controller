package etcd

import (
	"context"
	"errors"
	"time"

	etcdcv3 "go.etcd.io/etcd/clientv3"
)

const (
	priority    = 10  // default priority when nothing is set
	ttl         = 300 // default ttl when nothing is set
	etcdTimeout = 5 * time.Second
)

type EtcdClient struct {
	PathPrefix string
	Ctx        context.Context
	Cli        *etcdcv3.Client
}

var errKeyNotFound = errors.New("key not found")
var errParse = errors.New("parse etcd fail")

func New(ep []string) (*EtcdClient, error) {
	etcdcleint := new(EtcdClient)

	etcdCfg := etcdcv3.Config{
		Endpoints: ep,
	}
	ecli, err := etcdcv3.New(etcdCfg)
	if err != nil {
		return nil, err
	}
	etcdcleint.Cli = ecli
	return etcdcleint, nil
}
