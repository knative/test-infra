# Copyright 2022 The Knative Authors

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/usr/bin/env python
"""Mostly simple redirection app for knative.dev subdomains."""

import os
import re
import webapp2


class RedirectHandler(webapp2.RequestHandler):
  """Redirects URLs."""

  def get(self):
    original = self.request.uri
    domain = 'knative.dev'
    redirect_path = os.getenv('REDIR_TO', 'https://github.com/knative')
    where = original.find(domain)
    if where != -1:
      # Increment to beginning of domain, then end of domain,
      #  then one for the slash; okay if there is no slash,
      #  python just gives empty string for slices past the end
      extra_path = original[where+len(domain)+1:]
    else:  # Probably only use this section when hitting app URL directly
      match = re.match(r'https?://[^/]+/(.*)', original)
      if not match:
        # Failed to figure out the URL, redirect to the base page
        self.redirect(redirect_path)
        return
      extra_path = match.groups()[0] or ''
    # extra_path should never have a leading slash
    if extra_path and redirect_path[-1] != '/':
      self.redirect(redirect_path + '/' + extra_path)
    else:
      self.redirect(redirect_path + extra_path)


app = webapp2.WSGIApplication([(r'/.*', RedirectHandler),],
                              debug=True,
                              config={})
