using Newtonsoft.Json;

namespace emoGpacked.Models;

public class Server
{
	[JsonProperty("name")] public string Name { get; set; }
	[JsonProperty("_id")] public string Id { get; set; }
}