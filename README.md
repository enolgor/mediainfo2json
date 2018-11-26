# MediaInfo2Json

MediaInfo2Json is a tool to format the output of the [mediainfo tool](https://mediaarea.net/es/MediaInfo) into JSON in a more friendly format

Usage: `mediainfo2json [OPTIONS] PATH`

Extract mediainfo of PATH in json format, where PATH can be a directory or a file. If PATH is a directory it will extract the mediainfo of each subfile.

Options:

 - `-recursive, -r`: If PATH is a directory, list files recursively
 - `-output, -o`: Specify output file to write, otherwise print to standard output

