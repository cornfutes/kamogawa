import { test, expect } from '@playwright/test';

test('desktops', async ({ page }) => {
  // Makes sure landing page loads
  await page.goto('https://diceduckmonk.com');
  await expect(page).toHaveTitle(/DiceDuckMonk/);

  // Go from landing page to proceed to demo page
  const landingCta = page.locator('text=Try it out!');
  await landingCta.click();
  await expect(page).toHaveURL(/demo/);

  // Go from demo page to app
  const demoCta = page.locator('text=Proceed to demo');
  await demoCta.click();
  await expect(page).toHaveURL(/search\?q=test/);

  // inspect the search reviews page
  await expect(page.locator('#content')).toHaveText(/1 results/)
  await expect(page.locator('text=GCP Project')).toHaveText(/GCP Project \"diceduckmonk-test-project\"/)

  // Go to VMs page
  await page.locator('#nav-gcp').click();
  await page.locator('text=VMs').click();
  await expect(page).toHaveURL(/gce/);
  await expect(page.locator('text=Cache avoided')).toBeVisible();
  const cachedResultContent = await page.locator('text=Cache avoided').textContent();
  expect(cachedResultContent?.trim()).toEqual("Cache avoided 15 API calls");
  await expect(page).toHaveScreenshot({
    mask: [page.locator('text=Query took')],
  })
});
