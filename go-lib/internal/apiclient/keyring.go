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

package apiclient

import (
	"fmt"
	"github.com/ProtonMail/go-proton-api"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/sirupsen/logrus"
)

type UnlockedKeyRing struct {
	keyRing *crypto.KeyRing
	addrMap map[string]*crypto.KeyRing
	user    *proton.User
}

func NewUnlockedKeyRing(user *proton.User, addresses []proton.Address, keyPass []byte) (*UnlockedKeyRing, error) {
	userKR, err := user.Keys.Unlock(keyPass, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to unlock user keys: %w", err)
	}

	keyring := &UnlockedKeyRing{
		keyRing: userKR,
		addrMap: make(map[string]*crypto.KeyRing),
		user:    user,
	}

	for _, addr := range addresses {
		addrKR, err := addr.Keys.Unlock(keyPass, userKR)
		if err != nil {
			logrus.WithField("addressID", addr.ID).WithError(err).Warn("Failed to unlock address keys")
			continue
		}

		if addrKR.CountDecryptionEntities() == 0 {
			addrKR.ClearPrivateParams()
			logrus.WithField("addressID", addr.ID).Warn("Address keyring has no decryption entities")
			continue
		}

		keyring.addrMap[addr.ID] = addrKR
	}

	return keyring, nil
}

func (u *UnlockedKeyRing) Close() {
	for _, v := range u.addrMap {
		v.ClearPrivateParams()
	}
	u.addrMap = nil
	u.keyRing.ClearPrivateParams()
}

func (u *UnlockedKeyRing) GetAddrKeyRing(addrID string) (*crypto.KeyRing, bool) {
	kr, ok := u.addrMap[addrID]

	return kr, ok
}

func (u *UnlockedKeyRing) GetAddrKeyRingMap() map[string]*crypto.KeyRing {
	return u.addrMap
}
