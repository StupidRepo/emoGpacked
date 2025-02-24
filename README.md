# emoGpacked
A Go tool for uploading emoji.gg packs to a Revolt server.

## Usage
0. A Windows binary is available in the [releases][r] tab.
   - See the [building](#building) section for building the yourself.
1. Go to [emoji.gg/packs][p] and choose a pack.
2. Copy the pack's URL (e.g. `https://emoji.gg/pack/14495-kermit-pack`), and add `&type=json` to the end. 
   - e.g. `https://emoji.gg/pack/14495-kermit-pack&type=json`
3. Run the tool with [the flags](#flags) required.

## Flags
| Flag | Description                                                                                                     | Required |
|------|-----------------------------------------------------------------------------------------------------------------|----------|
| `-p` | The URL of the emoji.gg pack.                                                                                   | Yes      |
| `-s` | The server ID of the server to add the emojis to.                                                               | Yes      |
| `-t` | Your user login token.                                                                                          | Yes      |
| `-g` | Upload animated emojis (disabled by default due to weird Revolt bug - upload animated GIFs using the frontend). | No       |

## Example
```sh
# No GIF emojis
emoGpacked -p "https://emoji.gg/pack/14495-kermit-pack&type=json" -s "01JMX2H8HYBWXVGBDY5JB5S2EW" -t "user-token"

# GIF emojis
emoGpacked -p "https://emoji.gg/pack/14495-kermit-pack&type=json" -s "01JMX2H8HYBWXVGBDY5JB5S2EW" -t "user-token" -g
```

## Other APIs
If you know an emoji pack website which supplies data in the following JSON format, it will work with this tool.
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

[r]: https://github.com/StupidRepo/emoGpacked/releases
[p]: https://emoji.gg/packs