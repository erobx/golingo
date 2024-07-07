# Usage

This fork of the Duolingo library is not published to PyPi: you can't install it using pip, but you can copy it into your project if you also copy in the accompanying LICENSE.md.

The Duolingo API only allows JWT login: you can't log in using your username and password. To acquire your JWT from Duolingo, open up the network console and visit the Duolingo site. Be sure that you're logged in. Find a request from your browser to the Duolingo API: you should see a `jwt_token` cookie on that request.

## *class* Duolingo(username: str, \*, jwt: str | None = None)

```pycon
from duolingo import Duolingo

lingo = Duolingo(username="foo", jwt="bar")
```

### get_abbreviation_of(lang: str)

Get abbreviation of a language.

```pycon
>>> lingo.get_abbreviation_of(lang='High Valyrian')
'hv'
```

### get_daily_xp_progress()

Return the topics learned but not golden by a user in a language.

```pycon
>>> lingo.get_daily_xp_progress()
{'xp_goal': ..., ..., 'xp_today': ...}
```

### get_friends()

Get user’s friends.

```pycon
>>> lingo.get_friends()
[{'username': '...', ..., 'displayName': '...'}...]
```

### get_golden_topics(abbr)

Return the topics mastered (“golden”) by a user in a language.

```pycon
>>> lingo.get_golden_topics(abbr=MY_FAVORITE_LANGUAGE_ABBREVIATION)
[...]
```

### get_known_topics(abbr)

Return the topics learned by a user in a language.

```pycon
>>> lingo.get_known_topics(abbr=MY_FAVORITE_LANGUAGE_ABBREVIATION)
['...', '...', ...]
```

### get_known_words(abbr)

Get a list of all words learned by user in a language.

```pycon
>>> lingo.get_known_words(abbr=MY_FAVORITE_LANGUAGE_ABBREVIATION)
[...]
```

### get_language_details(language)

Get user’s status about a language.

```pycon
>>> lingo.get_language_details(language='High Valyrian')
{'streak': ..., ..., 'to_next_level': ...}
```

### get_language_from_abbr(abbr: str)

Get language full name from abbreviation.

```pycon
>>> lingo.get_language_from_abbr(abbr='hv')
'High Valyrian'
```

### get_language_progress(abbr: str)

Get information about user’s progression in a language.

```pycon
>>> lingo.get_language_progress(abbr='hv')
{'language_string': 'High Valyrian', ..., 'level': ...}
```

### get_languages(abbreviations: bool = False)

Get practiced languages.

```pycon
>>> lingo.get_languages()
['...', 'High Valyrian'...]
```

```pycon
>>> lingo.get_languages(abbreviations=True)
['...', 'hv'...]
```

### get_reviewable_topics(abbr)

Return the topics learned but not golden by a user in a language.

```pycon
>>> lingo.get_reviewable_topics(abbr=MY_FAVORITE_LANGUAGE_ABBREVIATION)
['...', '...', ...]
```

### get_streak_info()

Get user’s streak information.

```pycon
>>> lingo.get_streak_info()
{'daily_goal': ..., ..., 'streak_extended_today': ...}
```

### get_unknown_topics(abbr)

Return the topics learned by a user in a language.

```pycon
>>> lingo.get_unknown_topics(abbr=MY_FAVORITE_LANGUAGE_ABBREVIATION)
['...', '...', ...]
```

### get_user_info()

Get user’s information.

```pycon
>>> lingo.get_user_info()
{'username': ..., ..., 'ui_language': ...}
```

### get_vocabulary(abbr: str, source_language_abbr=None)

Get overview of user’s vocabulary in a language.

```pycon
>>> lingo.get_vocabulary(abbr=MY_FAVORITE_LANGUAGE_ABBREVIATION)
[{'text': ..., ...}, ...]
```
