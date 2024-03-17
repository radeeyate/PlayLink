# PlayLink

PlayLink is a protocol for the Playdate console that allows you to make arbitrary HTTP(S) requests.
It uses the Playdate serial interface to interact with the "server" running on a computer.

## Running the server

1. Clone the repository:

    Begin by cloning the PlayLink repository using Git:

    ```bash
    git clone https://github.com/radeeyate/PlayLink.git
    ```

2. **Install Go:**

   PlayLink's server component is written in Go. To proceed, you'll need Go installed on your system. You can download and install Go from the official website: [https://go.dev/doc/install](https://go.dev/doc/install)

3. **Build the Server:**

   Navigate to the PlayLink server directory:

   ```bash
   cd PlayLink/server
   ```

   Then, build the server executable using the following commands:

   ```bash
   go get
   go build
   ```

   This will generate a binary file. The filename will vary depending on your operating system:

   - **Windows:** PlayLink.exe (not officially tested yet)
   - **Linux:** PlayLink

   On Linux, grant executable permissions to the server file using chmod:

   ```bash
   chmod +x ./PlayLink
   ```

4. **Run the Server:**

   Execute the server using the following command:

   ```bash
   ./PlayLink
   ```

## Using the Lua Client

To leverage PlayLink's functionalities within your Playdate project, follow these steps:

1. **Copy Lua Files:**

   Copy the `b64.lua` and `playlink.lua` files from the PlayLink repository into your Playdate project's directory.

2. **Import Library:**

   Within your Playdate Lua code, include the PlayLink library using the following statement:

   ```lua
   local playlink = import "playlink"
   ```

3. **Initialize Library:**

   PlayLink requires initialization before you can use its functions. Call the following code snippet to initialize the library:

   ```lua
   playlink.init()
   ```

   Currently, there are no error messages if the initialization fails. However, subsequent PlayLink functions will not work without successful initialization.

4. **Serial Message Processing:**

   To enable PlayLink to recognize serial messages from the Playdate, incorporate the following code into your project:

   ```lua
   function playdate.serialMessageReceived(message)
       playlink.process(message, request_return)
   end

   function request_return(body)
       -- Do something with the response body here
   end
   ```

   This code establishes a callback function that executes whenever the Playdate receives a serial message. The `playlink.process` function handles the message and calls the `request_return` function (or whichever function you specify as the second argument) once the request is complete. The `body` parameter within `request_return` holds the response data from the server.

**Note:** Presently, PlayLink only returns JSON responses. Support for other response formats will be added in future updates.
