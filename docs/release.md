# Release Process

Use the following steps to make a release and publish binaries.

1. Bump version in `README.md` and `docs/changelog.md`.
2. Commit changes and push to `main`.
3. Create an annotated git tag, for example:

```bash
git tag -a v1.1.0 -m "Release v1.1.0"
```

4. Build release binaries:

```bash
mkdir -p dist
GOOS=linux GOARCH=amd64 go build -o dist/preekeeper-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o dist/preekeeper-windows-amd64.exe .
GOOS=darwin GOARCH=amd64 go build -o dist/preekeeper-darwin-amd64 .
```

5. Push tags and branch:

```bash
git push origin main --follow-tags
```

6. Create GitHub Release (via web UI or using GitHub CLI):

```bash
gh release create v1.1.0 dist/preekeeper-* --title "v1.1.0" --notes "Release notes..."
```

7. Verify assets are available on the release page.

Notes: If you need cross-architecture builds (arm64) add GOARCH/GOOS targets appropriately.