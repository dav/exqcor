# Exqcor Apple Apps

This directory contains a multiplatform Xcode project for macOS and iOS writer apps.

## Generate the Xcode project

From this directory:

```
xcodegen generate
```

This produces `Exqcor.xcodeproj` with two app targets:

- `ExqcorMac` (macOS)
- `ExqcoriOS` (iOS)

## Running

1. Run the Rails server in `modern/` on `http://localhost:3000`.
2. Open `Exqcor.xcodeproj`.
3. Select the target (`ExqcorMac` or `ExqcoriOS`) and run.
4. Enter the server URL and load scripts.

## Notes

- The app uses bearer tokens if you set `APIClient.token` in the shared core code.
- `NSAllowsArbitraryLoads` is enabled for local development.
