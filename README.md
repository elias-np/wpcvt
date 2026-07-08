# wpcvt

wpcvt is a command line tool written in Go that converts common image formats, such as jpg and png, into webp, and the idea behind it is fairly simple, which is to offer a fast and direct way of doing that conversion without needing to open any image editor or depend on some online service to get it done. The project came out of a practical everyday need, since converting images to webp tends to be a recurring task for anyone working on website performance, given that webp usually produces a much smaller file while keeping an acceptable visual quality, and having that available straight from the terminal, integrated into scripts or into any build pipeline, ends up saving a fair amount of time compared to doing that process by hand.

The project is open source and the license is still to be decided, but the intention is to keep the code open so that anyone can use it, study it, suggest improvements or adapt it to their own workflow, and contributions are welcome once the initial structure of the repository is a bit more mature.

## How it works

The basic syntax planned for the first version is the following.

```
wpcvt image.jpg -q 85 output.webp
```

In this example, image.jpg is the input file that will be converted, the -q flag sets the compression quality applied while generating the webp, accepting a numeric value that usually ranges from 0 to 100, where the higher the number, the better the final quality and, consequently, the larger the resulting file, and output.webp is the name of the output file, which is an optional argument.

If the output file is not given when calling the command, wpcvt will automatically assume the same name as the input file, just swapping the extension to webp, so running wpcvt image.jpg -q 85 without specifying anything after the quality flag will generate a file called image.webp in the same folder, which makes everyday use much leaner when you just want to convert an image quickly without worrying about typing an output name every single time.

## Project status

This project is in its early stages of development, so the functionality described here represents the initial goal and not necessarily what is already implemented and working in the repository, and the idea is to evolve it gradually, starting with basic conversion of a single file and then possibly expanding into other possibilities, such as batch conversion of entire directories, support for other input formats besides jpg and png, and other flags that make sense for the workflow of whoever ends up using the tool day to day.
