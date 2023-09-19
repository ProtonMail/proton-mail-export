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

#include <iostream>
#include <sstream>
#include <string_view>

#include "tasks/task.hpp"
#include "tui_util.hpp"

constexpr const std::string_view kNetworkLostText = "Can't connect to proton servers. Retrying...";

void fillSpaces(size_t previousLength, size_t currentLength) {
    for (size_t i = previousLength; i < currentLength; i++) {
        std::cout << ' ';
    }
}

template <class R>
R runTask(const TaskAppState& state, Task<R>& task) {
    auto desc = task.description();
    auto future = std::async(std::launch::async, [&]() -> R { return task.run(); });
    auto spinner = CliSpinner();

    do {
        if (state.shouldQuit()) {
            task.cancel();
        } else {
            if (state.networkLost()) {
                std::cout << '\r' << spinner.next() << " " << kNetworkLostText;
                fillSpaces(kNetworkLostText.length(), desc.length());
                std::cout << std::flush;
            } else {
                std::cout << '\r' << spinner.next() << " " << desc;
                fillSpaces(desc.length(), kNetworkLostText.length());
                std::cout << std::flush;
            }
        }
    } while (future.wait_for(std::chrono::milliseconds(500)) != std::future_status::ready);
    std::cout << "\n" << std::flush;

    if constexpr (std::is_void_v<R>) {
        future.get();
    } else {
        return future.get();
    }
}

template <class R>
R runTaskWithProgress(const TaskAppState& state, TaskWithProgress<R>& task) {
    auto future = std::async(std::launch::async, [&]() -> R { return task.run(); });
    auto spinner = CliSpinner();
    auto progressBar = CLIProgressBar();
    const auto progressBarLen = progressBar.value().length();
    do {
        if (state.shouldQuit()) {
            task.cancel();
        } else {
            const float progress = task.pollProgress();
            progressBar.update(progress);

            if (state.networkLost()) {
                std::cout << '\r' << spinner.next() << " " << kNetworkLostText << std::flush;
                fillSpaces(kNetworkLostText.length(), progressBarLen);
                std::cout << std::flush;
            } else {
                std::cout << '\r' << progressBar.value();
                fillSpaces(progressBarLen, kNetworkLostText.length());
                std::cout << std::flush;
            }
        }
    } while (future.wait_for(std::chrono::milliseconds(0)) != std::future_status::ready);
    std::cout << "\n" << std::flush;

    if constexpr (std::is_void_v<R>) {
        future.get();
    } else {
        return future.get();
    }
}
