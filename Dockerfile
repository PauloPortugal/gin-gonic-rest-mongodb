FROM golang:1.16.5 as development

WORKDIR /app
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy app files
COPY . .
# Install Reflex for development
RUN go install github.com/cespare/reflex@latest
# Expose port
EXPOSE 8080
EXPOSE 8081

