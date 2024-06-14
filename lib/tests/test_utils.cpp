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

#include "test_utils.h"
#include <random>

//****************************************************************************************************************************************************
/// \brief Returns a random string composed of lowercase roman alphabet letters.
///
/// \param[in] length The length of the string.
/// \return A string of length length composed of characters in the range.
//****************************************************************************************************************************************************
std::string randomLowerCaseLetters(size_t const length) {
    static std::random_device rd;
    static std::mt19937 generator(rd());
    static std::uniform_int_distribution<> distribution('a', 'z');

    std::string result;
    result.reserve(length);
    for (size_t i = 0; i < length; i++) {
        result += static_cast<char>(distribution(generator));
    }

    return result;
}

//****************************************************************************************************************************************************
/// \brief Creates a temporary folder.
/// \return The path of the temporary folder.
//****************************************************************************************************************************************************
std::filesystem::path createTempDir() {
    std::filesystem::path const tempDir = std::filesystem::temp_directory_path();
    std::filesystem::path path;
    while (true) {
        path = tempDir / randomLowerCaseLetters(16);
        if (!std::filesystem::exists(path)) {
            break;
        }
    }

    if (!std::filesystem::create_directories(path)) {
        throw std::runtime_error("Could not create temporary folder");
    }

    return path;
}

//****************************************************************************************************************************************************
//
//****************************************************************************************************************************************************
ScopedTempFolder::ScopedTempFolder() : path_(createTempDir()) {}

//****************************************************************************************************************************************************
//
//****************************************************************************************************************************************************
ScopedTempFolder::~ScopedTempFolder() {
    std::filesystem::remove_all(path_);
}

//****************************************************************************************************************************************************
/// \return the path of the temporary folder.
//****************************************************************************************************************************************************
std::filesystem::path ScopedTempFolder::getPath() const {
    return path_;
}