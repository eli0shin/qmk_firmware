#! bin/bash
./qmk_macro_mapper/macro_mapper ./bare_layout_34.json ./macros.json ./34_built_layout.json
qmk flash -kb eli_3x5 ./34_built_layout.json
