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


#ifndef ET_OPERATION_H
#define ET_OPERATION_H

#include <string>

extern std::string backupStr;
extern std::string restoreStr;

//****************************************************************************************************************************************************
/// \brief Enumeration for the operation to perform.
//****************************************************************************************************************************************************
enum class EOperation {
    Backup = 0,
    Restore = 1,
    Unknown = 2,
};

EOperation stringToOperation(std::string_view operationString); ///< Converts a string to an operation.

#endif // ET_OPERATION_H
