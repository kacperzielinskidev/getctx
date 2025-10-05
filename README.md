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

Choose the installation method that suits you best. Using a package manager like Homebrew or Scoop is recommended for easy installation and automatic updates.

---

### macOS

The recommended way to install `getctx` on macOS is by using the [Homebrew](https://brew.sh/) package manager.

1.  **Add the Tap (one-time setup):**
    First, you need to add the repository containing the installation formula.

    ```sh
    brew tap kacperzielinskidev/tap
    ```

2.  **Install `getctx`:**
    Now you can install the program.

    ```sh
    brew install getctx
    ```

**Updating `getctx` in the future:**

```sh
brew upgrade getctx
```

---

### Linux

You have a few options for installing on Linux.

#### Option A: Debian, Ubuntu, and derivatives (via .deb package)

This method uses your system's native package manager (`dpkg`/`apt`) and is recommended for Debian-based distributions. The script below will automatically download and install the latest `.deb` package.

```sh
# Download the latest .deb package for your architecture (amd64 example)
curl -sL -o /tmp/getctx.deb $(curl -s https://api.github.com/repos/kacperzielinskidev/getctx/releases/latest | grep "browser_download_url.*amd64.deb" | cut -d '"' -f 4)

# Install the package
sudo dpkg -i /tmp/getctx.deb
```

_(Note: For ARM-based systems like Raspberry Pi, replace `amd64.deb` in the script with `arm64.deb`)_

#### Option B: Using Homebrew on Linux

If you already use [Homebrew on Linux (Linuxbrew)](https://docs.brew.sh/Homebrew-on-Linux), you can follow the same instructions as for macOS:

```sh
brew tap kacperzielinskidev/tap
brew install getctx
```

---

### Windows (Scoop)

If you are a Windows user, the easiest way to install and manage `getctx` is with the [Scoop](https://scoop.sh/) package manager.

1.  **Add the Bucket (one-time setup):**
    Open PowerShell and add the repository containing the app manifests.

    ```powershell
    scoop bucket add kacperzielinskidev https://github.com/kacperzielinskidev/scoop-bucket.git
    ```

2.  **Install `getctx`:**
    Now, install the package.

    ```powershell
    scoop install getctx
    ```

**Updating `getctx` in the future:**

```powershell
scoop update getctx
```

---

### Go Install

```sh
go install github.com/kacperzielinskidev/getctx/cmd/getctx@latest
```

The `getctx` binary will be placed in your `$GOPATH/bin` directory.

---

### Manual Installation

You can also download a pre-compiled binary for your operating system directly from the GitHub Releases page.

1.  Go to the [**Latest Release page**](https://github.com/kacperzielinskidev/getctx/releases/latest).
2.  Download the appropriate archive (`.zip` or `.tar.gz`) for your OS and architecture.
3.  Extract the archive.
4.  Move the `getctx` (or `getctx.exe`) executable to a directory in your system's `PATH` (e.g., `/usr/local/bin` on Linux or a dedicated folder on Windows).

