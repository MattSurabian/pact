# Pact
An experimental cryptographic application for encrypting and decrypting data which implements the experimental cryptographic 
library [MSG](https://github.com/MattSurabian/msg).

## Getting Started
There are a couple of ways to get started with Pact. The easiest method is to [download a compiled binary for your platform](https://github.com/MattSurabian/pact/releases), copy it into your path (try `/usr/bin` or `/usr/local/bin` on Linux/Mac) and just start using it.
If you're interested in compiling it yourself, here's how to do it:

1. Clone this repository anywhere you like.
1. Run `make`. This should create a `pact` binary which you can use directly like so: `./pact`; you can also copy the binary to `/usr/bin/` or `/usr/local/bin` to use it directly.
1. Run `pact config` to generate a config file. This will also generate a keypair if one does not already exist, and configure the `self` pact which will contain your own public key.
1. Running `pact list` will show all your pacts and the public keys they contain. A fresh configuration will only have a `self` pact
1. You should now be able to run the create and read methods! Try it: `pact create self "some data to encrypt" | pact read`

Of note, the `Makefile` and vendoring script are provided for convenience, using them is not mandatory. This package is able to be built using the standard
Golang workflow.

## What Is Pact?
Pact is a CLI application that uses [NaCl](http://nacl.cr.yp.to/) to securely distribute a key capable of decrypting an AES-256-GCM([Galois/Counter Mode](https://en.wikipedia.org/wiki/Galois/Counter_Mode)) ciphertext.
This allows data to be shared securely between many parties without the need for out of band secret sharing, something neither technology is capable of on its own. The ciphertext Pact outputs
is the concatenation of the AES-256-GCM ciphertext with fixed size repeating blocks of NaCl ciphertext which contains the key necessary for decryption of the original message.

If you're curious to take a deeper dive into the encryption, pact offloads all of that logic to the [MSG library](https://github.com/MattSurabian/msg), effort was made to extensively document the [MSG library](https://github.com/MattSurabian/msg) via code comments and test cases.

### Why Not Just Use PGP?
Frankly, you probably should. This project is an experiment aimed at making NaCl easier to use in a multi-party environment. Pact's main benefits are small keys courtesy of NaCl and it's lack of reliance on RSA.
It also aims to be marginally easier to use.

### How Does Pact Secure A Message
When Pact encrypts a message it does so using AES-256 in [Galois/Counter Mode](https://en.wikipedia.org/wiki/Galois/Counter_Mode)
with a [randomly generated nonce and key](https://github.com/MattSurabian/msg/blob/master/entropy.go#L25-L37). 
The AES-256-GCM key used is then encrypted with each pact member's public key (use `pact list` to see a list of pacts and the keys they contain). 
That payload is then prefixed with the [fingerprint of the public key](https://github.com/MattSurabian/msg/blob/master/keys.go#L106-L115) used for encryption, so on decryption the recipient can 
immediately know [which chunk of bytes to decrypt](https://github.com/MattSurabian/msg/blob/master/decrypter.go#L43-L59)
in order to learn the key necessary to decrypt the original message.

### Isn't Combining Cryptographic Methods Insecure?
Combining, yes. Concatenating, no. We assume that both AES-256-GCM and NaCl are PRPs(pseudo-random-permutations) 
or at worst PRFs (pseudo-random-functions); which is to say the output they produce is sufficiently indistinguishable 
from actual random output. The concatenation of two pseudo-random blocks is itself pseudo random. All parallelizable 
crypto algorithms rely on this principal. Pact takes advantage of producing a psuedo-random block which can be intelligently 
sliced appart by an authorized recipient and securely decrypted.

## Usage
The following examples use the default "self" pact which is created on initial configuration, but any pact shown in `pact list` could be used.

### With Files
Since Pact is a CLI tool it plays well with typical console functionality like piping (`|`) and output redirection (`>`) making
file encryption and decryption relatively straightforward:

*Linux/Mac:* 

`cat [path-to-file] | pact create self > file.encrypted` and `cat file.encrypted | pact read > file.decrypted`

*Windows:* 

`type [path-to-file] | pact create self > file.encrypted` and `type file.encrypted | pact read > file.decrypted`


### With Strings
Pact is also capable of reading in a plain text message or ciphertext directly from its arguments:

*Linux/Mac/Windows:* `pact create self "This is a secret message only I can decrypt"` and `pact read "SOME-ENCRYPTED-CIPHER-TEXT"`

### Creating Pacts
Using the `self` pact to encrypt/decrypt data for yourself is all well and good, but eventually you'll want to share data with other people. To do so 
ask that person to download pact and run `pact config` then send the output of `pact key-export` to you.

Use `pact add-key` to create a new pact that contains their key. For this example we're creating a pact called `friends`.

```
pact add-key friends SOME-PUBLIC-KEY
```

or if they send a file:

*Linux/Mac:*

```
cat friendPub.key | pact add-key friends
```

*Windows:*

```
type friendPub.key | pact add-key friends
```

Once a pact is created you can encrypt data such that only members of that pact can decrypt it. Of note, unless you explicitly add your own public key 
to the pact `pact key-export | pact add-key [name-of-pact]` you will not be able to decrypt the ciphertext.

## Available Commands

```
$ ./pact -h
A CLI tool that uses NaCl and AES-256-GCM to facilitate multiparty
communication without the need for out of band secret sharing.
Usage: 
  pact [flags]
  pact [command]

Available Commands: 
  create      Outputs an encrypted ciphertext given a plain text message
  read        Outputs a plain text message given an encrypted ciphertext
  config      Generates a new configuration file
  key-gen     Creates new NaCl keys in the location specified by pact's configuration
  key-export  Outputs the user's public key encoded as base64 to STDOUT
  new         Creates a new pact
  rm          Completely removes an existing pact
  list        Lists existing pacts
  add-key     Adds a key to an existing pact or creates a new pact containing the key
  rm-key      Interactively removes a single key from an existing pact
  help        Help about any command

Flags:
  -h, --help=false: help for pact


Use "pact [command] --help" for more information about a command.

```

### create

```
$ ./pact create --help
Uses AES-256-GCM to encrypt a message with a randomly generated key 
from PBKDF2 and encrypts that secret key with the public key of each 
member of a pact. Base64 encoded encrypted ciphertext is sent to STDOUT.
The plain text can be piped into this command.

Usage: 
  pact create [pact-name] [plain-text]

```

### read

```
$ ./pact read --help
Uses NaCl to decrypt a key which can be used to decrypt the message 
which has been secured with AES-256-GCM encryption. The ciphertext can be piped 
into this command.

Usage: 
  pact read [ciphertext]

Flags:
  -h, --help=false: help for read

```

### config

```
$ ./pact config --help
Generates a new configuration file and will refuse to overwrite an existing one.

Usage: 
  pact config

```

### key-gen

```
$ ./pact key-gen --help
Generates an NaCl keypair and writes their base64
string representation to the paths specified in Pact's configuration.

Usage: 
  pact key-gen

```
### key-export

```
$ ./pact key-export --help                                                                          
Sends the user's public key encoded as base64 to STDOUT for easy distribution

Usage: 
  pact key-export

```

### new

```
$ ./pact new --help
Creates a new pact in the configuration file that keys can be added to with the add-key command

Usage: 
  pact new [pact-name]

```

### rm

```
$ ./pact rm --help                                                                                  
Removes an existing pact and all the keys it contains from the user's configuration file.

Usage: 
  pact rm [pact-name]

```

### list

```
$ ./pact list --help                                                                                
Outputs a list of existing pacts and the keys they contain.

Usage: 
  pact list

```

### add-key

```
$ ./pact add-key --help                                                                             
Adds the provided public key to the specified pact. A new pact will be created if necessary.
The public-key can be piped into this command.

Usage: 
  pact add-key [pact-name] [public-key]

```

### rm-key

```
$ ./pact rm-key --help                                                                              
Removes a single key from an existing pact using interactive prompts.

Usage: 
  pact rm-key [pact-name]

```

## Known Issues
 - MSG is still *very* experimental and has not yet been thoroughly peer reviewed. Every effort was made to correctly utilize NaCl and AES-256-GCM, but until it's reviewed it shouldn't 
 be trusted with anything critical.
 - At present, a user's NaCl keypair cannot be secured with a passphrase the way an RSA key can. If a user loses control of their keys they also lose control over any data those keys protect.
 
## Contributing

This repo is still very much experimental, so the more the merrier. While a `Makefile` and vendoring 
script are provided for user convenience it's recommended that contributors clone this into their 
Gopath per the standard Go workflow (`$GOPATH/src/github/mattsurabian/pact`). Contributing to Go projects 
from a fork can be more complicated than project's developed in other languages. Fortunately there are blog 
posts on the subject, like [Katarina Owen's piece about Contributing to Open Source Git Repositories in Go](https://splice.com/blog/contributing-open-source-git-repositories-go/).