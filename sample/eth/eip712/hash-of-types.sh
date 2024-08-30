#!/bin/bash

jj -un < schema.json | tr -d '\n' | openssl dgst -sha224
jj -un < schema2.json | tr -d '\n' | openssl dgst -sha224