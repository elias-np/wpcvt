# webpcvt

webpcvt is a command-line tool written in Go that converts common image formats, such as JPG and PNG, into WebP. The idea behind it is fairly simple: offer a fast and direct way of doing that conversion without needing to open an image editor or depend on some online service to get it done.

The project came out of a practical, everyday need, since converting images to WebP tends to be a recurring task for anyone working on website performance, given that WebP usually produces a much smaller file while keeping an acceptable visual quality, and having that available straight from the terminal, integrated into scripts or into any build pipeline, ends up saving a fair amount of time compared to doing that process by hand, one image at a time, through a GUI.

What makes webpcvt worth using, beyond just being "another WebP converter," comes down to a few practical points:

- **Simple** - no complicated flags to memorize, no configuration files, nothing to set up first. You point it at an image and it converts.
- **Portable** - it ships as a single binary, so there's nothing else to install on the target machine. No runtime, no separate libraries to pull in first. Download it, drop it somewhere in your `PATH`, and it works.
- **Fast** - being written in Go, it starts instantly and processes images without the overhead of a scripting runtime, which matters when it's called dozens or hundreds of times inside a batch job or a CI pipeline.
- **Practical** - sensible defaults (a default quality of 80, automatic output naming) mean the common case, "just convert this file," takes a single short command, without having to think about it.

## How it works

The basic syntax planned for the first version is the following.

```bash
webpcvt image.jpg -q 85 output.webp
```

In this example, `image.jpg` is the input file that will be converted, the `-q` flag sets the compression quality applied while generating the WebP file, accepting a numeric value that usually ranges from 0 to 100, where the higher the number, the better the final quality and, consequently, the larger the resulting file, and `output.webp` is the name of the output file, which is an optional argument.

If the output file is not given when calling the command, webpcvt will automatically assume the same name as the input file, just swapping the extension to `.webp`, so running `webpcvt image.jpg -q 85` without specifying anything after the quality flag will generate a file called `image.webp` in the same folder, which makes everyday use much leaner when you just want to convert an image quickly without worrying about typing an output name every single time.

## Project status

This project is in its early stages of development, so the functionality described here represents the initial goal and not necessarily everything that is already implemented and working in the repository, and the idea is to evolve it gradually, starting with basic conversion of a single file and then possibly expanding into other directions, such as batch conversion of entire directories, support for other input formats besides JPG and PNG, and other flags that make sense for the day-to-day workflow of whoever ends up using the tool.

The project is open source and the license is still to be decided, but the intention is to keep the code open so that anyone can use it, study it, suggest improvements or adapt it to their own workflow, and contributions are welcome once the initial structure of the repository is a bit more mature.

## What's included

- Lossy conversion with adjustable quality (`-q`). Default quality is 80 when `-q` is omitted.
- `-v` flag to print the installed version.
- Automatic output naming (changes the extension to `.webp` if no output name is provided).

## Installation

### Windows (64-bit)

1. Download `webpcvt-windows-amd64.zip`.
2. Extract `webpcvt.exe` and place it in a folder that is in your `PATH` (e.g. `C:\webpcvt`).
3. Open a new terminal and run `webpcvt -v` to verify.

### Linux (64-bit)

1. Download `webpcvt-linux-amd64.tar.gz`.
2. Extract and install:

```bash
tar -xzf webpcvt-linux-amd64.tar.gz
sudo mv webpcvt /usr/local/bin/
sudo chmod +x /usr/local/bin/webpcvt
```

3. Run `webpcvt -v` to verify.

## Usage

```bash
webpcvt input.png
```
Converts using the default quality (80) and saves the result as `input.webp` in the same folder.

```bash
webpcvt input.png output.webp
```
Converts and saves the result with a custom output name or path.

```bash
webpcvt input.png -q 90
```
Converts using a custom quality (0–100). This can be combined with a custom output name:

```bash
webpcvt input.jpg output.webp -q 50
```
