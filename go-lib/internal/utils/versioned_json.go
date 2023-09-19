// Copyright (c) 2023 Proton AG
//
// This file is part of Proton Export Tool.
//
// Proton Mail Bridge is free software: you can redistribute it and/or modify
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

package utils

import (
	"encoding/json"
	"errors"
)

type VersionedJSON[T any] struct {
	Version int
	Payload T
}

type versionOnly struct {
	Version int
}

var ErrVersionDoesNotMatch = errors.New("version does not match")

func NewVersionedJSON[T any](expectedVersion int, contents []byte) (*VersionedJSON[T], error) {
	var vOnly versionOnly

	if err := json.Unmarshal(contents, &vOnly); err != nil {
		return nil, err
	}

	if vOnly.Version != expectedVersion {
		return nil, ErrVersionDoesNotMatch
	}

	var v = new(VersionedJSON[T])
	if err := json.Unmarshal(contents, v); err != nil {
		return nil, err
	}

	return v, nil
}

func (v *VersionedJSON[T]) GetVersion() int {
	return v.Version
}

func (v *VersionedJSON[T]) GetPayload() T {
	return v.Payload
}

func (v *VersionedJSON[T]) SetPayload(version int, data T) {
	v.Version = version
	v.Payload = data
}

func (v *VersionedJSON[T]) ToBytes() ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func GenerateVersionedJSON[T any](version int, data T) ([]byte, error) {
	v := VersionedJSON[T]{
		Version: version,
		Payload: data,
	}

	return v.ToBytes()
}
