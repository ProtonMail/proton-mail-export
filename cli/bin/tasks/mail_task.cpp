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

#include "tasks/mail_task.hpp"
#include <etsession.hpp>
#include <iostream>

MailTask::MailTask(etcpp::Session& session, const std::filesystem::path& exportPath)
    : mExport(session.newExportMail(exportPath.u8string().c_str())) {}

void MailTask::start(std::atomic_bool& shouldQuit) {
    try {
        runBackground([&](float progress) -> Task::Result {
            if (shouldQuit) {
                return Task::Result::Cancel;
            }

            mProgressBar.update(progress);
            std::cout << '\r' << mProgressBar.value() << std::flush;

            return Task::Result::Continue;
        });
    } catch (const std::exception& e) {
        std::cout << std::endl;
        throw e;
    }

    std::cout << std::endl;
}

void MailTask::startTask() {
    mExport.start(*this);
}

void MailTask::cancelTask() {
    mExport.cancel();
}

void MailTask::onProgress(float progress) {
    sendUpdate(progress);
}
