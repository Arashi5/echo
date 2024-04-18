protoc  \
-Iscripts \
-Iapi/protobuf-spec \
--plugin=protoc-gen-custom=protoc-gen-custom \
--custom_out=./ \
api/protobuf-spec/*.proto
