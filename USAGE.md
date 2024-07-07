# Usage

This fork of the Duolingo library is not published to PyPi: you can’t install it using pip, but you can copy it into your project if you also copy in the accompanying LICENSE.md.

The Duolingo API only allows JWT login: you can’t log in using your username and password. To acquire your JWT from Duolingo, open up the network console and visit the Duolingo site. Be sure that you’re logged in. Find a request from your browser to the Duolingo API: you should see a jwt_token cookie on that request.

<a id="module-duolingo.duolingo"></a>

Unofficial API for duolingo.com

### *class* duolingo.duolingo.Duolingo(\*, jwt: str)

Bases: `object`

Create a connection to the Duolingo server, providing your JWT for authentication.

```pycon
>>> from duolingo import Duolingo
>>> lingo = duolingo.Duolingo(jwt=MY_DUOLINGO_JWT)
>>>
>>> lingo.get_languages()  # Issue an arbitrary command
['...', 'High Valyrian'...]
```

#### get_abbreviation_of(language: str)

Get abbreviation of a language.

```pycon
>>> lingo.get_abbreviation_of(language='High Valyrian')
'hv'
```

#### get_daily_xp_progress()

Return information about the user’s daily learning.

```pycon
>>> lingo.get_daily_xp_progress()
{'xp_goal': ..., ..., 'xp_today': ...}
```

#### get_friends()

Get user’s friends.

```pycon
>>> lingo.get_friends()
[{'username': '...', ..., 'displayName': '...'}...]
```

#### get_golden_topics(abbreviation: str)

Return the topics mastered (“golden”) by a user in a language.

```pycon
>>> lingo.get_golden_topics(abbreviation=MY_FAVORITE_LANGUAGE_ABBREVIATION)
[...]
```

#### get_known_topics(abbreviation: str)

Return the topics learned by a user in a language.

```pycon
>>> lingo.get_known_topics(abbreviation=MY_FAVORITE_LANGUAGE_ABBREVIATION)
['...', '...', ...]
```

#### get_known_words(abbreviation: str)

Get a list of all words learned by user in a language.

```pycon
>>> lingo.get_known_words(abbreviation=MY_FAVORITE_LANGUAGE_ABBREVIATION)
[...]
```

#### get_language_details(language: str)

Get user’s status about a language.

```pycon
>>> lingo.get_language_details(language='High Valyrian')
{'streak': ..., ..., 'to_next_level': ...}
```

#### get_language_from_abbr(abbreviation: str)

Get language full name from abbreviation.

```pycon
>>> lingo.get_language_from_abbr(abbreviation='hv')
'High Valyrian'
```

#### get_language_progress(abbreviation: str)

Get information about user’s progression in a language.

```pycon
>>> lingo.get_language_progress(abbreviation='hv')
{'language_string': 'High Valyrian', ..., 'level': ...}
```

#### get_languages(as_abbreviations: bool = False)

Get practiced languages.

```pycon
>>> lingo.get_languages()
['...', 'High Valyrian'...]
```

```pycon
>>> lingo.get_languages(as_abbreviations=True)
['...', 'hv'...]
```

#### get_reviewable_topics(abbreviation: str)

Return the topics learned but not golden by a user in a language.

```pycon
>>> lingo.get_reviewable_topics(abbreviation=MY_FAVORITE_LANGUAGE_ABBREVIATION)
['...', '...', ...]
```

#### get_streak_info()

Get user’s streak information.

```pycon
>>> lingo.get_streak_info()
{'daily_goal': ..., ..., 'streak_extended_today': ...}
```

#### get_unknown_topics(abbreviation: str)

Return the topics learned by a user in a language.

```pycon
>>> lingo.get_unknown_topics(abbreviation=MY_FAVORITE_LANGUAGE_ABBREVIATION)
['...', '...', ...]
```

#### get_user_info()

Get user’s information.

```pycon
>>> lingo.get_user_info()
{'username': ..., ..., 'ui_language': ...}
```

#### get_vocabulary(abbreviation: str, source_language_abbreviation: str | None = None)

Get overview of user’s vocabulary in a language.

```pycon
>>> lingo.get_vocabulary(abbreviation=MY_FAVORITE_LANGUAGE_ABBREVIATION)
[{'text': ..., ...}, ...]
```

Note that this will only retrieve any learned words from any \_completed_ lessons: new words from
incomplete lessons will not be returned.

* **Parameters:**
  * **abbreviation** – Language abbreviation of learning language
  * **source_language_abbreviation** – Language abbreviation of source language (default: user’s UI language)

### *exception* duolingo.duolingo.DuolingoException

Bases: `Exception`

### *class* duolingo.duolingo.Struct(\*\*entries)

Bases: `object`
