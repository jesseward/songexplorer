# songexplorer

![Song Explorer](http://lh3.googleusercontent.com/FOJvQH9QkBAyUsnefWKKY-F4JvSAQQ80FXX0PnJSiIXM9J0BSG0F4pWmBpwzAtmFTLeHL7HroAW_2Q "Song Explorer")

Google Home assistant for music recommendations, built with Go and Redis. Using Google Home and API.ai

In a nutshell this service consists of the following

* a very light Go web API (https://github.com/jesseward/songexplorer) 
* An agent created at api.ai
* All api.ai intents, with the exception of the 'help' intent  are answered by the songdiscover webhook (this Go app).
* A self hosted Go app. Its sole purpose is to massage communication between api.api and the last.fm API. Calls that result in a redis miss, pass through to the last.fm API.

# Example Google Home invocations 

* OK Google, let me talk to song explorer
* OK Google, ask song explorer about artist Nightmares on Wax
* OK Google, ask song explorer what artists are similar to Underground Resistance
* OK Google, ask song explorer what are the popular songs by Boards of Canada
* OK Google, ask song explorer what songs are similar to Caught Up by Metro Area
* OK Google, ask song explorer what can i do
