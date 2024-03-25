# BTC calc goes API - Why?
Because I buy bitcoin on a weekly basis, hence I need to calculate the value for
each member regularly.

## How does it work
1. Update the `data.json` file with the proper values (use `example_data.json`
   as example)
2. Run the `make check` command

## Build
This build process uses a python slim image, so the resulting image is as small
as possible (less than 200MB).

Build process is as simple as:
1. Update the `Makefile` with your desired docker repo
2. Run `make build`

## Current deployment
API currently available at
[btc-api.borisquiroz.dev](https://btc-api.borisquiroz.dev)
