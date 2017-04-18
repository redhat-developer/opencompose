# Change Log

## [v0.1.0](https://github.com/redhat-developer/opencompose/tree/v0.1.0) (2017-04-18)
[Full Changelog](https://github.com/redhat-developer/opencompose/compare/v0.1.0-alpha.0...v0.1.0)

**Closed issues:**

- opencompose generates deployment per container [\#86](https://github.com/redhat-developer/opencompose/issues/86)
- Change the file naming convention to \<service\_name\>-\<object\_type\> [\#65](https://github.com/redhat-developer/opencompose/issues/65)
- `opencompose convert -f -` should accept input from STDIN? [\#63](https://github.com/redhat-developer/opencompose/issues/63)
- opencompose version shows as -not\_a\_git\_tree, is this normal? [\#60](https://github.com/redhat-developer/opencompose/issues/60)
- Better naming for generated files [\#58](https://github.com/redhat-developer/opencompose/issues/58)
- No output on successul conversion [\#57](https://github.com/redhat-developer/opencompose/issues/57)
- Allow specifying file as URL [\#55](https://github.com/redhat-developer/opencompose/issues/55)
- Odd spacing on `opencompose validate --help` [\#38](https://github.com/redhat-developer/opencompose/issues/38)
- Odd version number [\#37](https://github.com/redhat-developer/opencompose/issues/37)
- Tooling for releases [\#27](https://github.com/redhat-developer/opencompose/issues/27)
- Ingress support [\#21](https://github.com/redhat-developer/opencompose/issues/21)
- Allow specifying replica count for service [\#19](https://github.com/redhat-developer/opencompose/issues/19)

**Merged pull requests:**

- fix replica test name string from "Valid" to "Invalid" [\#102](https://github.com/redhat-developer/opencompose/pull/102) ([containscafeine](https://github.com/containscafeine))
- Minor change in the description of a test case [\#100](https://github.com/redhat-developer/opencompose/pull/100) ([containscafeine](https://github.com/containscafeine))
- Write information to user about artifacts \(files\) being created [\#89](https://github.com/redhat-developer/opencompose/pull/89) ([tnozicka](https://github.com/tnozicka))
- Generate deployment only once [\#87](https://github.com/redhat-developer/opencompose/pull/87) ([surajssd](https://github.com/surajssd))
- Make Wordpress accessible by a loadbalancer [\#74](https://github.com/redhat-developer/opencompose/pull/74) ([tnozicka](https://github.com/tnozicka))
- add support for passing in input from STDIN [\#72](https://github.com/redhat-developer/opencompose/pull/72) ([containscafeine](https://github.com/containscafeine))
- Fix generated file names [\#70](https://github.com/redhat-developer/opencompose/pull/70) ([tnozicka](https://github.com/tnozicka))
- Add ingress support [\#69](https://github.com/redhat-developer/opencompose/pull/69) ([tnozicka](https://github.com/tnozicka))
- Spell correction in one of the errors [\#67](https://github.com/redhat-developer/opencompose/pull/67) ([surajssd](https://github.com/surajssd))
- Added support for providing replicas [\#66](https://github.com/redhat-developer/opencompose/pull/66) ([surajssd](https://github.com/surajssd))
- add URL support [\#64](https://github.com/redhat-developer/opencompose/pull/64) ([containscafeine](https://github.com/containscafeine))
- Updated .gitignore file [\#59](https://github.com/redhat-developer/opencompose/pull/59) ([surajssd](https://github.com/surajssd))
- Test for more Golang versions and master [\#56](https://github.com/redhat-developer/opencompose/pull/56) ([tnozicka](https://github.com/tnozicka))
- Two spaces not tabs for help examples [\#49](https://github.com/redhat-developer/opencompose/pull/49) ([zoidbergwill](https://github.com/zoidbergwill))
- Add Community section to README [\#46](https://github.com/redhat-developer/opencompose/pull/46) ([tnozicka](https://github.com/tnozicka))
- Clean up installation way for linux releases in README [\#42](https://github.com/redhat-developer/opencompose/pull/42) ([tnozicka](https://github.com/tnozicka))
- Add documentation links to README [\#41](https://github.com/redhat-developer/opencompose/pull/41) ([tnozicka](https://github.com/tnozicka))
- Fix up README.md [\#40](https://github.com/redhat-developer/opencompose/pull/40) ([cdrage](https://github.com/cdrage))
- Update the supported installation way in README [\#39](https://github.com/redhat-developer/opencompose/pull/39) ([tnozicka](https://github.com/tnozicka))
- Update README after our first release [\#35](https://github.com/redhat-developer/opencompose/pull/35) ([tnozicka](https://github.com/tnozicka))

## [v0.1.0-alpha.0](https://github.com/redhat-developer/opencompose/tree/v0.1.0-alpha.0) (2017-02-24)
**Closed issues:**

- Zsh completion support [\#28](https://github.com/redhat-developer/opencompose/issues/28)
- Extend generated service options [\#20](https://github.com/redhat-developer/opencompose/issues/20)
- make test runs tests for packages in vendor [\#11](https://github.com/redhat-developer/opencompose/issues/11)
- add README.md [\#4](https://github.com/redhat-developer/opencompose/issues/4)

**Merged pull requests:**

- Add release tooling [\#34](https://github.com/redhat-developer/opencompose/pull/34) ([tnozicka](https://github.com/tnozicka))
- Update go get command [\#33](https://github.com/redhat-developer/opencompose/pull/33) ([cdrage](https://github.com/cdrage))
- Add Zsh completion [\#32](https://github.com/redhat-developer/opencompose/pull/32) ([tnozicka](https://github.com/tnozicka))
- Extend generated service options to cover internal and external types [\#31](https://github.com/redhat-developer/opencompose/pull/31) ([tnozicka](https://github.com/tnozicka))
- fix extension for yaml files [\#30](https://github.com/redhat-developer/opencompose/pull/30) ([kadel](https://github.com/kadel))
- fix forgotten HostPort [\#26](https://github.com/redhat-developer/opencompose/pull/26) ([kadel](https://github.com/kadel))
- Remove hostport from opencompose definitions [\#18](https://github.com/redhat-developer/opencompose/pull/18) ([tnozicka](https://github.com/tnozicka))
- Create services and deployments [\#17](https://github.com/redhat-developer/opencompose/pull/17) ([tnozicka](https://github.com/tnozicka))
- update examples [\#16](https://github.com/redhat-developer/opencompose/pull/16) ([kadel](https://github.com/kadel))
- OpenCompose file structure documentation [\#15](https://github.com/redhat-developer/opencompose/pull/15) ([kadel](https://github.com/kadel))
- Fix parser to conform to spec 0.1-dev [\#14](https://github.com/redhat-developer/opencompose/pull/14) ([tnozicka](https://github.com/tnozicka))
- Fix convert [\#13](https://github.com/redhat-developer/opencompose/pull/13) ([tnozicka](https://github.com/tnozicka))
- Run tests only for our code, not vendor [\#12](https://github.com/redhat-developer/opencompose/pull/12) ([tnozicka](https://github.com/tnozicka))
- Add readme [\#10](https://github.com/redhat-developer/opencompose/pull/10) ([pradeepto](https://github.com/pradeepto))
- Add Travis CI tests [\#9](https://github.com/redhat-developer/opencompose/pull/9) ([tnozicka](https://github.com/tnozicka))
- Fix Errorf format in encoding [\#8](https://github.com/redhat-developer/opencompose/pull/8) ([tnozicka](https://github.com/tnozicka))
- Specify /bin/bash as the primary shell for Makefile. This fixes array… [\#7](https://github.com/redhat-developer/opencompose/pull/7) ([tnozicka](https://github.com/tnozicka))
- Fix formatting [\#6](https://github.com/redhat-developer/opencompose/pull/6) ([tnozicka](https://github.com/tnozicka))
- Add formatting directives to Makefile [\#5](https://github.com/redhat-developer/opencompose/pull/5) ([tnozicka](https://github.com/tnozicka))
- rename imports \(tnozicka-\>redhat-developer\) [\#3](https://github.com/redhat-developer/opencompose/pull/3) ([kadel](https://github.com/kadel))
- Add vendoring [\#2](https://github.com/redhat-developer/opencompose/pull/2) ([tnozicka](https://github.com/tnozicka))
- update example [\#1](https://github.com/redhat-developer/opencompose/pull/1) ([kadel](https://github.com/kadel))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*