#!/bin/bash
############################################################################
# Copyright 2020 IBM Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
############################################################################

# Scan source files and check that the copyright header has the correct information.
# Run this script from the repository root dir as "scripts/copyright-check.sh".
# Input: list of files to scan (optional)
#        - this is the name of a file containing a list of file names to be scanned
#        - run "git diff --name-only master..HEAD --diff-filter=d" to build the list of files
#          that have changed in the pull request (PR)
#        - run "git diff --name-only <SHA> HEAD --diff-filter=d" to build the list of files
#          that have changed since the previous release. 
#          <SHA> is the SHA of the last commit in the previous release.
#
#        - Example:
#            scripts/copyright-check.sh ./copyright-list.txt
#            - if no filename is passed to this script, all the files in the repo will be scanned
#            - sample contents for copyright-list.txt:
#                scripts/startup.sh
#                lib/icp-util.js
#                app.js
#
##################################################################################################
# Search for "SET YEAR HERE" to set the current copyright year.
# Search for "SET EXCLUDED DIRS1 HERE" to add to the list of dirs that are not scanned.
# Search for "SET EXCLUDED DIRS2 HERE" to add to the list of dirs that are not scanned.
# Search for "SET FILE_TYPES HERE" to add to the list of filetypes that are scanned.
##################################################################################################

# "Copyright $YEAR IBM Corporation"
# $YEAR can be 1 year (2019) or 2 years separated by a comma and an optional space (2018,2019).
# the first year can be any value.
# ***** SET YEAR HERE ***** (don't set it manually, use CURR_YEAR)
CURR_YEAR=$(date +"%Y")
CHECK1="Copyright ([0-9]{4},)*[ ]*(${CURR_YEAR}) IBM Corporation"
CHECK2="Licensed under the Apache License, Version 2.0"
CHECK3="distributed on an \"AS IS\" BASIS"
CHECK4="WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND"

# array of lines to scan for
LIC_ARY=("$CHECK1" "$CHECK2" "$CHECK3" "$CHECK4")
LIC_ARY_SIZE=${#LIC_ARY[@]}

# Used to signal an exit
ERROR=0

# determine which files should be checked
FILE_NAME=$1
if [ "$FILE_NAME" == "" ] ; then
   # get list of files
   echo "building list of files"
   # Ignore .FILENAME types. ' ! -path "./doc/*" ' identifies a directory to be skipped
   # ***** SET EXCLUDED DIRS1 HERE *****
   FILE_LIST=$(find . -type f ! -iname ".*" ! -path "./.git/*" ! -path "./build/*" ! -path "./common/*" ! -path "./deploy/*" ! -path "./docs/*" ! -path "./hack/*")
else
   # use provided list of files
   echo "using provided list of files"
   FILE_LIST=$(cat "$FILE_NAME")
fi

echo -e "##### Copyright check started #####\n"
# Loop through all files
for f in $FILE_LIST; do
  # get the topmost dir in the file's path
  PARENT_DIR=$(dirname "$f" | cut -d "/" -f1)

  # skip non-existent files (shouldn't happen)
  if [ ! -f "$f" ]; then
    printf " Skipping non-existent file %s . . .\n" "$f"
    continue
  fi

  # ***** SET EXCLUDED DIRS2 HERE *****
  # similar to "EXCLUDED DIRS1" but covers the case where a file list created by 'git diff' was passed in
  case "${PARENT_DIR}" in
  	build | common | deploy | docs | hack)
      printf " Skipping %s because its parent dir is excluded. . .\n" "$f"
      continue
  	  ;;
  	*)
      printf " Checking filetype . . . "
      ;;
  esac

  # If a file doesn't have a file extension, like "Makefile", 
  # FILE_TYPE will be set to the filename so it will be skipped as expected.
  # If a file begins with a dot, like ".dockerignore", 
  # FILE_TYPE will be set to the filename without the dot so it will be skipped as expected.
  FILE_TYPE=$(basename "${f##*.}")
  # ***** SET FILE_TYPES HERE *****
  case "${FILE_TYPE}" in
  	sh | go)
      printf " Scanning %s . . . " "$f"
  	  ;;
  	*)
      printf " Skipping %s . . .\n" "$f"
      continue
      ;;
  esac

  # Read the first 20 lines, most Copyright headers use the first 15 lines.
  HEADER=$(head -20 "$f")

  # Check for all copyright lines
  for i in $(seq 0 $((LIC_ARY_SIZE))); do
    # Add a status message of OK, if all copyright lines are found
    if [ "$i" -eq "${LIC_ARY_SIZE}" ]; then
      printf "OK\n"
    else
      # Validate the copyright line being checked is present
      if ! [[ "$HEADER" =~ ${LIC_ARY[$i]} ]]; then
        printf " Missing copyright\n   >>Could not find %s in the file %s\n" "[${LIC_ARY[$i]}]" "$f"
        ERROR=1
        break
      fi
    fi
  done
done

echo -e "\n##### Copyright check finished, ReturnCode: ${ERROR} #####"
exit $ERROR
