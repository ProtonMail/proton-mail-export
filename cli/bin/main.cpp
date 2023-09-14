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
#include <optional>
#include <string>

#include <fmt/format.h>
#include <cxxopts.hpp>

#if !defined(_WIN32)
#include <termios.h>
#include <unistd.h>
#include <csignal>
#else
#include <windows.h>
#endif

#include <etconfig.hpp>
#include <etlog.hpp>
#include <etsession.hpp>
#include <etutil.hpp>

void setStdinEcho(bool enable = true) {
#ifdef WIN32
    HANDLE hStdin = GetStdHandle(STD_INPUT_HANDLE);
    DWORD mode;
    GetConsoleMode(hStdin, &mode);

    if (!enable) {
        mode &= ~ENABLE_ECHO_INPUT;
    } else {
        mode |= ENABLE_ECHO_INPUT;
    }

    SetConsoleMode(hStdin, mode);
#else
    struct termios tty;
    tcgetattr(STDIN_FILENO, &tty);
    if (!enable)
        tty.c_lflag &= ~ECHO;
    else
        tty.c_lflag |= ECHO;

    (void)tcsetattr(STDIN_FILENO, TCSANOW, &tty);
#endif
}

constexpr int kNumInputRetries = 3;

std::string readText(std::string_view label) {
    for (int i = 0; i < kNumInputRetries; i++) {
        std::string result;
        std::cout << label << ": " << std::flush;
        std::getline(std::cin, result);

        if (result.empty()) {
            std::cerr << "Value can't be empty" << std::endl;
            continue;
        }

        return result;
    }

    throw std::runtime_error(fmt::format("Failed read value for '{}'", label));
}

std::filesystem::path readPath(std::string_view label) {
    std::string result;

    for (int i = 0; i < kNumInputRetries; i++) {
        result.clear();
        std::cout << label << ": " << std::flush;
        std::getline(std::cin, result);

        if (result.empty()) {
            std::cerr << "Value can't be empty" << std::endl;
            continue;
        }

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

        return utf8Path;
    }

    throw std::runtime_error(fmt::format("Failed read value for '{}'", label));
}

std::string readSecret(std::string_view label) {
    struct PasswordScope {
        PasswordScope() { setStdinEcho(false); }
        ~PasswordScope() { setStdinEcho(true); }
    };

    PasswordScope pscope;

    for (int i = 0; i < kNumInputRetries; i++) {
        std::string result;
        std::cout << label << ": " << std::flush;
        std::getline(std::cin, result);

        if (result.empty()) {
            std::cerr << "Value can't be empty" << std::endl;
            continue;
        }

        return result;
    }

    throw std::runtime_error(fmt::format("Failed read value for '{}'", label));
}

template <class F>
std::string getCLIValue(cxxopts::ParseResult& parseResult,
                        const char* argKey,
                        std::optional<const char*> envVariable,
                        F fallback) {
    static_assert(std::is_invocable_r_v<std::string, F>);

    std::string result;
    if (parseResult.count(argKey)) {
        result = parseResult[argKey].as<std::string>();
    }
    if (!result.empty()) {
        return result;
    }

    if (envVariable) {
        auto envVar = std::getenv(*envVariable);
        if (envVar != nullptr && std::strlen(envVar) != 0) {
            return envVar;
        }
    }

    return fallback();
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

int main(int argc, const char** argv) {
    std::filesystem::path execPath;
    try {
        execPath = etcpp::getExecutableDir();
    } catch (const std::exception& e) {
        std::cerr << "Failed to get executable directory: " << e.what() << std::endl;
        std::cerr << "Will user working directory instead" << std::endl;
    }

    try {
        execPath.append("logs");
        auto logScope = etcpp::LogScope(execPath);

        const char* helpText = "Proton Data Exporter v{}";

        cxxopts::Options options("proton-export", fmt::format(helpText, et::VERSION_STR));

        options.add_options()("e,export-dir", "Export directory", cxxopts::value<std::string>())(
            "p,password", "User's password (can also be set with env var ET_USER_PASSWORD)",
            cxxopts::value<std::string>())(
            "m,mbox-password",
            "User's mailbox password when using 2 Password Mode (can also be set with env var "
            "ET_USER_MAILBOX_PASSWORD)",
            cxxopts::value<std::string>())(
            "t,totp", "User's TOTP 2FA code (can also be set with env var ET_TOTP_CODE)",
            cxxopts::value<std::string>())(
            "u,user", "User's account/email (can also be set with env var ET_USER_EMAIL",
            cxxopts::value<std::string>())("h,help", "Show help");

        auto argParseResult = options.parse(argc, argv);

        if (argParseResult.count("help")) {
            std::cout << options.help() << std::endl;
            return EXIT_SUCCESS;
        }

#if !defined(_WIN32)
        signal(SIGINT, [](int) { onSignalCancel(); });
#else
        if (SetConsoleCtrlHandler(
                [](DWORD ctrlType) -> BOOL {
                    switch (ctrlType) {
                        case CTRL_C_EVENT:
                            onSignalCancel();
                            return TRUE;
                        default:
                            return FALSE;
                    }
                },
                TRUE) == FALSE) {
            std::cerr << "Failed to instal console control hanlder" << std::endl;
            return EXIT_FAILURE;
        }
#endif
        auto session = etcpp::Session(et::DEFAULT_API_URL);

        etcpp::Session::LoginState loginState = etcpp::Session::LoginState::LoggedOut;

        while (loginState != etcpp::Session::LoginState::LoggedIn) {
            if (gShouldQuit) {
                return EXIT_SUCCESS;
            }

            switch (loginState) {
                case etcpp::Session::LoginState::LoggedOut: {
                    const auto username = getCLIValue(argParseResult, "user", "ET_USER_EMAIL",
                                                      []() { return readText("Username"); });
                    if (gShouldQuit) {
                        return EXIT_SUCCESS;
                    }

                    const auto password =
                        getCLIValue(argParseResult, "password", "ET_USER_PASSWORD",
                                    []() { return readSecret("Password"); });

                    try {
                        loginState = session.login(username.c_str(), password);
                    } catch (const etcpp::SessionException& e) {
                        std::cerr << "Failed to login" << e.what() << std::endl;
                        return EXIT_FAILURE;
                    }
                    break;
                }
                case etcpp::Session::LoginState::AwaitingTOTP: {
                    const auto totp = getCLIValue(argParseResult, "totp", "ET_TOTP_CODE",
                                                  []() { return readSecret("TOTP Code"); });
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
                    const auto mboxPassword =
                        getCLIValue(argParseResult, "mbox-password", "ET_USER_MAILBOX_PASSWORD",
                                    []() { return readSecret("Mailbox Password"); });
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

            std::filesystem::path exportPath;
            if (argParseResult.count("export-dir")) {
                exportPath =
                    std::filesystem::u8path(argParseResult["export-dir"].as<std::string>());
            }
            if (exportPath.empty()) {
                exportPath = readPath("Export Path");
            }

            auto exportMail = session.newExportMail(exportPath.u8string().c_str());
            auto cb = StdOutExportMailCallback{};

            try {
                exportMail.start(cb);
            } catch (const etcpp::ExportMailException& e) {
                std::cerr << "Failed to export: " << e.what() << std::endl;
                return EXIT_FAILURE;
            }
        }

    } catch (const std::exception& e) {
        std::cerr << "Encountered unexpected error: " << e.what() << std::endl;
        return EXIT_FAILURE;
    }

    return EXIT_SUCCESS;
}
