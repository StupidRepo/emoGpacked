using emoGpacked.Models;

namespace emoGpacked.APIs;

public interface IEmojiApi
{
	public string APIName { get; }
	
	public Emoji[] GetEmojisFromPack(string packId);
	public Emoji GetEmojiFromId(string emojiId);
}