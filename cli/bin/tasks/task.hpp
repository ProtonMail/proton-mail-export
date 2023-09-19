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

template <class R>
class Task {
   protected:
    Task() = default;

   public:
    virtual ~Task() = default;
    Task(const Task<R>&) = delete;
    Task<R>& operator=(const Task<R>&) = delete;
    Task(Task<R>&&) = delete;
    Task<R>& operator=(const Task<R>&&) = delete;

    virtual R run() = 0;

    virtual void cancel() = 0;

    virtual std::string_view description() const = 0;
};

template <class R>
class TaskWithProgress : public Task<R> {
   private:
    std::mutex mMutex;
    std::condition_variable mCond;
    float mProgress;

   protected:
    TaskWithProgress() = default;

   public:
    virtual ~TaskWithProgress() = default;

    float pollProgress() {
        std::unique_lock lockScope(mMutex);
        mCond.wait_for(lockScope, std::chrono::milliseconds(500));
        return mProgress;
    }

   protected:
    void updateProgress(float progress) {
        std::unique_lock lockScope(mMutex);
        mProgress = progress;
        mCond.notify_one();
    }
};

class TaskAppState {
   public:
    virtual ~TaskAppState() = default;

    virtual bool shouldQuit() const = 0;

    virtual bool networkLost() const = 0;
};