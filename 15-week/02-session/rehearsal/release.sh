#!/usr/bin/env bash
# Rehearsal of the release of the session: cut release/3.0.0 from main, apply one commit per story with cherry-pick -x,
# run the gate, merge to main, tag v3.0.0, deploy it, and roll back to the previous tag without a rebuild.
# Everything happens in a temporary directory with a local bare "remote". No project repository is touched.
set -u
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT
PASS=0
FAIL=0
check() { # check <description> <command...>
  local d=$1
  shift
  if "$@" >/dev/null 2>&1; then
    PASS=$((PASS + 1))
    echo "  ok   $d"
  else
    FAIL=$((FAIL + 1))
    echo "  FAIL $d"
  fi
}
refuses() { # refuses <description> <command...>: the command must fail
  local d=$1
  shift
  if "$@" >/dev/null 2>&1; then
    FAIL=$((FAIL + 1))
    echo "  FAIL $d (it was accepted)"
  else
    PASS=$((PASS + 1))
    echo "  ok   $d"
  fi
}

git init -q --bare "$WORK/remote.git"
git clone -q "$WORK/remote.git" "$WORK/repo" 2>/dev/null
cd "$WORK/repo" || exit 1
git config user.name rehearsal
git config user.email rehearsal@example.invalid
git checkout -q -b main

# the gate: the "pipeline". It fails when a file named BROKEN exists.
cat >gate.sh <<'EOF'
#!/usr/bin/env bash
[ -f BROKEN ] && { echo "gate: tests failed"; exit 1; }
echo "gate: green"
EOF
echo "0.1.0" >VERSION
git add . && git commit -q -m "chore: initial state of main"
git tag -a v2.0.0 -m "Release v2.0.0: previous release"
git push -q origin main --tags 2>/dev/null

git checkout -q -b develop
for hu in hu-001 hu-002 hu-003; do
  echo "$hu" >"$hu.txt"
  git add . && git commit -q -m "feat: add $hu"
done
git push -q origin develop 2>/dev/null
mapfile -t COMMITS < <(git log --reverse --format=%H main..develop)

echo "1. cut release/3.0.0 from main and apply one commit per story"
git checkout -q main
git checkout -q -b release/3.0.0
for c in "${COMMITS[@]}"; do git cherry-pick -x "$c" >/dev/null 2>&1; done
check "three story commits on the release branch" test "$(git rev-list --count main..release/3.0.0)" = 3
check "every commit carries the cherry picked from trail" test "$(git log main..release/3.0.0 --format=%B | grep -c 'cherry picked from commit')" = 3

echo "2. the gate decides"
touch BROKEN
refuses "a red gate blocks the tag" bash -c 'bash gate.sh && git tag -a v3.0.0 -m x'
check "no tag exists after the red gate" test -z "$(git tag -l v3.0.0)"
rm BROKEN
check "the gate is green again" bash gate.sh

echo "3. merge to main, tag and publish"
git checkout -q main
git merge -q --no-ff release/3.0.0 -m "chore: merge release/3.0.0"
echo "3.0.0" >VERSION
git add VERSION && git commit -q -m "chore: set the version to 3.0.0"
git tag -a v3.0.0 -m "Release v3.0.0: stories 001 to 003"
git push -q origin main v3.0.0 2>/dev/null
check "v3.0.0 is an annotated tag" test "$(git cat-file -t v3.0.0)" = tag
check "the remote has the tag" git ls-remote --exit-code --tags origin v3.0.0
refuses "a published tag is not moved (a second tag with the same name)" git tag -a v3.0.0 -m again
check "main contains the three stories" test -f hu-003.txt

echo "4. deploy the tag and roll back to the previous one (the same artifact, no rebuild)"
deploy() { git checkout -q "$1" && cp VERSION "$WORK/deployed" && git checkout -q main; }
deploy v3.0.0
check "the deployed version is 3.0.0" grep -qx 3.0.0 "$WORK/deployed"
deploy v2.0.0
check "the rollback deployed 0.1.0 (the content of v2.0.0)" grep -qx 0.1.0 "$WORK/deployed"
check "the rollback did not create a commit" test "$(git rev-parse main)" = "$(git rev-parse origin/main)"

echo
echo "checks passed: $PASS, failed: $FAIL"
[ "$FAIL" -eq 0 ]
