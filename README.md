# Image Saving Service

This application provides an HTTP endpoint to save multiple images from given URLs into a specified directory structure. It is built using [Go](https://golang.org/) and the [Gin](https://github.com/gin-gonic/gin) web framework.

## Overview

The service exposes a `POST /save-images` endpoint. You can send a JSON payload containing an array of image information objects. Each object specifies an image source URL and a desired folder path relative to the `dist` directory. The application will then download each image and save it into the appropriate folder.

By default, the application listens on port `8080`.

## Features

- **Concurrent Downloads**: The service processes up to 25 images concurrently for faster handling.
- **Customizable File Structure**: Each request can target different folders and image names.
- **Automatic Directory Creation**: If a specified folder does not exist, it will be created automatically within the `dist` directory.
- **Graceful Error Handling**: Errors are logged, but the endpoint will continue processing other images.

## Requirements

- **Go 1.18+** (or later)
- A working internet connection (to download images from the provided URLs)
- Network permissions to run on the specified port (default: `8080`)

## Installation

1. **Clone the repository**:
    ```bash
    git clone <REPO_URL>
    ```
   
2. **Navigate to the project directory**:
    ```bash
    cd <project-directory>
    ```

3. **Download dependencies**:
    ```bash
    go mod tidy
    ```

4. **Build the application**:
    ```bash
    go build -o image-downloader
    ```

## Running the Service

Once built, simply run the executable:

```bash
./image-downloader
```

This will start the server on `http://localhost:8080`.

## Usage Example

To save images, send a `POST` request to `http://localhost:8080/save-images` with a JSON body. For example, using `curl`:

```bash
curl -X POST http://localhost:8080/save-images \
  -H "Content-Type: application/json" \
  -d '[
    {
      "folderPath": "my-images/animals",
      "imageURL": "https://example.com/cat.jpg"
    },
    {
      "folderPath": "my-images/landscapes",
      "imageURL": "https://example.com/mountain.png",
      "imageName": "mountain_view.png"
    }
  ]'
```

### Request Body Fields

- **folderPath** (string, required): The relative path from `dist` where the image will be saved. E.g., `my-images/animals`.
- **imageURL** (string, required): The URL of the image to be downloaded.
- **imageName** (string, optional): The desired filename for the saved image. If not provided, the file name will be extracted from the URL.

## Directory Structure

After running the service and making requests, the directory structure might look like this:

```
dist/
├── my-images/
│   ├── animals/
│   │   └── cat.jpg
│   └── landscapes/
│       └── mountain_view.png
```
