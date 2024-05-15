"""Unofficial API for duolingo.com"""
import re
import json
import random
import time
from datetime import datetime, timedelta
from json import JSONDecodeError

import deprecated
import requests

__version__ = "0.5.4"
__author__ = "Kartik Talwar"
__email__ = "hi@kartikt.com"
__url__ = "https://github.com/KartikTalwar/duolingo"


class Struct:
    def __init__(self, **entries):
        self.__dict__.update(entries)


class DuolingoException(Exception):
    pass


class AlreadyHaveStoreItemException(DuolingoException):
    pass


class InsufficientFundsException(DuolingoException):
    pass


class CaptchaException(DuolingoException):
    pass


class OtherUserException(DuolingoException):
    """
    This exception is raised when set_username() has been called to get info on another user, but a method has then
    been used which cannot give data on that new user.
    """
    pass


class Duolingo(object):
    USER_AGENT = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 " \
                 "Safari/537.36"

    def __init__(self, username, password=None, *, jwt=None, session_file=None):
        """
        :param username: Username to use for duolingo
        :param password: Password to authenticate as user.
        :param jwt: Duolingo login token. Will be checked and used if it is valid.
        :param session_file: File path to a file that the session token can be stored in, to save repeated login
        requests.
        """
        self.username = username
        self._original_username = username
        self.password = password
        self.session_file = session_file
        self.session = requests.Session()
        self.leader_data = None
        self.jwt = jwt

        if password and (not jwt and not session_file):
            raise DuolingoException(
                "The Duolingo API no longer supports password login. You must provide a JWT or a session file."
            )

        self.session.cookies.set_cookie(
            requests.cookies.create_cookie("jwt_token", self.jwt)
        )

        if password or jwt or session_file:
            self._login()
        else:
            raise DuolingoException("Password, jwt, or session_file must be specified in order to authenticate.")

        self.user_data = Struct(**self._get_data())
        self.voice_url_dict = None

    def _make_req(self, url, data=None):
        headers = {}
        if self.jwt is not None:
            headers['Authorization'] = 'Bearer ' + self.jwt
        headers['User-Agent'] = self.USER_AGENT
        req = requests.Request('POST' if data else 'GET',
                               url,
                               json=data,
                               headers=headers,
                               cookies=self.session.cookies)
        prepped = req.prepare()
        resp = self.session.send(prepped)
        if resp.status_code == 403 and resp.json().get("blockScript") is not None:
            raise CaptchaException(
                "Request to URL: {}, using user agent {}, was blocked, and requested a captcha to be solved. "
                "Try changing the user agent and logging in again.".format(url, self.USER_AGENT)
            )
        return resp

    def _login(self):
        """
        Authenticate through ``https://www.duolingo.com/login``.
        """
        if self.jwt is None:
            self._load_session_from_file()
        if self._check_login():
            return True
        self.jwt = None

        login_url = "https://www.duolingo.com/login"
        data = {"login": self.username, "password": self.password}
        request = self._make_req(login_url, data)
        attempt = request.json()

        if "failure" not in attempt:
            self.jwt = request.headers['jwt']
            self._save_session_to_file()
            return True

        raise DuolingoException("Login failed")

    def _load_session_from_file(self):
        if self.session_file is None:
            return
        try:
            with open(self.session_file, "r") as f:
                self.jwt = json.load(f).get("jwt_session")
        except (OSError, JSONDecodeError):
            return

    def _save_session_to_file(self):
        if self.session_file is not None:
            with open(self.session_file, "w") as f:
                json.dump({"jwt_session": self.jwt}, f)

    def _check_login(self):
        resp = self._make_req(self.get_user_url())
        return resp.status_code == 200

    def get_user_url_by_id(self, fields=None):
        if fields is None:
            fields = []
        url = 'https://www.duolingo.com/2017-06-30/users/{}'.format(self.user_data.id)
        fields_params = requests.utils.requote_uri(','.join(fields))
        if fields_params:
            url += '?fields={}'.format(fields_params)
        return url

    def get_user_url(self):
        return "https://duolingo.com/users/%s" % self.username

    def set_username(self, username):
        self.username = username
        self.user_data = Struct(**self._get_data())


    def buy_item(self, item_name, abbr):
        url = 'https://www.duolingo.com/2017-06-30/users/{}/shop-items'
        url = url.format(self.user_data.id)

        data = {'itemName': item_name, 'learningLanguage': abbr}
        request = self._make_req(url, data)

        """
        status code '200' indicates that the item was purchased
        returns a text like: {"streak_freeze":"2017-01-10 02:39:59.594327"}
        """

        if request.status_code == 400:
            resp_json = request.json()
            if resp_json.get("error") == "ALREADY_HAVE_STORE_ITEM":
                raise AlreadyHaveStoreItemException("Already equipped with {}.".format(item_name))
            if resp_json.get("error") == "INSUFFICIENT_FUNDS":
                raise InsufficientFundsException("Insufficient funds to purchase {}.".format(item_name))
            raise DuolingoException(
                "Duolingo returned an unknown error while trying to purchase {}: {}".format(
                    item_name, resp_json.get("error")
                )
            )
        if not request.ok:
            # any other error:
            raise DuolingoException("Not possible to buy item.")

    def buy_streak_freeze(self):
        """
        figure out the users current learning language
        use this one as parameter for the shop
        """
        lang = self.get_abbreviation_of(self.get_user_info()['learning_language_string'])
        if lang is None:
            raise DuolingoException('No learning language found')
        try:
            self.buy_item('streak_freeze', lang)
            return True
        except AlreadyHaveStoreItemException:
            return False

    def buy_weekend_amulet(self):
        """
        figure out the users current learning language
        use this one as parameter for the shop
        """
        lang = self.get_abbreviation_of(self.get_user_info()['learning_language_string'])
        if lang is None:
            raise DuolingoException('No learning language found')
        try:
            self.buy_item('weekend_amulet', lang)
            return True
        except AlreadyHaveStoreItemException:
            return False

    def _switch_language(self, lang):
        """
        Change the learned language with
        ``https://www.duolingo.com/switch_language``.

        :param lang: Wanted language abbreviation (example: ``'fr'``)
        :type lang: str
        """
        data = {"learning_language": lang}
        url = "https://www.duolingo.com/switch_language"
        request = self._make_req(url, data)

        try:
            parse = request.json()['tracking_properties']
            if parse['learning_language'] == lang:
                self.user_data = Struct(**self._get_data())
        except ValueError:
            raise DuolingoException('Failed to switch language')

    def get_data_by_user_id(self, fields=None):
        """
        Get user's data from ``https://www.duolingo.com/2017-06-30/users/<user_id>``.
        """
        if fields is None:
            fields = []
        get = self._make_req(self.get_user_url_by_id(fields))
        if get.status_code == 404:
            raise DuolingoException('User not found')
        else:
            return get.json()

    def _get_data(self):
        """
        Get user's data from ``https://www.duolingo.com/users/<username>``.
        """
        get = self._make_req(self.get_user_url())
        if get.status_code == 404:
            raise Exception('User not found')
        else:
            return get.json()

    @staticmethod
    def _make_dict(keys, array):
        data = {}

        for key in keys:
            if type(array) == dict:
                data[key] = array[key]
            else:
                data[key] = getattr(array, key, None)

        return data

    @staticmethod
    def _compute_dependency_order_func(skills):
        # Create dictionary:
        skills_dict = {}
        for skill in skills:
            skills_dict[skill['name']] = skill
        # Get ordinal for all dependencies
        for skill in skills:
            skill['dependency_order'] = Duolingo._get_skill_ordinal(skills_dict, skill, [])

    @staticmethod
    def _get_skill_ordinal(skills_dict, skill, breadcrumbs):
        # If name is already in breadcrumbs, we've found a loop
        if skill['name'] in breadcrumbs:
            raise DuolingoException("Loop encountered: {}".format(breadcrumbs + [skill['name']]))
        # If order already set for this skill, return it
        if "dependency_order" in skill:
            return skill["dependency_order"]
        # If no dependencies, set order on this skill to 1
        if not skill['dependencies_name']:
            skill['dependency_order'] = 1
            return 1
        # Calculate order based on order of dependencies
        new_breadcrumbs = breadcrumbs + [skill['name']]
        order = 1 + max(
            [
                Duolingo._get_skill_ordinal(
                    skills_dict,
                    skills_dict[name],
                    new_breadcrumbs
                )
                for name in skill['dependencies_name']
            ]
        )
        skill["dependency_order"] = order
        return order

    def get_settings(self):
        """Get user settings."""
        keys = ['notify_comment', 'deactivated', 'is_follower_by',
                'is_following']

        return self._make_dict(keys, self.user_data)

    def get_languages(self, abbreviations=False):
        """
        Get practiced languages.

        :param abbreviations: Get language as abbreviation or not
        :type abbreviations: bool
        :return: List of languages
        :rtype: list of str
        """
        data = []

        for lang in self.user_data.languages:
            if lang['learning']:
                if abbreviations:
                    data.append(lang['language'])
                else:
                    data.append(lang['language_string'])
        return data

    def get_language_from_abbr(self, abbr):
        """Get language full name from abbreviation."""
        for language in self.user_data.languages:
            if language['language'] == abbr:
                return language['language_string']
        return None

    def get_abbreviation_of(self, name):
        """Get abbreviation of a language."""
        for language in self.user_data.languages:
            if language['language_string'].lower() == name.lower():
                return language['language']
        return None

    def get_language_details(self, language):
        """Get user's status about a language."""
        for lang in self.user_data.languages:
            if language == lang['language_string']:
                return lang

        return {}

    def get_user_info(self):
        """Get user's informations."""
        fields = ['username', 'bio', 'id', 'num_following', 'cohort',
                  'language_data', 'num_followers', 'learning_language_string',
                  'created', 'contribution_points', 'gplus_id', 'twitter_id',
                  'admin', 'invites_left', 'location', 'fullname', 'avatar',
                  'ui_language']

        return self._make_dict(fields, self.user_data)

    def get_streak_info(self):
        """Get user's streak informations."""
        fields = ['daily_goal', 'site_streak', 'streak_extended_today']
        return self._make_dict(fields, self.user_data)

    def _is_current_language(self, abbr):
        """Get if user is learning a language."""
        return abbr in self.user_data.language_data.keys()

    def get_calendar(self, language_abbr=None):
        """Get user's last actions."""
        if language_abbr:
            if not self._is_current_language(language_abbr):
                self._switch_language(language_abbr)
            return self.user_data.language_data[language_abbr]['calendar']
        else:
            return self.user_data.calendar

    def get_language_progress(self, lang):
        """Get informations about user's progression in a language."""
        if not self._is_current_language(lang):
            self._switch_language(lang)

        fields = ['streak', 'language_string', 'level_progress',
                  'num_skills_learned', 'level_percent', 'level_points',
                  'points_rank', 'next_level', 'level_left', 'language',
                  'points', 'fluency_score', 'level']

        return self._make_dict(fields, self.user_data.language_data[lang])
    def get_following(self):    #todo write tests for this function
        return self._make_req(f"https://www.duolingo.com/2017-06-30/friends/users/{self.user_data.id}/following").json()
    def get_friends(self):
        following = self.get_following()
        friends = []
        for follower in following["following"]["users"]:
            if follower["isFollowing"]:
                friends.append(
                    {"username": follower["username"], "id": follower["userId"],"points": follower["totalXp"],"avatar": follower["picture"],"displayName": follower["displayName"]})

        return friends

    def get_known_words(self, lang):
        """Get a list of all words learned by user in a language."""
        words = []
        for topic in self.user_data.language_data[lang]['skills']:
            if topic['learned']:
                words += topic['words']
        return list(set(words))

    def get_learned_skills(self, lang):
        """
        Return the learned skill objects sorted by the order they were learned
        in.
        """
        skills = [
            skill for skill in self.user_data.language_data[lang]['skills']
        ]

        self._compute_dependency_order_func(skills)

        return [skill for skill in
                sorted(skills, key=lambda skill: skill['dependency_order'])
                if skill['learned']]

    def get_known_topics(self, lang):
        """Return the topics learned by a user in a language."""
        return [topic['title']
                for topic in self.user_data.language_data[lang]['skills']
                if topic['learned']]

    def get_unknown_topics(self, lang):
        """Return the topics remaining to learn by a user in a language."""
        return [topic['title']
                for topic in self.user_data.language_data[lang]['skills']
                if not topic['learned']]

    def get_golden_topics(self, lang):
        """Return the topics mastered ("golden") by a user in a language."""
        return [topic['title']
                for topic in self.user_data.language_data[lang]['skills']
                if topic['learned'] and topic['strength'] == 1.0]

    def get_reviewable_topics(self, lang):
        """Return the topics learned but not golden by a user in a language."""
        return [topic['title']
                for topic in self.user_data.language_data[lang]['skills']
                if topic['learned'] and topic['strength'] < 1.0]

    def get_translations(self, words, source=None, target=None):
        """
        Get words' translations from
        ``https://d2.duolingo.com/api/1/dictionary/hints/<source>/<target>?tokens=``<words>``

        :param words: A single word or a list
        :type: str or list of str
        :param source: Source language as abbreviation
        :type source: str
        :param target: Destination language as abbreviation
        :type target: str
        :return: Dict with words as keys and translations as values
        """
        if not source:
            source = self.user_data.ui_language
        if not target:
            target = list(self.user_data.language_data.keys())[0]

        list_segments = self._segment_translations_list(words)
        results = dict()
        for segment in list_segments:
            results = {**results, **self._get_raw_translations(segment, source, target)}
        return results

    def _segment_translations_list(self, words):
        # These seem to be the length limits before Duolingo's API rejects the request
        word_count_limit = 2000
        word_json_limit = 12800

        # Handy internal function
        def is_word_list_valid(word_list):
            return (
                    len(word_list) < word_count_limit
                    and len(json.dumps(word_list)) < word_json_limit
            )

        # Fast return for simple lists
        if is_word_list_valid(words):
            return [words]
        # Start building segments until they trip the limits
        segments = []
        segment = []
        for word in words:
            if not is_word_list_valid(segment + [word]):
                segments.append(segment)
                segment = []
            segment.append(word)
        segments.append(segment)
        return segments

    def _get_raw_translations(self, words, target, source):
        word_parameter = json.dumps(words, separators=(',', ':'))
        url = "https://d2.duolingo.com/api/1/dictionary/hints/{}/{}?tokens={}" \
            .format(target, source, word_parameter)

        request = self.session.get(url)
        try:
            return request.json()
        except ValueError:
            raise DuolingoException('Could not get translations')

    def get_vocabulary(self, language_abbr, source_language_abbr=None):
        """Get overview of user's vocabulary in a language.

        :param language_abbr: Language abbreviation of learning language
        :param source_language_abbr: Language abbreviation of source language (default: user's UI language)
        """
        if self.username != self._original_username:
            raise OtherUserException("Vocab cannot be listed when the user has been switched.")

        if language_abbr and not self._is_current_language(language_abbr):
            self._switch_language(language_abbr)

        current_courses = self.get_data_by_user_id()["currentCourse"]["pathSectioned"]
        progressed_skills = []
        for section in current_courses:
            completedUnits = section["completedUnits"]
            units = section["units"]
            for i in range(completedUnits):
                unit = units[i]
                levels = unit["levels"]
                for l in levels:
                    level_type = l["type"]
                    # unit review doesnt contain new words
                    if level_type in ["chest", "unit_review"]:
                        continue
                    pathLevelClientData = l["pathLevelClientData"]
                    finishedSessions = l["finishedSessions"]
                    if "skillId" in pathLevelClientData:
                        skillId = pathLevelClientData["skillId"]
                        new_obj = {
                            "finishedLevels": 1,
                            "finishedSessions": finishedSessions,
                            "skillId": {
                                "id": skillId
                            }
                        }
                        progressed_skills.append(new_obj)
                    elif "skillIds" in pathLevelClientData:
                        skillIds = pathLevelClientData["skillIds"]
                        for skillId in skillIds:
                            new_obj = {
                                "finishedLevels": 1,
                                "finishedSessions": finishedSessions,
                                "skillId": {
                                    "id": skillId
                                }
                            }
                            progressed_skills.append(new_obj)

        # updated URL, default language to be english,
        current_index = 0
        data = []
        if source_language_abbr is None:
            source_language_abbr = self.user_data.ui_language
        while True:
            overview_url = f"https://www.duolingo.com/2017-06-30/users/{self.user_data.id}/courses/{language_abbr}/{source_language_abbr}/learned-lexemes?sortBy=ALPHABETICAL&startIndex={current_index}"
            overview_request = self._make_req(overview_url, data={
                "lastTotalLexemeCount": 0,
                "progressedSkills": progressed_skills
            })
            overview = overview_request.json()
            learnedLexemes = overview['learnedLexemes']
            data.extend(learnedLexemes)
            pagination = overview['pagination']
            totalLexemes = pagination['totalLexemes']
            if len(data) >= totalLexemes:
                break
            # its not my database so i am being wasteful :)
            nextStartIndex = pagination['nextStartIndex']
            current_index = nextStartIndex
        return data

    _cloudfront_server_url = None
    _homepage_text = None

    @property
    def _homepage(self):
        if self._homepage_text:
            return self._homepage_text
        homepage_url = "https://www.duolingo.com"
        request = self._make_req(homepage_url)
        self._homepage_text = request.text
        return self._homepage

    @property
    def _cloudfront_server(self):
        if self._cloudfront_server_url:
            return self._cloudfront_server_url

        server_list = re.search(r'//.+\.cloudfront\.net', self._homepage)
        self._cloudfront_server_url = "https:{}".format(server_list.group(0))

        return self._cloudfront_server_url

    _tts_voices = None

    def _process_tts_voices(self):
        voices_js = re.search(r'duo\.tts_multi_voices = {.+};',
                              self._homepage).group(0)

        voices = voices_js[voices_js.find("{"):voices_js.find("}") + 1]
        self._tts_voices = json.loads(voices)

    def get_language_voices(self, language_abbr=None):
        if not language_abbr:
            language_abbr = list(self.user_data.language_data.keys())[0]
        voices = []
        if not self._tts_voices:
            self._process_tts_voices()
        for voice in self._tts_voices[language_abbr]:
            if voice == language_abbr:
                voices.append('default')
            else:
                voices.append(voice.replace('{}/'.format(language_abbr), ''))
        return voices

    def get_audio_url(self, word, language_abbr=None):
        """
        Get the audio URL for a given word in the user's vocabulary.
        :param word: The word to retrieve the audio URL for.
        :param language_abbr: (Optional) The language abbreviation for the vocabulary.
        :return: The audio URL for the word.
        :raises DuolingoException: If the word is not found in the user's vocabulary.
        """
        if not language_abbr:
            language_abbr = self.get_languages(True)[0]
        for vocabulary in self.get_vocabulary(language_abbr):
            if vocabulary['text'] == word:

                return vocabulary['audioURL']
        raise DuolingoException("Word not found in vocabulary of user. Please make sure to use the correct language abbreviation. If language_abbr=None the first possible language is used")



    def _add_token_list_to_voice_url_dict(self, lang_abbr, token_list):
        for token in token_list:
            if isinstance(token, list):
                self._add_token_list_to_voice_url_dict(lang_abbr, token)
            if isinstance(token, dict) and token.get("tts") and token.get("value"):
                self._add_to_voice_url_dict(lang_abbr, token['value'], token['tts'])

    def _add_to_voice_url_dict(self, lang_abbr, word, url):
        word = word.lower()
        if word not in self.voice_url_dict[lang_abbr]:
            self.voice_url_dict[lang_abbr][word] = set()
        self.voice_url_dict[lang_abbr][word].add(url)

    def get_related_words(self, word, language_abbr=None):
        overview = self.get_vocabulary(language_abbr)

        for word_data in overview['vocab_overview']:
            if word_data['normalized_string'] == word.lower():
                related_lexemes = word_data['related_lexemes']
                return [w for w in overview['vocab_overview']
                        if w['lexeme_id'] in related_lexemes]

    def get_word_definition_by_id(self, lexeme_id):
        """
        Get the dictionary entry from
        ``https://www.duolingo.com/api/1/dictionary_page?lexeme_id=``<lexeme_id>``

        :param lexeme_id: Identifier of the word
        :type: str
        :return: The dictionary entry for the given word
        """
        url = "https://www.duolingo.com/api/1/dictionary_page?lexeme_id=%s" % lexeme_id

        request = self.session.get(url)

        try:
            return request.json()
        except:
            raise Exception('Could not get word definition')

    def get_daily_xp_progress(self):
        daily_progress = self.get_data_by_user_id(["xpGoal", "xpGains", "streakData"])

        if not daily_progress:
            raise DuolingoException(
                "Could not get daily XP progress for user \"{}\". Are you logged in as that user?".format(self.username)
            )

        # xpGains lists the lessons completed on the last day where lessons were done.
        # We use the streakData.updatedTimestamp to get the last "midnight", and get lessons after that.
        reported_timestamp = daily_progress['streakData']['updatedTimestamp']
        reported_midnight = datetime.fromtimestamp(reported_timestamp)
        midnight = datetime.fromordinal(datetime.today().date().toordinal())

        # Sometimes the update is marked into the future. When this is the case
        # we fall back on the system time for midnight.
        time_discrepancy = min(midnight - reported_midnight, timedelta(0))
        update_cutoff = round((reported_midnight + time_discrepancy).timestamp())

        lessons = [lesson for lesson in daily_progress['xpGains'] if
                   lesson['time'] > update_cutoff]

        return {
            "xp_goal": daily_progress['xpGoal'],
            "lessons_today": lessons,
            "xp_today": sum(x['xp'] for x in lessons)
        }


attrs = [
    'settings', 'languages', 'user_info', 'streak_info',
    'calendar', 'language_progress', 'known_words',
    'learned_skills', 'known_topics', 'vocabulary',"friends",
]

for attr in attrs:
    getter = getattr(Duolingo, "get_" + attr)
    prop = property(getter)
    setattr(Duolingo, attr, prop)
