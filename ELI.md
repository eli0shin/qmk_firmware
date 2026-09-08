# Eli's keyboard layouts

Run all commands from the repository root.

## 34-key split

- Keyboard: `eli_3x5`
- Edit the key layout: [`bare_layout_34.json`](bare_layout_34.json)
- Edit shared macros: [`macros.json`](macros.json)
- Generate the firmware configuration and flash the keyboard:

  ```sh
  ./build_34.sh
  ```

The command generates `34_built_layout.json`. Do not edit that generated file.

## 46-key Dactyl Manuform

- Keyboard: `handwired/dactyl_manuform/4x5`
- Edit the key layout: [`bare_layout.json`](bare_layout.json)
- Edit shared macros: [`macros.json`](macros.json)
- Generate the firmware configuration and flash the keyboard:

  ```sh
  ./build.sh
  ```

The command generates `dactyl_built_layout.json`. Do not edit that generated file.

## Layout file format

Layers are objects with explicit number keys. Each layer separates the left and right hands into visual rows:

```json
"layers": {
  "4": {
    "left": { "rows": [["KC_NO"], ["KC_LCTL"]] },
    "right": { "rows": [["KC_1"], ["KC_4"]] }
  }
}
```

Use raw QMK keycodes such as `KC_LGUI`, `KC_TRNS`, and `LT(1,KC_ESC)`. Layer numbers must be contiguous and start at `0`. Both hands must have matching row counts and row widths.

The Dactyl layout omits its final two unusable rows. The macro mapper uses `unused_rows` to add those positions as `KC_NO` in the generated QMK file.

## Macros

Both build scripts run the macro mapper from `qmk_macro_mapper`. In a layout file, use `ANY(M*KEYCODE)` to refer to a macro whose `keycode` is defined in `macros.json`.
