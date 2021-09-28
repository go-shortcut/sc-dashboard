# shortcut-story-marker

adds Labels to Shortcut stories related to the release branch

# build

```shell
docker login docker.io
TAGVER=v0.0.3
DOCKERID=...
docker build -f build/Dockerfile -t $DOCKERID/shortcut-story-marker:$TAGVER .
docker push                         $DOCKERID/shortcut-story-marker:$TAGVER
```

# workflow example

.github/workflows/shortcut-story-marker-main.yml
```yaml
name: GitHub Actions shortcut-story-marker

on:
  pull_request:
    types: [closed]
    branches: [main]

jobs:
  job-shortcut-story-marker:
    if: github.event.pull_request.merged == true
    runs-on: ubuntu-latest
    env:
      GITHUB_ACCESS_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      PULL_NUMBER:         ${{ github.event.pull_request.number }}
      SHORTCUT_API_TOKEN:  ${{ secrets.SHORTCUT_API_TOKEN }}
      SHORTCUT_ADD_LABEL:  "main"
      SHORTCUT_DEL_LABEL:  "premain"
    steps:
      - name: get env
        run: env
      - name: do all things
        uses: docker://<DOCKERID>/shortcut-story-marker:v0.0.3

```
