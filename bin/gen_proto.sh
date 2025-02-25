#!/bin/bash

root="$(realpath "$(dirname "$(dirname "$0")")")"
source_dir_path="$root/proto/harmonify/movie_reservation_system"
output_dir_paths=("$root/pkg/proto")

proto_files_path=("$source_dir_path"/**/*.proto)
go_opt_args="--go_opt=paths=source_relative"
go_grpc_opt_args="--go-grpc_opt=paths=source_relative"

for proto_file_path in "${proto_files_path[@]}"; do
	proto_dir_name="$(basename "$(dirname "$proto_file_path")" )"
	proto_file_name="$(basename "$proto_file_path")"

	source_relative_proto_file_path=""$proto_dir_name"/"$proto_file_name""
	go_pkg_name=""$proto_dir_name"_proto"

	go_opt_args+=",M"$source_relative_proto_file_path"=./"$go_pkg_name""
	go_grpc_opt_args+=",M"$source_relative_proto_file_path"=./"$go_pkg_name""
done

# clear

for output_dir_path in ${output_dir_paths[@]}; do
	rm -rf "$output_dir_path"
	mkdir -p "$output_dir_path"

	echo -e "Executing: protoc \n \
		--proto_path="${source_dir_path}" \n \
		--go_out="${output_dir_path}" \n \
		--go-grpc_out="${output_dir_path}" \n \
		"${go_opt_args// /}" \n \
		"${go_grpc_opt_args// /}" \n \
		"${proto_files_path[@]}""
	echo
	
	protoc \
		--proto_path="${source_dir_path}" \
		--go_out="${output_dir_path}" \
		--go-grpc_out="${output_dir_path}" \
		"${go_opt_args// /}" \
		"${go_grpc_opt_args// /}" \
		"${proto_files_path[@]}"

	echo -e "\n\n"
done
