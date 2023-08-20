# Use the official k6 image as the base image
FROM grafana/k6:latest

# Set the working directory inside the container
WORKDIR /scripts

# Copy your k6 script(s) into the container
COPY ./src/ ./

# # Run the k6 script when the container starts
# CMD ["run", "create-transaction.ts"]