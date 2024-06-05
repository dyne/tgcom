<div align="center">

# tgcom

### toggle comments in source files

</div>

<p align="center">
  <a href="https://dyne.org">
    <img src="https://files.dyne.org/software_by_dyne.png" width="170">
  </a>
</p>


---
<br><br>

`tgcom` is a command-line tool designed to comment, uncomment, and toggle comments in source code files. It supports multiple languages including JavaScript, Go, and Bash, and can be extended to support more. The tool can handle single lines, ranges of lines, and a mix of both. It also supports handling streams from stdin and processes multiple files and ranges in one command.

## tgcom Features


- **Comment/Uncomment/Toggle Comments**: Operate on single lines, ranges, or a mixture of both.
- **Multi-language Support**: Supports JavaScript, Go, Bash, and can be extended to other languages.
- **File Handling**: Works with filenames or streams from stdin.
- **Backup Creation**: Automatically creates a backup before modifying a file.
- **Performance**: Fast and efficient, does not load the entire file into memory.
- **Labels for Sections**: Supports labels for commenting sections in the style of heredocs.


<br>

<div id="toc">

### 🚩 Table of Contents

- [💾 Install](#-install)
- [🎮 Quick start](#-quick-start)
- [📜 Requirements](#-requirements)
- [🚑 Community & support](#-community--support)
- [🐋 Docker](#-docker)
- [😍 Acknowledgements](#-acknowledgements)
- [👤 Contributing](#-contributing)
- [💼 License](#-license)

</div>

***
## 💾 Install

```bash
go get github.com/dyne/tgcom
```


**[🔝 back to top](#toc)**

***
## 🎮 Quick start

### Basic Command

```sh
tgcom --file <filename> --line <line_number> --action <comment|uncomment|toggle>
```

### Examples
Comment a Single Line
```sh
tgcom --file main.go --line 10 --action comment
```

Uncomment a Range of Lines
```sh
tgcom --file main.go --lines 10-20 --action uncomment
```

Toggle Comments on Multiple Files and Lines
```sh
tgcom --files main.go:10-20,script.sh:4,index.html:#<p>,#</p> --action toggle
```

Using Stdin
```sh
cat main.go | tgcom --line 10 --action comment
```

Using Labels for Sections
```sh
tgcom --file main.go --start-label START --end-label END --action comment
```


**[🔝 back to top](#toc)**

***
## 📜 Requirements

1. **Language Support**:
    - The software must support at least JavaScript, Go, and Bash for commenting/uncommenting/toggling lines.
    - It should be extensible to support additional programming languages.

2. **File Handling**:
    - Accept filenames as input and work on streams from stdin.
    - Replace the files in place after making changes.
    - Create a backup of the original file before making any changes.
    - Provide a dry-run option to print the changes to stdout instead of modifying the files in place.

3. **Commenting Functionality**:
    - Comment, uncomment, and toggle comments for:
        - Single lines
        - Ranges of lines
        - A mixture of single lines and ranges
    - Accept labels for commenting sections, such as heredocs with a start keyword and an end keyword in source files.

4. **Performance**:
    - The tool must be fast and efficient.
    - Avoid loading the entire file into memory.

5. **User Interface**:
    - Provide a command-line interface (CLI) for user interaction.
    - Allow the CLI to handle multiple files and complex input specifications, such as `main.js:10-20 script.sh:4 index.html:#<p>,#</p>`.

6. **Testing**:
    - Include test units to ensure reliability and correctness of the software.
    - Tests should cover various scenarios and edge cases.

**[🔝 back to top](#toc)**

***
## 🚑 Community & support

**[📝 Documentation](#toc)** - Getting started and more.

**[🌱 Ecosystem](https://github.com/dyne/ecosystem)** - Plugins, resources, and more.

**[🚩 Issues](../../issues)** - Bugs end errors you encounter using tgcom.

**[💬 Discussions](../../discussions)** - Get help, ask questions, request features, and discuss tgcom.

**[[] Matrix](https://socials.dyne.org/matrix)** - Hanging out with the community.

**[🗣️ Discord](https://socials.dyne.org/discord)** - Hanging out with the community.

**[🪁 Telegram](https://socials.dyne.org/telegram)** - Hanging out with the community.

**[📖 Example](https://github.com/tgcom/example)** - An example repository that uses tgcom.

**[🔝 back to top](#toc)**

***
## 🐋 Docker

Please refer to [DOCKER PACKAGES](../../packages)


**[🔝 back to top](#toc)**

***
## 😍 Acknowledgements

<a href="https://dyne.org">
  <img src="https://files.dyne.org/software_by_dyne.png" width="222">
</a>


Copyleft 🄯 2023 by [Dyne.org](https://www.dyne.org) foundation, Amsterdam

Designed, written and maintained by Puria Nafisi Azizi.


**[🔝 back to top](#toc)**

***
## 👤 Contributing

Please first take a look at the [Dyne.org - Contributor License Agreement](CONTRIBUTING.md) then

1.  🔀 [FORK IT](../../fork)
2.  Create your feature branch `git checkout -b feature/branch`
3.  Commit your changes `git commit -am 'feat: New feature\ncloses #398'`
4.  Push to the branch `git push origin feature/branch`
5.  Create a new Pull Request `gh pr create -f`
6.  🙏 Thank you


**[🔝 back to top](#toc)**

***
## 💼 License
    tgcom - toggle comments in source files
    Copyleft 🄯 2023 Dyne.org foundation, Amsterdam

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.

**[🔝 back to top](#toc)**
