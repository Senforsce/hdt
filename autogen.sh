#!/usr/bin/env bash

aclocal -I m4 --install
autoreconf --install --force
gnulib-tool --import warnings
