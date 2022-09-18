# Week of 20220911 to 20220920

[x] YC application
[x] Sleep 
[x] Automated Screenshot Testing ( why? faster refactoring )
[x] cache 
  [x] fix GCE VMs page not caching case where project doesn't have GCE API enabled 
  [x] fix cache behavior for GAE 
  [x] cache bit on Project that no instances 
  [ ] TBD on using ListServices API
  App Engine states 
    [ ] zero state 
[ ] improve SERP ( search results page )
  [x] show total results count and duration 
  [ ] paginate. show 10 results per page.
[ ] always sort resorts to be predictable so screenshot test improves
[ ] search
  [x] evaluate issue with DB having duplicate 
  [ ] correctly implement multi-word search or go back to full-text search
  [ ] support searching last modified date 
[ ] DB 
  [ ] batch upsert
  [ ] understand how GORM upsert works.. works different in different environment 
  [ ] when people revoke, delete old projects and data
[x] eng prod
  [x] e2e + visual test setup
  [x] obsolete: gitlab ci/cd 
    [x] added config
    [x] obsolete: adding GCP service account json to Gitlab CI/Cd ( not showing up in settings )
  [x] add Github CI/CD
  [x] add makefile 
  [x] read local ENV variable for $THEME, so we dont have to manually change Docker Compose and avoid staging for git each time 
[x] caching model
  [x] by project? by gmail?
[x] try calling api with revoked refresh token and see what happens 
[x] prevent man in the middle attack on cookie
[ ] refactor
  [ ] all places where we parse response success and error response can be refactored.
    check responseStatus code and branch respectively
[x] derisk universality of Auth user ( only one globally refresh token )
  [x] local environment can break prod
    [x] add ability to disconnect google account, without revoking the refresh token
[ ] product 
  [ ] project details page
    [ ] apis page 
    [ ] fixed increment cache count when getting from cached 
[ ] add 404, 500 page 
[ ] look into this gin warning about proxies trusting all
[ ] look into removing plans.html and tbd.html as it seems like deadcode
[ ] upgrade to Go 1.19 from our host machine
[ ] bug 
  [ ] only landing page, GCP is selected. same ith login page. we need to add more bits
  [ ] route from search to account page
[ ] $onic 
  [ ] handle SPA navigation when left nav doesn't exist yet.
[ ] control panel
  [ ] hard refresh button 
