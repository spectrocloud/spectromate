#!/bin/sh
# Copyright (c) Spectro Cloud
# SPDX-License-Identifier: Apache-2.0

/usr/bin/spectromate 2>&1 | tee -a /var/log/spectromate.log
