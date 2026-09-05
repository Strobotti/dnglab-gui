# DNGLab GUI

A graphical user interface for dnglab — converts camera RAW files to Adobe DNG format.

## Prerequisites

- **Go 1.21 or later** — [https://go.dev/dl/](https://go.dev/dl/)
- **C compiler** — required by the [Fyne](https://fyne.io/) UI toolkit:
  - Linux: `gcc` (e.g. `sudo apt install gcc` on Debian/Ubuntu)
  - macOS: Xcode Command Line Tools (`xcode-select --install`)
- **dnglab CLI** — [https://github.com/dnglab/dnglab](https://github.com/dnglab/dnglab)

  The application searches for the `dnglab` binary in this order:
  1. The same directory as the `dnglab-gui` binary
  2. The system `PATH`

## Build

```sh
go build -o dnglab-gui ./cmd/dnglab-gui
```

## Run

Development (no prior build required):

```sh
go run ./cmd/dnglab-gui
```

After building:

```sh
./dnglab-gui
```

## Usage

The main window is divided into four sections:

### 1 Select images to convert

Pick a source folder containing RAW files. Enable **Include subfolders** to scan recursively. Enable **Skip existing DNGs** to skip files that already have a converted `.dng` counterpart in the output location.

### 2 Select location to save converted images

Choose where the DNG files are written. By default the converted files are saved to the same folder as the originals. Select a different output folder when you want to keep originals and conversions separate.

### 3 Conversion Options

| Option | Values |
|---|---|
| Compression | Lossless (default), Uncompressed |
| Crop mode | Best crop, Active area, None |
| Embed RAW | On / Off |
| Preview | On / Off |
| Thumbnail | On / Off |

### 4 Metadata

- **Artist** — written into the `Artist` DNG metadata field.

## Supported RAW Formats

| Manufacturer | Extensions |
|---|---|
| Canon | CR3, CR2, CRW |
| Nikon | NEF, NRW |
| Sony | ARW, SRF, SR2 |
| Fujifilm | RAF |
| Panasonic / Leica | RW2 |
| Olympus | ORF |
| Pentax / Ricoh | PEF |
| Others | 3FR, ARI, ERF, KDC, DCS, DCR, IIQ, MOS, MEF, MRW, SRW |

The actual list of supported formats depends on the version of `dnglab` installed on your system.

## License

MIT License — Copyright 2026 Juha Jantunen
