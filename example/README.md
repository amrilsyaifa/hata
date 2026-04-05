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

## Aliases

The config uses `aliases` to map full locale codes to short output filenames:

```yaml
aliases:
  en-US: en   # → locales/en.json
  id-ID: id   # → locales/id.json
```

The sheet columns always use the full code (`en-US`, `id-ID`). Aliases only affect output filenames.

## Export format

Controlled by `options.nested_json` in `i18n.config.yml`.

**Nested JSON** (`nested_json: true`, default): keys are expanded into objects.

**Flat JSON** (`nested_json: false`): keys stay as dot-notation strings.

## After `hata pull`

The `locales/` directory will be generated (using aliases as filenames):

```
locales/
├── en.json
├── id.json
└── ja-JP.json   ← no alias configured, uses full code
```

**Nested output** (`nested_json: true`):
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

**Flat output** (`nested_json: false`):
```json
{
  "auth.login": "Login",
  "auth.logout": "Logout",
  "auth.welcome": "Hello {{name}}",
  "home.title": "Welcome to our app"
}
```

## Migrating an existing project

If you already have `id.json` with many translations, import it in one command:

```bash
# Create key rows in sheet first
hata push

# Import existing translations into the matching column
hata import --file ./locales/id.json --lang id-ID
hata import --file ./locales/en.json --lang en-US
```

`import` accepts nested or flat JSON automatically and only updates rows that already exist in the sheet.

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
