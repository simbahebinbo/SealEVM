/*
 * Copyright 2020 The sealevm Authors
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package instructions

import (
	"github.com/simbahebinbo/sealevm/opcodes"
	"github.com/simbahebinbo/sealevm/storage"
)

func loadStorage() {
	instructionTable[opcodes.SLOAD] = opCodeInstruction{
		action:            sLoadAction,
		requireStackDepth: 1,
		enabled:           true,
	}

	instructionTable[opcodes.SSTORE] = opCodeInstruction{
		action:            sStoreAction,
		requireStackDepth: 2,
		enabled:           true,
		isWriter:          true,
	}
}

func sLoadAction(ctx *instructionsContext) ([]byte, error) {
	k := ctx.stack.Peek()

	v, err := ctx.storage.XLoad(ctx.environment.Contract.Namespace, k, storage.SStorage)
	if err != nil {
		return nil, err
	}

	k.Set(v.Int)
	return nil, nil
}

func sStoreAction(ctx *instructionsContext) ([]byte, error) {
	k := ctx.stack.Pop()
	v := ctx.stack.Pop()

	ctx.storage.XStore(ctx.environment.Contract.Namespace, k, v, storage.SStorage)
	return nil, nil
}
