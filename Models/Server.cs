using Newtonsoft.Json;

namespace emoGpacked.Models;

public class Server
{
	[JsonProperty("name")] public string Name { get; set; }
}