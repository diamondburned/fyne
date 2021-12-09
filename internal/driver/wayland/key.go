package wayland

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/wayland/internal/xkbcommon"
)

func symName(keysym uint32) fyne.KeyName {
	switch keysym {
	case xkbcommon.KEY_A, xkbcommon.KEY_a:
		return fyne.KeyA
	case xkbcommon.KEY_B, xkbcommon.KEY_b:
		return fyne.KeyB
	case xkbcommon.KEY_C, xkbcommon.KEY_c:
		return fyne.KeyC
	case xkbcommon.KEY_D, xkbcommon.KEY_d:
		return fyne.KeyD
	case xkbcommon.KEY_E, xkbcommon.KEY_e:
		return fyne.KeyE
	case xkbcommon.KEY_F, xkbcommon.KEY_f:
		return fyne.KeyF
	case xkbcommon.KEY_G, xkbcommon.KEY_g:
		return fyne.KeyG
	case xkbcommon.KEY_H, xkbcommon.KEY_h:
		return fyne.KeyH
	case xkbcommon.KEY_I, xkbcommon.KEY_i:
		return fyne.KeyI
	case xkbcommon.KEY_J, xkbcommon.KEY_j:
		return fyne.KeyJ
	case xkbcommon.KEY_K, xkbcommon.KEY_k:
		return fyne.KeyK
	case xkbcommon.KEY_L, xkbcommon.KEY_l:
		return fyne.KeyL
	case xkbcommon.KEY_M, xkbcommon.KEY_m:
		return fyne.KeyM
	case xkbcommon.KEY_N, xkbcommon.KEY_n:
		return fyne.KeyN
	case xkbcommon.KEY_O, xkbcommon.KEY_o:
		return fyne.KeyO
	case xkbcommon.KEY_P, xkbcommon.KEY_p:
		return fyne.KeyP
	case xkbcommon.KEY_Q, xkbcommon.KEY_q:
		return fyne.KeyQ
	case xkbcommon.KEY_R, xkbcommon.KEY_r:
		return fyne.KeyR
	case xkbcommon.KEY_S, xkbcommon.KEY_s:
		return fyne.KeyS
	case xkbcommon.KEY_T, xkbcommon.KEY_t:
		return fyne.KeyT
	case xkbcommon.KEY_U, xkbcommon.KEY_u:
		return fyne.KeyU
	case xkbcommon.KEY_V, xkbcommon.KEY_v:
		return fyne.KeyV
	case xkbcommon.KEY_W, xkbcommon.KEY_w:
		return fyne.KeyW
	case xkbcommon.KEY_X, xkbcommon.KEY_x:
		return fyne.KeyX
	case xkbcommon.KEY_Y, xkbcommon.KEY_y:
		return fyne.KeyY
	case xkbcommon.KEY_Z, xkbcommon.KEY_z:
		return fyne.KeyZ
	case xkbcommon.KEY_0:
		return fyne.Key0
	case xkbcommon.KEY_1:
		return fyne.Key1
	case xkbcommon.KEY_2:
		return fyne.Key2
	case xkbcommon.KEY_3:
		return fyne.Key3
	case xkbcommon.KEY_4:
		return fyne.Key4
	case xkbcommon.KEY_5:
		return fyne.Key5
	case xkbcommon.KEY_6:
		return fyne.Key6
	case xkbcommon.KEY_7:
		return fyne.Key7
	case xkbcommon.KEY_8:
		return fyne.Key8
	case xkbcommon.KEY_9:
		return fyne.Key9
	case xkbcommon.KEY_F1:
		return fyne.KeyF1
	case xkbcommon.KEY_F2:
		return fyne.KeyF2
	case xkbcommon.KEY_F3:
		return fyne.KeyF3
	case xkbcommon.KEY_F4:
		return fyne.KeyF4
	case xkbcommon.KEY_F5:
		return fyne.KeyF5
	case xkbcommon.KEY_F6:
		return fyne.KeyF6
	case xkbcommon.KEY_F7:
		return fyne.KeyF7
	case xkbcommon.KEY_F8:
		return fyne.KeyF8
	case xkbcommon.KEY_F9:
		return fyne.KeyF9
	case xkbcommon.KEY_F10:
		return fyne.KeyF10
	case xkbcommon.KEY_F11:
		return fyne.KeyF11
	case xkbcommon.KEY_F12:
		return fyne.KeyF12
	case xkbcommon.KEY_Escape:
		return fyne.KeyEscape
	case xkbcommon.KEY_Return:
		return fyne.KeyReturn
	case xkbcommon.KEY_Tab:
		return fyne.KeyTab
	case xkbcommon.KEY_BackSpace:
		return fyne.KeyBackspace
	case xkbcommon.KEY_Insert:
		return fyne.KeyInsert
	case xkbcommon.KEY_Delete:
		return fyne.KeyDelete
	case xkbcommon.KEY_Right:
		return fyne.KeyRight
	case xkbcommon.KEY_Left:
		return fyne.KeyLeft
	case xkbcommon.KEY_Down:
		return fyne.KeyDown
	case xkbcommon.KEY_Up:
		return fyne.KeyUp
	case xkbcommon.KEY_Page_Up:
		return fyne.KeyPageUp
	case xkbcommon.KEY_Page_Down:
		return fyne.KeyPageDown
	case xkbcommon.KEY_Home:
		return fyne.KeyHome
	case xkbcommon.KEY_End:
		return fyne.KeyEnd
	case xkbcommon.KEY_space:
		return fyne.KeySpace
	case xkbcommon.KEY_apostrophe:
		return fyne.KeyApostrophe
	case xkbcommon.KEY_comma:
		return fyne.KeyComma
	case xkbcommon.KEY_minus:
		return fyne.KeyMinus
	case xkbcommon.KEY_period:
		return fyne.KeyPeriod
	case xkbcommon.KEY_slash:
		return fyne.KeySlash
	case xkbcommon.KEY_backslash:
		return fyne.KeyBackslash
	case xkbcommon.KEY_leftanglebracket:
		return fyne.KeyLeftBracket
	case xkbcommon.KEY_rightanglebracket:
		return fyne.KeyRightBracket
	case xkbcommon.KEY_semicolon:
		return fyne.KeySemicolon
	case xkbcommon.KEY_equal:
		return fyne.KeyEqual
	case xkbcommon.KEY_asterisk:
		return fyne.KeyAsterisk
	case xkbcommon.KEY_plus:
		return fyne.KeyPlus
	case xkbcommon.KEY_grave:
		return fyne.KeyBackTick
	default:
		return fyne.KeyUnknown
	}
}
