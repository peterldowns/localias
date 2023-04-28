# This is based on the nixpkgs repository's `compose-xcodewrapper.nix` but it
# omits the iOS SDK's and simulator. It basically serves to ensure that
# there is an xcodebuild that works inside the developer shell.
{ stdenv, lib }:
{ version ? "11.1"
, allowHigher ? false
, xcodeBaseDir ? "/Applications/Xcode.app"
}:

assert stdenv.isDarwin;

stdenv.mkDerivation {
  pname = "xcode-wrapper${lib.optionalString allowHigher "-plus"}";
  inherit version;
  buildCommand = ''
    mkdir -p $out/bin
    cd $out/bin
    ln -s /usr/bin/xcode-select
    ln -s /usr/bin/security
    ln -s /usr/bin/codesign
    ln -s /usr/bin/xcrun
    ln -s /usr/bin/plutil
    ln -s /usr/bin/clang
    ln -s /usr/bin/lipo
    ln -s /usr/bin/file
    ln -s /usr/bin/rev
    ln -s "${xcodeBaseDir}/Contents/Developer/usr/bin/xcodebuild"
    cd ..
    # Check if we have the xcodebuild version that we want
    currVer=$($out/bin/xcodebuild -version | head -n1)
    ${if allowHigher then ''
      if [ -z "$(printf '%s\n' "${version}" "$currVer" | sort -V | head -n1)""" != "${version}" ]
    '' else ''
      if [ -z "$(echo $currVer | grep -x 'Xcode ${version}')" ]
    ''}
    then
        echo "We require xcodebuild version${
          if allowHigher then " or higher" else ""
        }: ${version}"
        echo "Instead what was found: $currVer"
        exit 1
    fi
  '';
}
