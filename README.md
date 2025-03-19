# emoGpacked
A C# tool for uploading emoji.gg packs to a Revolt server.

## Usage
0. No latest binaries because this is .NET and it's weird.
   - See the [building](#building) section for building yourself.
1. Run the binary. It will ask you for your **user token**. This is **NOT** a bot token!
   - You can get your user token by pressing `Ctrl+Shift+I` in the Revolt client, and pasting this code in the console:
     ```js
        const c = controllers.client.getReadyClient()
        console.log(c.session.token)
     ```
    - The token will be printed in the console. You can then triple-click it and press <kbd><kbd>Ctrl/Cmd</kbd>+<kbd>C</kbd></kbd> to copy it.
    - **Do not share your user token with anyone, and do not share the `.env` file with anyone!**
    - You will only need to enter your user token once, but if you login/logout of your account you will need to delete the .env file the script has generated.
2. You will get a menu of options to proceed through :)

## Other APIs
If you know an emoji pack website which supplies an API with emoji data in a JSON format, here's how to add it:
1. Fork this repository
2. Make a new class, of which extends `IEmojiAPI`, in the `APIs` folder. 
   - Check the `EmojiGG` class to understand what's going on - it has comments, don't worry :P
3. Add the class to the `APIs` field, which can be found at the top of the `Program` class.
4. Test it out, make sure it works, and then create a PR!

## Building
1. Install .NET 8.
2. Clone the repository.
3. Idk - Microsoft made it too complicated. Just do a `dotnet run` lmao.

[nightly]: https://nightly.link/StupidRepo/emoGpacked/workflows/build/c#?preview
[p]: https://emoji.gg/packs