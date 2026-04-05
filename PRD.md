# HATA

> A Golang-based CLI tool for synchronizing i18n data between local files and Google Sheets.

---

## 1. 📌 Overview

**Hata** is a command-line tool built with Golang that bridges your codebase's translation keys and Google Sheets, enabling a smooth collaboration workflow between developers and non-technical translators.

**Core capabilities:**

- Sync translation keys from code → Google Sheet
- Manage translations directly via Google Sheet
- Generate per-language JSON files from the sheet
- Support interpolation syntax (`{{name}}`)

---

## 2. 🎯 Objectives

### Primary Goals

- Simplify the translation workflow between developers and non-developers
- Eliminate hardcoded translations in the codebase
- Provide a single, flexible source of truth for all translation content

### Success Metrics

- Developers can sync i18n data with ≤ 1 command
- Translators only need to edit Google Sheets — no code knowledge required
- No manual editing of translation files

---

## 3. 👥 Target Users

### 👨‍💻 Developer
- Uses the CLI to sync and generate translation files
- Maintains `base.json` as the source of truth for keys

### 🌍 Translator / Non-Developer
- Edits translations directly in Google Sheets
- No CLI or coding knowledge required

---

## 4. 🧱 Core Concepts

### 4.1 Base File (`base.json`)

The source of truth for all translation keys. Uses a flat key format:

```json
{
  "auth.login": "Login",
  "auth.welcome": "Hello {{name}}"
}
```

### 4.2 Google Sheet Structure

| key | en | id |
|---|---|---|
| `auth.login` | Login | Masuk |
| `auth.welcome` | Hello {{name}} | Halo {{name}} |

- `key` — unique identifier for each translation entry
- Additional columns — each represents a language code

### 4.3 Generated Output

Files are generated under the configured output directory:

```
/locales/en.json
/locales/id.json
```

Example output (nested JSON):

```json
{
  "auth": {
    "welcome": "Hello {{name}}"
  }
}
```

---

## 5. ⚙️ Features

### 5.1 CLI Commands

#### `init`

Initializes the project configuration interactively.

**Flow:**
1. Select authentication method:
   - Service Account
   - OAuth
2. Configure:
   - Sheet ID
   - Supported languages (e.g., `en`, `id`)
3. Generate config file:

```yaml
# i18n.config.yml
project_id: my-i18n
sheet_id: xxx
auth:
  type: service_account # or oauth
  credentials_path: .i18n/credentials.json

languages:
  - en
  - id

output_path: ./locales
base_file: ./base.json
```

---

#### `push`

Syncs keys from `base.json` → Google Sheet.

**Behavior:**
- Adds new keys to the sheet
- Does **not** overwrite existing translations
- Does **not** delete keys already in the sheet

---

#### `pull`

Generates JSON translation files from the sheet.

**Behavior:**
- Reads all rows from the sheet
- Converts flat keys → nested JSON structure
- Generates one file per language

---

#### `diff` *(optional but recommended)*

Compares `base.json` against the sheet and reports discrepancies.

**Example output:**

```
Missing in sheet:
  - auth.register

Unused in base:
  - old.key
```

---

## 6. 🔐 Authentication

### 6.1 Service Account *(Default)*

**Flow:**
1. Create a service account in Google Cloud Platform
2. Download the JSON credential file
3. Share the sheet with the service account email

**Pros:**
- Simple setup
- Automation-friendly (CI/CD compatible)

### 6.2 OAuth *(Optional)*

**Flow:**
1. CLI triggers the login process
2. Opens browser for user authentication
3. Token is saved to `.i18n/token.json`

**Use cases:**
- Multi-user environments
- When sharing a sheet manually is not desired

---

## 7. 🗂️ Config File

**File:** `i18n.config.yml`

```yaml
project_id: my-project

sheet:
  id: "1abc123xyz"
  name: "Translations"

auth:
  type: service_account
  credentials_path: ".i18n/credentials.json"
  token_path: ".i18n/token.json"

languages:
  - en
  - id

paths:
  base: "./base.json"
  output: "./locales"

options:
  nested_json: true
  sort_keys: true
  keep_unused: true
```

---

## 8. 🔄 Data Flow

```
base.json
   ↓  push
Google Sheet
   ↓  pull
locales/*.json
   ↓
Application
```

---

## 9. 🧠 Key Behaviors

### 9.1 Flat → Nested Key Conversion

**Input key:**
```
auth.login
```

**Output JSON:**
```json
{
  "auth": {
    "login": "..."
  }
}
```

### 9.2 Interpolation Support

**Format:**
```
Hello {{name}}
```

**Rules:**
- Interpolation tokens are **not** parsed or transformed by the CLI
- They are passed through as-is to the output files
- Validation is optional: can warn on token mismatch across languages

### 9.3 Safe Sync Rules

| Action | Behavior |
|---|---|
| Add key | ✅ Allowed |
| Update existing key | ❌ Not allowed (no overwrite) |
| Delete key | ❌ Not allowed |
| Missing translation | ⚠️ Warning emitted |

---

## 10. ⚠️ Edge Cases

| Scenario | Handling |
|---|---|
| Duplicate key in sheet | Error — aborts operation |
| Missing language column | Skip column with warning |
| Empty sheet | Fail-safe — no files generated |
Invalid JSON structure → validation error
11. 🧪 Validation
Saat pull
Check semua language punya value
Check interpolation consistency:
en: Hello {{name}}
id: Halo
→ warning: missing {{name}}
12. 🏗️ Technical Architecture
CLI (Golang)
using cobra
structure folder follow this https://github.com/golang-standards/project-layout
Modules:

config loader (YAML)
Google Sheets client
JSON transformer
diff engine
13. 🚀 Future Enhancements
Web UI dashboard
Auto sync CI/CD
Translation status tracking (% complete)
Pluralization support
Namespaces (common.auth.login)
Versioning / history
14. 📦 Deliverables
MVP Scope
init
push
pull
service account auth
YAML config
nested JSON generator
15. 📌 Non-Goals (MVP)
Real-time sync
Web dashboard
Translation suggestion AI