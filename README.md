# hanabi-bot

This is a framework for creating [Hanabi](https://boardgamegeek.com/boardgame/98778/hanabi) bots that play with each other. It is written in [Go](https://golang.org/). Eventually, the idea is to have the option for players to play together with bots on [Hanabi Live](https://github.com/Zamiell/hanabi-live).

<br />

### Instructions for Coding a New Strategy

Copy the "strategy_dumb.go" file and rename it to "strategy_whatever.go". Fill in all of the functions. Then, add `whateverInit()` to the "strategy.go" file.

<br />

### Installation (for experts)

* The project doesn't require any special dependencies. Just clone the repo and go.

<br />

### Installation (for noobs/contributors)

Like many code projects, we use [golangci-lint](https://github.com/golangci/golangci-lint) to ensure that all of the code is written consistently and error-free. We ask that all pull requests pass our linting rules.

The following instructions will set up the development environment and the linter. This assumes you are on Windows and will be using Microsoft's [Visual Studio Code](https://code.visualstudio.com/), which is a very nice text editor that happens to be better than [Atom](https://atom.io/), [Notepad++](https://notepad-plus-plus.org/), etc. If you are using a different OS/editor, some adjustments will be needed (e.g. using `brew` on MacOS instead of `choco`).

Note that these steps require **an elevated (administrator) command-shell**.

* Install [Chocolatey](https://chocolatey.org/):
  * `@"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile -InputFormat None -ExecutionPolicy Bypass -Command "iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))" && SET "PATH=%PATH%;%ALLUSERSPROFILE%\chocolatey\bin"`
* Install [Git](https://git-scm.com/), [Golang](https://golang.org/), and [Visual Studio Code](https://code.visualstudio.com/):
  * `choco install git golang vscode -y`
* Configure Git:
  * `refreshenv`
  * `git config --global user.name "Your_Username"`
  * `git config --global user.email "your@email.com"`
  * `git config --global core.autocrlf false` <br />
  (so that Git does not convert LF to CRLF when cloning repositories)
  * `git config --global pull.rebase true` <br />
  (so that Git automatically rebases when pulling)
* Clone the repository:
  * `mkdir %GOPATH%\src\github.com\Zamiell`
  * `cd %GOPATH%\src\github.com\Zamiell`
  * `git clone git@github.com:Zamiell/hanabi-bot.git` <br />
  (or clone a fork, if you are doing development work)
  * `cd hanabi-bot`
* Install the Golang project dependencies:
  * `cd src`
  * `go get -u -v ./...`
  * `cd ..`
* Install the Golang linter:
  * `go get -u -v "github.com/golangci/golangci-lint/cmd/golangci-lint"`
* Install the VSCode extension for Golang:
  * `code --install-extension "ms-vscode.Go"`
* Import a solid set of starting VSCode user settings:
  * `copy "install\settings.json" "%APPDATA%\Code\User\settings.json"` <br />
  (feel free to tweak this file to your liking)
* Open VSCode using the cloned repository as the project folder:
  * `code .`
* Test the Golang linter:
  * On the left pane, navigate to and open "src\action.go".
  * In the bottom-right-hand corner, click on "Analysis Tools Missing" and then on "Install". You will know that it has finished once it displays: "All tools successfully installed."
  * Add a new line of "asdf" somewhere, save the file, and watch as some "Problems" appear in the bottom pane.
