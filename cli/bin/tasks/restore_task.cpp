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

#include "tasks/restore_task.hpp"
#include <etsession.hpp>
#include <iostream>

RestoreTask::RestoreTask(etcpp::Session& session, const std::filesystem::path& exportPath) :
    mRestore(session.newExportRestore(exportPath.u8string().c_str())) {}

void RestoreTask::onProgress(float progress) {
    updateProgress(progress);
}

void RestoreTask::run() {
    mRestore.start(*this);
}

void RestoreTask::cancel() {
    mRestore.cancel();
}

std::string_view RestoreTask::description() const {
    return "Restore Mail";
}
