package services

import (
	"fmt"
	"path"
	"time"

	"keyayun.com/seal-kafka-runner/pkg/client"
)

type carsService struct {
	client *client.SealClient
	Test   string
}

func NewCarsService() Service {
	return &carsService{
		Test: "111",
	}
}

func (c *carsService) Name() string {
	return "cars"
}

func (c *carsService) Scope() []string {
	return []string{}
}

func (c *carsService) Categories() []string {
	return []string{}
}

func (c *carsService) Version() string {
	return DefaultTaskVersion
}

func (c *carsService) Params() map[string][]client.Param {
	return map[string][]client.Param{
		DefaultTaskVersion: nil,
	}
}

func (c *carsService) DocTypes() client.DocDefs {
	return client.GetDocDefs([]string{})
}

func (c *carsService) RootDir() string {
	return path.Join("/mnt", c.Name(), DefaultTaskVersion)
}

func (c *carsService) Triggers() dict {
	return nil
}

func (c *carsService) RunJob(msg []byte) error {
	fmt.Println("carsService msg: ", string(msg))
	time.Sleep(10 * time.Second)
	return nil
}
