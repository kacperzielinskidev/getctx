# Your Program Name

A brief description of your program. Write one or two sentences here explaining what your tool does.

## Installation

Below are the instructions on how to download and install the program so that it is globally available on your system.

---

### 🐧 Linux

Installation on Linux involves downloading the archive, extracting it, giving the binary file execution permissions, and moving it to the `/usr/local/bin` directory, which is typically in the system's `PATH`.

1.  **Download the archive**
    Open a terminal and use the `curl` or `wget` command to download the latest version of the program in `.tar.gz` format.

    ```sh
    # Using curl
    curl -LO https://github.com/kacperzielinskidev/getctx/releases/download/v1.0.0/getctx_v1.0.0_linux_amd64.tar.gz

    # Or using wget
    wget https://github.com/kacperzielinskidev/getctx/releases/download/v1.0.0/getctx_v1.0.0_linux_amd64.tar.gz
    ```

2.  **Extract the archive**
    Use the `tar` command to extract the contents of the archive.

    ```sh
    tar -xzf archive-name.tar.gz
    ```

3.  **Grant execution permissions**
    After extracting, you will find the binary file. Grant it execution permissions.

    ```sh
    chmod +x your-program-name
    ```

4.  **Move the file to `/usr/local/bin`**
    Moving the binary file to this directory will make it accessible from anywhere in the system. You will need administrator privileges (`sudo`).

    ```sh
    sudo mv your-program-name /usr/local/bin/
    ```

5.  **Done!**
    You can now run the program by typing its name from any location in the terminal.

    ```sh
    your-program-name --version
    ```

---

### 🪟 Windows

On Windows, the process involves downloading the `.zip` archive, extracting it, and then adding the folder containing the `.exe` file to the system's `PATH` environment variable.

1.  **Download the archive**
    Download the latest `.zip` file from the [Releases](URL_TO_YOUR_REPOSITORY/releases) section on GitHub.

2.  **Extract the archive**
    Right-click the downloaded `.zip` file and select "Extract All...". Follow the on-screen instructions to extract the contents to a new folder.

3.  **Create a destination folder and move the file**

    - Create a dedicated, permanent folder for the program, for example `C:\Program Files\YourProgram`.
    - Move the extracted `your-program-name.exe` file into this newly created folder.

4.  **Add the folder to the PATH environment variable**
    This allows Windows to find your program from any command line.

    - Press the `Windows` key and type "environment variables".
    - Select "Edit the system environment variables".
    - In the new window, click the "Environment Variables..." button.
    - In the "System variables" section, find and select the `Path` variable, then click "Edit...".
    - Click "New" and paste the path to the folder you created, i.e., `C:\Program Files\YourProgram`.
    - Confirm all open windows by clicking "OK".

5.  **Restart your terminal**
    Close any open terminal windows (CMD or PowerShell) and open a new one. Changes to the `PATH` variable require a new terminal session to take effect.

6.  **Done!**
    You can now run the program from any location.

    ```shell
    your-program-name.exe --version
    ```
