#!/usr/bin/env sh
if [ ! -n "$PROMTAIL_YAML_BASE64" ]; then
  echo "Expecting promtail.yaml base64 encoded in PROMTAIL_YAML_BASE64"
  exit 1
fi

echo "$PROMTAIL_YAML_BASE64" | base64 -d > /promtail/promtail.yaml

/sidecars/bin/promtail -config.file=/promtail/promtail.yaml
