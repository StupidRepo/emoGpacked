using Newtonsoft.Json;
#pragma warning disable CS8618

namespace emoGpacked.Models;

public class ServerSentErrorException(ServerError serverError, string extra) : Exception
{
	public override string Message { get; } = $"Revolt returned an error: {serverError.Type} ({extra})";
	public override string? Source { get; set; } = serverError.Location;
}