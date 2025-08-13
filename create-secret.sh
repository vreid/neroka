#!/usr/bin/env bash

if [ -f .env ]; then
    set -a
    # shellcheck source=/dev/null
    source .env
    set +a
fi

kubectl delete secret neroka-secret || true
kubectl create secret generic neroka-secrets \
    --from-literal="mastodon_server=${MASTODON_SERVER}" \
    --from-literal="mastodon_access_token=${MASTODON_ACCESS_TOKEN}" \
    --from-literal="anthropic_api_key=${ANTHROPIC_API_KEY}" \
    --from-literal="anthropic_base_url=${ANTHROPIC_BASE_URL}" \
    --from-literal="deepseek_api_key=${DEEPSEEK_API_KEY}" \
    --from-literal="deepseek_base_url=${DEEPSEEK_BASE_URL}" \
    --from-literal="mistral_api_key=${MISTRAL_API_KEY}" \
    --from-literal="mistral_base_url=${MISTRAL_BASE_URL}" \
    --from-literal="openrouter_api_key=${OPENROUTER_API_KEY}" \
    --from-literal="openrouter_base_url=${OPENROUTER_BASE_URL}"
