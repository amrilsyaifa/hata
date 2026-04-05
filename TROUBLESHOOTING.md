# Troubleshooting

Common errors you may encounter when setting up or running `hata`.

---

## Error 400: redirect_uri_mismatch

**Full message:**
```
Error 400: redirect_uri_mismatch
The redirect URI in the request did not match a registered redirect URI.
```

**Cause:** The redirect URI used by `hata` during OAuth (`http://localhost:8085`) is not registered in your Google Cloud OAuth client.

**Fix:**
1. Go to [Google Cloud Console](https://console.cloud.google.com) → **APIs & Services** → **Credentials**
2. Click your OAuth 2.0 Client ID
3. Under **Authorized redirect URIs**, click **Add URI**
4. Enter exactly: `http://localhost:8085`
5. Click **Save**

> Changes may take a few minutes to propagate. Wait 1–2 minutes then retry.

---

## Error 403: access_denied

**Full message:**
```
Error 403: access_denied
Test Workshop GDG has not completed the Google verification process.
The app is currently being tested and can only be accessed by developer-approved testers.
```

**Cause:** Your OAuth consent screen is in **Testing** mode, which restricts access to explicitly approved test users. The Google account you're signing in with has not been added as a tester.

**Fix:**
1. Go to [Google Cloud Console](https://console.cloud.google.com) → **APIs & Services** → **OAuth consent screen**
2. Scroll down to the **Test users** section
3. Click **+ Add Users**
4. Enter the Gmail address you use to authenticate
5. Click **Save**

Then retry — authentication will now succeed for that account.

> Apps in Testing mode support up to 100 test users. Tokens expire after 7 days. For personal or team use, staying in Testing mode and adding team emails is sufficient.

---

## Error 403: The caller does not have permission

**Full message:**
```
googleapi: Error 403: The caller does not have permission, forbidden
```

**Cause:** The Google account that authenticated via OAuth does not have edit (or view) access to the Google Sheet specified in your `i18n.config.yml`.

**Fix:**
1. Open your Google Sheet
2. Click **Share** (top right)
3. Enter the Gmail address you authenticated with
4. Set the role to **Editor**
5. Click **Send**

Then retry `./hata push` — the cached token will be reused.

---

## Error 400: Unable to parse range

**Full message:**
```
googleapi: Error 400: Unable to parse range: 'SheetName'!A1:Z1, badRequest
```

**Cause:** The `sheet.name` value in your `i18n.config.yml` does not match the actual tab name in your Google Sheet.

**Fix:**
1. Open your Google Sheet
2. Check the tab name at the bottom (e.g., `Sheet1`, `Translations`, `amril`)
3. Update your config to match exactly:

```yaml
sheet:
  id: "YOUR_SHEET_ID"
  name: "Sheet1"   # must match the tab name exactly (case-sensitive)
```

---

## Token expired or invalid

**Symptom:** Authentication worked before but now fails with a token error.

**Fix:** Delete the cached token file and re-authenticate:

```bash
rm .i18n/token.json
./hata push
```

A new browser window will open for re-authentication.

---

## Credentials file not found

**Full message:**
```
unable to read credentials file: open .i18n/credentials.json: no such file or directory
```

**Fix:** Make sure your credentials file exists at the path specified in `auth.credentials_path` in your config. Download it from:

- [Google Cloud Console](https://console.cloud.google.com) → **APIs & Services** → **Credentials** → your OAuth Client → **Download JSON**

Rename it to match your config path (e.g., `.i18n/credentials.json`).
