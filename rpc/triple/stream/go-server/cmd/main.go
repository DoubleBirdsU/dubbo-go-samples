/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"strings"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	triple "dubbo.apache.org/dubbo-go/v3/protocol/triple/triple_protocol"
	"dubbo.apache.org/dubbo-go/v3/server"

	greet "github.com/apache/dubbo-go-samples/rpc/triple/stream/proto"

	"github.com/dubbogo/gost/log/logger"
)

type GreetTripleServer struct{}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: req.Name}
	return resp, nil
}

func (srv *GreetTripleServer) GreetStream(ctx context.Context, stream greet.GreetService_GreetStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if triple.IsEnded(err) {
				break
			}
			logger.Errorf("triple BidiStream recv error: %s", err)
			return err
		}
		if err := stream.Send(&greet.GreetStreamResponse{Greeting: req.Name}); err != nil {
			logger.Errorf("triple BidiStream send error: %s", err)
			return err
		}
	}
	return nil
}

func (srv *GreetTripleServer) GreetClientStream(ctx context.Context, stream greet.GreetService_GreetClientStreamServer) (*greet.GreetClientStreamResponse, error) {
	var reqs []string
	for stream.Recv() {
		reqs = append(reqs, stream.Msg().Name)
	}
	if stream.Err() != nil && !triple.IsEnded(stream.Err()) {
		logger.Errorf("triple ClientStream recv err: %s", stream.Err())
		return nil, stream.Err()
	}
	resp := &greet.GreetClientStreamResponse{
		Greeting: strings.Join(reqs, ","),
	}

	return resp, nil
}

func (srv *GreetTripleServer) GreetServerStream(ctx context.Context, req *greet.GreetServerStreamRequest, stream greet.GreetService_GreetServerStreamServer) error {
	for i := 0; i < 5; i++ {
		if err := stream.Send(&greet.GreetServerStreamResponse{Greeting: req.Name}); err != nil {
			logger.Errorf("triple ServerStream send err: %s", err)
			return err
		}
	}
	return nil
}

const (
	GroupVersionIdentifier = "g1v1"
)

type GreetTripleServerGroup1Version1 struct{}

func (g *GreetTripleServerGroup1Version1) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: GroupVersionIdentifier + req.Name}
	return resp, nil
}

func (g *GreetTripleServerGroup1Version1) GreetStream(ctx context.Context, stream greet.GreetService_GreetStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if triple.IsEnded(err) {
				break
			}
			logger.Errorf("triple BidiStream recv error: %s", err)
			return err
		}
		if err := stream.Send(&greet.GreetStreamResponse{Greeting: GroupVersionIdentifier + req.Name}); err != nil {
			logger.Errorf("triple BidiStream send error: %s", err)
			return err
		}
	}
	return nil
}

func (g *GreetTripleServerGroup1Version1) GreetClientStream(ctx context.Context, stream greet.GreetService_GreetClientStreamServer) (*greet.GreetClientStreamResponse, error) {
	var reqs []string
	for stream.Recv() {
		reqs = append(reqs, GroupVersionIdentifier+stream.Msg().Name)
	}
	if stream.Err() != nil && !triple.IsEnded(stream.Err()) {
		logger.Errorf("triple ClientStream recv err: %s", stream.Err())
		return nil, stream.Err()
	}
	resp := &greet.GreetClientStreamResponse{
		Greeting: strings.Join(reqs, ","),
	}

	return resp, nil
}

func (g *GreetTripleServerGroup1Version1) GreetServerStream(ctx context.Context, req *greet.GreetServerStreamRequest, stream greet.GreetService_GreetServerStreamServer) error {
	for i := 0; i < 5; i++ {
		if err := stream.Send(&greet.GreetServerStreamResponse{Greeting: GroupVersionIdentifier + req.Name}); err != nil {
			logger.Errorf("triple ServerStream send err: %s", err)
			return err
		}
	}
	return nil
}

func main() {
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20000),
		),
	)

	if err != nil {
		panic(err)
	}
	if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		panic(err)
	}
}
