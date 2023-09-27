# accuknox



# Running Your Go Application with Docker and `docker-compose.yaml`

This guide will walk you through the process of running your Go application using Docker and `docker-compose.yaml`. We assume that the Dockerfiles necessary for your application are already present in the repository.

## Prerequisites

Before you begin, ensure that you have the following prerequisites installed on your system:

- Docker: [Install Docker](https://docs.docker.com/get-docker/)
- Docker Compose: [Install Docker Compose](https://docs.docker.com/compose/install/)

## Clone the Repository

Clone your Go application's repository to your local machine.

```bash
git clone git@github.com:faheem-fk/accuknox.git
cd accuknox
```


## Build and Run

To build and run your Go application using Docker and `docker-compose`, follow these steps:

1. Open a terminal and navigate to your project directory containing the `docker-compose.yaml` file.

2. Build the Docker image:

   ```bash
   docker-compose build
   ```

3. Run the Docker container:

   ```bash
   docker-compose up
   ```

   This will start your Go application within a Docker container.

4. Access your Go application in a web browser or via an HTTP client:

   ```
   http://localhost:8080
   ```

   Your Go application should now be running and accessible at the specified port.

## Stopping the Application

To stop the running Docker container, press `Ctrl+C` in the terminal where it is running, or run the following command in the project directory:

```bash
docker-compose down
```

This will stop and remove the Docker container.

---

That's it! You've successfully set up and run your Go application using Docker and `docker-compose`. You can now easily share and deploy your application as a Docker container.
