#!/bin/bash
source ~/py_envs/bin/activate
detect-secrets scan --update .secrets.baseline --exclude-files ".secrets.baseline|requirements.txt|go.mod|go.sum|pom.xml|build.gradle|package-lock.json|yarn.lock|Cargo.lock|deno.lock|composer.lock|Gemfile.lock|Pipfile.lock"
exit $?
