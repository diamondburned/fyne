package xkbcommon

// #cgo pkg-config: xkbcommon
// #include <xkbcommon/xkbcommon.h>
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

// KeyMap is the xkbcommon key map.
type KeyMap struct {
	*keyMap
}

type keyMap struct {
	context *C.struct_xkb_context
	keymap  *C.struct_xkb_keymap
	state   *C.struct_xkb_state
}

var errUnknown = errors.New("xkbcommon error")

// Parse parses the keymap from the given data.
func Parse(keymapData []byte) (*KeyMap, error) {
	context := C.xkb_context_new(C.XKB_CONTEXT_NO_FLAGS)
	if context == nil {
		return nil, errUnknown
	}

	keymap := C.xkb_keymap_new_from_buffer(
		context,
		(*C.char)(unsafe.Pointer(&keymapData[0])),
		C.size_t(len(keymapData)),
		C.XKB_KEYMAP_FORMAT_TEXT_V1,
		C.XKB_KEYMAP_COMPILE_NO_FLAGS)
	if keymap == nil {
		return nil, errUnknown
	}

	state := C.xkb_state_new(keymap)
	if state == nil {
		return nil, errUnknown
	}

	km := &keyMap{
		context, keymap, state,
	}
	runtime.SetFinalizer(km, func(km *keyMap) {
		C.xkb_state_unref(km.state)
		C.xkb_keymap_unref(km.keymap)
		C.xkb_context_unref(km.context)
	})

	return &KeyMap{km}, nil
}

// UTF8 returns the UTF-8 string obtained from pressing a particular key in a
// given keyboard state.
func (k *KeyMap) UTF8(scancode uint32) string {
	var buf [128]byte

	n := C.xkb_state_key_get_utf8(
		k.state,
		C.xkb_keycode_t(scancode),
		(*C.char)(unsafe.Pointer(&buf[0])),
		C.size_t(128))
	if n < 128 {
		return string(buf[:n])
	}

	varbuf := make([]byte, n)
	C.xkb_state_key_get_utf8(
		k.state,
		C.xkb_keycode_t(scancode),
		(*C.char)(unsafe.Pointer(&varbuf[0])),
		C.size_t(n))
	return string(varbuf)
}

// OneSym gets the single keysym obtained from pressing a particular key in a
// given keyboard state.
func (k *KeyMap) OneSym(scancode uint32) uint32 {
	return uint32(C.xkb_state_key_get_one_sym(k.state, C.xkb_keycode_t(scancode)))
}

// UpdateMask updates a keyboard state from a set of explicit masks.
func (k *KeyMap) UpdateMask(
	depressedMods, latchedMods, lockedMods,
	depressedLayout, latchedLayout, lockedLayout uint32,
) {
	C.xkb_state_update_mask(
		k.state,
		C.xkb_mod_mask_t(depressedMods),
		C.xkb_mod_mask_t(latchedMods),
		C.xkb_mod_mask_t(lockedMods),
		C.xkb_layout_index_t(depressedLayout),
		C.xkb_layout_index_t(latchedLayout),
		C.xkb_layout_index_t(lockedLayout),
	)
}
