# webpcvt

webpcvt is a command-line tool written in Go. It converts JPG and PNG images to WebP.

To be precise about what that means: webpcvt does not implement its own WebP encoder. It vendors libwebp, the official WebP library maintained by Google, and calls it directly to do the actual conversion. What webpcvt adds is the packaging around that. A single binary, sensible defaults, and a command-line interface built for one specific job: converting images to WebP from a terminal, a script, or a CI pipeline, without asking you to install or compile libwebp yourself first.

That distinction matters. The encoding logic, the compression algorithm, the quality behavior, all of that comes from libwebp itself. webpcvt's job is to make that library convenient to reach for in a very specific context: everyday, script-friendly, one-command image conversion.

## Why this exists

Converting images to WebP is a recurring task for anyone working on website performance. WebP usually produces a much smaller file while keeping acceptable visual quality.  Having it available straight from the terminal, wired into scripts or a build pipeline, saves real time.

## What webpcvt actually does

- **Wraps libwebp.** The conversion itself is handled by libwebp, vendored directly into the binary. webpcvt does not reimplement WebP encoding; it exposes libwebp through a simple CLI.
- **Simple.** No complicated flags to memorize, no configuration files, nothing to set up first. Point it at an image and it converts.
- **Portable.** Because libwebp is vendored in, the result is a single binary. Nothing else to install on the target machine, no separate libwebp package, no runtime.
- **Fast.** Go starts instantly, and libwebp does the encoding at native speed. There's no scripting-runtime overhead, which matters when the tool runs dozens or hundreds of times inside a batch job.
- **Practical.** Sensible defaults (quality 80, automatic output naming) mean the common case takes one short command.

## How it works

```bash
webpcvt image.jpg -q 85 output.webp
```

`image.jpg` is the input file. The `-q` flag sets the quality passed straight to libwebp's encoder, a number from 0 to 100. Higher means better quality and a bigger file. `output.webp` is optional.

If you skip the output name, webpcvt reuses the input name and swaps the extension to `.webp`. So `webpcvt image.jpg -q 85` on its own produces `image.webp` in the same folder. One less thing to type every time.

## Project status

The plan is to grow it step by step: single-file conversion first, then batch conversion of whole directories, support for more input formats, and other flags that make sense day to day.

## What's included

- Lossy conversion through libwebp, with adjustable quality (`-q`). Default is 80 when `-q` is omitted.
- Automatic output naming (swaps the extension to `.webp` when no output name is given).
- Batch conversion of a whole directory (`-r`/`--recursive` to include subdirectories), converting matching files concurrently.

## Installation

### Windows (64-bit)

1. Download [webpcvt-windows-amd64.zip](https://github.com/elias-np/webpcvt/releases/download/v0.1.0/webpcvt-windows-amd64.zip).
2. Extract webpcvt.exe to a folder of your choice (e.g. C:\webpcvt) and add that folder to your PATH.
3. Open a new terminal and run `webpcvt -v` to check it's working.

### Linux (64-bit)

1. Download [webpcvt-linux-amd64.tar.gz](https://github.com/elias-np/webpcvt/releases/download/v0.1.0/webpcvt-linux-amd64.tar.gz)
2. Extract and install:

```bash
tar -xzf webpcvt-linux-amd64.tar.gz
sudo mv webpcvt /usr/local/bin/
sudo chmod +x /usr/local/bin/webpcvt
```

3. Run `webpcvt -v` to check it's working.

Because libwebp is vendored into the binary, no extra install step is needed for it. Whatever platform build you download already has it built in.

## Usage

```bash
webpcvt input.png
```
Converts with the default quality (80) and saves the result as `input.webp` in the same folder.

```bash
webpcvt input.png output.webp
```
Converts and saves the result with a custom output name or path.

```bash
webpcvt input.png -q 90
```
Converts with a custom quality (0 to 100). Can be combined with a custom output name:

```bash
webpcvt input.jpg output.webp -q 50
```
<<<<<<< HEAD

```bash
webpcvt ./photos
```
Converts every `.jpg`, `.jpeg` and `.png` file directly inside `./photos`, each saved next to its original as `.webp`. Subdirectories are left untouched.

```bash
webpcvt ./photos -r
```
Same as above, but recurses into every subdirectory of `./photos`, converting each image next to itself.

```bash
webpcvt ./photos ./photos-webp -r -q 90
```
Recursively converts `./photos` into `./photos-webp` instead of in place. If images live in more than one subdirectory, webpcvt asks once whether to mirror `./photos`'s folder structure inside `./photos-webp` or flatten everything into a single folder.

If any target `.webp` file already exists, webpcvt asks once whether to skip those files or overwrite (reconvert) them, then converts the rest of the batch concurrently.
=======
>>>>>>> 74a3b73d289412e4b091244b872f33bfbd6ebc61
