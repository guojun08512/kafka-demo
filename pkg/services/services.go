package services

import (
	"keyayun.com/seal-kafka-runner/pkg/client"
	"keyayun.com/seal-kafka-runner/pkg/utils"
)

const DefaultTaskVersion = "v0.1.0"

type (
	dict  = utils.Dict
	slice = utils.Slice
)

type Service interface {
	Name() string
	Scope() []string
	Categories() []string
	Version() string
	Params() map[string][]client.Param
	DocTypes() client.DocDefs
	RootDir() string
	Triggers() dict
	RunJob(b []byte) error
}
