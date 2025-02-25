using Newtonsoft.Json;
#pragma warning disable CS8618

namespace emoGpacked.Models;

public class ServerError
{
	[JsonProperty("type")]
	public string Type { get; set; }
	
	[JsonProperty("location")]
	public string Location { get; set; }
}