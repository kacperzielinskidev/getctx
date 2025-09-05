# Technical Documentation / LLM Context: `getctx` Tool

## 1. Purpose and Core Functionality

`getctx` (Get Context) is a Command-Line Interface (CLI) tool written in Go. Its primary purpose is to allow a user to **interactively select files and folders** from the filesystem and then **concatenate the contents** of all selected (and eligible) text files into a single, large output file (defaults to `context.txt`).

The main use case is to quickly aggregate source code context from a project, which can then be pasted into AI tools (LLMs), bug reports, or documentation.

## 2. Project Architecture

The project is structured into several files, each with a clearly defined responsibility, following the **Separation of Concerns** principle:

- **`main.go`**: **Entrypoint**. It is responsible for:

  - Parsing command-line flags (e.g., `-o` for the output file).
  - Initializing and running the TUI.
  - Passing the results from the TUI to the context-building logic.
  - Handling top-level errors.

- **`tui.go`**: **Terminal User Interface (TUI)**.

  - All interactive logic resides here.
  - Built using the **`bubbletea`** library and the **Model-View-Update** architectural pattern.
  - **Model (`model` struct):** Holds the entire state of the interface: the current path, the list of files/folders, the cursor position, and a map of selected items.
  - **View (`View()` method):** Renders the model's state to the terminal screen, using styles from the `lipgloss` library.
  - **Update (`Update()` method):** Handles all user input (key presses) and modifies the model's state accordingly.

- **`context_builder.go`**: **Business Logic**.

  - Contains the "brain" of the application that runs after the TUI exits.
  - Responsible for processing the list of selected paths, filtering out binary files, and building the final `context.txt` file.
  - Prints detailed logs to the console about which files are being added and which are being skipped.

- **`fs_utils.go`**: **File System Utilities**.

  - A collection of general-purpose helper functions for file and folder operations.
  - **`discoverFiles`**: Recursively scans the given paths and returns a list of all found files, classified as either text or binary.
  - **`isTextFile`**: Detects if a given file is a text file (using a MIME type-based heuristic).

- **`ui.go`**: **UI Definitions**.

  - A central place for managing the application's look and feel.
  - Defines icons (emoji), colors, and complex styles (`lipgloss.Style`) used in both the TUI and the logs.
  - Styles are organized into nested structs (`TUIStyles`, `TUIListStyles`, `TUILogStyles`) for better readability and scalability.

- **`keybindings.go`**: **Keybinding Definitions**.
  - A central place that defines all keyboard shortcuts used in the application as constants (`const`).
  - Eliminates "magic strings" in the TUI logic and makes reconfiguring keys easy.

## 3. Key Features and Logic

- **Navigation:** The user navigates the filesystem using the up/down arrow keys. `Enter` opens a directory, and `Backspace` goes to the parent directory.
- **Selection:**
  - `Spacebar`: Toggles the selection for the individual file or folder under the cursor.
  - `CTRL+A`: Toggles the selection for all items in the current view.
- **Binary File Filtering:** The tool automatically detects and skips binary files (e.g., `.exe`, images, compiled artifacts), including only the content of text files in the final output. The user is informed about any skipped files.
- **Program Exit:**
  - `q`: Exits the interface and **initiates the build process** for the `context.txt` file from the selected items.
  - `CTRL+C`: **Cancels the operation**. Exits the program without saving the file.
- **Styling:** The application makes extensive use of the `lipgloss` library for styling (colors, bolding) in both the interactive TUI and the final log output, ensuring a consistent and modern look.

## 4. External Dependencies

- `github.com/charmbracelet/bubbletea`: The foundation of the TUI.
- `github.com/charmbracelet/lipgloss`: For styling text in the terminal.
