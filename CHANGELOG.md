# Changelog

# [unreleased]

## <!-- 29 -->👷 CI/CD

- Use `git-cliff` for changelog generation ([b399796](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/b39979683f2c076c1504d78e320bcf8120700855))  - (dstrates)

## <!-- 3 -->📚 Documentation

- Update docs after release v1.1.0 ([#31](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/31)) ([7da6415](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/7da6415203ea585293a258b64983ca05845eca38))  - (dstrates)

## <!-- 30 -->📝 Other

- Merge b39979683f2c076c1504d78e320bcf8120700855 into 7da6415203ea585293a258b64983ca05845eca38
 ([2596ccc](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/2596cccb5c22ac9f0183723759bed8a02db7ea59))  - (dstrates)

# [1.1.0](https://github.com/cultureamp/terraform-provider-schemaregistry/compare/v1.0.1...v1.1.0) - (2024-08-12)

## <!-- 0 -->🚀 Features

- Enforce uppercase `schema_type` values for consistency ([#30](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/30)) ([443985b](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/443985bacd7bba6f81ca28f28f9279a3a5f17342))  - (dstrates)

# [1.0.1](https://github.com/cultureamp/terraform-provider-schemaregistry/compare/v1.0.0...v1.0.1) - (2024-07-29)

## <!-- 1 -->🐛 Bug Fixes

- Add support for schema json normalization ([#27](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/27)) ([c4a30bc](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/c4a30bc380dcb4e4a356aaeb023170c2b772985e))  - (dstrates)

## <!-- 7 -->⚙️ Miscellaneous Tasks

- **deps:** Bump the all-minor-updates group with 5 updates ([#26](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/26)) ([321e03d](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/321e03d108ff638386baea7bbb893009333bace4))  - (dependabot[bot])
- **deps:** Configure dependabot groups ([#25](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/25)) ([fbf48c8](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/fbf48c89a5daf1b90e511a4d5ceb4e1591b1b339))  - (dstrates)
- **deps:** Enable dependabot updates ([#19](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/19)) ([7d6790e](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/7d6790e6a51f8ce7dd87d28a0146a200e8cdc876))  - (dstrates)

# [1.0.0](https://github.com/cultureamp/terraform-provider-schemaregistry/tree/v1.0.0) - (2024-07-25)

## <!-- 0 -->🚀 Features

- Initialize repository and implement provider ([#2](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/2)) ([43d407c](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/43d407c497f50dabc1f07251bd7ab8bdd37a6522))  - (dstrates)

## <!-- 1 -->🐛 Bug Fixes

- Improve support for schema references ([#8](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/8)) ([2e2fb02](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/2e2fb023000b16282720ee1e06a983ddb10ef667))  - (dstrates)

## <!-- 14 -->🎉 Initial Commit

- Initial commit ([e237b95](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/e237b95929356f980428e98eb7021786d6136c7d))  - (dstrates)

## <!-- 2 -->🚜 Refactor

- Rename provider from `schema-registry` to `schemaregistry` ([#18](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/18)) ([8e2121c](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/8e2121c4921eeb86e1111608c5fb8a39534bb596))  - (dstrates)
- Remove specific version reference in goreleaser ([#13](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/13)) ([0cc8fa7](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/0cc8fa7de2cc7ffd6698e64fb9b07558ed7bd0f1))  - (Elliot Schot)
- Rename `kafka-schema-registry` to `schema-registry` as the provider is kafka agnostic ([#12](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/12)) ([ffd72af](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/ffd72af7b575e1d04ca53d7612d78b26ee80329f))  - (Elliot Schot)

## <!-- 29 -->👷 CI/CD

- Update release workflow triggers ([#14](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/14)) ([c80d39d](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/c80d39d35b3f945295db6211bad8615868123b08))  - (dstrates)
- Fix go-semantic-release version ([#10](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/10)) ([61f7b1e](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/61f7b1e955b13f82bede856e0da79107e252cb68))  - (dstrates)
- Fix go-semantic-release version ([#9](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/9)) ([4759ad4](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/4759ad4d7fdaecb512c5f1fb1576ab985efd586b))  - (dstrates)
- Add semantic-release-go, scorecard and changelog automation ([#7](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/7)) ([1461bfd](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/1461bfd81c7cf2ac2d2140bbdd9d7887311f9fab))  - (dstrates)
- Update release and remove generate docs ([#3](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/3)) ([ea79869](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/ea79869cf84776cd66efe75d1adf16e06743eafb))  - (dstrates)

## <!-- 3 -->📚 Documentation

- Create SECURITY.md ([#16](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/16)) ([b3be8b2](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/b3be8b26b2afc6658222f8dd1e4f094d435cbc64))  - (dstrates)

## <!-- 7 -->⚙️ Miscellaneous Tasks

- **deps:** Update go toolchain version to `1.N.P` ([#17](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/17)) ([14c6cf8](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/14c6cf840bfa76439f34c879a0ea15fa4eb68f5c))  - (dstrates)
- **deps:** Replace renovate with dependabot ([#15](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/15)) ([f017e2c](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/f017e2ce7f6bd7a870efd1cd6bc9cfb50c51da6b))  - (dstrates)
- Update Makefile and README ([#11](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/11)) ([abcb3c0](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/abcb3c0bac6896e4a9b259f94653b26d639a13ee))  - (dstrates)
- Add validation for schema registry URL ([#6](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/6)) ([761cacb](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/761cacb9346b8c70ce0945ff988139d13364c2c3))  - (dstrates)
- Update descriptions and errors for clarity ([#5](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/5)) ([16429d5](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/16429d58062fa6f6c6c4b11cdb0d38dc65d04f5c))  - (dstrates)
- Fix typo in workflow, redefine variables and update error message ([#4](https://github.com/cultureamp/terraform-provider-schemaregistry/issues/4)) ([d4930fd](https://github.com/cultureamp/terraform-provider-schemaregistry/commit/d4930fdc953d329cc8c60c940eb5445c73af43cb))  - (dstrates)

<!-- generated by git-cliff -->
