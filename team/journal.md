

** 2022 Sep 16 
  - Added e2e test automation w/ Playwright. For now, just-ad hoc against prod.
    - took less than half a day
    - probably comparable to the product of Japan's hot test automation startup.
      - Goes to show lack of hard tech in Japan
    - Playwright debug mode is awesome. Can interactiely step through and debug.
    - Masks: can mask out UI elements whether it is for confidentiality reasons
      or avoiding flakiness ( for example, when the specific values in the UI 
        between run )s
    - note: Playwright developed by Microsoft
      - ironically, there are some issues on Windows machines
        - in WSL ( unix simulator ), compatibility issues with x11 preventing 
          debug mode of Playwright
  - Abandoned Gitlab for Github
    - Gitlab experience was mediocre thus far. 
    - Final nail in the coffin was CI/CD bug or shadow ban
    - Product-wise, Gitlab only has 400 free CI/CD minutes a month versus
      2,000-4,000 on Github. Assuming 2 CI/CD minutes per build, and a
      conservative 500+ commits a month, Gitlab just doesn't cut it.
    - No CI/CD integration with GCP
    - It's ironic because Gitlab is pushing really hard on CI/CD but it 
      feels not very well thought out.
  - Configured CI/CD ( deploy to GCR )
    - Could not configure Gitlab
      - originally, account was missing credit card and couldn't use CI/CD
      - after adding credit card, still couldn't test CI/CD
        - could not see "settings > CI/CD" alluded to be docs. seem to be user 
          set bit, as David could see but Louis could not. we decided to just
          abandon Gitlab altogether
    - Got Google Cloud Build working, but preferred Github or anyone else.
      - Opaque billing and business model
      - Proprietary configuration format with less community support
      - slow. 2min+ builds
      - mediocre UI
      - small nuanced of affecting CloudRun UI
    - Went with Github Actions
  - Settled on Github Actions
  - GCP IAM
    - Finally understand:
      - Roles consist of permissions, which are the atoms
      - Can define custom roles 
      - Security theater. InfoSec risk.
      - GCP DevRel docs refer to internal role name, not searchable in UI
    - very terrible UI:
      - Service Accounts Page
        - Details page has no information about the service account's role....
        - The "Manage access" account CTA does not populate base on selected
          service accounts / "principals"
        - No link or way to edit roles. Must go to IAM
      - IAM Page
        - Can associate role to principals
        - UI tools allow search by Role Name, i.e. "AAM Viewer:,  not By Role
          ID such as roles/dialogflow.aamViewer, which GCP's devrel for this
          alludes to
  - GCP Supply Chain Risk
    - Bad GCP DevRel docs poses supply chain risk
      1. On GitHub, there are predefined templates. One of which is deploy to CloudRun.
      2. This template includes a build step google-github-actions/auth
        - https://github.com/google-github-actions/auth
      3. GCP DevRel added template, but it's broken. Docs are distracting, too.
        - turns out the file they checked in was missing special key 
          - people at GCP do not check their work or deploy to GCP even 
      4. first, we need to add 'roles/iam.serviceAccountTokenCreator' to
        the service account to auth 
      5. was getting empty access_token values back. turns out we need to 
         add special property "token_format: 'access_token'". 
      6. then push to GCR failed. Error with missing 'storage.buckets.get' 
        permissions. It turns out we need the '' was insufficient. We originally
        gave the service account 'cloud-platform' ( reaad/write) permissions,
        which is far too wide and poses risk of compromised deleted our project,
        for example.
      7. The 'Container Registry Service Agent' lacked the permissions.
         This is a very misleading role. In the end, we didn't even need this 
         role. "Storage Admin" was sufficient.
  - Github 
    - user-owned repos can only have one owner, everyone else is maintainer
    - maintainer cannot see settings such as secrets used by Github Actions
    - Created Org 
      - transfered ownership from repo to org
        - super dodgy transfer mechanism. if typo, can send repo to unintended
          party. lacking two-way handshake mechanism for proof of ownership.
      - made David and Louis owners of the org
      - Public page, but members private by default ( good )
    - Observed they are selling teams and other products now
  - Cleanup GCP
    - removed old GCR repos in linear-cinema-360910 (kanazawa2)
      - Goal: save server costs 
      - product observations:
        - noticed if you delete all images, repo disappears  
        - noticed you can't delete entire repo. must delete image by image
      - ~100mb/image. We removed maybe 50 images.
    - Remove service accounts
      - Goal: infosec 
    - Is it safe to rename project without breaking things? 🤷🏻