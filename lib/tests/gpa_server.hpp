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

#pragma once

#include <exception>
#include <string>
#include <string_view>
#include <vector>

#include "etgpa.h"

class GPAException final : public std::exception {
   private:
    friend class Session;
    std::string mWhat;

   public:
    GPAException(std::string_view what);
    const char* what() const noexcept;
};

class GPAServer {
   private:
    gpaServer* mServer;

   public:
    GPAServer();
    ~GPAServer();

    std::string createUser(const char* email, const char* password, std::string& outAddrID);

    std::string url() const;

    std::vector<std::string> createTestMessages(const char* userID,
                                                const char* addrRD,
                                                const char* email,
                                                const char* password,
                                                int count);

    GPAServer(const GPAServer&) = delete;
    GPAServer(GPAServer&&) = delete;
    GPAServer& operator=(const GPAServer&) = delete;
    GPAServer& operator=(GPAServer&&) = delete;
};