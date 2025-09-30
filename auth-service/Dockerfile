#  User Arguments
ARG USERNAME=appuser
ARG USER_UID=1001
ARG USER_GID=1001

# First Stage
FROM golang:1.25-alpine AS builder 

# Pass User Arguments
ARG USERNAME
ARG USER_UID
ARG USER_GID

# Create the user and group 
RUN addgroup -g ${USER_GID} ${USERNAME} && \
    adduser -u ${USER_UID} -G ${USERNAME} -S ${USERNAME}

RUN echo "Building as user ${USERNAME}"

# Install Delve for debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# set build directory
WORKDIR /build

COPY go.mod .
COPY go.sum .  
RUN go mod download

COPY . .

# Build the app with debugging flags
RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o ./auth-service ./cmd/main.go

#Secodn Stage
FROM scratch

# Pass User Arguments
ARG USERNAME
ARG USER_UID
ARG USER_GID

# Coppy user group files
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# set workdirectory
WORKDIR /app
COPY --from=builder /build/auth-service ./auth-service
COPY --from=builder /go/bin/dlv /usr/local/bin/

# Switch to the non-root user
USER ${USERNAME}

CMD [ "/app/auth-service" ]

