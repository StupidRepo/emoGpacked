using emoGpacked.Models;
using Newtonsoft.Json;

namespace emoGpacked;

public class RevoltClient
{
	public readonly string Token;

	public const string ApiBase = "https://app.revolt.chat/api/";
	public const string AutumnApiBase = "https://autumn.revolt.chat/";

	public readonly HttpClient RequestClient = new();
	
	public static RevoltClient Instance { get; private set; } = null!;

	public RevoltClient(string token)
	{
		Instance = this;
		Token = token;
	}

	public ApiResponse<T> GetEndpoint<T>(string endpoint)
	{
		return Get<T>(ApiBase + endpoint);
	}
	
	public ApiResponse<T> Get<T>(string url)
	{
		var request = new HttpRequestMessage(HttpMethod.Get, url);
		return DoRequest<T>(request);
	}

	public ApiResponse<T> DoRequest<T>(HttpRequestMessage requestMessage)
	{
		requestMessage.Headers.Add("X-Session-Token", Token);
    
		var response = RequestClient.SendAsync(requestMessage).Result;
		if (!response.IsSuccessStatusCode)
		{
			var eResult = response.Content.ReadAsStringAsync().Result;
			return new ApiResponse<T>
			{
				Error = JsonConvert.DeserializeObject<ServerError>(eResult)
			};
		}
    
		var content = response.Content.ReadAsStringAsync().Result;
		return new ApiResponse<T>
		{
			Result = JsonConvert.DeserializeObject<T>(content, new JsonSerializerSettings
			{
				CheckAdditionalContent = false
			})
		};
	}
	
	public ApiResponse<T> UploadToAutumn<T>(string endpoint, string filePath, string fileName)
	{
		var request = new HttpRequestMessage(HttpMethod.Post, AutumnApiBase + endpoint);
		request.Content = new MultipartFormDataContent
		{
			{ new ByteArrayContent(File.ReadAllBytes(filePath)), "file", fileName }
		};
		return DoRequest<T>(request);
	}
}

public class ApiResponse<T>
{
	public T? Result { get; set; }
	public ServerError? Error { get; set; }
	public bool IsSuccess => Error == null;
}