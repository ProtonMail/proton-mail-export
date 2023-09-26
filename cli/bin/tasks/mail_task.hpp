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

#include <etexport_mail.hpp>
#include <filesystem>

#include "tasks/task.hpp"
#include "tui_util.hpp"

class MailTask final : public TaskWithProgress<void>, etcpp::ExportMailCallback {
   private:
    etcpp::ExportMail mExport;
    CLIProgressBar mProgressBar;

   public:
    MailTask(etcpp::Session& session, const std::filesystem::path& exportPath);
    ~MailTask() override = default;
    MailTask(const MailTask&) = delete;
    MailTask(MailTask&&) = delete;
    MailTask& operator=(const MailTask&) = delete;
    MailTask& operator=(MailTask&&) = delete;

    void run() override;

    void cancel() override;

    std::string_view description() const override;

    inline std::filesystem::path getExportPath() const { return mExport.getExportPath(); }

   private:
    void onProgress(float progress) override;
};