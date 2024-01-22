# Cache Container Images Action

This GitHub Action caches container images pulled from the Docker Registry. Normally, each time we run a workflow, we must pull images freshly, even if the images change infrequently.

This action solves this problem by first pulling the images and then immediately adding them to a tar archive. The archive is then stored in the cache using @actions/cache. Next time the workflow runs, the container images are instead retrieved from the cache in the tar archive and extracted to the default location for container images.

This action only works for the podman engine. Since both Podman and Docker follow the OCI standard, the same images which we build and run with docker will also run with podman. The learning curve for podman is extremely shallow. The CLI commands work exactly the same as with docker. So `docker run --rm ubuntu:latest` will work just the same as `podman run --rm ubuntu:latest`

## Time (and money) savings

GitHub Action runners for Linux are billed at $0.008 per minute. The less time it takes to run a workflow, the more we're able to run workflows. We also get faster feedback.

When pulling selenium/standalone-chrome:latest, a 1.3GB container image, my average time savings was 30 seconds per run. When there's a cache miss, on average, it takes 23 seconds longer to archive and store the pulled images in the cache. This is not counting the average 31 seconds to pull the image when there's a cache miss. The break even period comes with just one cached workflow run.

## Limitations

- Only works with podman, not docker.
- Not yet tested on macos-latest. Probably won't work on windows-latest.
- Only works with images pulled from Docker's registry, but this could be expanded to support other registries or images built inside a workflow.

## Why Podman and not Docker

I have been experimenting with different container runtimes, including Podman, which now is becoming a viable free alternative to Docker Desktop on macOS. But the main reason for caching images pulled with podman is that it was just easier to cache images pulled with podman than with docker.

GitHub Action runners come preinstalled with some base Docker images, and this added up to a lot of space and a lot of time to archive them. It will also take some work to extract just the pulled images from /var/lib/docker, as well as repositories.json, and archive only what is pulled.  There is also information about docker load and docker save, which imports and exports tar archives of images, but I haven't looked into this yet.

## Usage

Add this to your workflow file, and replace the required field, images, with ones you'll be using in your workflow. In the example below, we plan to cache selenium/node-chrome:latest and selenium/hub:latest, so we declare this in the images block:

```yaml
      - name: Cache Container Images
        id: cache-container-images
        uses: jamesmortensen/cache-container-images-action@v1
        with:
          images: |
            selenium/node-chrome:latest
            selenium/hub:latest
```

When new images are pushed to the container registry at Docker Hub, then when the action runs, it will pull fresh images. As long as you're using the "latest" tag, the system will pull the latest images and then cache them until new ones are pushed to the registry.  If you use immutable tags, then the images will theoretically remain cached until something happens to clear them from [actions/cache](https://github.com/actions/cache).

If you need to force a cache flush, change the prefix-key, 'podman-cache' by default, to any other value:

```yaml
      - name: Cache Container Images
        id: cache-container-images
        uses: jamesmortensen/cache-container-images-action@v1
        with:
          prefix-key: 'afdafdasfds'  # Optional: this can be anything you want. Change it to force a cache flush.
          images: |
            selenium/node-chrome:4.1.2-20220130
            selenium/node-firefox:4.1.2-20220130
            selenium/node-edge:4.1.2-20220130
            selenium/hub:4.1.2-20220130
```

In your workflows, you may want to perform some actions only in the event of a cache hit or a cache miss. You can do this like so:

```yaml
      - name: Run this step only if container images were found in the cache
        if: ${{ steps.cache-container-images.outputs.cache-hit == 'true' }}
        shell: bash
        run: echo "Container images were found in the cache..."

      - name: Run this step only if container images were NOT found in the cache
        if: ${{ steps.cache-container-images.outputs.cache-hit != 'true' }}
        shell: bash
        run: echo "Container images were NOT found in the cache..."
```

NOTE:  Be sure to use `!= 'true'` as actions/cache does not set any cache-hit value if there's a cache miss.


See the [test-action.yml](https://github.com/jamesmortensen/cache-container-images-action/blob/master/.github/workflows/demo-action.yml) workflow file for test results and examples. 

## License 

Copyright (c) James Mortensen, 2022, 2024 MIT License
