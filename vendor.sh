#!/usr/bin/env bash
set -e

# Downloads dependencies into _vendor/ directory
mkdir -p _vendor
cd _vendor

clone() {
	vcs=$1
	pkg=$2
	rev=$3

	pkg_url=https://$pkg
	target_dir=src/$pkg

	echo -n "Getting dependency -> $pkg @ $rev: "

	if [ -d $target_dir ]; then
		echo -n 'removing old version, '
		rm -fr $target_dir
	fi

	echo -n 'cloning, '
	case $vcs in
		git)
			git clone --quiet --no-checkout $pkg_url $target_dir
			( cd $target_dir && git reset --quiet --hard $rev )
			;;
		hg)
			hg clone --quiet --updaterev $rev $pkg_url $target_dir
			;;
	esac

	echo -n 'removing VCS hidden files, '
	( cd $target_dir && rm -rf .{git,hg} )

	echo done
}

# List Project Dependencies
clone hg code.google.com/p/go.crypto 69e2a90ed92d

clone git github.com/mattsurabian/msg c329a42586fca968e152a235c3a155b10819fa78
clone git github.com/mitchellh/cli 8230c3f351c1efa17429df4e771ab8dcd67ff4bd
clone git github.com/andrew-d/go-termutil 91702f30b7f6d63f574b486457bae6acb1534dce
clone git github.com/rakyll/globalconf 415abc325023f1a00cd2d9fa512e0e71745791a2
clone git github.com/rakyll/goini 907cca0f578a5316fb864ec6992dc3d9730ec58c
