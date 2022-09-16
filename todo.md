# Week of 20220911 to 20220920

[x] YC application
[x] Sleep 
[ ] Automated Screenshot Testing ( why? faster refactoring )
[ ] cache 
  [x] fix GCE VMs page not caching case where project doesn't have GCE API enabled 
  [ ] fix cache behavior for GAE 
[ ] improve SERP ( search results page )
  [x] show total results count and duration 
  [ ] paginate. show 10 results per page.
[ ] search
  [ ] evaluate issue with DB having duplicate 
  [ ] correctly implement multi-word search or go back to full-text search
  [ ] support searching last modified date 
[ ] DB 
  [ ] batch upsert
  [ ] understand how GORM upsert works.. works different in different environment 
  [ ] when people revoke, delete old projects and data
[ ] caching model
  [ ] by project? by gmail?
[x] try calling api with revoked refresh token and see what happens 
[x] prevent man in the middle attack on cookie
[ ] refactor
  [ ] all places where we parse response success and error response can be refactored.
    check responseStatus code and branch respectively