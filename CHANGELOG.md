# Changelog

All notable changes to this project will be documented in this file.

## [1.4.1] - 2025-05-27

### 🐛 Bug Fixes

- Update version state and references logic (#99)

### ⚙️ Miscellaneous Tasks

- *(release)* Update changelog (#96)

## [1.4.0] - 2025-05-23

### 🚀 Features

- *(schema)* Add retry backoff to `CreateSchema` (#95)

### ⚙️ Miscellaneous Tasks

- *(release)* Update changelog (#94)

## [1.3.1] - 2025-05-21

### 🐛 Bug Fixes

- Update schema definition on read (#85)

### ⚙️ Miscellaneous Tasks

- *(release)* Update changelog (#88)

## [1.3.0] - 2025-05-02

### 🚀 Features

- Use a custom client for cookie awareness (#89)

### 📚 Documentation

- Add openssf badge (#54)
- Update attribute references (#60)

### ⚙️ Miscellaneous Tasks

- Update changelog workflow triggers (#47)
- Implement additional workflow checks (#50)
- Update golangci linters (#56)
- Ungroup dependency updates (#57)
- Address golangci linting issues (#63)
- Re-enable markdownlint workflow with changed-files bump (#86)
- Update PR title validation workflow (#87)

## [1.2.2] - 2024-09-05

### 🐛 Bug Fixes

- Handle undefined subject compatibility levels (#45)

## [1.2.1] - 2024-08-27

### 🐛 Bug Fixes

- Resolve schema normalization (#43)

### ⚙️ Miscellaneous Tasks

- Use `git-cliff` for changelog generation (#34)
- Refactor changelog workflow (#38)
- *(release)* Update changelog (#40)

## [1.2.0] - 2024-08-13

### 🚀 Features

- Add support for optionally hard deleting schemas (#33)

### 📚 Documentation

- Update docs after release v1.1.0 (#31)

## [1.1.0] - 2024-08-12

### 🚜 Refactor

- [**breaking**] Enforce uppercase `schema_type` values for consistency (#30)

## [1.0.1] - 2024-07-29

### 🐛 Bug Fixes

- Add support for schema json normalization (#27)

## [1.0.0] - 2024-07-25

### 🚜 Refactor

- [**breaking**] Rename provider from `schema-registry` to `schemaregistry` (#18)

## [0.0.1] - 2024-07-22

### 🚀 Features

- Initialize repository and implement provider (#2)

### 🐛 Bug Fixes

- Improve support for schema references (#8)

### 🚜 Refactor

- [**breaking**] Rename `kafka-schema-registry` to `schema-registry` as the provider is kafka agnostic (#12)
- Remove specific version reference in goreleaser (#13)

### 📚 Documentation

- Create SECURITY.md (#16)

### ⚙️ Miscellaneous Tasks

- Update release and remove generate docs (#3)
- Fix typo in workflow, redefine variables and update error message (#4)
- Update descriptions and errors for clarity (#5)
- Add validation for schema registry URL (#6)
- Add semantic-release-go, scorecard and changelog automation (#7)
- Fix go-semantic-release version (#9)
- Fix go-semantic-release version (#10)
- Update Makefile and README (#11)
