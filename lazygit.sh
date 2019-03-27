#!/usr/bin/env bash
cd HG
git add .
git commit -a -m "$1"
git push
