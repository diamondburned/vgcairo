{ pkgs ? import <nixpkgs> {} }:

let gotk4-nix = pkgs.fetchFromGitHub {
		owner = "diamondburned";
		repo  = "gotk4-nix";
		rev   = "2c031f93638f8c97a298807df80424f68ffaac76";
		hash  = "sha256:0lpbnbzl1sc684ypf6ba5f8jnj6sd8z8ajs0pa2sqi8j9w0c87b0";
	};

in import "${gotk4-nix}/shell.nix" {
	base = {
		pname = "vgcairo";
		version = "dev";
	};
}
