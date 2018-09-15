#!/usr/bin/env python

"""Simple redirection handler to Gubernator's GitHub service."""

import webapp2

class GitHubRedirect(webapp2.RequestHandler):
  def get(self):
    self.redirect("https://github-dot-knative-tests.appspot.com" + self.request.path_qs)

app = webapp2.WSGIApplication([(r'/.*', GitHubRedirect),], debug=True, config={})
