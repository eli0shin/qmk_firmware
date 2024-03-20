#! bin/bash

qmk compile ./dactyl_layout.json
qmk flash .build/handwired_dactyl_manuform_4x5_handwired_dactyl_manuform_4x5_layout_mine.hex
