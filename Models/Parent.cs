using Newtonsoft.Json;

namespace emoGpacked.Models;

public class Parent
{
	[JsonProperty("type")] public string Type { get; set; }
	[JsonProperty("id")] public string Id { get; set; }
}