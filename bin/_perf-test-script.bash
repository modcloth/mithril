#!/usr/bin/env bash

main() {
  case "$1" in
    -h|--help)
      usage
      exit 0
  esac

  if [ -z "$1" ] ; then
    echo "Not enough arguments" >&2
    exit 1
  fi

  local TARGET="$(echo "$1" | sed -E s/^http[s]?:[/]*// | sed -E s/:.*$// | sed -E s/[/]$//)"
  local HTTP_PORT="${2:-"80"}"
  local HTTPS_PORT="${3:-"443"}"

  local EXAMPLE_ROUTE="${ROUTE:-'foo/bar'}"

  curl_command  "http://$TARGET:$HTTP_PORT" 'GET'
  curl_command  "http://$TARGET:$HTTP_PORT/$EXAMPLE_ROUTE" 'POST'
  curl_command  "https://$TARGET:$HTTPS_PORT" 'GET'
  curl_command  "https://$TARGET:$HTTPS_PORT/$EXAMPLE_ROUTE" 'POST'


  for qty in 1 3 10 50 100 250 500 ; do
    wrk_command "http://$TARGET:$HTTP_PORT" "$qty" "$qty" '1s' 'GET'
    wrk_command "http://$TARGET:$HTTP_PORT/$EXAMPLE_ROUTE" "$qty" "$qty" '1s' 'POST'

    wrk_command "https://$TARGET:$HTTPS_PORT" "$qty" "$qty" '1s' 'GET'
    wrk_command "https://$TARGET:$HTTPS_PORT/$EXAMPLE_ROUTE" "$qty" "$qty" '1s' 'POST'
  done

  tar -czf "/tmp/performance_results.tar.gz" "$RESULTS_DIR" >/dev/null 2>&1
}

# curl_command TARGET REQ_TYPE OUTFILE
curl_command() {
  local TARGET="$1"
  local REQ_TYPE="${2:-"GET"}"
  local OUTFILE="${3:-"curl-$RANDOM-$(date +"%s")"}.txt"

  local CMD="curl -w\"@\$CURL_FORMAT_FILE\" -X${REQ_TYPE} -o/dev/null -s -d\"\$(sample_data)\" \"$TARGET\" >> \"$RESULTS_DIR/$OUTFILE\""
  echo "$CMD" > "$RESULTS_DIR/$OUTFILE"
  echo "$CMD"
  eval "$CMD"
}

# wrk_command TARGET CONNECTIONS THREADS DURATION REQ_TYPE OUTFILE
wrk_command() {
  if ! which wrk >/dev/null ; then
    sudo apt-get install libssl-dev -y
    sudo git clone https://github.com/wg/wrk.git
    pushd wrk
    make >/dev/null
    mv {,/usr/local/bin/}wrk
    popd
    rm -rf wrk
  fi

  local TARGET="$1"
  local CONNECTIONS="${2:-"10"}"
  local THREADS="${3:-"1"}"
  local DURATION="${4:-"10s"}"
  local REQ_TYPE="${5:-"GET"}"
  local OUTFILE="${6:-"wrk-$RANDOM-$(date +"%s")"}.txt"

  local CMD="wrk -c${CONNECTIONS} -t${THREADS} -d${DURATION} -M${REQ_TYPE} --body \"\$(sample_data)\" --latency \"$TARGET\" >> \"$RESULTS_DIR/$OUTFILE\""
  echo "$CMD" > "$RESULTS_DIR/$OUTFILE"
  echo "$CMD"
  eval $CMD
}

sample_data() {
  cat <<EOF
{\n
"field1": 1234567890,\n
"field2": "bogus data",\n
"field3": "bogus data",\n
"field4": "bogus data",\n
"field5": "bogus data",\n
"field6": "bogus data",\n
"field7": "bogus data",\n
"field8": "bogus data",\n
"field9": "bogus data",\n
"field10": "bogus data",\n
"field11": "bogus data",\n
"field12": "bogus data",\n
"field13": "bogus data",\n
"field14": "bogus data",\n
"field15": "bogus data"
}
EOF
}

curl_format_file() {
  cat <<EOF
\n
   time_namelookup:  %{time_namelookup}\n
time to first byte:  %{time_starttransfer}\n
      time_connect:  %{time_connect}\n
   time_appconnect:  %{time_appconnect}\n
  time_pretransfer:  %{time_pretransfer}\n
     time_redirect:  %{time_redirect}\n
time_starttransfer:  %{time_starttransfer}\n
------------------------------------------\n
        time_total:  %{time_total}\n
       size_upload:  %{size_upload}\n
\n
EOF
}

export RESULTS_DIR="/tmp/performance_results"
rm -rf "$RESULTS_DIR"
mkdir -p "$RESULTS_DIR"
export CURL_FORMAT_FILE="$(mktemp -t curl_format_file.XXXXX)"
curl_format_file > "$CURL_FORMAT_FILE"

trap "rm -f $CURL_FORMAT_FILE" EXIT SIGTERM SIGINT

# This line should be written programmatically
# main "$@"
