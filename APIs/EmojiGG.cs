using emoGpacked.Models;

namespace emoGpacked.APIs;

public class EmojiGG : IEmojiApi
{
	public string APIName => "emoji.gg";
	private string APIUrl => "https://emoji.gg/";
	
	public Emoji[] GetEmojisFromPack(string packId)
	{
		var response = RevoltClient.Instance.Get<EmojiGGPacksResponse>($"{APIUrl}pack/{packId}&type=json");
		var emotes = response.Result!.Emotes;
		
		foreach (var emote in emotes)
		{
			emote.Name = GetCleanEmojiName(emote.Name);
			emote.FileInfo = Utils.GenerateEmojiFileInfo(emote);
		}
		return emotes;
	}

	public Emoji GetEmojiFromId(string emojiId)
	{
		throw new NotSupportedException("There is currently no way to get a single emoji by ID from emoji.gg, due to their very tiny API. Therefore, you cannot use this method.");
	}

	private string GetCleanEmojiName(string dirtyBoy)
	{
		var split = dirtyBoy.Split('-');
		return string.Join('-', split[1..]);
	}
}

public class EmojiGGPacksResponse
{
	public Emoji[] Emotes;
}