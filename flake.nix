{
  # TODO: change this when creating a new flake
  description = "Nix Flake-based FHS development environment for protho";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }: flake-utils.lib.eachSystem flake-utils.lib.allSystems (system:
    let
      pkgs = import nixpkgs { inherit system; };
    in
    rec {
      # TODO: change these when creating a new flake
      flakeName = "protho";
      flakeIsFHS = false;
      flakePackages = pkgs: with pkgs; [
        bash
        git
        go
      ];
      flakeShellHook = ''
        go version
      '';

      devShell = pkgs.mkShell rec {
        fhsName = "${flakeName}-fhs-env";
        fhsScriptName = "${flakeName}-fhs-env-script";

        # this is needed to debug:
        # https://nixos.wiki/wiki/Go
        hardeningDisable = [ "fortify" ];
        packages = [
          # if it's an FHS environment, we need to start bash again manually
          (pkgs.writeShellScriptBin fhsScriptName (flakeShellHook + (if flakeIsFHS then "bash" else "")))
        ] ++ (if flakeIsFHS then [
          (pkgs.buildFHSEnv {
            name = fhsName;
            targetPkgs = flakePackages;

            # this whole shenanigan was made because putting
            # the raw commands in the "runScript" just wouldn't
            # work
            runScript = fhsScriptName;
          })
        ] else flakePackages pkgs);

        # TODO: all this mess creates a level-4 shell (`echo $SHLVL`)
        shellHook = ''
          ${if flakeIsFHS then fhsName else fhsScriptName}
          # the "exit 0" is necessary for not having to type in
          # "exit" another time after quitting the nix develop;
          # this is only necessary when using the FHS environment
          ${if flakeIsFHS then "exit 0" else ""}
        '';
      };
    }
  );
}
