using Newtonsoft.Json;
#pragma warning disable CS8618

namespace emoGpacked.Models;

public class ServerSentErrorException(ServerError serverError) : Exception
{
	public override string Message { get; } = serverError.Type;
	public override string? Source { get; set; } = serverError.Location;
}