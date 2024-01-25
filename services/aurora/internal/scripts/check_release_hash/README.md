# check_release_hash

Docker image for comparing releases hash to a local builds hash.

## Usage

1. Build the image. Optionally pass `--build-arg` with `golang` image you want to use for compilation. See `Dockerfile` for a default value.
2. `docker run -e "TAG=aurora-vX.Y.Z"  -e "PACKAGE_VERSION=X.Y.Z-BUILD_ID" check_release_hash`
3. Compare the hashes in the output. `released` directory contains packages from GitHub, `dist` are the locally built packages.