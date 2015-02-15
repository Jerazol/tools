#!/bin/bash
/usr/bin/find $1 -type f -exec stat --printf "%s\n" "{}" \;