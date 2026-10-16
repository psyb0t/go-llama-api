# Releases and docs

Refines `~/.claude/rule-details/versioning-a-project.md` and `~/.claude/rule-details/writing-docs.md`.

- **REL1** The git tag is the version. A regenerated client that renames or removes an exported name bumps the minor version before 1.0 and says so in the CHANGELOG; everything else is a patch.
- **REL2** Every release prepends a `## vX.Y.Z, YYYY-MM-DD` section to `CHANGELOG.md`, naming the spec version it was generated from. The annotated tag message is `vX.Y.Z`, a blank line, then that section.
- **REL3** `.agents/.codex-plugin/plugin.json` carries `version` and moves to the release version in the same commit. `.agents/.claude-plugin/plugin.json` has no `version`.
- **REL4** The agent skill in `.agents/skills/go-llama-api/` changes together with any generated method or type it names, and every example in it must compile against the code.
- **REL5** Commits, tags and pushes go through `git-update.sh` only. Tag pushes publish the skill to ClawHub through the `clawhub` job in `pipeline.yml`.
