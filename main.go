// Copyright 2015 The blockchainrpc Authors
// This file is part of the blockchainrpc library.
//
// The blockchainrpc library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The blockchainrpc library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the blockchainrpc library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/lzxm160/blockchainrpc/rpc"
)

func main() {
	server := rpc.NewServer()
	defer server.Stop()
	service := new(testService)

	if err := server.RegisterName("test", service); err != nil {
		fmt.Println(err)
		return
	}

	listener, err := net.Listen("tcp", "127.0.0.1:8545")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	go server.ServeListener(listener)

	select {}
}

type testService struct{}

type Args struct {
	S string
}

type Result struct {
	String string
	Int    int
	Args   *Args
}

func (s *testService) NoArgsRets() {}

func (s *testService) Echo(str string, i int, args *Args) Result {
	return Result{str, i, args}
}

func (s *testService) EchoWithCtx(ctx context.Context, str string, i int, args *Args) Result {
	return Result{str, i, args}
}

func (s *testService) Sleep(ctx context.Context, duration time.Duration) {
	time.Sleep(duration)
}

func (s *testService) Rets() (string, error) {
	return "", nil
}

func (s *testService) InvalidRets1() (error, string) {
	return nil, ""
}

func (s *testService) InvalidRets2() (string, string) {
	return "", ""
}

func (s *testService) InvalidRets3() (string, string, error) {
	return "", "", nil
}

func (s *testService) CallMeBack(ctx context.Context, method string, args []interface{}) (interface{}, error) {
	c, ok := rpc.ClientFromContext(ctx)
	if !ok {
		return nil, errors.New("no client")
	}
	var result interface{}
	err := c.Call(&result, method, args...)
	return result, err
}

func (s *testService) CallMeBackLater(ctx context.Context, method string, args []interface{}) error {
	c, ok := rpc.ClientFromContext(ctx)
	if !ok {
		return errors.New("no client")
	}
	go func() {
		<-ctx.Done()
		var result interface{}
		c.Call(&result, method, args...)
	}()
	return nil
}

func (s *testService) Subscription(ctx context.Context) (*rpc.Subscription, error) {
	return nil, nil
}
