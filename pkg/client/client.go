/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package client

import (
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/farisekananda/kitexcall/pkg/config"
	"github.com/farisekananda/kitexcall/pkg/errors"
)

type Client interface {
	Init(conf *config.Config) error
	Call() error
	HandleBizError(bizErr kerrors.BizStatusErrorIface) error
	Output() error
	GetResponse() interface{}
	GetMetaBackward() map[string]string
}

func InvokeRPC(conf *config.Config) (Client, error) {
	var c Client
	switch conf.Type {
	case config.Thrift:
		c = NewThriftGeneric()
	case config.Protobuf:
		c = NewPbGeneric()
	default:
		c = NewThriftGeneric()
	}

	if err := c.Init(conf); err != nil {
		return nil, errors.New(errors.ClientError, "Client init failed: %v", err)
	}

	if err := c.Call(); err != nil {
		// Handle Biz error
		bizErr, isBizErr := kerrors.FromBizStatusError(err)
		if isBizErr {
			if err := c.HandleBizError(bizErr); err != nil {
				return nil, errors.New(errors.OutputError, "BizError parse error: %v", err)
			}
			return c, nil
		}
		return nil, errors.New(errors.ServerError, "RPC call failed: %v", err)
	}

	if err := c.Output(); err != nil {
		return nil, errors.New(errors.OutputError, "Response parse error: %v", err)
	}

	return c, nil
}
