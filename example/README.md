# Example Project

This folder shows the expected file structure for using `hata` in a project.

## Files

| File | Description |
|------|-------------|
| `base.json` | Your source of truth. Supports both flat (`"auth.login": "Login"`) and nested JSON formats. |
| `i18n.config.yml` | Example config — copy to your project root and fill in your Sheet ID. |

## Google Sheet structure

After running `hata push`, your sheet will look like:

| key | base | en-US | id-ID | ja-JP |
|-----|------|-------|-------|-------|
| auth.login | Login | Login | Masuk | ログイン |
| auth.logout | Logout | Logout | Keluar | ログアウト |
| auth.welcome | Hello {{name}} | Hello {{name}} | Halo {{name}} | こんにちは {{name}} |
| home.title | Welcome to our app | Welcome to our app | Selamat datang | ようこそ |

> The `base` column is filled automatically by `hata push`. Translators only edit the language columns (`en-US`, `id-ID`, …).

## After `hata pull`

The `locales/` directory will be generated:

```
locales/
├── en-US.json
├── id-ID.json
└── ja-JP.json
```

Each file is nested JSON:

```json
{
  "auth": {
    "login": "Login",
    "logout": "Logout",
    "welcome": "Hello {{name}}"
  },
  "home": {
    "title": "Welcome to our app"
  }
}
```

## OAuth setup

See the full [OAuth setup guide](../README.md#option-b-oauth-personal-use) in the main README.

Quick checklist:
- [ ] Google Cloud project created
- [ ] Google Sheets API enabled
- [ ] OAuth consent screen configured (add your email as test user)
- [ ] OAuth client created (Web app type)
- [ ] `http://localhost:8085` added to Authorized redirect URIs
- [ ] credentials JSON downloaded to `.i18n/credentials.json`
- [ ] Run `./hata push` — browser opens, authorize, done
