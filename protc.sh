#!/bin/bash

protoc mojoroutes/routes.proto --go_out=plugins=grpc:.
