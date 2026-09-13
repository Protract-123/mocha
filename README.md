# mocha

**AI Disclaimer**: I used AI for code review, and for minor things like porting tests from the original scoop. All code was reviewed by me, and a large majority written by me.

mocha is a [Scoop](https://github.com/ScoopInstaller/Scoop) alternative written in Go. I have the following goals for this project: 

- 100% backward compatibility with existing Scoop packages and buckets
- Act as a drop-in replacement for existing Scoop deployments
- Be faster than Scoop for all use cases
- Add new features at a fast pace
- Create a more maintainable codebase

> [!CAUTION]
> mocha is in **very** early development, and isn't suitable for users who want a stable experience.

## Installation

**Requirements:** Windows 10/11, [Git](https://git-scm.com/downloads)

mocha may need the following extra tools on your PATH:

- [7-Zip](https://www.7-zip.org/) (`7z.exe`) for archives other than `.zip`, like `.7z` and `.tar.gz`
- [innounp](https://innounp.sourceforge.net/) (`innounp.exe`) for apps packaged as InnoSetup installers
- Java, Python or PowerShell for `.jar`, `.py` and `.ps1` shims. The runner for each can be changed in the [configuration](#shim)

Some of these tools can be installed through mocha once it has been set up.

### Install from release

Prebuilt binaries are available on the [Releases](https://github.com/Protract-123/mocha/releases) page. Download the binary for your architecture and place it in a folder on your PATH.

### Install with Go

**Requirements:** [Go 1.27+](https://go.dev/doc/install), [Git](https://git-scm.com/downloads)

```shell
go install github.com/Protract-123/mocha@latest
```

This places `mocha.exe` in `$GOPATH\bin`, which is `%USERPROFILE%\go\bin` by default.

### Build from source

**Requirements:** [Go 1.27+](https://go.dev/doc/install), [Git](https://git-scm.com/downloads)

```shell
git clone https://github.com/Protract-123/mocha.git
cd mocha
go build -ldflags="-s -w" -trimpath .
```

Running the commands above will output a `mocha.exe` binary in the repository directory, which can then be moved anywhere.

## Getting Started

mocha keeps everything in `%USERPROFILE%\mocha` by default. To use a different folder, set the `MOCHA_DIR` environment variable 
to an absolute path.

### Set Up the Shim Binary

mocha uses a shim to put installed apps on your PATH, and won't install apps until one is set up. List the available releases, 
then pick a language (e.g. `zig` or `rust`). A comparison between the different languages can be found at [ScoopInstaller/Shim](https://github.com/ScoopInstaller/Shim).
You can provide a version if you want to install a specific one, or just provide a language for the latest version.

```shell
mocha shim binary releases
mocha shim binary setup <language> [version]
```

### Add Shims to PATH

To make installed apps accessible, add the shims folder (`%USERPROFILE%\mocha\shims` by default) to your user `PATH`
environment variable. Currently, mocha doesn't do this for you.

### Add a Bucket

Buckets from Scoop's known buckets list (see `mocha bucket known`) can be added by name, and any other bucket can be added with its git URL:

```shell
mocha bucket add main
mocha bucket add mybucket https://github.com/user/my-bucket
```

### Install an App

You can now search for an app in the added buckets through `mocha search`, and install an app through `mocha install`.

```shell
mocha search ripgrep
mocha install ripgrep
```

If an app adds to `PATH` or sets environment variables, the changes only apply to apps opened after the installation.

## Usage

```shell
mocha <command> [options]
```

Apps are written as `[bucket/]app[@version]`, e.g. `git`, `main/git` or `bat@0.26.1`. If no bucket is given, mocha uses 
the first bucket that has a manifest for the app.

At this time the `mocha install` command takes in the version parameter, but always installs the latest version. 

To view all available commands, run `mocha -h`. To see a command's subcommands and options, run `mocha <command> -h`.

## Directory Layout

mocha stores its files inside the mocha directory:

| Path                   | Contents                                                                    |
| ---------------------- | --------------------------------------------------------------------------- |
| `apps\<app>\<version>` | Installed files for each version of an app                                  |
| `apps\<app>\current`   | Junction to the active version, used by shims, shortcuts and `Path` entries |
| `persist\<app>`        | Data kept between versions of an app (e.g. config files)                    |
| `buckets\<bucket>`     | Cloned bucket repositories                                                  |
| `cache`                | Downloaded files                                                            |
| `shims`                | Shims for installed apps and ones created with `shim add`                   |
| `shim.exe`             | Shim binary set up with `shim binary setup`                                 |
| `known_buckets.json`   | Scoop's list of known buckets, refreshed by `mocha update`                  |

Start Menu shortcuts are created in `%APPDATA%\Microsoft\Windows\Start Menu\Programs\Mocha Apps`.

## Configuration

mocha is configured through a TOML file named `mocha.toml`. Running `mocha config` opens it in the editor set by `EDITOR` or `VISUAL`, 
or the default app for `.toml` files if neither is set. 

If no config file exists yet when running `mocha config`, mocha writes the default one to the mocha directory first. mocha uses the first config file it finds 
in these locations:

| Priority | Path                                     |
| -------- | ---------------------------------------- |
| 1        | `<mocha directory>\mocha.toml`           |
| 2        | `%APPDATA%\mocha\mocha.toml`             |
| 3        | `%XDG_CONFIG_HOME%\mocha\mocha.toml`     |
| 4        | `%USERPROFILE%\.config\mocha\mocha.toml` |

Any setting left out of the file uses its default value. The default configuration can be found in [`config/default.toml`](config/default.toml).

### Cat

| Field     | Description                                                                                                                                                                                                 |
| --------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `command` | Command used by `mocha cat` to show a manifest, run through `cmd.exe`. It must contain `[path]`, which is replaced with the manifest's path (e.g. `bat [path]`). Leave it empty to print the manifest as is |

### Colors

| Field     | Default   | Description                   |
| --------- | --------- | ----------------------------- |
| `success` | `magenta` | Color of success messages     |
| `error`   | `red`     | Color of error messages       |
| `warning` | `yellow`  | Color of warning messages     |
| `info`    | `blue`    | Color of information messages |

The available colors are `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan` and `white`, plus `bright-red`, `bright-green`, `bright-yellow`, `bright-blue`, `bright-magenta`, `bright-cyan` and `bright-white`. If any color is invalid, mocha shows a warning and uses the default colors.

### Shim

| Field               | Default          | Description                      |
| ------------------- | ---------------- | -------------------------------- |
| `java-runner`       | `java.exe`       | Program used to run `.jar` shims |
| `powershell-runner` | `powershell.exe` | Program used to run `.ps1` shims |
| `python-runner`     | `python.exe`     | Program used to run `.py` shims  |

Each runner can be a name on PATH (e.g. `pwsh.exe`) or a full path. mocha looks up the runner when it creates a shim, so existing shims keep the old runner until they're recreated.

## License

mocha is released under the [MIT License](LICENSE).  
`fileops/junction.go` and `fileops/shortcut.go` are adapted from MIT licensed projects, and their license notices are included in those files.

## Acknowledgments

- [Scoop](https://github.com/ScoopInstaller/Scoop) - Inspired this project. mocha uses its manifest and bucket format, and its list of known buckets
- [ScoopInstaller/Shim](https://github.com/ScoopInstaller/Shim) - Provides the shim binaries
- [nyaosorg/go-windows-junction](https://github.com/nyaosorg/go-windows-junction) - Junction creation code
- [jxeng/shortcut](https://github.com/jxeng/shortcut) - Shortcut creation code
