/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package command

import (
	"configcenter/src/common/blog"
	"configcenter/src/storage/mongodb"
	"configcenter/src/storage/rpc"
	"configcenter/src/storage/tmserver/core"
	"configcenter/src/storage/types"
)

func init() {
	core.GCommands.SetCommand(types.OPDistinctCode, &distinct{})
}

var _ core.SetDBProxy = (*distinct)(nil)

type distinct struct {
	dbProxy mongodb.Client
}

func (d *distinct) SetDBProxy(db mongodb.Client) {
	d.dbProxy = db
}

func (d *distinct) Execute(ctx core.ContextParams, decoder rpc.Request) (*types.OPReply, error) {
	msg := types.OPDistinctOperation{}
	reply := &types.OPReply{}
	reply.RequestID = ctx.Header.RequestID
	if err := decoder.Decode(&msg); nil != err {
		reply.Message = err.Error()
		return reply, err
	}
	blog.V(4).Infof("[MONGO OPERATION] %+v", &msg)

	var targetCol mongodb.CollectionInterface
	if nil != ctx.Session {
		targetCol = ctx.Session.Collection(msg.Collection)
	} else {
		targetCol = d.dbProxy.Collection(msg.Collection)
	}

	results, err := targetCol.Distinct(ctx, msg.Field, msg.Filter)
	reply.DistinctRes = results
	if nil == err {
		reply.Success = true
	} else {
		blog.ErrorJSON("distinct execute error.  err: %s, raw data: %s, rid:%s", err.Error(), msg, msg.RequestID)
		reply.Message = err.Error()
	}

	return reply, err
}