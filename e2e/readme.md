# Overview

This directory houses "End-to-End" test automation.

When ran, tests performs the following:

  - opens a web browser ( headless by default )
  - looks for stuff ( including fuzzy matching text )
  - can assert on the text of stuff
  - clicks stuff
  - grabs screenshot ( temporarily stored in 'test-results' directory )
  - performs of current screenshot DIFF against golden snapshot

## Installation 

```
$PATH_TO_REPO/e2e $ npm i 
$PATH_TO_REPO/e2e $ npm i playwright --g
```
## Run 

Run test in headless browser mode

```
$PATH_TO_REPO/playwright $ npx playwright test
```

Run test in debug mode ( which launches a step-by-step visual debugger )

```
$PATH_TO_REPO/playwright $ npx playwright test --debug
```
