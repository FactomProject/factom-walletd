# factom-walletd Docker Helper

The factom-walletd Docker Helper is a simple tool to help build and run factom-walletd as a container

## Prerequisites

You must have at least Docker v17 installed on your system.

Having this repo cloned helps too 😇

## Build
From wherever you have cloned this repo, run

`docker build -t factom-walletd_container .`

(yes, you can replace **factom-walletd_container** with whatever you want to call the container.  e.g. **factom-walletd**, **foo**, etc.)

#### Cross-Compile
To cross-compile for a different target, you can pass in a `build-arg` as so

`docker build -t factom-walletd_container --build-arg GOOS=darwin .`

## Run
#### No Persistence
**WARNING: DO NOT USE FOR REAL_WORLD STUFF** : When this container goes down, *you lose all data*. This is, literally, for testing purposes.

`docker run --rm -p 8089:8089 factom-walletd_container`
  
* This will start up **factom-walletd** with no flags.
* **When the container terminates, all data will be lost**
* **Note** - In the above, replace **factom-walletd_container** with whatever you called it when you built it - e.g. **factom-walletd**, **foo**, etc.

#### With Persistence
1. `docker volume create factom-walletd_volume`
2. `docker run --rm -v $(PWD)/factomd.conf:/source -v factom-walletd_volume:/destination busybox /bin/cp /source /destination/factomd.conf`
3. `docker run --rm -p 8089:8089 -v factom-walletd_volume:/root/.factom/m2 factom-walletd_container`

* This will start up **factom-walletd** with no flags.
* When the container terminates, the data will remain persisted in the volume **factom-walletd_volume**
* The above copies **factom-walletd.conf** from the local directory into the container. Put _your_ version in there, or change the path appropriately.
* **Note**.  In the above
   * replace **factom-walletd_container** with whatever you called it when you built it - e.g. **factom-walletd**, **foo**, etc.
   * replace **factom-walletd_volume** with whatever you might want to call it - e.g. **myvolume**, **barbaz**, etc.

#### Additional Flags
In all cases, you can startup with additional flags by passing them at the end of the docker command, e.g.

`docker run --rm -p 8089:8089 factom-walletd_container -p 9999`


## Copy
So yeah, you want to get your binary _out_ of the container. To do so, you basically mount your target into the container, and copy the binary over, like so


`docker run --rm --entrypoint='' -v <FULLY_QUALIFIED_PATH_TO_TARGET_DIRECTORY>:/destination factom-walletd_container /bin/cp /go/bin/factom-walletd /destination`

e.g.

`docker run --rm --entrypoint='' -v /tmp:/destination factom-walletd_container /bin/cp /go/bin/factom-walletd /destination`

which will copy the binary to `/tmp/factom-walletd`

**Note** : You should replace ** factom-walletd_container** with whatever you called it in the **build** section above  e.g. **factom-walletd**, **foo**, etc.

#### Cross-Compile
If you cross-compiled to a different target, your binary will be in `/go/bin/<target>/factom-walletd`.  e.g. If you built with `--build-arg GOOS=darwin`, then you can copy out the binary with

`docker run --rm --entrypoint='' -v <FULLY_QUALIFIED_PATH_TO_TARGET_DIRECTORY>:/destination factom-walletd_container /bin/cp /go/bin/darwin_amd64/factom-walletd /destination`

e.g.

`docker run --rm --entrypoint='' -v /tmp:/destination factom-walletd_container /bin/cp /go/bin/darwin_amd64/factom-walletd /destination` 

which will copy the darwin_amd64 version of the binary to `/tmp/factom-walletd`

**Note** : You should replace ** factom-walletd_container** with whatever you called it in the **build** section above  e.g. **factom-walletd**, **foo**, etc.
