#!/bin/bash
# https://gist.github.com/judy2k/7656bfe3b322d669ef75364a46327836
export $(egrep -v '^#' .env | xargs)
