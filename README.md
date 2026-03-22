# REG

REG is short for **Receiver of Expired Guides**.

This repository automates workflow actions against the sica.sunagro.gob.ve web app, including login flow and subsequent receipt/inventory features

### Required environment variables

- `.env` (recommended) has:
  - `REG_TEST_USERNAME`
  - `REG_TEST_PASSWORD`

### Network flag for integration tests

To avoid running browser tests unintentionally (offline or pure unit test passes), we guard E2E with:

- `REG_E2E=1`

### Headless mode on/off flag

You can control headless mode with:

- `REG_HEADLESS=1` (default if not set)
- `REG_HEADLESS=0` or `REG_HEADLESS=false` (for debugging with visible browser)

If not set, Rod suite is skipped.

### Run tests

```bash
# set credentials in .env (or environment)
# run integration tests (rod-based) explicitly:
REG_E2E=1 go test ./... -v

# or on Windows PowerShell:
$env:REG_E2E = "1"
go test ./... -v
```

## Notes

- `testutil/rod_suite.go` keeps the browser session shared across suite and page per test.

