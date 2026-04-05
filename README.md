# HATA

```
██╗  ██╗ █████╗ ████████╗ █████╗ 
██║  ██║██╔══██╗╚══██╔══╝██╔══██╗
███████║███████║   ██║   ███████║
██╔══██║██╔══██║   ██║   ██╔══██║
██║  ██║██║  ██║   ██║   ██║  ██║
╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝
```

> A lightweight **Translation Management System (TMS)** — sync i18n translation keys between your codebase and Google Sheets with a single CLI command.

[![Go Version](https://img.shields.io/badge/go-1.25+-blue)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Author](https://img.shields.io/badge/author-amrilsyaifa-orange)](https://github.com/amrilsyaifa)

---

## Overview

**Hata** is a lightweight translation management tool ✅ that bridges the gap between your codebase's translation keys and your translators. No more copy-pasting between JSON files and spreadsheets — developers push keys from `base.json` to a shared Google Sheet, translators fill in the translations without touching any code, and developers pull the completed translations back as ready-to-use JSON files.

If you're looking for a simple Translation Management System (TMS) that works with Google Sheets, `hata` is built for you.

```
base.json  →  (push)  →  Google Sheet  →  (pull)  →  locales/*.json  →  Your App
```

---

## Features

- **`init`** — Interactive setup wizard with locale picker, alias config, and export format selection
- **`push`** — Sync keys from `base.json` → Google Sheet, filling the `base` column; never touches language columns
- **`pull`** — Generate per-language JSON files from the sheet (flat or nested output)
- **`diff`** — Show what's out of sync between `base.json` and the sheet
- **`import`** — Bulk-import an existing locale JSON file into a sheet column (great for migrating existing projects)
- Supports **Service Account** and **OAuth** authentication
- Flat and nested `base.json` input support (both flattened to dot-notation keys)
- **Nested or flat JSON output** — controlled by `nested_json` option in config
- **Language aliases** — export as `en.json` / `id.json` instead of `en-US.json` / `id-ID.json`
- Interpolation passthrough (`Hello {{name}}`)
- Interactive locale selector with search/filter (250+ locales)

---

## Installation

### Option 1: Install via `go install`

```bash
go install github.com/amrilsyaifa/hata@latest
```

### Option 2: Build from source

```bash
git clone https://github.com/amrilsyaifa/hata.git
cd hata
go build -o hata .

# Move to your PATH (macOS/Linux)
mv hata /usr/local/bin/hata
```

### Option 3: Download binary

Download the latest binary from the [Releases](https://github.com/amrilsyaifa/hata/releases) page.

---

## Google Sheets Setup

### Option A: Service Account (recommended for teams/CI)

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a project → Enable **Google Sheets API**
3. Create a **Service Account** → Download the JSON credentials
4. Share your Google Sheet with the service account's email address (Editor access)

### Option B: OAuth (personal use)

#### Step 1 — Create OAuth credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a project → Enable **Google Sheets API**
3. Go to **APIs & Services → OAuth consent screen**
   - User Type: **External**
   - Fill in App name and support email → Save
   - Under **Test users** → add your Google email → Save
4. Go to **APIs & Services → Credentials → + Create Credentials → OAuth client ID**
   - Application type: **Web application** (or Desktop app)
   - Name: anything (e.g. `hata`)
   - Under **Authorized redirect URIs** → click **Add URI** → enter:
     ```
     http://localhost:8085
     ```
   - Click **Create**
5. Click **Download JSON** → save as `.i18n/credentials.json` in your project

#### Step 2 — Run hata

```bash
mkdir -p .i18n
# Move your downloaded credentials file
mv ~/Downloads/client_secret_*.json .i18n/credentials.json

# Update i18n.config.yml
hata init
# Choose: OAuth, credentials path: .i18n/credentials.json, token path: .i18n/token.json
```

#### Step 3 — First run (browser opens automatically)

```bash
./hata push
```

1. Browser opens → log in with your Google account
2. Grant access to Google Sheets
3. Browser shows **"Authorization successful!"** — close it
4. Token saved to `.i18n/token.json` automatically
5. All future runs reuse the cached token — no browser needed

> **Note:** If you see `redirect_uri_mismatch`, make sure `http://localhost:8085` is added to your OAuth client's **Authorized redirect URIs** in Google Cloud Console.

---

## Quick Start

### 1. Initialize your project

```bash
hata init
```

This will guide you through:
- Project ID
- Google Sheet ID (from the sheet URL: `https://docs.google.com/spreadsheets/d/YOUR_SHEET_ID/edit`)
- Auth method (Service Account or OAuth)
- Language selection (interactive picker — Space to select, Enter to confirm)
- **Short aliases** for each locale (e.g. `en` for `en-US`, `id` for `id-ID`) — used as output filenames
- Base file and output directory paths
- **Export format** — nested JSON or flat JSON

It generates `i18n.config.yml` in your project root.

### 2. Create your base file

`base.json` is your source of truth. Both **flat** and **nested** formats are supported:

**Flat:**
```json
{
  "auth.login": "Login",
  "auth.logout": "Logout",
  "auth.welcome": "Hello {{name}}",
  "home.title": "Welcome to our app"
}
```

**Nested:**
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

Both produce the same dot-notation keys in the sheet (`auth.login`, `home.title`, etc.).

### 3. Push keys to Google Sheet

```bash
hata push
```

New keys are appended to the sheet with the **base** text pre-filled. Language columns (`en-US`, `id-ID`, …) are left empty for translators to fill in. Existing rows are never deleted or overwritten.

### 4. Translators fill in the sheet

| key | base | en-US | id-ID |
|---|---|---|---|
| auth.login | Login | Login | Masuk |
| auth.logout | Logout | Logout | Keluar |
| auth.welcome | Hello {{name}} | Hello {{name}} | Halo {{name}} |
| home.title | Welcome to our app | Welcome to our app | Selamat datang |

> The `base` column is managed by `hata push`. Translators only edit the language columns.

### 5. Pull translations into your project

```bash
hata pull
```

Generates one JSON file per language. If aliases are configured (`en-US → en`), files are named `locales/en.json`, `locales/id.json`, etc. Output format (nested or flat) is controlled by `nested_json` in config.

**Nested output** (`nested_json: true`, default):
```json
{
  "auth": {
    "login": "Login",
    "logout": "Logout",
    "welcome": "Hello {{name}}"
  },
  "home": {
    "subtitle": "Get started below",
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
  "home.subtitle": "Get started below",
  "home.title": "Welcome to our app"
}
```

### 6. Migrate an existing locale file into the sheet

If you already have `locales/id.json` with many translations, import it directly:

```bash
# First push your keys to create sheet rows
hata push

# Then import the existing translations into the id-ID column
hata import --file ./locales/id.json --lang id-ID
hata import --file ./locales/en.json --lang en-US
```

- Accepts **nested or flat** JSON (auto-flattened)
- Updates only keys that already exist in the sheet
- Prints a list of any keys not yet in the sheet

### 7. Check for drift


```bash
hata diff
```

```
Missing in sheet:
  - auth.register

Unused in base:
  - old.key
```

---

## Config Reference

`i18n.config.yml`:

```yaml
project_id: my-project

sheet:
  id: "1abc123xyz..."   # from your Google Sheet URL
  name: "Translations"  # sheet tab name

auth:
  type: service_account          # or "oauth"
  credentials_path: ".i18n/credentials.json"
  token_path: ".i18n/token.json" # used by oauth only

languages:
  - en-US
  - id-ID

# Optional: short aliases used as output filenames when pulling.
# Sheet columns still use the full locale code (en-US, id-ID).
aliases:
  en-US: en   # → locales/en.json
  id-ID: id   # → locales/id.json

paths:
  base: "./base.json"
  output: "./locales"

options:
  nested_json: true   # true = nested JSON output, false = flat key output
  sort_keys: true     # sort keys alphabetically in output
  keep_unused: true   # keep stale keys in sheet (don't auto-delete)
```

> **Security:** Add `.i18n/` to your `.gitignore` to avoid committing credentials.

---

## Command Reference

| Command | Description |
|---|---|
| `hata init` | Interactive setup — creates `i18n.config.yml` |
| `hata push` | Sync new keys from `base.json` to sheet |
| `hata pull` | Download translations from sheet → locale JSON files |
| `hata diff` | Show keys that are out of sync |
| `hata import -f <file> -l <lang>` | Import an existing locale file into a sheet column |
| `hata clean` | Interactively remove stale keys from Google Sheet |

### `hata import` flags

| Flag | Short | Description |
|---|---|---|
| `--file` | `-f` | Path to the locale JSON file (required) |
| `--lang` | `-l` | Locale code as it appears in the sheet header, e.g. `id-ID` (required) |

```bash
hata import --file ./locales/id.json --lang id-ID
hata import -f ./locales/en.json -l en-US
```

### `hata clean`

Finds keys that exist in the Google Sheet but are no longer present in your `base.json`. Presents an interactive multi-select list, then permanently deletes the confirmed rows from the sheet.

```bash
hata clean
```

**Flow:**

1. Reads all keys from the sheet and compares against `base.json`.
2. Displays a multi-select list of stale keys (Space to toggle, Enter to confirm).
3. Asks for a final confirmation before deleting.
4. Deletes selected rows from the sheet in a single batch update.

> **Note:** Deletions cannot be undone. Review the key list carefully before confirming.

---

## Integration Guides

### React / Next.js

Install a compatible i18n library (e.g., `react-i18next` or `next-intl`):

```bash
npm install react-i18next i18next
```

Add a `hata pull` step to your build/dev workflow:

```bash
# package.json
{
  "scripts": {
    "i18n:pull": "hata pull",
    "dev": "hata pull && next dev",
    "build": "hata pull && next build"
  }
}
```

Load the generated JSON files:

```js
// i18n.js
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import en from './locales/en.json';
import id from './locales/id.json';

i18n.use(initReactI18next).init({
  resources: { en: { translation: en }, id: { translation: id } },
  lng: 'en',
  fallbackLng: 'en',
});

export default i18n;
```

Use in components:

```jsx
import { useTranslation } from 'react-i18next';

function LoginButton() {
  const { t } = useTranslation();
  return <button>{t('auth.login')}</button>;
}
```

Use with Next.js `next-intl`:

```js
// next.config.js
import createNextIntlPlugin from 'next-intl/plugin';
const withNextIntl = createNextIntlPlugin();
export default withNextIntl({});
```

```js
// messages/en.json  ← point hata output here
// app/[locale]/layout.tsx
import { NextIntlClientProvider } from 'next-intl';
```

---

### React Native

```bash
npm install i18next react-i18next
```

Copy `locales/` into your project (e.g., `src/locales/`), then update `i18n.config.yml`:

```yaml
paths:
  output: "./src/locales"
```

```js
// src/i18n.js
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import en from './locales/en.json';
import id from './locales/id.json';

i18n.use(initReactI18next).init({
  resources: { en: { translation: en }, id: { translation: id } },
  lng: 'en',
  fallbackLng: 'en',
  interpolation: { escapeValue: false },
});

export default i18n;
```

```jsx
// App.js
import './src/i18n';
import { useTranslation } from 'react-i18next';

export default function App() {
  const { t } = useTranslation();
  return <Text>{t('home.title')}</Text>;
}
```

---

### Vue.js

```bash
npm install vue-i18n
```

```js
// src/i18n.js
import { createI18n } from 'vue-i18n';
import en from './locales/en.json';
import id from './locales/id.json';

export default createI18n({
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en, id },
});
```

```html
<!-- Component -->
<template>
  <button>{{ $t('auth.login') }}</button>
</template>
```

---

### Golang

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
)

func loadLocale(lang string) (map[string]interface{}, error) {
    data, err := os.ReadFile(fmt.Sprintf("locales/%s.json", lang))
    if err != nil {
        return nil, err
    }
    var messages map[string]interface{}
    return messages, json.Unmarshal(data, &messages)
}
```

Or use a library like [`go-i18n`](https://github.com/nicksnyder/go-i18n):

```bash
go get github.com/nicksnyder/go-i18n/v2
```

```go
bundle := i18n.NewBundle(language.English)
bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
bundle.LoadMessageFile("locales/en.json")
bundle.LoadMessageFile("locales/id.json")
```

---

### Ruby on Rails

Add a `hata pull` step to your CI or `Rakefile`, then copy the output to your Rails locale path:

```yaml
# i18n.config.yml — point output to Rails locale dir
paths:
  output: "./config/locales"
```

Rails expects YAML by default. To use JSON, add `i18n-js` or convert with a Rake task:

```ruby
# Rakefile
task :pull_translations do
  system('hata pull')
end
```

Alternatively use `hata pull` and process with a script to convert JSON → YAML:

```ruby
require 'json'
require 'yaml'

Dir['config/locales/*.json'].each do |f|
  lang = File.basename(f, '.json')
  data = JSON.parse(File.read(f))
  File.write("config/locales/#{lang}.yml", { lang => data }.to_yaml)
end
```

---

## CI/CD Integration

### GitHub Actions

```yaml
# .github/workflows/i18n.yml
name: Sync Translations

on:
  schedule:
    - cron: '0 8 * * 1'  # Every Monday at 8am
  workflow_dispatch:

jobs:
  pull:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install hata
        run: go install github.com/amrilsyaifa/hata@latest

      - name: Write credentials
        run: |
          mkdir -p .i18n
          echo '${{ secrets.GOOGLE_CREDENTIALS }}' > .i18n/credentials.json

      - name: Pull translations
        run: hata pull

      - name: Commit updated locales
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git add locales/
          git diff --cached --quiet || git commit -m "chore: update translations"
          git push
```

> Store your service account credentials JSON as a GitHub secret named `GOOGLE_CREDENTIALS`.

---

## CLI Reference

```
hata [command] [flags]

Commands:
  init        Initialize hata configuration interactively
  push        Sync translation keys from base.json to Google Sheet
  pull        Generate per-language JSON files from Google Sheet
  diff        Show differences between base.json and Google Sheet
  help        Help about any command

Flags:
  --config string   Config file path (default "i18n.config.yml")
  -h, --help        Help for hata
```

---

## Project Structure

```
your-project/
├── base.json           ← Your translation keys (source of truth)
├── i18n.config.yml     ← Hata configuration
├── .i18n/
│   └── credentials.json  ← Google credentials (gitignored)
└── locales/            ← Generated output (can be gitignored)
    ├── en.json
    └── id.json
```

---

## Author

**Amril Syaifa**
- GitHub: [@amrilsyaifa](https://github.com/amrilsyaifa)
- Repository: [github.com/amrilsyaifa/hata](https://github.com/amrilsyaifa/hata)

---

## License

MIT © [Amril Syaifa](https://github.com/amrilsyaifa)
