using emoGpacked.Models;

namespace emoGpacked.APIs;

public interface IEmojiApi
{
	public string APIName { get; } // Displayed to user
	
	public Emoji[] GetEmojisFromPack(string packId);
	public Emoji GetEmojiFromId(string emojiId);
}