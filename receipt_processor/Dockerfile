FROM golang:1.23

# Set the working directory
WORKDIR /app

# Copy files to the container
COPY . .

# Install dependencies and build the app
RUN go build -o main .

# Expose the application on port 8080
EXPOSE 8080

# Command to run the application
CMD ["./main"]
