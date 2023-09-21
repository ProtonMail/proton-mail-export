// Copyright (c) 2023 Proton AG
//
// This file is part of Proton Export Tool.
//
// Proton Mail Bridge is Free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Proton Mail Bridge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Proton Export Tool.  If not, see <https://www.gnu.org/licenses/>.

package internal

import "sync"

type Handle[T any] int

// HandleMap exists to bridge the GAP between C interface and the Go Code. Since we are not allowed to pass
// go pointers to the "outside" world, we pass handle instead. On each C call, we then Resolve the handle.
type HandleMap[T any] struct {
	sync      sync.RWMutex
	instances []*T
}

func NewHandleMap[T any](capacity int) *HandleMap[T] {
	return &HandleMap[T]{
		instances: make([]*T, 0, capacity),
	}
}

func (a *HandleMap[T]) Alloc(i *T) Handle[T] {
	a.sync.Lock()
	defer a.sync.Unlock()

	a.instances = append(a.instances, i)

	// Return index +1 so we can still do null checks with c pointers.
	return Handle[T](len(a.instances))
}

func (a *HandleMap[T]) Free(h Handle[T]) {
	a.sync.Lock()
	defer a.sync.Unlock()

	index := h - 1

	a.instances[index] = nil
}

func (a *HandleMap[T]) Resolve(h Handle[T]) (*T, bool) {
	a.sync.RLock()
	defer a.sync.RUnlock()

	index := h - 1
	instance := a.instances[index]

	return instance, instance != nil
}
