#! bin/bash
./qmk_macro_mapper/macro_mapper ./bare_layout.json ./macros.json ./dactyl_built_layout.json
qmk compile ./dactyl_build_layout.json
qmk flash .build/handwired_dactyl_manuform_4x5_handwired_dactyl_manuform_4x5_layout_mine.hex
