# Usage

The Duolingo API only allows JWT login: you can’t log in using your username and password. To acquire your JWT from Duolingo, open up the network console and visit the Duolingo site. Be sure that you’re logged in. Find a request from your browser to the Duolingo API: you should see a jwt_token cookie on that request.

<a id="module-duolingo.duolingo"></a>

Unofficial API for duolingo.com

#### func NewDuolingo(username, token, url, abbr string) *Duolingo

Creates a connection to the Duolingo server, providing your JWT for authentication.

#### func GetKnownWords() map[string]struct{}

Returns a set of all words learned by a user in a language.

#### func GetVocab() []interface{}

Returns a slice of maps containing info on a user's vocab learned in a language.
