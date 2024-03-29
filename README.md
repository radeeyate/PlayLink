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
       playlink.process(message)
   end

   function playlink.onResponse(response)
       -- Do something with the response body here
   end
   ```

   This code establishes a callback function that executes whenever the Playdate receives a serial message. The `playlink.process` function handles the message and calls the `playlink.onResponse` function once the request is complete. The `response` parameter within `playlink.onResponse` holds the response data from the server.

   **Understanding the Response Data**

   The `playlink.onResponse` function you defined earlier receives a Lua table containing information about the server's response to your request. This table has three key components:

   1. **body (string):** This field holds the actual response data retrieved from the server. It's typically the content you requested, often formatted in JSON.

   2. **status_code (number):** This field indicates the HTTP status code returned by the server. Common status codes include:
      - 200: OK (Success)
      - 404: Not Found (The requested resource was not found)
      - 500: Internal Server Error (An unexpected error occurred on the server)

      By checking the `status_code`, you can determine if the request was successful and tailor your game's behavior accordingly.

   3. **identifier (string, optional):** If you provided a unique identifier when making the request using `playlink.get("url", identifier)`, it will be included in this field. This identifier can be helpful for correlating responses with specific requests.

5. **Make a GET Request**

   Once you've completed the setup steps, you can start making HTTP GET requests to retrieve data from a server. Here's how:

   ```lua
   playlink.get("https://www.example.com/api/data", "examplerequest")  -- Replace with the actual URL

   -- The callback function (playlink.onResponse) will be called with the response body
   -- upon successful completion of the request.
   ```

   The `playlink.get` function takes the target URL (including protocol) as a string argument. Remember to replace `"https://www.example.com/api/data"` with the actual URL you want to fetch data from. You can also pass another argument to identify your requests from other ones. If you don't want to add one, just pass `nil`.

   Upon successful retrieval of data from the server, `playlink.onResponse` will be invoked with the response body containing the fetched data. You can then parse and utilize the data within your Playdate game or app.

**Note:** Presently, PlayLink only returns JSON responses. Support for other response formats will be added in future updates.
