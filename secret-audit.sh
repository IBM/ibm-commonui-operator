#!/bin/bash
source ~/py_envs/bin/activate
detect-secrets audit .secrets.baseline
