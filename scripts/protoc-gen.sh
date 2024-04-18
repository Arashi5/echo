docker run --rm -v $(pwd):/defs namely/protoc-all \
    -d api/protobuf-spec \
    -i vendor \
    -i scripts \
    -o . \
    -l go

