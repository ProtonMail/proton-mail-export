// Copyright (c) 2024 Proton AG
//
// This file is part of Proton Mail Bridge.
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
// along with Proton Mail Bridge. If not, see <https://www.gnu.org/licenses/>.

#include "operation.h"

std::string backupStr = "backup";
std::string restoreStr = "restore";

//****************************************************************************************************************************************************
/// \param[in] operationStr The string representing the operation.
/// \return The operation matching the string.
//****************************************************************************************************************************************************
EOperation stringToOperation(std::string_view const operationStr) {
    if (operationStr == backupStr) {
        return EOperation::Backup;
    }

    if (operationStr == restoreStr) {
        return EOperation::Restore;
    }

    return EOperation::Unknown;
}
