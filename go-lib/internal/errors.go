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

package internal

import (
	"errors"
	"fmt"

	"github.com/ProtonMail/go-proton-api"
)

func MapError(err error) error {
	if err == nil {
		return nil
	}

	var protonErr *proton.APIError
	if errors.As(err, &protonErr) {
		return fmt.Errorf(protonErr.Message)
	}

	return err
}
