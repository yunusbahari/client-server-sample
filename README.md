# Client Server Application Sample

Simple client-server app in golang with docker & k8s deployment example

## Requirement for Development

- go ~> 1.11.5
- make
- docker

## Build image

### Client

- compile the source
    ```
    cd client/
    make compile
    ```

- build the image (go to root folder project)

    `docker build -t 'your-image-tag/client:v1' client/`

### Server

- compile the source

    ```
    cd server/
    make compile
    ```

- build the image (go to root folder project)

    `docker build -t 'your-image-tag/server:v1' server/`

## Running Application (via Docker)

    `docker run --rm -d -p 9999:9999 -p 13000:13000 your-image-tag:v1`

### Client

    `docker run --rm -d -p 9998:9998 'your-image-tag/client:v1'`

### Server

    `docker run --rm -d -p 9999:9999 -p 13000:13000 your-image-tag/server:v1`

*Parameter explanation*
- --rm `cleanup container when it stopped`
- -d `run the container in the background`
- -p <host>:<container> `forward traffic from host port to container port`
