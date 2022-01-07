using Microsoft.JSInterop;

namespace Client.Logic
{
    public class AudioService
    {
        private readonly IJSRuntime _jsRuntime;

        public AudioService(IJSRuntime runtime) { 
            _jsRuntime = runtime;
        }
        public async Task PlaySound(string name)
        {
            await _jsRuntime.InvokeAsync<string>("PlayAudio", name);
        }
    }
}
