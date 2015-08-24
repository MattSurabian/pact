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
clone git go.googlesource.com/crypto 81bf7719a6b7ce9b665598222362b50122dfc13b golang.org/x/crypto

clone git github.com/mattsurabian/msg d2a9d565127023e7a954262d624c13eb97e8c564
clone git github.com/spf13/cobra e4993076d845b7127c760e2e57a4984c166c6c05
clone git github.com/spf13/viper 9fca10189b1307bba68b2cd487dd93da0bfbda06
clone git github.com/andrew-d/go-termutil 009166a695a2f516c749a26b4ac1f183d89aa336
clone git github.com/mitchellh/go-homedir 56f508a88415ab57e596a176f0789ede8f790903

## Nested Deps
clone git github.com/BurntSushi/toml 056c9bc7be7190eaa7715723883caffa5f8fa3e4
clone git github.com/inconshreveable/mousetrap 76626ae9c91c4f2a10f34cad8ce83ea42c93bb75
clone git github.com/kr/pretty e6ac2fc51e89a3249e82157fa0bb7a18ef9dd5bb
clone git github.com/magiconair/properties 359442d561ca28acd0fe503aa9f075f505bc9ed0
clone git github.com/mitchellh/mapstructure 2caf8efc93669b6c43e0441cdc6aed17546c96f3
clone git github.com/spf13/cast 4d07383ffe94b5e5a6fa3af9211374a4507a0184
clone git github.com/spf13/jwalterweatherman 3d60171a64319ef63c78bd45bd60e6eab1e75f8b
clone git github.com/spf13/pflag 67cbc198fd11dab704b214c1e629a97af392c085
clone git github.com/kr/text e373e137fafd8abd480af49182dea0513914adb4
clone git gopkg.in/yaml.v2 7ad95dd0798a40da1ccdff6dff35fd177b5edf40
clone git github.com/cpuguy83/go-md2man 71acacd42f85e5e82f70a55327789582a5200a90
clone git github.com/russross/blackfriday 8cec3a854e68dba10faabbe31c089abf4a3e57a6
clone git github.com/shurcooL/sanitized_anchor_name 244f5ac324cb97e1987ef901a0081a77bfd8e845