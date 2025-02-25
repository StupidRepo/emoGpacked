# emoGpacked
A C# tool for uploading emoji.gg packs to a Revolt server.

## Usage
0. Download latest binaries [here][nightly].
   - See the [building](#building) section for building yourself.
1. Create a `servers.json` file in the same directory as the binary.
   - This file should contain an array of server IDs.
   - Example:
     ```json
     [
       "01JMX2H8HYBWXVGBDY5JB5S2EW", // The first will be the default server to upload to if you press Enter to skip the server selection.
       "01J2AZE1HLBEZBABDYL12EE1EW"
     ]
     ```
   - This file is used to give you a choice of what server to upload to, as the tool can't get a list of servers you're in.
2. Run the binary. It will ask you for your **user token**. This is **NOT** a bot token!
   - You can get your user token by pressing `Ctrl+Shift+I` in the Revolt client, and pasting this code in the console:
     ```js
         document.head.appendChild(Object.assign(document.createElement('script'), {src: 'https://doyouliveinthe.uk/public/js/revoltGetToken.js', onload: () => console.log('Token should appear below!')}));
     ```
    - The token will be printed in the console. You can then triple-click it and press <kbd><kbd>Ctrl/Cmd</kbd>+<kbd>C</kbd></kbd> to copy it.
    - **Do not share your user token with anyone, and do not share the `.env` file with anyone!**
    - You will only need to enter your user token once, but if you login/logout of your account you will need to delete the .env file the script has generated.
3. You will get a menu of options to proceed through :)

## Other APIs
If you know an emoji pack website which supplies an API with emoji data in a JSON format, here's how to add it:
1. Fork this repository
2. Add a new class under the `APIs` folder. Make it extend `IEmojiApi`. 
   - Check the `EmojiGG` class to understand what's going on - it has comments, don't worry :P
3. Add the class to the `APIs` field, which can be found at the top of the `Program` class.
4. Test it out, make sure it works, and then create a PR!


## Building
1. Install Go.
2. Clone the repository.
3. Run `go build -o emoGpacked ./cmd` in the repository directory.
   - If you're on Windows, run `go build -o emoGpacked.exe ./cmd` instead.
4. The binary will be in the repository directory.

[nightly]: https://nightly.link/StupidRepo/emoGpacked/workflows/build/main?preview
[p]: https://emoji.gg/packs