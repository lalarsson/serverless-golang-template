#!/usr/bin/env bash

dep ensure

echo "Compiling functions to bin/ ..."

rm -rf bin/

cd src/
for f in **/*[!_test].go; do
  filename="${f%.go}"
  if CGO_ENABLED=1 go build -race -o "../bin/$filename" ${f}; then
    echo "✓ Compiled $filename"
  else
    echo "✕ Failed to compile $filename!"
    exit 1
  fi
done

echo "Done."
