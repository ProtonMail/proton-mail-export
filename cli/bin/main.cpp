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

#include <string>

#include <ftxui/component/component.hpp>
#include <ftxui/component/component_options.hpp>
#include <ftxui/component/screen_interactive.hpp>
#include <ftxui/dom/elements.hpp>

#include <etconfig.hpp>
#include <etcpp.hpp>

int main() {
    auto session = etcpp::Session(et::DEFAULT_API_URL);

    std::string email;
    std::string password;
    auto inputEmail = ftxui::Input(&email, "Email");

    auto screen = ftxui::ScreenInteractive::TerminalOutput();
    auto exit = screen.ExitLoopClosure();

    auto passwordOptions = ftxui::InputOption{};
    passwordOptions.password = true;
    auto inputPassword = ftxui::Input(&password, "password", passwordOptions);
    auto loginButton = ftxui::Button("Login", [&]() {
        try {
            auto loginState = session.login(email.c_str(), password.c_str());
            if (loginState != etcpp::Session::LoginState::LoggedIn) {
                std::cerr << "Login requires mores steps" << std::endl;
            }
        } catch (const etcpp::Exception& e) {
            std::cerr << "Error: " << e.what() << std::endl;
        }
        exit();
    });

    auto components = ftxui::Container::Vertical({
        inputEmail,
        inputPassword,
        loginButton,
    });

    auto renderer = ftxui::Renderer(components, [&] {
        return ftxui::vbox({
                   ftxui::text("Login into your proton account"),
                   ftxui::separator(),
                   ftxui::hbox(ftxui::text("Email   : "), inputEmail->Render()),
                   ftxui::hbox(ftxui::text("Password: "), inputPassword->Render()),
                   ftxui::separator(),
                   loginButton->Render(),
               }) |
               ftxui::border;
    });

    screen.Loop(renderer);
    return EXIT_SUCCESS;
}