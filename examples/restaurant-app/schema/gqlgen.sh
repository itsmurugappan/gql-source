#!/bin/bash

go get github.com/99designs/gqlgen
cd examples/restaurant-app/schema
gqlgen schema.graphql
