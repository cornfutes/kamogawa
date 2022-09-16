# TODO

  [ ] display warning if regex 
  [ ] add UI codes 
- [ ] make top global search actually do its job
  - [x] make search the /search page responds given query param. don't worry 
        about the interactive UI search box in nav bar. i will hook it up.
  - [ ] for each entity retrieved from the network, add it to the DB 
  - [x] in DB, start with just two colums: id (string)  kind (string)
  - [x] figure out how to do full text search with postgres
  - [x] update /search?q=foobar web page route handler
      - [x] extract query param 'q'
      - [x] call function to get list of results for 'q', render results to html 
      - [x] DAVID: replace the fake with real one. see /Users/louis/workspace/kamogawa/handler/search.go
      search for 'foobar' in postgres
  - [x] don't worry about pagination or length. just render a html list
- [x] add page to export data ( SQLLite. DuckDB comin soon)
- [ ] export CSV, SQLite, DuckDb ( future )
- [ ] refactor the ad-hoc retrieval of user from DB on demand.
      for demo, we can load it on server start to speed up time.
- [ ] demo: use in memory data for all API data requests for test user.
- [x] boilerplate the security and vulnerability disclosure page 
- [ ] add a Google Survey / Forms to request access
- [x] fix: "you have no projects" when missing scope 

Low priority
- [x] minor bug. when logged in, /login route give immpression you are logged out

auth 
- [x] route to login page if people's session are invalid or expired. Currently do not differentiate.

production
- [x] embed assets so we can deploy verified
- [x] productionize
  - [x] hook domain name up with app
  - [x] automated deployment or atleast self contained file
        infeasible with Gitlab 
  - [x] document process, where to find things in the CSP 
  - [x] refactor environmental assumptions to ENV 
- [ ] load test
  - [ ] perform load test 
  - [ ] consider having fake static website, cant login / not dynamic, in case of HN hug of death
  - [ ] then host the real app on app.ListVMs.com or something
- [ ] logs and monitoring 
- [x] parameterize dev environment flag properly ... its showing in prod 
- [x] Minify CSS file. Goal: reduce bundle size
- [x] Prune CSS. Goal: reduce bundle size.
  - [x] move inline styles to class/ID selector and style.css
  - [x] refactor redundant / overlapping CSS rules
- [x] Enable GZIP on files. Goal: reduce bundle size.
  - [x] works on HTML file
  - [x] works on media? hit or miss. TXT, but not CSS. why?
    - [x] fixed gzip not being applied to  asset.Config registered HTML handlers
    - [x] audit all assets get gzipped ( doesn't apply to images )
    - [ ] not caching
- [x] fix resources not being cached 

quality
- [x] be more robust about missing config 
  - [ ] nice to have: use proper ENV manager and consolidate scattered logic
- [x] handle edge cases
  - [x] handle edge case around partial and no oauth2 scopes
  - [x] display the delta between needed and missing scopes
    - [x] can provide more granulartiy
  - [x] handle scenario where GCP not yet authorized.
- [x] fix "zero state" for q=''
- [x] display error if q.length < 3

qa 
- [ ] ensure all pages covered
  - [x] show release notes
- [x] add HN news account settings-like page to /account page. 
  - [x] show account tier
  - [ ]  add account tier page

polish 
- [x] add some CSS polish
- [x] mark / indicate the current page on the left nav
- [ ] make header sticky, scroll beneath 
- [ ] add breadcrumbs
- [x] instruct user to login when clicking global search
  - [x] improve UI component beyond an alert. 
- [x] make sure font size is normalized across pages
- [x] obsolete. move system status, release notes, demo to left size
- [x] fix top padding of footer 
- [x] maintain client state across across page navs utilize url param
- [x] obsolete. perform a Lighthouse metrics test and improve
- [x] make search and 'launch search' CTA keyboard accessible

meta 
- [x] set up email. team@ListVMs.com

easter egg 
- [x] log stuff to console 


