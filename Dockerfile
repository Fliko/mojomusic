FROM heroku/heroku:18-build as build

COPY . /app
WORKDIR /app

# Setup buildpack
RUN mkdir -p /tmp/buildpack/heroku/go /tmp/build_cache /tmp/env
RUN curl https://codon-buildpacks.s3.amazonaws.com/buildpacks/heroku/go.tgz | tar xz -C /tmp/buildpack/heroku/go
#
# #Execute Buildpack
RUN STACK=heroku-18 /tmp/buildpack/heroku/go/bin/compile /app /tmp/build_cache /tmp/env
#
FROM golang:1.12.1-stretch
RUN go get -u github.com/bwmarrin/dca/cmd/dca
# # Prepare final, minimal image
FROM heroku/heroku:18
#
COPY --from=build /app /app
ENV HOME /app
WORKDIR /app
RUN useradd -m heroku
USER heroku
RUN sudo apt-get install -y youtube-dl
CMD /app/bin/mojoMusic
