using Newtonsoft.Json;
using Spectre.Console;
using Emoji = emoGpacked.Models.Emoji;

namespace emoGpacked.APIs;

public class EmojiGG : IEmojiAPI
{
	public string APIName => "emoji.gg"; // Displayed to user
	private string APIUrl => "https://emoji.gg/";
	
	public string? InputText => null; // Displayed to user
	
	public Emoji[] GetEmojisFromPack(string packId, RevoltAutumnClient autumnClient)
	{
		var response = autumnClient.Get<EmojiGGPacksResponse>($"{APIUrl}pack/{packId}&type=json");
		var emotes = response.Result!.Emotes;
		
		foreach (var emote in emotes)
		{
			emote.Name = GetCleanEmojiName(emote.Name);
			
			emote.FileInfo = Utils.GenerateEmojiFileInfo(emote); // Generate file info for the emoji - you need to do this before returning out of this function.
			emote.FormFileName = GetCleanEmojiName(emote.Url.Split('/').Last()); // The file name used in the form data.
		}
		return emotes;
	}

	public Emoji GetEmojiFromId(string emojiId, RevoltAutumnClient autumnClient)
	{
		// so.... emoji.gg has one endpoint but it returns a list of every emoji ever made on the site since it first launched (about 5K at the time of writing this comment)
		// why the f**k they would do that, i have no idea. but it's all we got to work with :sob:
		
		// the url in question btw: https://emoji.gg/api/
		// so we're gonna have to get all the emojis and then filter them by the ID
		
		// this only responds with 5K emojis wtf piss off emoji.gg
		// var response = RevoltClient.Instance.Get<EmojiGGEmote[]>($"{APIUrl}api/");
		// var emotes = response.Result!;
		//
		// var emote = emotes.FirstOrDefault(e => e.Slug == emojiId);
		// if (emote == null)
		// {
		// 	throw new Exception("Emoji not found! Maybe it's too old...");
		// }
		//
		// // Ask the user for the name of the emoji
		// var name = AnsiConsole.Ask("What do you want the emoji to be called? (lowercase alphabet only)", GetCleanEmojiName(emote.Title));
		//
		// var emoji = new Emoji
		// {
		// 	Name = GetCleanEmojiName(name),
		// 	Url = emote.Image,
		// 	
		// 	FormFileName = emote.Image.Split('/').Last()
		// };
		//
		// emoji.FileInfo = Utils.GenerateEmojiFileInfo(emoji);
		//
		// return emoji;
		
		throw new NotSupportedException("This API does not support getting a single emoji by ID. Emoji.gg does not have an easy way to implement this, so it's not supported.");
	}

	// THIS IS REQUIRED because Revolt SUPPORTS ONLY lowercase alphabet-only names for emojis.
	private string GetCleanEmojiName(string dirtyBoy)
	{
		var split = dirtyBoy.ToLower().Split('-');
		split = split.Select(s => new string(s.Where(ch => !char.IsNumber(ch)).ToArray())).ToArray();
		
		return string.Join(null, split);
	}
}

public class EmojiGGPacksResponse
{
	public Emoji[] Emotes;
}

public class EmojiGGEmote
{
	[JsonProperty("title")] public string Title;
	[JsonProperty("slug")] public string Slug;
	[JsonProperty("image")] public string Image;
}