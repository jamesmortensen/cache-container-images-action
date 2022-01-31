# Cache Container Images Action

This GitHub Action caches container images pulled from the Docker Registry. Normally, each time we run a workflow, we must pull images freshly, even if the images change infrequently.

This action solves this problem by first pulling the images and then immediately adding them to a tar archive. The archive is then stored in the cache using @actions/cache. Next time the workflow runs, the container images are instead retrieved from the cache in the tar archive and extracted to the default location for container images.

This action only works for the podman engine. Since both Podman and Docker follow the OCI standard, the same images which we build and run with docker will also run with podman. The learning curve for podman is extremely shallow. The CLI commands work exactly the same as with docker. So `docker run --rm ubuntu:latest` will work just the same as `podman run --rm ubuntu:latest`

## Time (and money) savings

GitHub Action runners for Linux are billed aat $0.008 per minute. The less time it takes to run a workflow, the more we're able to run workflows. We also get faster feedback.

When pulling selenium/standalone-chrome:latest, a 1.3GB container image, my average time savings was 30 seconds per run. When there's a cache miss, on average, it takes 23 seconds longer to archive and store the pulled images in the cache. This is not counting the average 31 seconds to pull the image when there's a cache miss. The break even period comes with just one cached workflow run.

## Limitations

- Only works with podman, not docker.
- Not yet tested on macos-latest. Probably won't work on windows-latest.
- Only works with images pulled from Docker's registry, but this could be expanded to support other registries or images built inside a workflow.

## Why Podman and not Docker

I have been experimenting with different container runtimes, including Podman, which now is becoming a viable free alternative to Docker Desktop on macOS. But the main reason for caching images pulled with podman is that it was just easier to cache images pulled with podman than with docker.

GitHub Action runners come preinstalled with some base Docker images, and this added up to a lot of space and a lot of time to archive them. It will also take some work to extract just the pulled images from /var/lib/docker, as well as repositories.json, and archive only what is pulled.  There is also information about docker load and docker save, which imports and exports tar archives of images, but I haven't looked into this yet.

## Usage

Add this to your workflow file, and replace the images with ones you'll be using in your workflow. In the example below, we plan to cache selenium/node-chrome:latest and selenium/hub:latest:

```yaml
      - name: Cache Container Images
        id: cache-container-images
        uses: jamesmortensen/cache-container-images-action@v1
        with:
          runtime: podman
          images: |
            selenium/node-chrome:latest
            selenium/hub:latest
```

See the [demo-action.yml](https://github.com/jamesmortensen/cache-container-images-action/blob/master/.github/workflows/demo-action.yml) workflow file for a simple example. 

## License 

Copywright (c) James Mortensen, 2022 MIT License
