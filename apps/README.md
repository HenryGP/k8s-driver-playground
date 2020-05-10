# App instructions

The small code samples inside this directory are meant to be used for testing an app running inside a Docker container. 


## Uploading the built image to Docker Hub
1. Login to Docker Hub using the CLI:
   ```
   docker login
   ```
1. When building the image, set a tag that contains your Docker Hub username as a prefix:
   ```
   docker build -t <docker_username>/<image_name>:<version>
   ```
   For example:
   ```
   docker build -t enrique/demo-app:latest
   ```
1. Push the image to Docker Hub:
   ```
   docker push <docker_username>/<image_name>:<version>
   ```
   Matching the tag from the previous step. If the repository doesn't exist already, it will be created in Docker Hub and associated to your account.
1. View the image details by navigating to the following URL:
   ```
   https://hub.docker.com/repository/docker/<username>/<image_name>
   ```

## Uploading the built image into a registry in AWS ECR
1. Create a registry entry in AWS ECR using the `aws` CLI:
   ```
   aws ecr create-repository --repository-name <repository_name>
   ```
   The resulting output should be similar to the following:
   ```
   {
    "repository": {
        "repositoryArn": "arn:aws:ecr:<aws-region>:<registry-id>:repository/docker-demo",
        "registryId": <registry-id>,
        "repositoryName": <repository_name>,
        "repositoryUri": "<registry-id>.dkr.ecr.<aws-region>.amazonaws.com/docker-demo",
        "createdAt": "2020-04-28T11:56:48+01:00",
        "imageTagMutability": "MUTABLE",
        "imageScanningConfiguration": {
            "scanOnPush": false
        }
      }
    }
   ```
1. Build the image in your local Docker Daemon by running `docker build` make sure to assign a tag:
   ```
   docker build <path_to_Dockerfile> -t <registry-id>.dkr.ecr.<aws-region>.amazonaws.com/<repository-name>:<version>
   ```
   `registry-id`, `aws-region` and `repository-name` are obtained from the previous step of this procedure. `version` will be an additional tag for versioning the image. 
1. Authenticate with ECR using Docker:
   ```
   aws ecr get-login-password --region us-east-1 | docker login --username `aws-username` --password-stdin <registry-id>.dkr.ecr.<aws-region>.amazonaws.com/<repository-name>
   ```
1. Push the image to the repository:
   ```
   docker push <registry-id>.dkr.ecr.<aws-region>.amazonaws.com/<repository-name>:<version>
   ```
1. Verify that the image has been successfully pushed to ECR:
   ```
   aws ecr list-images --repository-name <repository-name>
   ```
   The result should be similar to the following:
   ```
   {
     "imageIds": [
        {
            "imageDigest": "sha256:##########",
            "imageTag": <version>
        }
     ]
   }
   ```            