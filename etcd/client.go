package etcd

import (
	"context"
	"errors"
	"fmt"
	"time"

	etcdcv3 "github.com/coreos/etcd/clientv3"
)

const (
	priority    = 10  // default priority when nothing is set
	ttl         = 300 // default ttl when nothing is set
	etcdTimeout = 5 * time.Second
)

type Client struct {
	PathPrefix string
	Ctx        context.Context
	Cli        *etcdcv3.Client
}

var errKeyNotFound = errors.New("key not found")
var errParse = errors.New("parse etcd fail")

func New(ep []endPoints) *Client,error{
	cli, err := newEtcdClient(ep);
	if err != nil {
		return nil,err
	}
	return cli,nil
}

func newEtcdClient(endpoints []string) (*etcdcv3.Client, error) {
	etcdCfg := etcdcv3.Config{
		Endpoints: endpoints,
	}
	cli, err := etcdcv3.New(etcdCfg)
	if err != nil {
		return nil, err
	}
	return cli, nil
}
