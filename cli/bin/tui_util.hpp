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

#include <functional>
#include <string>

void setStdinEcho(bool enable = true);

bool isStdoutTerminal();

bool isStdInTerminal();

/// Install a signal handler.
/// ThreadSafety: This may get called from any thread at any time.
bool registerCtrlCSignalHandler(std::function<void()>&& handler);

class CliSpinner {
private:
    int mState = 0;

public:
    CliSpinner() = default;
    char next();
};

class CLIProgressBar {
private:
    int mActiveBars = -1;
    std::string mValue;

public:
    CLIProgressBar();
    void update(float progress);

    inline std::string_view value() const { return mValue; }
};
