using System.ComponentModel;
using System.Text;
using DotNetEnv;
using emoGpacked.APIs;
using EnumsNET;
using Spectre.Console;

namespace emoGpacked;

internal class Program
{
	enum UploadProcessType
	{
		[Description("Upload Emoji Pack")]
		EmojiPack,
		
		[Description("Upload Single Emoji")]
		SingleEmoji
	}
	
	public static RevoltClient? Client;
	
	public static IEmojiApi[] APIs = [
		new EmojiGG()
	];

	// ReSharper disable once FunctionNeverReturns
	private static void Main(string[] args)
	{
		Env.Load();
		
		// Try to get the token from the environment
		var token = Environment.GetEnvironmentVariable("TOKEN");
		while (string.IsNullOrEmpty(token))
		{
			// Ask for the token
			token = AnsiConsole.Ask<string>("Enter your Revolt token: ");

			if(string.IsNullOrEmpty(token))
			{
				Console.WriteLine("Token cannot be empty.");
				continue;
			}
			
			// Save the token to .env
			var env = File.OpenWrite(".env");
			var writer = new StreamWriter(env);
			
			writer.WriteLine("TOKEN=" + token);
			writer.Close();
			
			Console.WriteLine("Token saved to .env");
		}
		
		Client = new RevoltClient(token);
		
		AnsiConsole.MarkupLine("[bold green]Successfully made client![/]");
		
		// Ask which API to use
		var api = AnsiConsole.Prompt(
			new SelectionPrompt<IEmojiApi>()
				.Title("Which API do you want to use?")
				.AddChoices(APIs)
				.UseConverter(api => api.APIName)
		);

		// var servers = Utils.GetServers();

		var serverId = "01JMYEVRSYXG98MT2504QAPQ2R";
		
		// Ask the user if they want
		while (true)
		{
			DoAskLoop(api, serverId);
		}
	}

	private static void DoAskLoop(IEmojiApi api, string serverId)
	{
		// Ask for the process type
		var processType = AnsiConsole.Prompt(
			new SelectionPrompt<UploadProcessType>()
				.Title("What do you want to do?")
				.AddChoices(Enum.GetValues(typeof(UploadProcessType)).Cast<UploadProcessType>())
				.UseConverter(type => type.AsString(EnumFormat.Description) ?? type.ToString())
		);
			
		// Ask for the ID of the pack/emoji
		var sb = new StringBuilder();
		sb.Append("Enter the ID of the ");
		sb.Append(processType == UploadProcessType.SingleEmoji ? "emoji" : "emoji pack");
		sb.Append(": ");
		
		var id = AnsiConsole.Ask<string>(sb.ToString());

		switch (processType)
		{
			case UploadProcessType.EmojiPack:
				Utils.ProcessEmojiPack(api, id, serverId);
				break;
			case UploadProcessType.SingleEmoji:
				Utils.ProcessSingleEmoji(api, id, serverId);
				break;
			default:
				throw new ArgumentOutOfRangeException();
		}
	}
}