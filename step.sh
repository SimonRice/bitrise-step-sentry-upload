#!/bin/bash
set -ex

sentry-cli --auth-token ${sentry_auth_token} --url ${sentry_url} upload-dif --org ${sentry_org_slug} --project ${sentry_project_slug} ${sentry_dsym_path}
