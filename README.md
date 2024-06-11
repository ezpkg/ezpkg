<p align="center">
<a href="https://ezpkg.io">
<img alt="gopherz" src="./_/gopherz.png" style="width:420px">
</a>
</p>

# ezpkg.io

Collection of packages and tools to make writing Go code easier.

## Introduction

As I work on various Go projects, I often find myself creating utility functions, extending existing packages, or developing packages to solve specific problems. Moving from one project to another, I usually have to copy or rewrite these solutions. So I created this repository to have all these utilities and packages in one place. Hopefully, you'll find them useful as well.

These packages aim to enhance the functionality of the standard library and other popular packages. They are intended to be used together with other packages rather than replacing them. The APIs are designed based on my experience working with Go, focusing on simplicity and ease of use. I will try to follow best practices in Go, but not always. I also tend to choose a more performance implementation if possible.

If you have any questions or suggestions, feel free to [create new issues](https://github.com/ezpkg/ezpkg/issues), [open pull requests](https://github.com/ezpkg/ezpkg/pulls), or [post discussions](https://github.com/ezpkg/ezpkg/discussions).

## Packages

| Package                                     | Description                                                                                   |
|---------------------------------------------|-----------------------------------------------------------------------------------------------|
| [bytez](https://ezpkg.io/bytez)             | Extends the standard library [bytes](https://pkg.go.dev/bytes) with additional functions.     |
| [colorz](https://ezpkg.io/colorz)           | Working with colors in terminal.                                                              |
| [diffz](https://ezpkg.io/diffz)             | Comparing and displaying differences between two strings.                                     |
| [errorz](https://ezpkg.io/errorz)           | Utilities for working with errors: Must, Validate, multi-errors, etc.                         |
| [mapz](https://ezpkg.io/mapz)               | Extends the package [golang.org/x/exp/maps](https://pkg.go.dev/golang.org/x/exp/maps).        |
| [slicez](https://ezpkg.io/slicez)           | Extends the standard library [slices](https://pkg.go.dev/slices) with additional functions.   |
| [stacktracez](https://ezpkg.io/stacktracez) | Get stack trace for using in errors and logs.                                                 |
| [stringz](https://ezpkg.io/stringz)         | Extends the standard library [strings](https://pkg.go.dev/strings) with additional functions. |
| [testingz](https://ezpkg.io/testingz)       | Utilities for testing.                                                                        |
| [typez](https://ezpkg.io/typez)             | Generic functions for working with types.                                                     |
| [unsafez](https://ezpkg.io/unsafez)         | Convert bytes to strings using unsafe pointer.                                                |

### Installation

Download each package into your project using `go get`:

```sh
go get -u ezpkg.io/errorz
go get -u ezpkg.io/stringz
...
```

## Versioning

All packages are released together with the same version number to simplify management, as they often call each other. When the API evolves, the version number is incremented for all packages.

- Minor version **v0.1.0** → **v0.2.0**: Add new features or functions. May introduce a couple of (breaking) API changes.
- Patch version **v0.1.0** → **v0.1.1**: Fix bugs or improve performance. May add new (small) functions.

**Release Frequency**: New versions are released periodically as new features are developed or bugs are fixed. All packages will have their versions updated at the same time.

_**NOTE:** This project is in its early stages. While all packages are usable, the API may change over time._

## Project Structure

### ezpkg/ezpkg

The main repository that contains source code and tests of all packages.

### ezpkg/PKGNAME

Each package code is copied from the main repository to its own repository. This setup allows for importing only the necessary packages into a project. To minimize dependencies, the tests are removed, ensuring that users only need to download the package code without including the testing packages.

## Build and Test

### Setup local environment

- Create a root directory `ezpkg_root`.
- Clone this repository to `ezpkg_root/ezpkg`.
- Install [direnv](https://direnv.net/) to load environment variables from `.envrc`.

```sh
mkdir ezpkg_root      # create a root directory
git clone https://github.com/ezpkg/ezpkg ezpkg_root/ezpkg
cd ezpkg_root/ezpkg
direnv allow
```

### Test all packages

```sh
cd ezpkg_root/ezpkg
run testall
```

### Generate all packages

The following commands will copy each package code to the `ezpkg_root/ztarget` directory.

```sh
cd ezpkg_root/ezpkg
run pkgall
```

## FAQ

### Why should I NOT use these packages?

- **More dependencies**: These packages will add more dependencies to your project.
- **Early development**: This project is in its early stages and will have API changes. There are other packages that are more mature and offer more features.
- **Customization**: Sometimes, writing your own code allows for better customization. You can also copy code from these packages and modify it to fit your specific needs.

### Why should I use these packages?

- You find yourself copying the same code over and over.
- You are starting a new project and want some simple and easy-to-use packages.
- You are learning Go and want to see how some common tasks are implemented.

### Where are the tests?

- The tests are in the main repository `ezpkg/ezpkg`. You can run the tests for all packages using the `run testall` command.
- In each package repository, the tests are removed to minimize dependencies.

## Author

<a href="https://olivernguyen.io"><img alt="olivernguyen.io" src="https://olivernguyen.io/_/badge.png" height="28px"></a>&nbsp;&nbsp;[![github](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/iOliverNguyen)
