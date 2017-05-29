# Change Log

## [v0.2.0](https://github.com/redhat-developer/opencompose/tree/v0.2.0) (2017-05-29)
[Full Changelog](https://github.com/redhat-developer/opencompose/compare/v0.1.0...v0.2.0)

**Closed issues:**

- In container support for specifying command [\#129](https://github.com/redhat-developer/opencompose/issues/129)
- Change syntax for referencing higher value fields [\#119](https://github.com/redhat-developer/opencompose/issues/119)
- Changing env definition [\#118](https://github.com/redhat-developer/opencompose/issues/118)
- check the gofmt and other go related errors for every PR [\#110](https://github.com/redhat-developer/opencompose/issues/110)
- sometimes test in pkg/util takes more time [\#94](https://github.com/redhat-developer/opencompose/issues/94)
- why is services a list/slice and not map [\#78](https://github.com/redhat-developer/opencompose/issues/78)
- Using "-" for stdout? [\#61](https://github.com/redhat-developer/opencompose/issues/61)
- Support for storage \(volume mounts\) [\#25](https://github.com/redhat-developer/opencompose/issues/25)
- Labels support [\#23](https://github.com/redhat-developer/opencompose/issues/23)

**Merged pull requests:**

- add validation for environment variables [\#140](https://github.com/redhat-developer/opencompose/pull/140) ([containscafeine](https://github.com/containscafeine))
- BREAKING CHANGE\(spec\): Change volumeName to volumeRef [\#136](https://github.com/redhat-developer/opencompose/pull/136) ([surajssd](https://github.com/surajssd))
- refactor env definition from a=b to name: a,value: b [\#121](https://github.com/redhat-developer/opencompose/pull/121) ([containscafeine](https://github.com/containscafeine))
- fix indentation in wordpress examples for env [\#117](https://github.com/redhat-developer/opencompose/pull/117) ([containscafeine](https://github.com/containscafeine))
- Update readme instructions regarding installation [\#111](https://github.com/redhat-developer/opencompose/pull/111) ([cdrage](https://github.com/cdrage))
- Add volumes implementation support [\#88](https://github.com/redhat-developer/opencompose/pull/88) ([surajssd](https://github.com/surajssd))
- implement service level labels [\#82](https://github.com/redhat-developer/opencompose/pull/82) ([containscafeine](https://github.com/containscafeine))
- Proposal for volume/storage support for opencompose [\#79](https://github.com/redhat-developer/opencompose/pull/79) ([surajssd](https://github.com/surajssd))

## [v0.1.0](https://github.com/redhat-developer/opencompose/tree/v0.1.0) (2017-04-19)
[Full Changelog](https://github.com/redhat-developer/opencompose/compare/v0.1.0-alpha.0...v0.1.0)

**Closed issues:**

- opencompose generates deployment per container [\#86](https://github.com/redhat-developer/opencompose/issues/86)
- Change the file naming convention to \<service\_name\>-\<object\_type\> [\#65](https://github.com/redhat-developer/opencompose/issues/65)
- opencompose version shows as -not\_a\_git\_tree, is this normal? [\#60](https://github.com/redhat-developer/opencompose/issues/60)
- Allow specifying file as URL [\#55](https://github.com/redhat-developer/opencompose/issues/55)
- Odd spacing on `opencompose validate --help` [\#38](https://github.com/redhat-developer/opencompose/issues/38)
- Odd version number [\#37](https://github.com/redhat-developer/opencompose/issues/37)
- Allow specifying replica count for service [\#19](https://github.com/redhat-developer/opencompose/issues/19)
- `opencompose convert -f -` should accept input from STDIN? [\#63](https://github.com/redhat-developer/opencompose/issues/63)
- Better naming for generated files [\#58](https://github.com/redhat-developer/opencompose/issues/58)
- No output on successul conversion [\#57](https://github.com/redhat-developer/opencompose/issues/57)
- Tooling for releases [\#27](https://github.com/redhat-developer/opencompose/issues/27)
- Ingress support [\#21](https://github.com/redhat-developer/opencompose/issues/21)
- Extend generated service options [\#20](https://github.com/redhat-developer/opencompose/issues/20)
- add README.md [\#4](https://github.com/redhat-developer/opencompose/issues/4)

**Merged pull requests:**

- 0.1.0 release [\#109](https://github.com/redhat-developer/opencompose/pull/109) ([cdrage](https://github.com/cdrage))
- Add release script [\#108](https://github.com/redhat-developer/opencompose/pull/108) ([cdrage](https://github.com/cdrage))
- fix replica test name string from "Valid" to "Invalid" [\#102](https://github.com/redhat-developer/opencompose/pull/102) ([containscafeine](https://github.com/containscafeine))
- Minor change in the description of a test case [\#100](https://github.com/redhat-developer/opencompose/pull/100) ([containscafeine](https://github.com/containscafeine))
- Generate deployment only once [\#87](https://github.com/redhat-developer/opencompose/pull/87) ([surajssd](https://github.com/surajssd))
- Make Wordpress accessible by a loadbalancer [\#74](https://github.com/redhat-developer/opencompose/pull/74) ([tnozicka](https://github.com/tnozicka))
- add support for passing in input from STDIN [\#72](https://github.com/redhat-developer/opencompose/pull/72) ([containscafeine](https://github.com/containscafeine))
- Spell correction in one of the errors [\#67](https://github.com/redhat-developer/opencompose/pull/67) ([surajssd](https://github.com/surajssd))
- Added support for providing replicas [\#66](https://github.com/redhat-developer/opencompose/pull/66) ([surajssd](https://github.com/surajssd))
- Updated .gitignore file [\#59](https://github.com/redhat-developer/opencompose/pull/59) ([surajssd](https://github.com/surajssd))
- Two spaces not tabs for help examples [\#49](https://github.com/redhat-developer/opencompose/pull/49) ([zoidbergwill](https://github.com/zoidbergwill))
- Add Community section to README [\#46](https://github.com/redhat-developer/opencompose/pull/46) ([tnozicka](https://github.com/tnozicka))
- Clean up installation way for linux releases in README [\#42](https://github.com/redhat-developer/opencompose/pull/42) ([tnozicka](https://github.com/tnozicka))
- Add documentation links to README [\#41](https://github.com/redhat-developer/opencompose/pull/41) ([tnozicka](https://github.com/tnozicka))
- Fix up README.md [\#40](https://github.com/redhat-developer/opencompose/pull/40) ([cdrage](https://github.com/cdrage))
- Update the supported installation way in README [\#39](https://github.com/redhat-developer/opencompose/pull/39) ([tnozicka](https://github.com/tnozicka))
- Update README after our first release [\#35](https://github.com/redhat-developer/opencompose/pull/35) ([tnozicka](https://github.com/tnozicka))
- Write information to user about artifacts \(files\) being created [\#89](https://github.com/redhat-developer/opencompose/pull/89) ([tnozicka](https://github.com/tnozicka))
- Fix generated file names [\#70](https://github.com/redhat-developer/opencompose/pull/70) ([tnozicka](https://github.com/tnozicka))
- Add ingress support [\#69](https://github.com/redhat-developer/opencompose/pull/69) ([tnozicka](https://github.com/tnozicka))
- add URL support [\#64](https://github.com/redhat-developer/opencompose/pull/64) ([containscafeine](https://github.com/containscafeine))
- Test for more Golang versions and master [\#56](https://github.com/redhat-developer/opencompose/pull/56) ([tnozicka](https://github.com/tnozicka))
- Add release tooling [\#34](https://github.com/redhat-developer/opencompose/pull/34) ([tnozicka](https://github.com/tnozicka))
- Extend generated service options to cover internal and external types [\#31](https://github.com/redhat-developer/opencompose/pull/31) ([tnozicka](https://github.com/tnozicka))
- OpenCompose file structure documentation [\#15](https://github.com/redhat-developer/opencompose/pull/15) ([kadel](https://github.com/kadel))
- Add readme [\#10](https://github.com/redhat-developer/opencompose/pull/10) ([pradeepto](https://github.com/pradeepto))

## [v0.1.0-alpha.0](https://github.com/redhat-developer/opencompose/tree/v0.1.0-alpha.0) (2017-02-24)
**Closed issues:**

- Zsh completion support [\#28](https://github.com/redhat-developer/opencompose/issues/28)
- make test runs tests for packages in vendor [\#11](https://github.com/redhat-developer/opencompose/issues/11)

**Merged pull requests:**

- Update go get command [\#33](https://github.com/redhat-developer/opencompose/pull/33) ([cdrage](https://github.com/cdrage))
- Add Zsh completion [\#32](https://github.com/redhat-developer/opencompose/pull/32) ([tnozicka](https://github.com/tnozicka))
- fix extension for yaml files [\#30](https://github.com/redhat-developer/opencompose/pull/30) ([kadel](https://github.com/kadel))
- fix forgotten HostPort [\#26](https://github.com/redhat-developer/opencompose/pull/26) ([kadel](https://github.com/kadel))
- Remove hostport from opencompose definitions [\#18](https://github.com/redhat-developer/opencompose/pull/18) ([tnozicka](https://github.com/tnozicka))
- Create services and deployments [\#17](https://github.com/redhat-developer/opencompose/pull/17) ([tnozicka](https://github.com/tnozicka))
- update examples [\#16](https://github.com/redhat-developer/opencompose/pull/16) ([kadel](https://github.com/kadel))
- Fix parser to conform to spec 0.1-dev [\#14](https://github.com/redhat-developer/opencompose/pull/14) ([tnozicka](https://github.com/tnozicka))
- Fix convert [\#13](https://github.com/redhat-developer/opencompose/pull/13) ([tnozicka](https://github.com/tnozicka))
- Run tests only for our code, not vendor [\#12](https://github.com/redhat-developer/opencompose/pull/12) ([tnozicka](https://github.com/tnozicka))
- Add Travis CI tests [\#9](https://github.com/redhat-developer/opencompose/pull/9) ([tnozicka](https://github.com/tnozicka))
- Fix Errorf format in encoding [\#8](https://github.com/redhat-developer/opencompose/pull/8) ([tnozicka](https://github.com/tnozicka))
- Specify /bin/bash as the primary shell for Makefile. This fixes array… [\#7](https://github.com/redhat-developer/opencompose/pull/7) ([tnozicka](https://github.com/tnozicka))
- Fix formatting [\#6](https://github.com/redhat-developer/opencompose/pull/6) ([tnozicka](https://github.com/tnozicka))
- Add formatting directives to Makefile [\#5](https://github.com/redhat-developer/opencompose/pull/5) ([tnozicka](https://github.com/tnozicka))
- rename imports \(tnozicka-\>redhat-developer\) [\#3](https://github.com/redhat-developer/opencompose/pull/3) ([kadel](https://github.com/kadel))
- Add vendoring [\#2](https://github.com/redhat-developer/opencompose/pull/2) ([tnozicka](https://github.com/tnozicka))
- update example [\#1](https://github.com/redhat-developer/opencompose/pull/1) ([kadel](https://github.com/kadel))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*