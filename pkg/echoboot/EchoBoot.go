// Package main with for obtaining echo capabilities
package main

import (
	"context"

	"capnproto.org/go/capnp/v3/rpc"
	ecgo "github.com/hiveot/echorpc/capnp/go"
)

// EchoCapClient is a client to obtain echo capabilities
type EchoCapClient struct {
	connection    *rpc.Conn
	ctx           context.Context
	ctxCancel     context.CancelFunc
	echocapClient ecgo.EchoBootstrap
}

func (ecc *EchoCapClient) GetEcho() (*ecgo.EchoService, error) {
	var err error

	resp, release := ecc.echocapClient.GetEcho(ecc.ctx,
		func(ecgo.EchoBootstrap_getEcho_Params) error {
			return nil
		})

	result, err := resp.Struct()
	if err != nil {
		return nil, err
	}

	defer release()
	echoSvc := result.Service()
	return &echoSvc, err
}
