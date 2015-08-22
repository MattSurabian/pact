#!/usr/bin/env bash
set -e

# Downloads dependencies into _vendor/ directory
mkdir -p _vendor
cd _vendor

clone() {
	vcs=$1
	pkg=$2
	rev=$3
	dest=$4

	pkg_url=https://$pkg
	target_dir=src/$pkg

	if [ -n "$dest" ]; then
	    echo "EMPTY"
		target_dir=src/$dest
	fi

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
	esac

	echo -n 'removing VCS hidden files, '
	( cd $target_dir && rm -rf .{git} )

	echo done
}

# List Project Dependencies
clone git go.googlesource.com/crypto 1351f936d976c60a0a48d728281922cf63eafb8d golang.org/x/crypto

clone git github.com/mattsurabian/msg d2a9d565127023e7a954262d624c13eb97e8c564
clone git github.com/spf13/cobra c55cdf33856a08e4822738728b41783292812889
clone git github.com/spf13/viper 2abb1bebfde865b0bb6bb7ada5be63ec78527fa6
clone git github.com/andrew-d/go-termutil 91702f30b7f6d63f574b486457bae6acb1534dce
