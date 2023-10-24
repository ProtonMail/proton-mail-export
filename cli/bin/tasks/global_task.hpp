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

#include <et.hpp>
#include <string>
#include <type_traits>

#include "tasks/task.hpp"

template <class R>
class GlobalTask : public Task<R> {
   protected:
    std::string mDesc;
    etcpp::GlobalScope& mScope;

   protected:
    GlobalTask(etcpp::GlobalScope& scope, std::string_view desc) : mDesc(desc), mScope(scope) {}

   public:
    virtual ~GlobalTask() override = default;

    void cancel() override {}

    std::string_view description() const override { return mDesc; }
};

class NewVersionCheckTask final : public GlobalTask<bool> {
   public:
    NewVersionCheckTask(etcpp::GlobalScope& scope, std::string_view desc)
        : GlobalTask<bool>(scope, desc) {}

    ~NewVersionCheckTask() override = default;

    bool run() override;
};
