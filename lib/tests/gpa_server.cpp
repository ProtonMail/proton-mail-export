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

#include "gpa_server.hpp"

GPAServer::GPAServer() : mServer(gpaServerNew()) {}

GPAServer::~GPAServer() {
    gpaServerDelete(mServer);
}

std::string GPAServer::createUser(const char* email, const char* password) {
    char* userID = nullptr;

    if (gpaServerCreateUser(mServer, email, password, &userID) != GPA_SERVER_STATUS_OK) {
        throw GPAException("Failed to create user");
    }

    auto result = std::string(userID);
    free(userID);

    return result;
}

std::string GPAServer::url() const {
    char* outURL = nullptr;
    if (gpaServerGetURL(mServer, &outURL) != GPA_SERVER_STATUS_OK) {
        throw GPAException("Failed to get server url");
    }

    auto result = std::string(outURL);
    free(outURL);

    return result;
}

GPAException::GPAException(std::string_view what) : mWhat(what) {}

const char* GPAException::what() const noexcept {
    return mWhat.c_str();
}