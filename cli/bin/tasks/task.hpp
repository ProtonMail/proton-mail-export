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

#include <condition_variable>
#include <future>
#include <mutex>
#include <thread>

/// CLI background task helper.
template <class P>
class Task {
    static_assert(std::is_default_constructible_v<P> && std::is_move_assignable_v<P>);

   protected:
    std::mutex mMutex;
    std::condition_variable mCond;
    std::once_flag mCancelCall;
    P mValue;

   public:
    Task() = default;
    virtual ~Task() = default;
    Task(const Task&) = delete;
    Task& operator=(const Task&) = delete;
    Task(Task&&) = delete;
    Task& operator=(const Task&&) = delete;

    /// Start a task. If the task should cancel, simply set shouldQuit to true.
    virtual void start(std::atomic_bool& shouldQuit) = 0;

   protected:
    enum class Result { Continue, Cancel };

    /// Run a task in the background and periodically run F on the main thread so that we can
    /// produce an update while the task is running. If the task needs to be cancelled, call to F
    /// should return Result::Cancel and the task cancellation code will be cancelled once. Note: It
    /// is safe to return cancel more than once, this helper class ensures cancel is only called
    /// once.
    template <class F>
    void runBackground(F&& f) {
        static_assert(std::is_invocable_r_v<Result, F, const P&>);
        auto future = std::async(std::launch::async, [&] { startTask(); });
        do {
            std::unique_lock lockScope(mMutex);
            if (mCond.wait_for(lockScope, std::chrono::milliseconds(500)) ==
                std::cv_status::no_timeout) {
                if (f(mValue) == Result::Cancel) {
                    std::call_once(mCancelCall, [&] { cancelTask(); });
                }
            }
        } while (future.wait_for(std::chrono::seconds(0)) != std::future_status::ready);

        return future.get();
    }

    virtual void startTask() = 0;

    virtual void cancelTask() = 0;

    void sendUpdate(P p) {
        std::unique_lock lockScope(mMutex);
        mValue = std::move(p);
        mCond.notify_one();
    }

    void sendUpdate(P&& p) {
        std::unique_lock lockScope(mMutex);
        mValue = std::move(p);
        mCond.notify_one();
    }
};
