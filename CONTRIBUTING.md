# Contributing

## Principles

This is a library used by a small number of people, which gets touched infrequently only when Duolingo makes a change to their service which breaks the existing code.

The following principles help keep the code as easy to maintain and quick to review as possible, with minimal spin-up time despite looking at it infrequently:
  * Keep the functionality simple: whenever possible let the calling code do the complicated stuff.
  * Don't expose bare responses. Instead, return specifically typed dictionaries which we can hope to fill out consistently even when Duolingo makes changes.
  * Return the minimum needed: more values returned means more that can be broken.
  * Add a docstring which demonstrates how to use any public method; we use autodoc to convert these into helpful documentation and docstring tests to ensure that our documentation is accurate.

## Running the Tests

You'll need to define two environment variables to run the project's integration tests:
  * `MY_DUOLINGO_JWT`: A Duolingo JWT token.
  * `MY_FAVORITE_LANGUAGE_ABBREVIATION`: A language abbreviation for a language that you've got some lessons/words learned in. For example `es` for Spanish.
  * Additionally you must add the course **High Valyrian** in Duolingo in order for the integration tests to work, although you don't need to do any study in the this course.

You can then run: `$ python -m unittest discover -s . -p tests.py`

## Extra Instructions

To rebuild the docs after making any public API changes:
  * Install the documentation requirements: `$ pip install sphinx sphinx-markdown-builder`
  * Build the docs: `$ ./docs/make.bat markdown`, or perhaps `$ (cd docs; make markdown)` for Mac/Linux?
  * Copy the docs into their needed place: `$ cp docs/_build/markdown/index.md ./USAGE.md`

