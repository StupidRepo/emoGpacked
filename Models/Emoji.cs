using Newtonsoft.Json;

namespace emoGpacked.Models;

public class Emoji
{
	[JsonProperty("name")] public string Name { get; set; } // Name of the emoji (e.g. "smile")
	[JsonProperty("url")] public string Url { get; set; } // URL to the emoji
	
	public string FormFileName { get; set; } // File name of the emoji used for form uploads
	public FileInfo? FileInfo { get; set; } // File info of the emoji

	public class UploadResponse
	{
		[JsonProperty("id")] public string Id { get; set; }
	}

	public class FinaliseCommon
	{
		[JsonProperty("name")] public string Name { get; set; }
		[JsonProperty("parent")] public Parent Parent { get; set; }
	}
}