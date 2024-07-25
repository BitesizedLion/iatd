# Internet Archive Torrent Downloader (IATD)

Internet Archive Torrent Downloader (IATD) is a command-line tool designed to automate the download of torrent files from archive.org collections. Given a list of identifiers, IATD will concurrently download the torrent files.

IATD is useful for researchers, archivists, or anyone who needs to bulk download torrent files from archive.org collections.

## Features

- Concurrent downloading with configurable number of workers.
- Skips already existing torrent files to avoid duplicate downloads. (Only checks file name, not contents)
- Logs all failed downloads in `failed_downloads.log` (Currently not configurable)

## Installation

> [!IMPORTANT]
> You will need to have Golang installed.

1. Clone the repository:
    ```sh
    git clone https://github.com/BitesizedLion/iatd.git
    cd iatd
    ```

2. Build the project:
    ```sh
    go build .
    ```

OR:

Install with `go install`:
```sh
go install github.com/BitesizedLion/iatd
```

OR:

Download a build from the [releases](https://github.com/BitesizedLion/iatd/releases)

## Usage
> [!TIP]
> You can get this by using [Advanced Search](https://archive.org/advancedsearch.php) and selecting "CSV format"
1. Prepare a CSV file (e.g., `search.csv`) with the following format:
    ```
    identifier
    winampskins_Sunico
    winampskin_05PROAmp
    winampskin_08_Meaningless_Shallow
    ...
    ```

2. Run the downloader:
    ```sh
    ./iatd -csv search.csv -workers 10
    ```
    Optional flags:
    - `-csv`: Path to the CSV file containing the archive.org identifiers. Default is `search.csv`.
    - `-workers`: Number of concurrent workers. Default is `10`.

> All torrent files will be downloaded to `torrents/`, this is currently not configurable.
