# zscripts

This directory contains scripts to build and manage the ezpkg.io project. It's intended to be used internally by the author. The code may be changed at any time.

### Project Structure

The project are organized in the following structure:

```
ezpkg_root/
  ├─ ezpkg/               # github.com/ezpkg/ezpkg
  │   ├─ errorz/          # source code and tests for errorz package
  │   ├─ stringz/         # source code and tests for stringz package
  │   ├─ ...              # ...
  │   ├─ zscripts/        # scripts and tools for managing ezpkg.io
  │   │   └─ ezrun        # the main code for building & managing ezpkg.io
  │   └─ run              # the entrypoint script, which wraps zscripts/ezrun
  ├── ztarget/            # root directory for generated packages
  │   ├─ errorz/          # generated errorz package (without tests)
  │   ├─ stringz/         # generated stringz package (without tests)
 ...  └─ ...              # ...
```
