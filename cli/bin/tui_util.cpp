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

#include "tui_util.hpp"

#include <fmt/format.h>
#include <atomic>
#include <cmath>
#include <cstdio>

#if !defined(_WIN32)
#include <termios.h>
#include <unistd.h>
#include <csignal>
#else
#include <io.h>
#include <windows.h>
#endif

void setStdinEcho(bool enable) {
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

static inline bool isTerm(FILE* f) {
#if defined(_WIN32)
    return _isatty(_fileno(f)) == 1;
#else
    return isatty(fileno(f)) == 1;
#endif
}

bool isStdoutTerminal() {
    return isTerm(stdout);
}

bool isStdInTerminal() {
    return isTerm(stdin);
}

static std::function<void()> gSignalHandler = []() {};

bool registerCtrlCSignalHandler(std::function<void()>&& handler) {
#if !defined(_WIN32)
    if (signal(SIGINT, [](int) { gSignalHandler(); }) == SIG_ERR) {
        return false;
    }
#else
    if (SetConsoleCtrlHandler(
            [](DWORD ctrlType) -> BOOL {
                switch (ctrlType) {
                    case CTRL_C_EVENT:
                        gSignalHandler();
                        return TRUE;
                    default:
                        return FALSE;
                }
            },
            TRUE) == FALSE) {
        return false;
    }
#endif
    gSignalHandler = handler;
    return true;
}

char CliSpinner::next() {
    constexpr const int kMaxSpinStates = 4;
    constexpr const char kSpinStates[kMaxSpinStates] = {
        '-',
        '\\',
        '|',
        '/',
    };

    auto curIdx = mState;
    mState = (mState + 1) % kMaxSpinStates;
    return kSpinStates[curIdx];
}

CLIProgressBar::CLIProgressBar() {
    update(0);
}

void CLIProgressBar::update(float progress) {
    constexpr int kNumBars = 50;
    mValue.reserve(kNumBars);

    const int activeBars = std::ceil((progress * float(kNumBars)) / 100.0f);
    if (mActiveBars != activeBars) {
        mValue.clear();
        if (progress < 10.0f) {
            fmt::format_to(std::back_inserter(mValue), "[0{:.2f}%]", double(progress));
        } else {
            fmt::format_to(std::back_inserter(mValue), "[{:.2f}%]", double(progress));
        }
        mValue.push_back('[');
        for (int i = 0; i < kNumBars; i++) {
            if (i < activeBars) {
                mValue.push_back('|');
            } else {
                mValue.push_back(' ');
            }
        }
        mValue.push_back(']');
    }
}
