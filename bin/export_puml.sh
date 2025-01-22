#!/bin/bash

if [[ -z "$(which plantuml)" ]]; then
	echo "plantuml is not installed"
	exit 1
fi

root_dir="$(realpath "$(dirname "$(dirname "$0")")")"
src_dir=""$root_dir"/docs/design"
dest_dir=""$root_dir"/.dist/docs/design"

# Delete all output files in destination directory
find "$dest_dir" -name "*.png" -name "*.md" -exec rm -f {} +

# Remove empty directories
find "$dest_dir" -type d -empty -delete

# Generate output files
find "$src_dir" -type f -name "*.puml" | while read -r file; do
	# Determine the relative path of the file
	rel_path="${file#$src_dir/}"
	# Determine the target directory for the exported diagram
	target_dir="$dest_dir/$(dirname "$rel_path")"
	# Create the target directory if it doesn't exist
	mkdir -p "$target_dir"
	# Run PlantUML and export the diagram to the target directory
	echo "Exporting file: "$file""
	plantuml -o "$target_dir" "$file" "$([[ "$1" = "true" ]] && echo '-darkmode')" && echo "Success"
	echo -e "\n==========================\n"
done

# Copy markdown files to destination
find "$src_dir" -type f -name "*.md" | while read -r file; do
    # Compute the relative path of the file
    rel_path="${file#$src_dir/}"
    # Determine the target directory
    target_dir="$dest_dir/$(dirname "$rel_path")"
    # Create the target directory if it doesn't exist
    mkdir -p "$target_dir"
    # Copy the file to the target directory
    cp "$file" "$target_dir"
done
