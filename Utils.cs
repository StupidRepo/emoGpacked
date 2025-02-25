using System.Security.Cryptography;
using System.Text;
using emoGpacked.APIs;
using emoGpacked.Models;
using Newtonsoft.Json;
using Spectre.Console;
using Emoji = emoGpacked.Models.Emoji;

namespace emoGpacked;

internal static class Utils
{
	internal static readonly HttpClient DLClient = new();
	
	public static readonly string TempDir = Path.GetTempPath();

	internal static Server[] GetServers()
	{
		if(!File.Exists("servers.json"))
		{
			var file = File.CreateText("servers.json");
			file.Write("[]");
			file.Close();

			return [];
		}
		
		var serversJson = File.ReadAllText("servers.json");
		var serverIds = JsonConvert.DeserializeObject<string[]>(serversJson) ?? [];
		var servers = new List<Server>();
		
		foreach (var serverId in serverIds)
		{
			var server = RevoltClient.Instance.Get<Server>("https://api.revolt.chat/servers/" + serverId);
			if (!server.IsSuccess)
			{
				throw new ServerSentErrorException(server.Error!);
			}
			
			servers.Add(server.Result!);
		}
		
		return servers.ToArray();
	}
	
	internal static void ProcessSingleEmoji(IEmojiApi api, string emojiId, string serverId)
	{
		try
		{
			var emoji = api.GetEmojiFromId(emojiId);

			if (emoji.FileInfo == null)
			{
				throw new ArgumentNullException(nameof(emoji.FileInfo));
			}

			DownloadEmoji(emoji);

			var response = RevoltClient.Instance.UploadToAutumn<Emoji.UploadResponse>("emojis", emoji.FileInfo.FullName, emoji.FormFileName);
			if (!response.IsSuccess)
			{
				throw new ServerSentErrorException(response.Error!);
				// Stfu, IntelliJ. If you were smart, you would know IsSuccess is checking if Error is null already.
				// F**k off!
			}

			var autumnId = response.Result!.Id;
			FinaliseEmoji(serverId, autumnId, emoji);

			AnsiConsole.MarkupLine($"[bold green]Successfully uploaded {emoji.Name}![/]");
		}
		catch (Exception e)
		{
			AnsiConsole.MarkupLine($"[bold red]Error trying to process emoji {emojiId}: {e.Message}[/]");
		}
	}
	
	internal static void ProcessEmojiPack(IEmojiApi api, string packId, string serverId)
	{
		try
		{
			var emojis = api.GetEmojisFromPack(packId);
			if (emojis.Any(emoji => emoji.FileInfo == null))
			{
				throw new ArgumentNullException(nameof(emojis));
			}

			DownloadEmojis(emojis);

			foreach (var emoji in emojis)
			{
				try
				{
					AnsiConsole.MarkupLine($"[bold]Processing {emoji.Name}...[/]");

					var response =
						RevoltClient.Instance.UploadToAutumn<Emoji.UploadResponse>("emojis", emoji.FileInfo!.FullName, emoji.FormFileName);
					if (!response.IsSuccess)
					{
						throw new ServerSentErrorException(response.Error!);
					}

					var autumnId = response.Result!.Id;

					FinaliseEmoji(serverId, autumnId, emoji);

					AnsiConsole.MarkupLine($"[bold green]Successfully uploaded {emoji.Name}![/]");
				} catch (Exception e)
				{
					AnsiConsole.MarkupLine($"[bold red]Error trying to process emoji {emoji.Name}: {e.Message}[/]");
				}
			}
		} 
		catch (Exception e)
		{
			AnsiConsole.MarkupLine($"[bold red]Error trying to process pack {packId}: {e.Message}[/]");
		}
	}
	
	// AT THIS POINT we expect the Emoji to have an associated FileInfo.
	internal static void DownloadEmoji(Emoji emoji)
	{
		if (emoji.FileInfo == null)
		{
			throw new ArgumentNullException(nameof(emoji.FileInfo));
		}

		var request = new HttpRequestMessage(HttpMethod.Get, emoji.Url);
		var response = DLClient.SendAsync(request).Result;
		
		var content = response.Content.ReadAsByteArrayAsync().Result;
		
		File.WriteAllBytes(emoji.FileInfo.FullName, content);
	}

	internal static void DownloadEmojis(Emoji[] emojis)
	{
		foreach (var emoji in emojis)
		{
			DownloadEmoji(emoji);
		}
	}
	
	internal static void FinaliseEmoji(string serverId, string autumnId, Emoji emoji)
	{
		var finalise = new Emoji.FinaliseCommon
		{
			Name = emoji.Name,
			Parent = new Parent
			{
				Type = "Server",
				Id = serverId
			}
		};
		
		var json = JsonConvert.SerializeObject(finalise);
		var content = new StringContent(json, Encoding.UTF8, "application/json");
		
		var request = new HttpRequestMessage(HttpMethod.Put, RevoltClient.ApiBase + "custom/emoji/" + autumnId)
		{
			Content = content
		};
		
		var response = RevoltClient.Instance.DoRequest<Emoji.FinaliseCommon>(request);
		if (!response.IsSuccess)
		{
			throw new ServerSentErrorException(response.Error!);
		}
		
		AnsiConsole.MarkupLine("[bold green]Successfully finalised emoji![/]");
	}

	public static FileInfo GenerateEmojiFileInfo(Emoji emoji)
	{
		return new FileInfo(Path.Combine(TempDir, emoji.Url.Split('/').Last()));
	}
}