#!/usr/bin/env bash
protoc calcpb/calc.proto --go_out=plugins=grpc:.