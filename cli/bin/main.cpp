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

#include <atomic>
#include <filesystem>
#include <iostream>
#include <string>

#if !defined(_WIN32)
#include <csignal>
#else
#include <windows.h>
#endif

#include <etconfig.hpp>
#include <etsession.hpp>

std::string readText(std::string_view label) {
    std::string result;
    std::cout << label << ": " << std::flush;
    std::getline(std::cin, result);

    return result;
}

std::string readPath(std::string_view label) {
    std::string result;

    while (true) {
        result.clear();
        std::cout << label << ": " << std::flush;
        std::getline(std::cin, result);
        // Ensure the path is converted from utf8 to native type.
        auto utf8Path = std::filesystem::u8path(result);

        try {
            if (std::filesystem::exists(utf8Path) && !std::filesystem::is_directory(utf8Path)) {
                std::cerr << "Path is not a directory" << std::endl;
                continue;
            }
        } catch (const std::exception& e) {
            std::cerr << "Failed to check export path:" << e.what() << std::endl;
            continue;
        }

        // Ensure we get utf8 string back;
        return utf8Path.u8string();
    }
}

std::string readSecret(std::string_view label) {
    std::string result;
    std::cout << label << ": " << std::flush;
    std::getline(std::cin, result);

    return result;
}

static std::atomic_bool gShouldQuit = std::atomic_bool(false);

class StdOutExportMailCallback final : public etcpp::ExportMailCallback {
   public:
    etcpp::ExportMailCallback::Reply onProgress(float progress) override {
        printf("Export Mail Progress: %.02f", double(progress));
        std::cout << std::endl;

        if (gShouldQuit) {
            return etcpp::ExportMailCallback::Reply::Cancel;
        }

        return etcpp::ExportMailCallback::Reply::Continue;
    }
};

void onSignalCancel() {
    std::cout << std::endl << "Received Ctrl+C, exiting as soon as possible" << std::endl;
    gShouldQuit.store(true);
}

int main() {
#if !defined(_WIN32)
    signal(SIGINT, [](int) { onSignalCancel(); });
#else
    SetConsoleCtrlHandler(
        [](DWORD ctrlType) {
            switch (ctrlType) {
                case CTRL_C_EVENT:
                    onSignalCancel();
                    break;
                default:
                    break;
            }
        },
        1);
#endif
    auto session = etcpp::Session(et::DEFAULT_API_URL);

    etcpp::Session::LoginState loginState = etcpp::Session::LoginState::LoggedOut;

    while (loginState != etcpp::Session::LoginState::LoggedIn) {
        if (gShouldQuit) {
            return EXIT_SUCCESS;
        }

        switch (loginState) {
            case etcpp::Session::LoginState::LoggedOut: {
                auto username = readText("Username");
                if (gShouldQuit) {
                    return EXIT_SUCCESS;
                }

                auto password = readSecret("Password");
                if (gShouldQuit) {
                    return EXIT_SUCCESS;
                }

                try {
                    loginState = session.login(username.c_str(), password);
                } catch (const etcpp::SessionException& e) {
                    std::cerr << "Failed to login" << e.what() << std::endl;
                    return EXIT_FAILURE;
                }
                break;
            }
            case etcpp::Session::LoginState::AwaitingTOTP: {
                auto totp = readSecret("TOTP Code");
                if (gShouldQuit) {
                    return EXIT_SUCCESS;
                }
                try {
                    loginState = session.loginTOTP(totp.c_str());
                } catch (const etcpp::SessionException& e) {
                    std::cerr << "Failed to submit totp code:" << e.what() << std::endl;
                    return EXIT_FAILURE;
                }
                break;
            }
            case etcpp::Session::LoginState::AwaitingHV: {
                std::cerr << "Not yet implemented" << std::endl;
                return EXIT_FAILURE;
            }
            case etcpp::Session::LoginState::AwaitingMailboxPassword: {
                auto mboxPassword = readSecret("Mailbox Password");
                if (gShouldQuit) {
                    return EXIT_SUCCESS;
                }

                try {
                    loginState = session.loginMailboxPassword(mboxPassword);
                } catch (const etcpp::SessionException& e) {
                    std::cerr << "Failed to set mailbox password" << e.what() << std::endl;
                    return EXIT_FAILURE;
                }
                break;
            }
            default:
                std::cerr << "Unknown login state" << std::endl;
                return EXIT_FAILURE;
        }

        auto exportPath = readPath("Export Path");

        auto exportMail = session.newExportMail(exportPath.c_str());
        auto cb = StdOutExportMailCallback{};

        try {
            exportMail.start(cb);
        } catch (const etcpp::ExportMailException& e) {
            std::cerr << "Failed to export: " << e.what() << std::endl;
            return EXIT_FAILURE;
        }
    }

    return EXIT_SUCCESS;
}