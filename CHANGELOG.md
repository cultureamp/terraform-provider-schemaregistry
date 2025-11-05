# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### <!-- 7 -->⚙️ Miscellaneous Tasks
- Ci: update build and release pipeline by @dstrates in [#185](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/185)


## [1.5.1] - 2025-11-05


### <!-- 1 -->🐛 Bug Fixes
- Fix: only call `CreateSchema` on changes and update `ModifyPlan` by @dstrates in [#183](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/183)

### <!-- 7 -->⚙️ Miscellaneous Tasks
- Ci: fix pipeline dependencies by @dstrates in [#133](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/133)
- Ci: centralise github workflows by @dstrates in [#129](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/129)
- Chore(release): update changelog by @github-actions[bot] in [#128](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/128)


## [1.5.0] - 2025-06-26


### <!-- 0 -->🚀 Features
- Feat: update `ModifyPlan` with semantic lookup by @dstrates in [#127](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/127)

### <!-- 7 -->⚙️ Miscellaneous Tasks
- Chore(release): update changelog by @github-actions[bot] in [#112](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/112)


## [1.4.3] - 2025-06-19


### <!-- 1 -->🐛 Bug Fixes
- Fix: add `ModifyPlan` to suppress schema diffs by @dstrates in [#104](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/104)

### <!-- 7 -->⚙️ Miscellaneous Tasks
- Ci: allow `workflow_dispatch` trigger for releases by @dstrates in [#123](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/123)
- Ci: fix release conditional and scope permissions by @dstrates in [#122](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/122)
- Ci: enable merge group required checks by @dstrates in [#118](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/118)
- Ci: enable dependabot for github-actions by @dstrates in [#106](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/106)
- Chore(release): update changelog by @github-actions[bot] in [#103](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/103)


## [1.4.2] - 2025-06-03


### <!-- 1 -->🐛 Bug Fixes
- Fix: revert update schema definition on read by @dstrates in [#101](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/101)

### <!-- 7 -->⚙️ Miscellaneous Tasks
- Chore(release): update changelog by @github-actions[bot] in [#100](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/100)


## [1.4.1] - 2025-05-27


### <!-- 1 -->🐛 Bug Fixes
- Fix: update version state and references logic by @dstrates in [#99](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/99)

### <!-- 7 -->⚙️ Miscellaneous Tasks
- Chore(release): update changelog by @github-actions[bot] in [#96](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/96)


## [1.4.0] - 2025-05-23


### <!-- 0 -->🚀 Features
- Feat(schema): add retry backoff to `CreateSchema` by @dstrates in [#95](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/95)

### <!-- 7 -->⚙️ Miscellaneous Tasks
- Chore(release): update changelog by @github-actions[bot] in [#94](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/94)


## [1.3.1] - 2025-05-21


### <!-- 1 -->🐛 Bug Fixes
- Fix: update schema definition on read by @dstrates in [#85](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/85)

### <!-- 7 -->⚙️ Miscellaneous Tasks
- Chore(release): update changelog by @github-actions[bot] in [#88](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/88)


## [1.3.0] - 2025-05-02


### <!-- 0 -->🚀 Features
- Feat: use a custom client for cookie awareness by @dstrates in [#89](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/89)

### <!-- 10 -->💼 Other
- Removing vulnerable Github Action package from  tj-actions by @mr-joshcrane in [#80](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/80)
- Removing vulnerable Github Action package from  tj-actions by @mr-joshcrane

### <!-- 3 -->📚 Documentation
- Docs: update attribute references by @dstrates in [#60](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/60)
- Docs: add openssf badge by @dstrates in [#54](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/54)

### <!-- 7 -->⚙️ Miscellaneous Tasks
- Ci: update PR title validation workflow in [#87](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/87)
- Ci: re-enable markdownlint workflow with changed-files bump by @dstrates in [#86](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/86)
- Chore: address golangci linting issues by @dstrates in [#63](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/63)
- Ci: ungroup dependency updates by @dstrates in [#57](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/57)
- Ci: update golangci linters by @dstrates in [#56](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/56)
- Ci: implement additional workflow checks by @dstrates in [#50](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/50)
- Ci: update changelog workflow triggers by @dstrates in [#47](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/47)

## New Contributors
- @mr-joshcrane made their first contribution in [#80](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/80)

## [1.2.2] - 2024-09-05


### <!-- 1 -->🐛 Bug Fixes
- Fix: handle undefined subject compatibility levels by @dstrates in [#45](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/45)


## [1.2.1] - 2024-08-27


### <!-- 1 -->🐛 Bug Fixes
- Fix: resolve schema normalization by @dstrates in [#43](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/43)

### <!-- 7 -->⚙️ Miscellaneous Tasks
- Chore(release): update changelog by @github-actions[bot] in [#40](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/40)
- Ci: refactor changelog workflow by @dstrates in [#38](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/38)
- Ci: use `git-cliff` for changelog generation by @dstrates in [#34](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/34)

## New Contributors
- @github-actions[bot] made their first contribution in [#40](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/40)

## [1.2.0] - 2024-08-13


### <!-- 0 -->🚀 Features
- Feat: add support for optionally hard deleting schemas by @dstrates in [#33](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/33)

### <!-- 3 -->📚 Documentation
- Docs: update docs after release v1.1.0 by @dstrates in [#31](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/31)


## [1.1.0] - 2024-08-12


### <!-- 0 -->🚀 Features
- Feat: enforce uppercase `schema_type` values for consistency by @dstrates in [#30](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/30)


## [1.0.1] - 2024-07-29


### <!-- 1 -->🐛 Bug Fixes
- Fix: add support for schema json normalization by @dstrates in [#27](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/27)

## New Contributors
- @dependabot[bot] made their first contribution in [#26](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/26)

## [1.0.0] - 2024-07-25


### <!-- 0 -->🚀 Features
- Feat: initialize repository and implement provider by @dstrates in [#2](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/2)

### <!-- 1 -->🐛 Bug Fixes
- Fix: improve support for schema references by @dstrates in [#8](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/8)

### <!-- 10 -->💼 Other
- Initial commit by @dstrates

### <!-- 2 -->🚜 Refactor
- Refactor!: rename provider from `schema-registry` to `schemaregistry` by @dstrates in [#18](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/18)
- Refactor: remove specific version reference in goreleaser in [#13](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/13)
- Refactor!: rename `kafka-schema-registry` to `schema-registry` as the provider is kafka agnostic in [#12](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/12)

### <!-- 3 -->📚 Documentation
- Docs: create SECURITY.md by @dstrates in [#16](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/16)

### <!-- 7 -->⚙️ Miscellaneous Tasks
- Ci: update release workflow triggers by @dstrates in [#14](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/14)
- Chore: update Makefile and README by @dstrates in [#11](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/11)
- Ci: fix go-semantic-release version by @dstrates in [#10](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/10)
- Ci: fix go-semantic-release version by @dstrates in [#9](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/9)
- Ci: add semantic-release-go, scorecard and changelog automation by @dstrates in [#7](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/7)
- Chore: add validation for schema registry URL by @dstrates in [#6](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/6)
- Chore: update descriptions and errors for clarity by @dstrates in [#5](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/5)
- Chore: fix typo in workflow, redefine variables and update error message by @dstrates in [#4](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/4)
- Ci: update release and remove generate docs by @dstrates in [#3](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/3)

## New Contributors
- @dstrates made their first contribution in [#18](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/18)
- @ made their first contribution in [#13](https://github.com/cultureamp/terraform-provider-schemaregistry/pull/13)

