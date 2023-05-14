# Copyright (c) Spectro Cloud
# SPDX-License-Identifier: Apache-2.0

resource "random_password" "password" {
  length           = 12
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}