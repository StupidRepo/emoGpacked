using emoGpacked.Models;

namespace emoGpacked.APIs;

public class EmojiGG : IEmojiApi
{
	public string APIName => "emoji.gg"; // Displayed to user
	private string APIUrl => "https://emoji.gg/";
	
	public Emoji[] GetEmojisFromPack(string packId)
	{
		var response = RevoltClient.Instance.Get<EmojiGGPacksResponse>($"{APIUrl}pack/{packId}&type=json");
		var emotes = response.Result!.Emotes;
		
		foreach (var emote in emotes)
		{
			emote.Name = GetCleanEmojiName(emote.Name);
			
			emote.FileInfo = Utils.GenerateEmojiFileInfo(emote); // Generate file info for the emoji - you need to do this before returning out of this function.
			emote.FormFileName = GetCleanEmojiName(emote.Url.Split('/').Last()); // The file name used in the form data.
		}
		return emotes;
	}

	public Emoji GetEmojiFromId(string emojiId)
	{
		throw new NotSupportedException("There is currently no way to get a single emoji by ID from emoji.gg, due to their very tiny API. Therefore, you cannot use this method.");
	}

	// THIS IS REQUIRED because Revolt SUPPORTS ONLY lowercase alphabet-only names for emojis.
	private string GetCleanEmojiName(string dirtyBoy)
	{
		var split = dirtyBoy.Split('-');
		return string.Join(null, split[1..]);
	}
}

public class EmojiGGPacksResponse
{
	public Emoji[] Emotes;
}