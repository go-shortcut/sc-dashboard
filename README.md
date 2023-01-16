# shortcut-story-marker

adds Labels to Shortcut stories related to the release branch

# build

```shell
TAGVER=v0.0.5
DOCKERID=ghcr.io
OWNER=go-shortcut
echo $GITHUB_ACCESS_TOKEN | docker login ghcr.io -u $OWNER --password-stdin
docker build -f build/Dockerfile -t $DOCKERID/$OWNER/shortcut-story-marker:$TAGVER .
docker push                         $DOCKERID/$OWNER/shortcut-story-marker:$TAGVER
```

# workflow example

.github/workflows/shortcut-story-marker-main.yml
```yaml
name: GitHub Actions shortcut-story-marker

on:
  pull_request:
    types: [ labeled ]

jobs:
  job-shortcut-story-marker-b2b:
    if: github.event.label.name == 'need_shortcut_report'
    runs-on: ubuntu-latest
    env:
      GITHUB_ACCESS_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      PULL_NUMBER:         ${{ github.event.pull_request.number }}
      SHORTCUT_API_TOKEN:  ${{ secrets.SHORTCUT_API_TOKEN }}
      SHORTCUT_ADD_LABEL:  ${{ github.event.pull_request.head.ref }}
    steps:
      - name: get env
        run: env
      - name: do all things
        uses: docker://<DOCKERID>/<OWNER>/shortcut-story-marker:<TAGVER>

```
