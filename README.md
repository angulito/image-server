# Image Server

[![Build Status](https://magnum.travis-ci.com/wanelo/image-server.svg?token=xxYxjHDAXkDK41qZ1dqA&branch=master)](https://magnum.travis-ci.com/wanelo/image-server)

## Usage

The default port number is 7000

To start the server under development
```bash
make run
```

To run the compiled executable
```
./bin/image-server -e production
```

## Configuration

Example of json configuration on `config/production.json`
```json
{
  "server_port": "7000",
  "status_port": "7001",
  "source_domain": "http://cdn-s3-2.wanelo.com",
  "whitelisted_extensions": ["jpg", "png", "webp"],
  "maximum_width": 1000,
  "manta_base_path": "public/images/production",
  "manta_concurrency": 10,
  "default_quality": 75,
  "graphite_enabled": false,
  "graphite_port": 7777,
  "graphite_host": "localhost"
}

```

### Sample images

**Image Types**

    User Avatar
    http://localhost:7000/user/avatar/3589782/w500.jpg

    Product
    http://localhost:7000/product/image/10855050/x400.jpeg


**Dimensions**

    Maximum Width
    http://localhost:7000/user/avatar/3589782/w500.jpg

    Square
    http://localhost:7000/user/avatar/3589782/x600.jpg

    Rectangle
    http://localhost:7000/user/avatar/3589782/300x400.jpg

**Quality**

The default compression of the image can modified by appending `-q` and the desired quality `1-100`.

    Square with quality 50
    http://localhost:7000/user/avatar/3589782/x600-q50.jpg

## Error Handling

Few errors will cause the server to return error pages

- Source image is not found: NotFound (404)
- Image requested is larger than maximum_width: NotAcceptable (406)

## Development

Prerequisites

```bash
brew install imagemagick --with-webp
```

To download dependencies
```bash
make deps
```

## Compilation

`make` will build the executable under `./bin`
```bash
make
```

## Graphite Events

A local cache was not found and the image was processed. This also tracks count of images sent to remote store.
```
stats.image_server.image_request
```

In addition, the format is tracked (jpg, gif, webp)
```
stats.image_server.image_request.jpg
```

Request failed to return an image
```
stats.image_server.image_request_fail
```

Every download from original source, and a 404 was returned
```
stats.image_server.original_downloaded
```

The original image is not available, and a 404 was returned
```
stats.image_server.original_unavailable
```

## Benchmarks

Make sure your computer can handle enough simultaneous connections. MacOS X by default allows 128. Need a lot more!

```shell
$ sudo sysctl -w kern.ipc.somaxconn=2048
```

Also need to increase the limit of maximum open files

To find out the limits on your computer:
```shell
launchctl limit
```

Increase the limits!
```shell
launchctl limit maxfiles 400000 1000000
```

to increase them [permanently](https://coderwall.com/p/lfjoaq)

## Pending

### Required
- Optimize image generation. Make files smaller. Want to replicate all current configurations.
- Default background color needs to be white (for transparent gifs, etc)

### Operations/Deployment
- deploy to smartos
- chef cookbook for deployment
- configure fastly

### After Release

- Zero-downtime restart: http://rcrowley.org/talks/strange-loop-2013.html#27
- status page
  - current images processing count
  - current original download count
- keep track of the image dimension statistics by image type, dimension, and extension (product image, user avatar). When an image is requested, other popular sizes can be generated in the background after the request. Images created on the background should not count towards statistics.

### Needs discussion

- Ability to have manta jobs to create new versions. First list all product Ids already stored in manta. Map: split into smaller batches. A CLI will be needed for this.
- Split into subdirectories in manta? Currently all product Ids are on one level. Difficult to list files.

### Done
- ~~only allow whitelisted formats: jpg, png, webp~~
- ~~Allow to have variable compressions: x50-q60.jpg~~
- ~~graphite events (https://github.com/marpaia/graphite-golang)~~
  - ~~image processed~~
  - ~~image processed by extension~~
  - ~~original downloaded~~
  - ~~failed downloading from s3~~
  - ~~completed with errors~~
- ~~Limit the number of simultaneous manta uploads. Channels can be used instead of go routines.~~
- ~~Move default compression to configuration~~
- ~~save processed images into manta~~
- ~~error handling [done]~~
- ~~accept flags with environment `-e production` [done]~~
  - ~~default will be `development`~~
- ~~ability to overwrite the source of the image [done]~~
  - ~~by passing query parameter `source`~~
- ~~configuration options/file~~
  - ~~port number [done]~~
  - ~~status page port number [done]~~
  - ~~max dimensions [done]~~
  - ~~source domain [done]~~
  - ~~whitelist image extensions [done]~~
