Usage
=====

This fork of the Duolingo library is not published to PyPi: you can't install it using pip, but you can copy it into your project if you also copy in the accompanying LICENSE.md.

The Duolingo API only allows JWT login: you can't log in using your username and password. To acquire your JWT from Duolingo, open up the network console and visit the Duolingo site. Be sure that you're logged in. Find a request from your browser to the Duolingo API: you should see a jwt_token cookie on that request.

.. automodule:: duolingo.duolingo
   :members:
   :undoc-members:
   :show-inheritance:

