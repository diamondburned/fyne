{ pkgs ? import <nixpkgs> {} }:

let go = pkgs.go.overrideAttrs (old: {
	version = "1.17";
	src = builtins.fetchurl {
		url    = "https://golang.org/dl/go1.17.linux-amd64.tar.gz";
		sha256 = "sha256:0b9p61m7ysiny61k4c0qm3kjsjclsni81b3yrxqkhxmdyp29zy3b";
	};
	doCheck = false;
	patches = [
		# cmd/go/internal/work: concurrent ccompile routines
		(builtins.fetchurl "https://github.com/diamondburned/go/commit/4e07fa9fe4e905d89c725baed404ae43e03eb08e.patch")
		# cmd/cgo: concurrent file generation
		(builtins.fetchurl "https://github.com/diamondburned/go/commit/432db23601eeb941cf2ae3a539a62e6f7c11ed06.patch")
	];
});

in pkgs.mkShell {
	name = "fyne";

	buildInputs = with pkgs; [
		glfw
		glfw-wayland
		mesa
		wayland
		xorg.libX11
		xorg.libXcursor
		xorg.libXi
		xorg.libXinerama
		xorg.libXrandr
		libxkbcommon
		xlibs.libXext
	];

	nativeBuildInputs = [
		pkgs.pkg-config
		go
	];

	CGO_ENABLED = "1";
}
