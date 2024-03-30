#!/usr/bin/env bash

set -e

name="blarg"
path="$1"
binary="$2"

if [[ -z "$path" || -z "$binary" ]]; then
  echo "Usage: $0 <path> <binary>"
  exit 1
fi

# init
mkdir -p "$path"
cd "$path" || exit 1
"$binary" --silent init "$name" "github.com/test/$name"

# add
cd "$name" || exit 1
"$binary" --silent add model release name version
"$binary" --silent add group release --rest -m '*Release'
"$binary" --silent add route release additional -m POST
"$binary" --silent add group hello
"$binary" --silent add route hello world -p funky/world id count:int
"$binary" --silent add route hello new -m POST
"$binary" --silent add model hello world:string count:int
"$binary" --silent add model --struct metric time:time.Time key value job:*Job
"$binary" --silent add model --struct job time:time.Time name external_id:primitive.ObjectID
"$binary" --silent add event jobs event id job:*Job
"$binary" --silent add event reporter metric:*Metric -c '$(NAME).summary.report'
"$binary" --silent add event metrics time:time.Time key value -c 'metrics.report' --receiver -p Metric -t 'reporter'
"$binary" --silent add event flame -r time:time.Time download:float64 upload:float64
"$binary" --silent plugin enable cache
"$binary" --silent add worker ProcessRelease id
"$binary" --silent add queue downloads -c 3
"$binary" --silent add worker process_download id -s '0 0 11 * * *' -q downloads
"$binary" --silent add client go
"$binary" --silent generate
# "$binary" readme
# "$binary" routes
