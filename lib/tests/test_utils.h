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

#ifndef TEST_UTILS_H
#define TEST_UTILS_H

#include <filesystem>

//****************************************************************************************************************************************************
/// \brief A class providing an empty temporary folder that will be erased on destruction.
//****************************************************************************************************************************************************
class ScopedTempFolder {
public: // member functions.
    ScopedTempFolder(); ///< Default constructor.
    ScopedTempFolder(ScopedTempFolder const&) = delete; ///< Disabled copy-constructor.
    ScopedTempFolder(ScopedTempFolder&&) = delete; ///< Disabled assignment copy-constructor.
    ~ScopedTempFolder(); ///< Destructor.
    ScopedTempFolder& operator=(ScopedTempFolder const&) = delete; ///< Disabled assignment operator.
    ScopedTempFolder& operator=(ScopedTempFolder&&) = delete; ///< Disabled move assignment operator.

    std::filesystem::path getPath() const; ///< Return the path of the temporary folder.

private:
    std::filesystem::path const path_; ///< The path of the temporary folder.
};


#endif //TEST_UTILS_H