# emoGpacked
A Go tool for uploading emoji.gg packs to a Revolt server.

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
         document.head.appendChild(Object.assign(document.createElement('script'), {src: 'https://doyouliveinthe.uk/public/js/revoltGetToken.js', onload: () => console.log('Token should appear below!'), onerror: () => console.error('Failed to load the script.')}));
     ```
    - The token will be printed in the console. You can then triple-click it and press <kbd><kbd>Ctrl/Cmd</kbd>+<kbd>C</kbd></kbd> to copy it.
    - **Do not share your user token with anyone, and do not share the `.env` file with anyone!**
    - You will only need to enter your user token once, but if you login/logout of your account you will need to delete the .env file the script has generated.
3. Enter the URL of the emoji pack you want to upload.
   - You can find emoji packs from [emoji.gg/packs][p] or [other supported websites](#other-apis).
   - If you are using an emoji.gg pack link, copy the link and append `&type=json` to the end.
     - Example: `https://emoji.gg/pack/14495-kermit-pack&type=json`
4. A list will appear with the servers you have in the `servers.json` file.
   - Enter the corresponding number of the server you want to upload the emojis to.
   - If you want to upload to the default server, just press Enter.
5. The emojis will be uploaded to the server you selected.

## Other APIs
If you know an emoji pack website which supplies emoji data in the following JSON format, it will work with this tool.
```json
{
  // Doesn't matter if there are fields before/after this. Just as long as there is an array called "emotes".
  "emotes": [
    {
      "name": "23232-emoji-name", // Format = [ID]-[Name seperated by hyphens]
      "url": "https://example.com/emoji.png"
    }
  ]
}
```

## Building
1. Install Go.
2. Clone the repository.
3. Run `go build -o emoGpacked ./cmd` in the repository directory.
   - If you're on Windows, run `go build -o emoGpacked.exe ./cmd` instead.
4. The binary will be in the repository directory.

[nightly]: https://nightly.link/StupidRepo/emoGpacked/workflows/build/main?preview
[p]: https://emoji.gg/packs