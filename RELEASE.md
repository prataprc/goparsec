Build and Release Check list
============================

* Ensure coverage to be > 90%. Unit tests are not silver bullets to catch
  all bugs, but they can help maintain decent code-quality against
  future commits.
* Set up with https://coveralls.io for coverage analysis.
* Set up with https://travis-ci.org/ for continuous integration and
  enable daily CRON job.
* Spell check all `.md` files.

Repository structure
--------------------

* README.md file that includes -
  * Badges.
  * Bullet point of why and what of this package.
  * Quicklinks.
  * Short descriptions and details about this package.
  * Panic and recovery.
  * Examples.
  * External references.
  * How to contribute.
* RELEASE.md checklist.
* LICENSE under which this package is released on github.
* AUTHORS list of contributing authors, copyright is collectively
  held by authors.
* PITCHME files, if any, for package presentation.

Badges
------

* CI badge
* Coverage badge
* Godoc reference.
* Issue stats badge for response time.
  http://issuestats.com/github/prataprc/goparsec
* Sourcegraph for "used by projects" badge
  https://sourcegraph.com/github.com/prataprc/goparsec/-/badge.svg
* Report card.
  https://goreportcard.com/report/github.com/prataprc/goparsec
