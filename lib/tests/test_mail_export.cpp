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

#include <catch2/catch_test_macros.hpp>
#include <etsession.hpp>
#include <filesystem>
#include <fmt/format.h>
#include <iostream>

#include "gpa_server.hpp"
#include "test_data/ExportedData.h"
#include "test_utils.h"

class NullBackupCallback final : public etcpp::BackupCallback {
public:
    void onProgress(float) override {}
};

class ProgressCancelCallback final : public etcpp::BackupCallback {
private:
    etcpp::Backup& mE;
    float mCancelPercentage;
    bool mCancelled = false;

public:
    explicit ProgressCancelCallback(etcpp::Backup& e, float cancelPercentage) : etcpp::BackupCallback(), mE(e), mCancelPercentage(cancelPercentage) {}

    void onProgress(float p) override {
        if (p > mCancelPercentage && !mCancelled) {
            mCancelled = true;
            mE.cancel();
        }
    }
};

TEST_CASE("MailExport") {
    GPAServer server;

    const char* userEmail = "hello";
    const char* userPassword = "12345";

    std::string addrID;
    const auto userID = server.createUser(userEmail, userPassword, addrID);
    const auto url = server.url();

    auto session = etcpp::Session(url.c_str());
    {
        auto loginState = session.getLoginState();
        REQUIRE(loginState == etcpp::Session::LoginState::LoggedOut);
    }

    auto loginState = session.login(userEmail, userPassword);
    REQUIRE(loginState == etcpp::Session::LoginState::LoggedIn);

    std::vector<std::string> messageIDs;
    REQUIRE_NOTHROW(messageIDs = server.createTestMessages(userID.c_str(), addrID.c_str(), userEmail, userPassword, 50));

    time_t t = time(nullptr);

    auto tmpDir = std::filesystem::temp_directory_path();

    // Japanese text below to test unicode path handling on Win32.
    tmpDir /= std::filesystem::u8path("ことわざ") / std::to_string(t);

    std::filesystem::path exportDir{};
    {
        auto backup = session.newBackup(tmpDir.u8string().c_str());
        exportDir = backup.getExportPath();
        auto nullCallback = NullBackupCallback();
        REQUIRE_NOTHROW(backup.start(nullCallback));
    }

    for (const auto& msgID: messageIDs) {
        auto msgPath = exportDir / (msgID + ".eml");
        auto metadataPath = exportDir / (msgID + ".metadata.json");
        REQUIRE(std::filesystem::is_regular_file(msgPath));
        REQUIRE(std::filesystem::is_regular_file(metadataPath));
    }

    REQUIRE_FALSE(std::filesystem::exists(exportDir / "exportProgress.json"));
}


class NullRestoreCallback final : public etcpp::RestoreCallback {
public:
    void onProgress(float) override {}
};

TEST_CASE("MailRestore") {
    GPAServer server;

    const char* userEmail = "hello";
    const char* userPassword = "12345";

    std::string addrID;
    std::string const userID = server.createUser(userEmail, userPassword, addrID);
    std::string const url = server.url();

    auto session = etcpp::Session(url.c_str());
    auto loginState = session.getLoginState();

    REQUIRE(loginState == etcpp::Session::LoginState::LoggedOut);

    loginState = session.login(userEmail, userPassword);
    REQUIRE(loginState == etcpp::Session::LoginState::LoggedIn);

    ScopedTempFolder dir;
    std::cout << dir.getPath();
    createTestBackup(dir.getPath());

    etcpp::Restore restore = session.newRestore(dir.getPath().u8string().c_str());
    NullRestoreCallback nullCallback = NullRestoreCallback();
    REQUIRE_NOTHROW(restore.start(nullCallback));
}