package handler

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strconv"

	"github.com/vtxnetworks/network-controller/etcd"
)

var Pool = "/ippool"
var service = "/service/"
var max = 254

var ip = "192.168."

type Checker struct {
	name   string
	client *etcd.EtcdClient
}
type reslover struct {
	Ip []string `json:"ip"`
}

func (c *Checker) IPRegister() string {
	for {
		ip = ip + strconv.Itoa(rand.Intn(max)) + "." + strconv.Itoa(rand.Intn(max))
		fmt.Printf("new the %s \n", ip)
		if c.ipChecker() {
			c.serviceRegister()
			return ip
		}
	}
	return ""
}
func (c *Checker) serviceRegister() {
	response, err := c.client.GetValue(service + c.name)
	if err != nil {
		c.client.PutValue(service+c.name, "")
	}
	rData, err := c.client.ParseReslove(response, new(reslover))
	if err != nil {
		log.Fatalf("etcd server get %v cat not reslove data ERROR: %v", service+c.name, err)
	}
	fmt.Println(reflect.TypeOf(rData))
	rr, ok := (rData).(*reslover)
	if !ok {
		log.Fatalf("the data can not parse { ip:[\"xxxxx\"]})")
	}
	fmt.Printf("%v\n", rr)
	rr.Ip = append(rr.Ip, ip)
	c.client.PutValue(service+c.name, rr)
}

func (c *Checker) ipChecker() bool {
	response, err := c.client.GetValue(Pool)
	if err != nil {
		c.client.PutValue(Pool, "")
	}
	rData, err := c.client.ParseReslove(response, new(reslover))
	if err != nil {
		log.Fatalf("etcd server get %v cat not reslove data ERROR: %v", ip, err)
	}
	fmt.Println(reflect.TypeOf(rData))
	rr, ok := (rData).(*reslover)
	if !ok {
		log.Fatalf("the data can not parse { ip:[\"xxxxx\"]})")
	}
	fmt.Printf("%v\n", rr)
	if len(rr.Ip) == 0 {
		rr.Ip = append(rr.Ip, ip)
		c.client.PutValue(Pool, rr)
	} else {
		for _, data := range rr.Ip {
			if data == ip {
				return false
			}
			rr.Ip = append(rr.Ip, ip)
			c.client.PutValue(Pool, rr)
			return true
		}
	}
	return true
}

func InitHandler(ep []string, name string, ctx context.Context) (*Checker, error) {
	cli, err := etcd.New([]string{"http://192.168.33.10:32379"})
	cli.Ctx = ctx

	return &Checker{client: cli, name: name}, err
}
