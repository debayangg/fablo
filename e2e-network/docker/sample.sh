#!/bin/bash

set -e

TEST_TMP="$(rm -rf "$0.tmpdir" && mkdir -p "$0.tmpdir" && (cd "$0.tmpdir" && pwd))"
TEST_LOGS="$(mkdir -p "$0.logs" && (cd "$0.logs" && pwd))"
FABLO_HOME="$TEST_TMP/../../.."

export FABLO_HOME

FABLO_YAML="$FABLO_HOME/samples/sample.yaml"

# Step 1: Check if config file exists
if [[ ! -f "$FABLO_YAML" ]]; then
  echo "❌ File $FABLO_YAML not found!"
  exit 1
fi

# Step 4: Start the network
echo "🚀 Starting Fablo network..."
"$FABLO_HOME/fablo.sh" generate "$FABLO_YAML"
"$FABLO_HOME/fablo.sh" up

# Step 5: Verify Docker containers
echo "🔎 Checking running containers..."
docker ps --filter "name=org1" --filter "name=orderer" --format "table {{.Names}}\t{{.Status}}"

# Step 6: Test peer0.org1 connectivity (default port 7051)
echo "🌐 Checking peer0.org1..."
if nc -z localhost 7051; then
  echo "✅ peer0.org1 is reachable"
else
  echo "❌ peer0.org1 not reachable"
fi

echo "✅ Fablo network is up and validated."
