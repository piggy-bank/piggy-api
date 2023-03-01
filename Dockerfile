FROM golang

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

WORKDIR /src

# Fetch modules first. Module dependencies are less likely to change per build,
# so we benefit from layer caching
ADD ./go.mod ./go.sum* ./
RUN go mod download
# Import the remaining source from the context
COPY . ./
RUN go build -installsuffix cgo -ldflags '-s -w' -o ./app main.go

# Document the service listening port(s)
EXPOSE 8001

# Define the containers executable entrypoint
ENTRYPOINT ["./app"]