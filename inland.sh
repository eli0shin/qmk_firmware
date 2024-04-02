#! bin/bash

qmk compile ./inland_layout.json
qmk flash .build/inland_mk47_inland_mk47_eli.hex
