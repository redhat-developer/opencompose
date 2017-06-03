# Development Guide

## Building OpenCompose

Read about building OpenCompose [here](https://github.com/redhat-developer/opencompose#go-from-source).

## Workflow
### Fork the main repository

1. Go to https://github.com/redhat-developer/opencompose
2. Click the "Fork" button (at the top right)

### Clone your fork

The commands below require that you have $GOPATH. We highly recommended you put OpenCompose code into your $GOPATH.

```console
git clone https://github.com/$YOUR_GITHUB_USERNAME/opencompose $GOPATH/src/github.com/redhat-developer/opencompose
cd $GOPATH/src/github.com/redhat-developer/opencompose
git remote add upstream 'https://github.com/redhat-developer/opencompose'
```

### Create a branch and make changes

```console
git checkout -b myfeature
# Make your code changes
```

### Keeping your development fork in sync

```console
git fetch upstream
git rebase upstream/master
```

Note: If you have write access to the main repository at github.com/redhat-developer/opencompose, you should modify your git configuration so that you can't accidentally push to upstream:

```console
git remote set-url --push upstream no_push
```

### Committing changes to your fork

```console
git commit
git push -f origin myfeature
```

### Creating a pull request

1. Visit https://github.com/$YOUR_GITHUB_USERNAME/opencompose
2. Click the "Compare and pull request" button next to your "myfeature" branch.
3. Check out the pull request process for more details

## `glide`, `glide-vc` and dependency management

OpenCompose uses `glide` to manage dependencies and `glide-vc` to clean vendor directory.
They are not strictly required for building OpenCompose but they are required when managing dependencies under the `vendor/` directory.
If you want to make changes to dependencies please make sure that `glide` and `glide-vc` are installed and are in your `$PATH`.

### Installing glide

There are many ways to build and host golang binaries. Here is an easy way to get utilities like `glide` and `glide-vc` installed:

Ensure that Mercurial and Git are installed on your system. (some of the dependencies use the mercurial source control system).
Use `apt-get install mercurial git` or `yum install mercurial git` on Linux, or `brew.sh` on OS X, or download them directly.

```console
go get -u github.com/Masterminds/glide
go get github.com/sgotti/glide-vc
```

Check that `glide` and `glide-vc` commands are working.

```console
glide --version
glide-vc -h
```

### Using glide

#### Adding new dependency
1. Update `glide.yml` file.

  Add new packages or subpackages to `glide.yml` depending if you added whole new package as dependency or
  just new subpackage.

2. Run `make recreate-vendor` to get new dependencies.
 then run `make strip-vendor` to delete all unnecessary files from vendor.

3. Commit updated `glide.yml`, `glide.lock` and `vendor` to git.


#### Updating dependencies

1. Set new package version in  `glide.yml` file.

2. Run `make update-vendor` to update dependencies.
   Then run `make strip-vendor` to delete all unnecessary files from vendor.

