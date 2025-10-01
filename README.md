# Getctx

`getctx` is a powerful and interactive command-line (CLI) tool, written in Go, designed to streamline the process of gathering and consolidating source code and text files into a single, cohesive context file.

### What problem does it solve?

Have you ever found yourself manually copying and pasting content from multiple files just to provide context to a large language model (LLM) like GPT or Google AI Studio, create a detailed bug report, or prepare a code snippet for documentation? This process is often tedious, error-prone, and clutters your clipboard. You might forget key files or accidentally include irrelevant ones (like binaries or files from a `.git` directory).

### How can `getctx` help?

`getctx` solves this problem by providing a fast and intuitive Terminal User Interface (TUI) that allows you to navigate your filesystem, interactively select the files and folders you need, and automatically concatenate their contents into a single output file (defaults to `context.txt`). The tool intelligently filters out non-text files and ignores common, unnecessary directories like `node_modules` or `.git`, ensuring the final context is clean and relevant.

By automating the context-gathering process, `getctx` helps you:

- **Improve LLM Prompt Quality:** Provide large, clean, and well-organized code snippets to AI models to get better results.
- **Create Bug Reports Faster:** Instantly package all relevant files needed to reproduce an issue.
- **Streamline Documentation:** Easily aggregate source code for technical documentation or tutorials.

### Key Features

- **Interactive TUI:** Navigate your project using familiar keybindings in a modern and responsive terminal interface.
- **Live Filtering:** Instantly search and filter files and directories in the current view to find what you need quickly.
- **Direct Path Input:** Jump to any directory by typing or pasting its path directly in the terminal, complete with autocompletion support.
- **Smart Selection:** Select single files (`space`) or all visible items at once (`Ctrl+A`), even on a filtered list.
- **Intelligent Exclusion:** The tool automatically ignores irrelevant files and directories (e.g., `.git`, `node_modules`, `dist`, binaries, images) to keep your context clean.
- **Modern Look & Feel:** Built with `Bubble Tea` and `Lipgloss`, `getctx` offers a polished and enjoyable terminal experience.

## Installation

Below are the instructions on how to download and install the program so that it is globally available on your system.

---

### 🐧 Linux

Installation on Linux involves downloading the archive, extracting it, giving the binary file execution permissions, and moving it to the `/usr/local/bin` directory, which is typically in the system's `PATH`.

Installation on Linux involves cloning the repository, building the binary, giving it execution permissions, and moving it to the `/usr/local/bin` directory, which is typically in the system's `PATH`.

1.  **Clone the repository**
    Open a terminal and use the `git` command to clone the repository and navigate into the project directory.

    ```sh
    git clone https://github.com/kacperzielinskidev/getctx.git
    cd getctx
    ```

2.  **Build the program**
    Use the `make build` command to build the binary.

    ```sh
    make build
    ```

3.  **Grant execution permissions**
    After extracting, you will find the binary file. Grant it execution permissions.

    ```sh
    chmod +x bin/getctx
    ```

4.  **Move the file to `/usr/local/bin`**
    Moving the binary file to this directory will make it accessible from anywhere in the system. You will need administrator privileges (`sudo`).

    ```sh
    sudo mv bin/getctx /usr/local/bin/
    ```

5.  **Done!**
    You can now run the program by typing its name from any location in the terminal.

---

### 🪟 Windows

On Windows, the process involves downloading the `.zip` archive, extracting it, and then adding the folder containing the `.exe` file to the system's `PATH` environment variable.

1.  **Clone the repository**
    Open a terminal and use the `git` command to clone the repository and navigate into the project directory.

    ```powershell
    git clone https://github.com/kacperzielinskidev/getctx.git
    cd getctx
    ```

2.  **Build the program**
    You can build the executable in two ways:

    **A) Using `make` (if available)**
    If you have `make` installed on your system, run:

    ```powershell
    make build-windows
    ```

    This will create `getctx.exe` inside the `bin` directory.

    **B) Using the `go` command**
    If you don't have `make`, you can use the standard `go` command:

    ```powershell
    go build -o bin\getctx.exe ./cmd/getctx
    ```

    This will also create `getctx.exe` inside the `bin` directory.

3.  **Create a destination folder and move the file**

    - Create a dedicated, permanent folder for the program, for example `C:\Program Files\getctx`.
    - Move the extracted `getctx.exe` file into this newly created folder.

4.  **Add the folder to the PATH environment variable**
    This allows Windows to find your program from any command line.

    - Press the `Windows` key and type "environment variables".
    - In the Advanced tab, click the Environment Variables... button.
    - In the "System variables" section, find and select the `Path` variable, then click "Edit...".
    - Click "New" and paste the path to the folder you created, i.e., `C:\Program Files\getctx`.
    - Confirm all open windows by clicking "OK".

5.  **Restart your terminal**
    Close any open terminal windows (CMD or PowerShell) and open a new one. Changes to the `PATH` variable require a new terminal session to take effect.

6.  **Done!**
    You can now run the program from any location.
