#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
ruby $DIR/generate_redis_set_operations.rb | redis-cli --pipe
