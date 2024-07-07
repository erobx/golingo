import os
import unittest
import doctest

import duolingo

MY_DUOLINGO_JWT = os.environ['MY_DUOLINGO_JWT']
MY_FAVORITE_LANGUAGE_ABBREVIATION = os.environ['MY_FAVORITE_LANGUAGE_ABBREVIATION']


def load_tests(_, tests, __):
    tests.addTests(
        doctest.DocTestSuite(
            duolingo.duolingo,
            globs={
                'duolingo': duolingo,
                'MY_DUOLINGO_JWT': MY_DUOLINGO_JWT,
                'MY_FAVORITE_LANGUAGE_ABBREVIATION': MY_FAVORITE_LANGUAGE_ABBREVIATION,
                'lingo': duolingo.Duolingo(jwt=MY_DUOLINGO_JWT),
            },
            optionflags=doctest.NORMALIZE_WHITESPACE | doctest.ELLIPSIS,
        )
    )

    return tests


if __name__ == '__main__':
    unittest.main()
