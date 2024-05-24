#! bin/bash
./qmk_macro_mapper/macro_mapper ./bare_layout.json ./macros.json ./dactyl_built_layout.json
qmk flash ./dactyl_built_layout.json

