# Simple Go Webapp for Kubernetes Demos

This is a very simple webapp written in Go, used for Kubernetes demos. It serves a simple HTML page at the root path: `/`, showing the value of the HOSTNAME environment variable inside a colored bar.

The `/grid` path can generate a grid of iframes displaying the main page. The number of columns and rows in the grid is determined by the `cols` and `rows` query parameters respectively.

For example, to generate a 3x3 grid just use the following path: `http://<hostname>:<port>/grid?cols=3&rows=3`

## Building the Go Executable

On a POSIX-like environment you can simply run the `build.sh` script. It runs the `go build` command with the correct configuration to generate a statically-linked, Linux ELF executable, in the current directory, simply called `main`. Or you can run it directly like this:
```
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
```
Please note that this is not intended to be _installed_ in your local $GOPATH. Instead, it is used as the basis for building a Docker image as explained below.

## Creating a Docker Image

After building the Linux executable just run:
```
docker build -t <your_image_tag> .
```

To upload the new image to a public Docker Registry such as hub.docker.com push the image like this:
```
docker push <your_image_tag>
```

## Creating the Kubernetes Objects

Once your Docker image has been published you can use the sample descriptor files in the `kdemo` folder to create the relevant objects in a Kubernetes cluster. TO USE YOUR OWN IMAGE MAKE SURE YOU REPLACE THE `image` FIELD IN THE `kdemo-dep.yaml` FILE WITH THE TAG YOU USED ABOVE.
```
kubectl create -f kdemo/kdemo-dep.yaml
kubectl create -f kdemo/kdemo-svc.yaml
```
This will create a Deployment and related Service using a `NodePort`. Find the dynamically allocated Node port (with something like `kubectl describe svc/kdemo-svc`) to access the Service from outside the cluster. If you have an Ingress Controller configured in your cluster, you may find the `kdemo-ing.yaml` file useful in creating an Ingress for this Deployment.

## Who do I talk to?

[Isaac "Ike" Arias](mailto:ike@apprenda.com)