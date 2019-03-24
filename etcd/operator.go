package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	etcdcv3 "go.etcd.io/etcd/clientv3"
)

type reslove struct {
	Ip []string `json:"ip"`
}

func (c *EtcdClient) GetValue(path string) (*etcdcv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(c.Ctx, etcdTimeout)
	defer cancel()
	fmt.Printf("%s\n", path)
	fmt.Printf("====================\n")
	r, err := c.Cli.Get(ctx, path)
	if err != nil {
		fmt.Printf("Client.Get Error \n")
		fmt.Printf("%s", err)
		fmt.Printf("====================\n")
		return nil, err
	}
	if r.Count == 0 {
		fmt.Printf("count = 0\n")
		fmt.Printf("====================\n")
		return nil, errKeyNotFound
	}
	return r, nil
}

func (c *EtcdClient) PutValue(path string, value interface{}) (*etcdcv3.PutResponse, error) {
	ctx, cancel := context.WithTimeout(c.Ctx, etcdTimeout)
	defer cancel()
	fmt.Printf("%s\n", path)
	fmt.Printf("====================\n")
	fmt.Printf("%s\n", value)
	data, err := json.Marshal(value)
	if err != nil {
		fmt.Println("json err:", err)
		return nil, fmt.Errorf("json parse err:%s", err.Error())
	}
	r, err := c.Cli.Put(ctx, path, string(data))
	if err != nil {
		fmt.Printf("Client. Error \n")
		fmt.Printf("%s", err)
		fmt.Printf("====================\n")
		return nil, err
	}
	return r, nil
}

func (c *EtcdClient) ParseReslove(response *etcdcv3.GetResponse, data interface{}) (interface{}, error) {
	if response == nil {
		fmt.Printf("this etcd data is %s\n", response)
		return data, nil
	}
	fmt.Printf("this value is %s\n", response.Kvs[0].Value)
	fmt.Printf("====================\n")
	if err := json.Unmarshal(response.Kvs[0].Value, data); err != nil {
		return nil, fmt.Errorf("json parse err:%s", err.Error())
	}
	fmt.Println(reflect.TypeOf(data))
	fmt.Printf("this struct is %s\n", data)
	fmt.Printf("====================\n")
	return data, nil
}

func (c *EtcdClient) ParseData(data reslove) (interface{}, error) {
	fmt.Printf("this struct is %s\n", data)
	fmt.Printf("====================\n")
	reslove, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json err:", err)
		return nil, fmt.Errorf("json parse err:%s", err.Error())
	}
	fmt.Printf("this value is %s\n", reslove)
	fmt.Printf("====================\n")
	return string(reslove), nil
}
