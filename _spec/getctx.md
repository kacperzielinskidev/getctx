# Technical Documentation / LLM Context: `getctx` Tool

## 1. Purpose and Core Functionality

`getctx` (Get Context) is a Command-Line Interface (CLI) tool written in Go. Its primary purpose is to allow a user to **interactively select files and folders** from the filesystem and then **concatenate the contents** of all selected (and eligible) text files into a single, large output file (defaults to `context.txt`).

The main use case is to quickly aggregate source code context from a project, which can then be pasted into AI tools (LLMs), bug reports, or documentation.

## 2. Project Architecture

The project follows modern Go application design, emphasizing a clear **Separation of Concerns**. The core logic is managed by a central `App` orchestrator in the `core` package, with distinct packages for the CLI setup, TUI, business logic, filesystem abstraction, and other utilities.

- **`cmd/getctx/main.go`**: **Minimal Entrypoint**. Its sole responsibility is to call the `Run` function from the `internal/cli` package and handle any fatal errors at the highest level.

- **`internal/cli/cli.go`**: **CLI Layer**. This package is responsible for:

  - Parsing command-line flags (e.g., `-o`, `--debug`).
  - Initializing dependencies like the logger, configuration, and filesystem abstraction.
  - Instantiating the core application orchestrator (`core.App`) and the `ContextBuilder`.
  - Running the main application.

- **`internal/core/app.go`**: **Application Orchestrator**.

  - Contains the central `App` struct that encapsulates the application's lifecycle and dependencies (`logger`, `ContextBuilder`, `fsys`, etc.).
  - **`Run()`**: Executes the main application flow: starts the TUI, waits for it to exit, and then passes the final state to the `ContextBuilder` to perform the main business logic.

- **`internal/tui/`**: **Terminal User Interface (TUI) Core**.

  - All interactive logic resides in this package, built using the **`bubbletea`** library and the Model-View-Update pattern.
  - **`model.go`**: Defines the `Model` struct, which holds the entire TUI state, including the cursor, selected items, current path, and modes (e.g., input mode, filter mode).
  - **`view.go`**: Contains the `View()` method, which renders the UI based on the model's state.
  - **`update.go`**: Contains the `Update()` method, a key component that handles all user input and state changes. It performs dynamic layout calculation to ensure the UI remains responsive.
  - Integrates a **`viewport`** component to smoothly handle scrolling through long file lists.
  - **`theme.go`**, **`keymap.go`**, **`completions.go`**: Helper files that define visual styles, keybindings, and path auto-completion logic, respectively.

- **`internal/build/context_builder.go`**: **Business Logic**.

  - Contains the "brain" of the application that runs after the TUI exits.
  - It is responsible for processing the list of selected paths received from the TUI, discovering all nested files, filtering out non-text files, and building the final `context.txt`.

- **`internal/fs/`**: **File System Abstraction**.

  - Provides a clean `FileSystem` interface for improved testability (decoupling from the `os` package).
  - **`fsys.go`**: Defines the `FileSystem` interface.
  - **`os_fs.go`**: Provides the concrete implementation of the interface using standard library functions.
  - **`utils.go`**: Contains helper functions like `discoverFiles` (to recursively find eligible files) and `isTextFile` (to detect if a file is text-based).

- **`internal/config/config.go`**: **Application Configuration**.

  - Stores application-wide configuration, primarily the lists of excluded file names, folder names, and file extensions.

- **`internal/logger/logger.go`**: **Global Structured Logger**.
  - A dedicated, site-wide package for logging.
  - Logs messages in a structured **JSON format** to a `debug.log` file when enabled.
  - Includes log levels (DEBUG, INFO, WARN, ERROR) and a `name` field to identify the context of the log entry.

## 3. Key Features and Logic

- **File System Navigation:** The user navigates the filesystem with arrow keys (`handleMoveCursorUp`/`Down`). `Enter` (`handleEnterDirectory`) opens a directory, and `Backspace` (`handleNavigateToParent`) goes to the parent directory.

- **In-View Filtering (Search):**

  - **Activation:** Pressing `/` (`handleEnterFilterMode`) activates a filter input field.
  - **Live Filtering:** The list of currently visible items is filtered in real-time as the user types. The search is case-insensitive and operates purely in-memory for instant feedback.
  - **Interacting with Results:** Pressing `Enter` exits the text input mode but **keeps the view filtered**, allowing the user to navigate and select items from the search results using the standard keys (`Space`, `CTRL+A`).
  - **Clearing the Filter:** To restore the full directory view, the user can press `Escape` or `CTRL+C`. Any selections made on the filtered items will be preserved. Navigating to a new directory also clears the filter automatically.
  - **User Guidance:** A clear indicator `[Filtering by: "query"]` is displayed in the header to inform the user that their view is filtered.

- **Direct Path Input Mode:**

  - **Activation:** Pressing `CTRL+P` (`handleEnterPathInputMode`) activates a text input field.
  - **Functionality:** Allows the user to directly type or paste an absolute or relative path. Supports `~` as a shortcut for the user's home directory and provides **tab-completion**.
  - **User Guidance:** Clear, color-coded on-screen hints (`(enter: Confirm, esc/ctrl+c: Cancel)`) guide the user.
  - **Confirmation & Cancellation:** `Enter` (`handleConfirmPathChange`) attempts to navigate to the path. `Esc` or `CTRL+C` (`handleCancelPathChange`) exits the input mode without changes.
  - **Error Handling:** If an invalid path is entered, a non-disruptive error message appears directly below the input field.

- **Intelligent File Exclusion (Blacklist):** The tool maintains a configurable list of names and extensions to ignore (e.g., `.git`, `node_modules`, `.png`). Ignored items are visually dimmed and cannot be interacted with.

- **Selection:**

  - `Spacebar` (`handleSelectFile`): Toggles selection for a single item (works on both full and filtered lists).
  - `CTRL+A` (`handleSelectAllFiles`): Toggles selection for all _visible_ items (works on both full and filtered lists).

- **Dynamic & Responsive UI:** The TUI is fully responsive. The `viewport` ensures that lists of any length are scrollable, and the appearance of status messages or input fields correctly resizes the view without breaking the UI layout.

- **Structured Logging:** All significant application events, warnings, and errors are logged to `debug.log` in a machine-readable JSON format for easier debugging.

- **Program Exit & Cancellation:**

  - `q` (`handleConfirmAndExit`): Exits and initiates the build process with the current selections.
  - `CTRL+C`: This key is now context-aware:
    - If the view is **filtered**, it clears the filter.
    - If an **input field** (path or filter) is active, it cancels the input.
    - Otherwise, it clears all selections and exits the application (`handleCancelAndExit`).
  - `Escape`: Clears an active filter or cancels an input field.

- **Enhanced Styling:** The application uses `lipgloss` for a modern look. Keybinding hints in the help text are color-coded to improve usability.

## 4. External Dependencies

- `github.com/charmbracelet/bubbletea`: The TUI framework.
- `github.com/charmbracelet/bubbles/viewport`: The component for scrollable views.
- `github.com/charmbracelet/bubbles/textinput`: The component that provides text input fields for the 'Direct Path Input' and 'In-View Filtering' features.
- `github.com/charmbracelet/lipgloss`: The library for terminal styling.

## 5. Project structure

```GETCTX
├── _spec
│   ├── getctx.md
│   └── rules.md
├── .github
│   └── workflows
│       └── release.yml
├── cmd
│   └── getctx
│       └── main.go
├── internal
│   ├── build
│   │   └── context_builder.go
│   ├── cli
│   │   └── cli.go
│   ├── config
│   │   └── config.go
│   ├── core
│   │   └── app.go
│   ├── fs
│   │   ├── fsys.go
│   │   ├── os_fs.go
│   │   └── utils.go
│   ├── logger
│   │   └── logger.go
│   └── tui
│       ├── completions.go
│       ├── keymap.go
│       ├── model.go
│       ├── theme.go
│       ├── update.go
│       └── view.go
├── .gitignore
├── context.txt
├── go.mod
├── go.sum
└── Makefile
```
